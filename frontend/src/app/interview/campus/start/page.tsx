'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Typography, Row, Col, Card as AntCard, Space, Tag, Button, Input, Avatar, Progress, message } from 'antd';
import { AudioOutlined, CustomerServiceOutlined, QuestionCircleOutlined, SendOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

interface ConversationItem {
  type: 'question' | 'answer';
  content: string;
  index?: number;
  timestamp: number;
}

export default function CampusInterviewStartPage() {
  const [elapsed, setElapsed] = useState(0);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [questionText, setQuestionText] = useState<string>('');
  const [questionIndex, setQuestionIndex] = useState<number>(0);
  const [answeredCount, setAnsweredCount] = useState(0);
  const [answer, setAnswer] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const abortControllerRef = useRef<AbortController | null>(null);
  const [uploadPercent, setUploadPercent] = useState(0);
  const [starting, setStarting] = useState(false);
  const [waitingNextQuestion, setWaitingNextQuestion] = useState(false);
  const [conversationHistory, setConversationHistory] = useState<ConversationItem[]>([]);
  const chatContainerRef = useRef<HTMLDivElement>(null);
  const router = useRouter();

  useEffect(() => {
    const timer = setInterval(() => setElapsed(prev => prev + 1), 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    const params = (window as any).__interviewParams || (() => { try { return JSON.parse(sessionStorage.getItem('interviewParams') || 'null'); } catch { return null; } })();
    if (!params || !params.resume_id) {
      message.error('缺少面试参数或简历，请从表单页重新进入');
      return;
    }
    setStarting(true);
    const sanitize = (s: string) => s.replace(/[<>&"'`]/g, '');
    const requestBody = {
      type: String(params.type || '综合面试'),
      domain: String(params.domain || '校招'),
      difficulty: String(params.difficulty || 'easy'),
      position_name: String(params.position_name || ''),
      company_name: sanitize(String(params.company_name || '')),
      resume_id: Number(params.resume_id),
    };

    const abortController = new AbortController();
    abortControllerRef.current = abortController;

    const startInterview = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          message.error('请先登录后再开始面试');
          setStarting(false);
          return;
        }

        // 先测试后端服务是否可达
        console.log('[检测] 测试后端服务连接...');
        try {
          const testResponse = await fetch('http://localhost:8888/api/user/login', {
            method: 'OPTIONS',
            mode: 'cors',
          });
          console.log('[检测] 后端服务连接正常');
        } catch (e) {
          message.error('无法连接到后端服务，请确认后端服务是否运行在 http://localhost:8888');
          setStarting(false);
          console.error('[检测] 后端服务连接失败:', e);
          return;
        }

        // 使用JSON格式发送请求，包含resume_id
        let response;
        console.log('[面试启动] 使用JSON格式发送请求');
        console.log('[面试启动] 请求参数:', requestBody);
        
        try {
          response = await fetch('http://localhost:8888/api/mianshi/stream/start', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(requestBody),
            signal: abortController.signal,
            mode: 'cors',
          });
        } catch (headerError) {
          // 如果Authorization header方式失败，尝试使用URL参数
          console.log('[面试启动] 方案1失败，尝试方案2: 使用URL参数传递token');
          const urlWithToken = `http://localhost:8888/api/mianshi/stream/start?token=${encodeURIComponent(token)}`;
          
          response = await fetch(urlWithToken, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
            signal: abortController.signal,
            mode: 'cors',
          });
        }

        console.log('[面试启动] 收到响应:', response.status, response.statusText);

        if (!response.ok) {
          if (response.status === 401) {
            message.error('登录已过期，请重新登录');
          } else if (response.status === 404) {
            console.error('[面试启动] 404错误 - 接口不存在');
            message.error({
              content: '接口返回404，请在后端 middleware.go 中将 /api/mianshi/stream/start 添加到 jwtPublicRoutes',
              duration: 10,
            });
          } else {
            message.error(`面试启动失败：${response.status} ${response.statusText}`);
          }
          setStarting(false);
          return;
        }

        if (!response.body) {
          message.error('无法读取响应流');
          setStarting(false);
          return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
          const { done, value } = await reader.read();
          if (done) {
            setStarting(false);
            break;
          }

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n\n');
          buffer = lines.pop() || '';

          for (const block of lines) {
            // SSE 格式可能是 "event: xxx\ndata: {...}" 或仅 "data: {...}"
            const dataMatch = block.match(/^data:\s*(.+)$/m);
            if (dataMatch) {
              const json = dataMatch[1];
              try {
                const payload = JSON.parse(json);
                console.log('[SSE数据]', payload);
                if (payload?.type === 'session_id') {
                  const sid = payload.session_id || payload.data?.session_id || '';
                  console.log('[会话ID]', sid);
                  setSessionId(sid);
                  setStarting(false);
                } else if (payload?.type === 'start') {
                  const sid = payload.session_id || '';
                  console.log('[面试开始] session_id:', sid);
                  setSessionId(sid);
                  setStarting(false);
                } else if (payload?.type === 'question') {
                  const q = payload.data?.question_text || '';
                  const idx = payload.index || payload.data?.index || 0;
                  console.log('[问题]', q, 'index:', idx);
                  setQuestionText(q);
                  setQuestionIndex(idx);
                  setStarting(false);
                  setWaitingNextQuestion(false);
                  
                  // 添加问题到对话历史
                  setConversationHistory(prev => [
                    ...prev,
                    {
                      type: 'question',
                      content: q,
                      index: idx,
                      timestamp: Date.now()
                    }
                  ]);
                } else if (payload?.type === 'end' || payload?.type === 'complete') {
                  console.log('[面试结束]', payload);
                  message.info('面试已结束');
                  setStarting(false);
                }
              } catch (e) {
                console.error('解析SSE数据失败:', json, e);
              }
            }
          }
        }
      } catch (error: any) {
        if (error.name !== 'AbortError') {
          message.error('面试启动失败：网络错误');
          console.error('启动面试错误:', error);
        }
        setStarting(false);
      }
    };

    startInterview();

    return () => {
      abortController.abort();
      abortControllerRef.current = null;
    };
  }, []);
  
  // 自动滚动到底部
  useEffect(() => {
    if (chatContainerRef.current) {
      chatContainerRef.current.scrollTop = chatContainerRef.current.scrollHeight;
    }
  }, [conversationHistory, waitingNextQuestion]);
  
  const mm = String(Math.floor(elapsed / 60)).padStart(2, '0');
  const ss = String(elapsed % 60).padStart(2, '0');
  const percent = Math.min(100, answeredCount > 0 ? Math.round((answeredCount / Math.max(answeredCount, 1)) * 100) : 0);

  const onSubmit = async (act?: 'next' | 'quit') => {
    if (!sessionId) {
      message.warning('会话已失效，请重新开始面试');
      return;
    }
    
    const action = act || 'next';
    
    // 如果是结束面试
    if (action === 'quit') {
      try {
        const token = localStorage.getItem('token');
        if (token) {
          // 调用后端接口结束面试
          // 更新为新的结束面试接口
          await fetch('http://localhost:8888/api/mianshi/interview/end', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify({
              session_id: sessionId,
              // 结束面试可能不需要answer字段，仅传session_id即可
            }),
            mode: 'cors',
          });
        }
      } catch (e) {
        console.error('[结束面试] 请求失败:', e);
      }
      try { abortControllerRef.current?.abort(); } catch {}
      message.success('面试已结束，正在跳转...');
      router.push('/user/interviews');
      return;
    }
    
    // 验证答案不为空
    if (!answer.trim()) {
      message.warning('请输入答案后再提交');
      return;
    }
    
    setSubmitting(true);
    
    // 添加答案到对话历史
    setConversationHistory(prev => [
      ...prev,
      {
        type: 'answer',
        content: answer,
        timestamp: Date.now()
      }
    ]);
    
    setAnsweredCount(prev => prev + 1);
    const currentAnswer = answer;
    setAnswer('');
    setWaitingNextQuestion(true);
    
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        message.error('登录已过期，请重新登录');
        setSubmitting(false);
        setWaitingNextQuestion(false);
        return;
      }
      
      console.log('[提交答案] 调用submit/answer接口:', { session_id: sessionId, answer: currentAnswer });
      
      // 调用submit/answer接口提交答案
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      };
      
      // 更新为新的提交答案接口
      const response = await fetch('http://localhost:8888/api/mianshi/answer/submit', {
        method: 'POST',
        headers,
        body: JSON.stringify({
          session_id: sessionId,
          answer: currentAnswer
        }),
        mode: 'cors',
      });
      
      console.log('[提交答案] 响应状态:', response.status);
      
      if (!response.ok) {
        const errorText = await response.text();
        console.error('[提交答案] 错误响应:', errorText);
        message.error('答案提交失败，请重试');
        setWaitingNextQuestion(false);
        return;
      }
      
      // 提交成功，等待SSE流推送下一题
      console.log('[提交答案] 提交成功，等待SSE推送下一题');
      message.success('答案已提交，正在生成下一题...');
      
    } catch (error: any) {
      console.error('[提交答案] 异常:', error);
      message.error('答案提交失败：网络错误');
      setWaitingNextQuestion(false);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="container mx-auto px-4">
      <Space align="center" className="mb-2" wrap>
        <Title level={2} style={{ margin: 0 }}>综合面试 · 校招简历面试</Title>
        <Space size={16} className="ml-2">
          <span style={{ fontSize: 16 }}>已回答 {answeredCount} 题</span>
          <Space align="center" size={8}>
            <span style={{ fontSize: 16 }}>进度</span>
            <Progress percent={percent} size="small" style={{ width: 140 }} />
            <Tag color="green" style={{ fontSize: 16, padding: '4px 10px' }}>{percent}%</Tag>
          </Space>
          <Tag style={{ fontSize: 16, padding: '4px 10px' }}>{mm}:{ss}</Tag>
          <Button className="bg-green-500" type="primary" onClick={() => onSubmit('quit')}>结束面试</Button>
        </Space>
      </Space>

      <Row gutter={[24, 24]}>
        <Col xs={24} md={20}>
          <Space direction="vertical" style={{ width: '100%' }}>
            {/* 对话历史区域 */}
            <AntCard className="rounded-2xl" style={{ minHeight: '400px', maxHeight: '600px' }}>
              <div 
                ref={chatContainerRef}
                style={{ 
                  maxHeight: '550px', 
                  overflowY: 'auto',
                  paddingRight: '10px'
                }}
              >
                <Space direction="vertical" style={{ width: '100%' }} size={16}>
                  {conversationHistory.length === 0 && !starting && (
                    <div style={{ textAlign: 'center', color: '#999', padding: '40px 0' }}>
                      正在生成首题，请稍候…
                    </div>
                  )}
                  
                  {conversationHistory.map((item, index) => {
                    // 计算问题序号：统计当前项之前有多少个问题，然后+1
                    const questionNumber = conversationHistory
                      .slice(0, index + 1)
                      .filter(i => i.type === 'question').length;
                    
                    return (
                      <div key={index}>
                        {item.type === 'question' ? (
                          <Space align="start" style={{ width: '100%' }}>
                            <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                            <Space direction="vertical" style={{ flex: 1 }}>
                              <div className="bg-orange-50 rounded-2xl px-5 py-4 text-base">
                                {item.content}
                              </div>
                              <Text type="secondary">第 {questionNumber} 题</Text>
                            </Space>
                          </Space>
                        ) : (
                          <div style={{ width: '100%', display: 'flex', justifyContent: 'flex-end', alignItems: 'flex-start' }}>
                            <Space direction="vertical" style={{ alignItems: 'flex-end', maxWidth: '70%' }}>
                              <div 
                                className="bg-blue-50 rounded-2xl px-5 py-4 text-base" 
                                style={{ 
                                  backgroundColor: '#e6f7ff',
                                  wordBreak: 'break-word',
                                  whiteSpace: 'pre-wrap',
                                  textAlign: 'left'
                                }}
                              >
                                {item.content}
                              </div>
                              <Text type="secondary">你的回答</Text>
                            </Space>
                            <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=user" size={48} style={{ marginLeft: '12px' }} />
                          </div>
                        )}
                      </div>
                    );
                  })}
                  
                  {/* 等待下一题的提示 */}
                  {waitingNextQuestion && (
                    <Space align="start" style={{ width: '100%' }}>
                      <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                      <div className="bg-gray-100 rounded-2xl px-5 py-4 text-base" style={{ fontStyle: 'italic', color: '#999' }}>
                        正在生成下一题，请稍候...
                      </div>
                    </Space>
                  )}
                  
                  {/* 正在生成首题的提示 */}
                  {starting && conversationHistory.length === 0 && (
                    <Space align="start" style={{ width: '100%' }}>
                      <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                      <div className="bg-gray-100 rounded-2xl px-5 py-4 text-base" style={{ fontStyle: 'italic', color: '#999' }}>
                        正在生成首题，请稍候…
                      </div>
                    </Space>
                  )}
                </Space>
              </div>
            </AntCard>

            <AntCard className="rounded-2xl">
              <Space direction="vertical" style={{ width: '100%' }}>
                <Input.TextArea 
                  rows={8} 
                  placeholder={waitingNextQuestion ? "等待下一题..." : "请在此作答，建议结构化回答（背景/职责/挑战/成果/反思）"} 
                  maxLength={500} 
                  showCount 
                  value={answer} 
                  onChange={(e) => setAnswer(e.target.value)}
                  disabled={waitingNextQuestion || starting}
                />
                <Space align="center" size={16}>
                  <Tag style={{ fontSize: 14 }}>已回答 {answeredCount} 题</Tag>
                  <Tag style={{ fontSize: 14 }}>字数建议 80-300</Tag>
                  {waitingNextQuestion && <Tag color="processing" style={{ fontSize: 14 }}>等待下一题...</Tag>}
                </Space>
                <Space size={12}>
                  <Button disabled={waitingNextQuestion || starting}>求助</Button>
                  <Button 
                    type="primary" 
                    icon={<SendOutlined />} 
                    loading={submitting} 
                    disabled={!sessionId || waitingNextQuestion || starting || !answer.trim()} 
                    onClick={() => onSubmit()}
                  >
                    {waitingNextQuestion ? '等待中...' : '提交答案'}
                  </Button>
                </Space>
              </Space>
            </AntCard>
          </Space>
        </Col>
        <Col xs={0} md={4}>
          <Space direction="vertical" className="fixed right-6" size={16}>
            <Button shape="circle" icon={<CustomerServiceOutlined />} />
            <Button shape="circle" icon={<AudioOutlined />} />
            <Button shape="circle" icon={<QuestionCircleOutlined />} />
          </Space>
        </Col>
      </Row>
    </div>
  );
}
