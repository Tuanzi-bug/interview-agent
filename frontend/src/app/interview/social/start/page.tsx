'use client';

import { useEffect, useRef, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Space, Tag, Button, Input, Avatar, Progress, message } from 'antd';
import { AudioOutlined, CustomerServiceOutlined, QuestionCircleOutlined, SendOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

export default function SocialInterviewStartPage() {
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
  useEffect(() => {
    const timer = setInterval(() => setElapsed(prev => prev + 1), 1000);
    return () => clearInterval(timer);
  }, []);
  useEffect(() => {
    const params = (window as any).__interviewParams || (() => { try { return JSON.parse(sessionStorage.getItem('interviewParams') || 'null'); } catch { return null; } })();
    const resumeFile = (window as any).__interviewResume || null;
    if (!params || !resumeFile) {
      message.error('缺少面试参数，请从表单页重新进入');
      return;
    }
    setStarting(true);
    const formData = new FormData();
    const sanitize = (s: string) => s.replace(/[<>&"'`]/g, '');
    formData.append('title', String(params.title || '综合面试'));
    formData.append('type', String(params.type || '综合面试'));
    formData.append('domain', String(params.domain || '社招'));
    formData.append('difficulty', String(params.difficulty || 'easy'));
    formData.append('position_name', String(params.position_name || ''));
    formData.append('company_name', sanitize(String(params.company_name || '')));
    if (params.query) {
      formData.append('query', String(params.query));
    }
    const f = (resumeFile?.originFileObj || resumeFile) as Blob;
    const fname = (resumeFile?.name || 'resume.pdf');
    formData.append('resume', f, fname);

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

        // 尝试方案1: 使用Authorization header
        let response;
        console.log('[面试启动] 尝试方案1: 使用Authorization header');
        try {
          const headers: Record<string, string> = {};
          headers['Authorization'] = `Bearer ${token}`;
          
          console.log('[面试启动] 请求URL:', 'http://localhost:8888/api/interview/start/stream');
          
          response = await fetch('http://localhost:8888/api/interview/start/stream', {
            method: 'POST',
            headers,
            body: formData,
            signal: abortController.signal,
            mode: 'cors',
          });
        } catch (headerError) {
          // 如果Authorization header方式失败，尝试使用URL参数
          console.log('[面试启动] 方案1失败，尝试方案2: 使用URL参数传递token');
          const urlWithToken = `http://localhost:8888/api/interview/start/stream?token=${encodeURIComponent(token)}`;
          console.log('[面试启动] 请求URL:', urlWithToken);
          
          response = await fetch(urlWithToken, {
            method: 'POST',
            body: formData,
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
              content: '接口返回404，请在后端 middleware.go 中将 /api/interview/start/stream 添加到 jwtPublicRoutes',
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

          for (const line of lines) {
            if (line.trim().startsWith('data:')) {
              const json = line.trim().replace(/^data:\s*/, '');
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
                  console.log('[问题]', q);
                  setQuestionText(q);
                  setQuestionIndex(payload.index || 0);
                  setStarting(false);
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
  const mm = String(Math.floor(elapsed / 60)).padStart(2, '0');
  const ss = String(elapsed % 60).padStart(2, '0');
  const percent = Math.min(100, answeredCount > 0 ? Math.round((answeredCount / Math.max(answeredCount, 1)) * 100) : 0);
  const onSubmit = async (act?: 'next' | 'quit') => {
    if (!sessionId) return;
    const action = act ? act : (answer.trim() === '结束面试' ? 'quit' : 'next');
    setSubmitting(true);
    try {
      const token = localStorage.getItem('token');
      const headers: Record<string, string> = { 'Content-Type': 'application/json' };
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
      await fetch('http://localhost:8888/api/interview/submit/answer', {
        method: 'POST',
        headers,
        body: JSON.stringify({ session_id: sessionId, answer, action }),
        mode: 'cors',
      });
      if (action === 'next') {
        setAnsweredCount(prev => prev + 1);
        setAnswer('');
      } else {
        try { abortControllerRef.current?.abort(); } catch {}
      }
    } catch {
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="container mx-auto px-4">
      <Space align="center" className="mb-2" wrap>
        <Title level={2} style={{ margin: 0 }}>综合面试 · 社招简历面试</Title>
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
            <AntCard className="rounded-2xl">
              <Space align="start" style={{ width: '100%' }}>
                <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer" size={48} />
                <Space direction="vertical" style={{ width: '100%' }}>
                  <div className="bg-orange-50 rounded-2xl px-5 py-4 text-base">
                    {questionText || '正在生成首题，请稍候…'}
                  </div>
                  <Text type="secondary">第 {questionIndex || 1} 题</Text>
                </Space>
              </Space>
            </AntCard>

            <AntCard className="rounded-2xl">
              <Space direction="vertical" style={{ width: '100%' }}>
                <Input.TextArea rows={8} placeholder="请在此作答，建议结构化回答（背景/职责/挑战/成果/反思）" maxLength={500} showCount value={answer} onChange={(e) => setAnswer(e.target.value)} />
                <Space align="center" size={16}>
                  <Tag style={{ fontSize: 14 }}>已回答 {answeredCount} 题</Tag>
                  <Tag style={{ fontSize: 14 }}>字数建议 80-300</Tag>
                </Space>
                <Space size={12}>
                  <Button>求助</Button>
                  <Button type="primary" icon={<SendOutlined />} loading={submitting} disabled={!sessionId} onClick={() => onSubmit()}>提交答案</Button>
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