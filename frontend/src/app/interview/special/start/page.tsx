'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Typography,
  Row,
  Col,
  Card as AntCard,
  Space,
  Tag,
  Button,
  Input,
  Avatar,
  Progress,
  message,
} from 'antd';
import { AudioOutlined, CustomerServiceOutlined, QuestionCircleOutlined, SendOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

type SpecialInterviewParams = {
  domain: string;
  difficulty: string;
};

export default function SpecialInterviewStartPage() {
  const [elapsed, setElapsed] = useState(0);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [questionText, setQuestionText] = useState<string>('');
  const [questionIndex, setQuestionIndex] = useState<number>(0);
  const [answeredCount, setAnsweredCount] = useState(0);
  const [answer, setAnswer] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [starting, setStarting] = useState(false);
  const [waitingNextQuestion, setWaitingNextQuestion] = useState(false);
  const [conversationHistory, setConversationHistory] = useState<
    { type: 'question' | 'answer'; content: string; index?: number; timestamp: number }[]
  >([]);
  const [meta, setMeta] = useState<SpecialInterviewParams>({ domain: '-', difficulty: '-' });
  const abortControllerRef = useRef<AbortController | null>(null);
  const chatContainerRef = useRef<HTMLDivElement>(null);
  const specialParamsRef = useRef<SpecialInterviewParams | null>(null);
  const router = useRouter();

  useEffect(() => {
    const timer = setInterval(() => setElapsed((prev) => prev + 1), 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    const params: SpecialInterviewParams | null =
      (window as any).__specialInterviewParams ||
      (() => {
        try {
          return JSON.parse(sessionStorage.getItem('specialInterviewParams') || 'null');
        } catch {
          return null;
        }
      })();

    if (!params || !params.domain || !params.difficulty) {
      message.error('缺少专项面试参数，请重新选择专项方向');
      router.push('/interview/special');
      return;
    }

    const normalizedParams: SpecialInterviewParams = {
      domain: String(params.domain || '').trim() || 'java',
      difficulty: String(params.difficulty || '').trim() || '简单',
    };

    specialParamsRef.current = normalizedParams;
    setMeta(normalizedParams);
    setStarting(true);

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

        const requestBody = {
          domain: normalizedParams.domain,
          difficulty: normalizedParams.difficulty,
        };

        let response: Response;
        try {
          response = await fetch('http://localhost:8888/api/interview/special/stream', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify(requestBody),
            signal: abortController.signal,
            mode: 'cors',
          });
        } catch (headerError) {
          const fallbackUrl = `http://localhost:8888/api/interview/special/stream?token=${encodeURIComponent(token)}`;
          response = await fetch(fallbackUrl, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
            signal: abortController.signal,
            mode: 'cors',
          });
        }

        if (!response.ok) {
          if (response.status === 401) {
            message.error('登录已过期，请重新登录');
          } else {
            message.error(`专项面试启动失败：${response.status} ${response.statusText}`);
          }
          setStarting(false);
          return;
        }

        if (!response.body) {
          message.error('无法读取专项面试响应流');
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
            const dataMatch = block.match(/^data:\s*(.+)$/m);
            if (dataMatch) {
              const json = dataMatch[1];
              try {
                const payload = JSON.parse(json);
                if (payload?.type === 'session_id') {
                  const sid = payload.session_id || payload.data?.session_id || '';
                  setSessionId(sid);
                  setStarting(false);
                } else if (payload?.type === 'start') {
                  const sid = payload.session_id || '';
                  setSessionId(sid);
                  setStarting(false);
                } else if (payload?.type === 'question') {
                  const q = payload.data?.question_text || '';
                  const idx = payload.index || payload.data?.index || 0;
                  setQuestionText(q);
                  setQuestionIndex(idx);
                  setStarting(false);
                  setWaitingNextQuestion(false);
                  setConversationHistory((prev) => [
                    ...prev,
                    {
                      type: 'question',
                      content: q,
                      index: idx,
                      timestamp: Date.now(),
                    },
                  ]);
                } else if (payload?.type === 'end' || payload?.type === 'complete') {
                  message.info('专项面试已结束');
                  setStarting(false);
                }
              } catch (e) {
                console.error('解析专项面试 SSE 数据失败:', json, e);
              }
            }
          }
        }
      } catch (error: any) {
        if (error.name !== 'AbortError') {
          message.error('专项面试启动失败：网络错误');
          console.error('专项面试启动错误:', error);
        }
        setStarting(false);
      }
    };

    startInterview();

    return () => {
      abortController.abort();
      abortControllerRef.current = null;
    };
  }, [router]);

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
      message.warning('会话已失效，请重新开始专项面试');
      return;
    }

    const params = specialParamsRef.current;
    const action = act || 'next';

    if (action === 'quit') {
      try {
        const token = localStorage.getItem('token');
        if (token) {
          await fetch('http://localhost:8888/api/interview/submit/answer', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              session_id: sessionId,
              answer: '',
              action: 'quit',
              domain: params?.domain,
              difficulty: params?.difficulty,
            }),
            mode: 'cors',
          });
        }
      } catch (e) {
        console.error('[专项面试] 结束请求失败:', e);
      }
      try {
        abortControllerRef.current?.abort();
      } catch {}
      message.success('专项面试已结束，正在跳转...');
      router.push('/interview/special');
      return;
    }

    if (!answer.trim()) {
      message.warning('请输入答案后再提交');
      return;
    }

    setSubmitting(true);
    setConversationHistory((prev) => [
      ...prev,
      {
        type: 'answer',
        content: answer,
        timestamp: Date.now(),
      },
    ]);
    setAnsweredCount((prev) => prev + 1);
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

      const response = await fetch('http://localhost:8888/api/interview/submit/answer', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          session_id: sessionId,
          answer: currentAnswer,
          action: 'next',
          domain: params?.domain,
          difficulty: params?.difficulty,
        }),
        mode: 'cors',
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error('[专项面试] 提交答案失败:', errorText);
        message.error('答案提交失败，请重试');
        setWaitingNextQuestion(false);
        return;
      }

      message.success('答案已提交，正在生成下一题...');
    } catch (error: any) {
      console.error('[专项面试] 提交异常:', error);
      message.error('答案提交失败：网络错误');
      setWaitingNextQuestion(false);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="container mx-auto px-4">
      <Space align="center" className="mb-2" wrap>
        <Title level={2} style={{ margin: 0 }}>
          专项面试 · {meta.domain}（{meta.difficulty}）
        </Title>
        <Space size={16} className="ml-2" wrap>
          <span style={{ fontSize: 16 }}>已回答 {answeredCount} 题</span>
          <Tag color="blue" style={{ fontSize: 16, padding: '4px 10px' }}>方向：{meta.domain}</Tag>
          <Tag color="purple" style={{ fontSize: 16, padding: '4px 10px' }}>难度：{meta.difficulty}</Tag>
          <Space align="center" size={8}>
            <span style={{ fontSize: 16 }}>进度</span>
            <Progress percent={percent} size="small" style={{ width: 140 }} />
            <Tag color="green" style={{ fontSize: 16, padding: '4px 10px' }}>{percent}%</Tag>
          </Space>
          <Tag style={{ fontSize: 16, padding: '4px 10px' }}>{mm}:{ss}</Tag>
          <Button className="bg-green-500" type="primary" onClick={() => onSubmit('quit')}>
            结束面试
          </Button>
        </Space>
      </Space>

      <Row gutter={[24, 24]}>
        <Col xs={24} md={20}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <AntCard className="rounded-2xl" style={{ minHeight: '400px', maxHeight: '600px' }}>
              <div
                ref={chatContainerRef}
                style={{
                  maxHeight: '550px',
                  overflowY: 'auto',
                  paddingRight: '10px',
                }}
              >
                <Space direction="vertical" style={{ width: '100%' }} size={16}>
                  {conversationHistory.length === 0 && !starting && (
                    <div style={{ textAlign: 'center', color: '#999', padding: '40px 0' }}>
                      正在生成首题，请稍候…
                    </div>
                  )}

                  {conversationHistory.map((item, index) => {
                    const questionNumber = conversationHistory.slice(0, index + 1).filter((i) => i.type === 'question').length;
                    return (
                      <div key={`${item.type}-${item.timestamp}-${index}`}>
                        {item.type === 'question' ? (
                          <Space align="start" style={{ width: '100%' }}>
                            <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                            <Space direction="vertical" style={{ flex: 1 }}>
                              <div className="bg-orange-50 rounded-2xl px-5 py-4 text-base">{item.content}</div>
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
                                  textAlign: 'left',
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

                  {waitingNextQuestion && (
                    <Space align="start" style={{ width: '100%' }}>
                      <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                      <div className="bg-gray-100 rounded-2xl px-5 py-4 text-base" style={{ fontStyle: 'italic', color: '#999' }}>
                        正在生成下一题，请稍候...
                      </div>
                    </Space>
                  )}

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
                  placeholder={waitingNextQuestion ? '等待下一题...' : '请在此作答，建议结构化回答（背景/职责/挑战/成果/反思）'}
                  maxLength={500}
                  showCount
                  value={answer}
                  onChange={(e) => setAnswer(e.target.value)}
                  disabled={waitingNextQuestion || starting}
                />
                <Space align="center" size={16}>
                  <Tag style={{ fontSize: 14 }}>已回答 {answeredCount} 题</Tag>
                  <Tag style={{ fontSize: 14 }}>字数建议 80-300</Tag>
                  {waitingNextQuestion && (
                    <Tag color="processing" style={{ fontSize: 14 }}>
                      等待下一题...
                    </Tag>
                  )}
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
