'use client';

import { useEffect, useRef, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Space, Tag, Button, Input, Avatar, Progress } from 'antd';
import { AudioOutlined, CustomerServiceOutlined, QuestionCircleOutlined, SendOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

export default function CampusInterviewStartPage() {
  const [elapsed, setElapsed] = useState(0);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [questionText, setQuestionText] = useState<string>('');
  const [questionIndex, setQuestionIndex] = useState<number>(0);
  const [answeredCount, setAnsweredCount] = useState(0);
  const [answer, setAnswer] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const esRef = useRef<EventSource | null>(null);
  useEffect(() => {
    const timer = setInterval(() => setElapsed(prev => prev + 1), 1000);
    return () => clearInterval(timer);
  }, []);
  useEffect(() => {
    const url = 'http://localhost:8888/api/interview/start/stream';
    try {
      const es = new EventSource(url);
      esRef.current = es;
      es.onmessage = (evt) => {
        try {
          const payload = JSON.parse(evt.data);
          if (payload && payload.type === 'session_id') {
            setSessionId(payload.session_id || payload.data?.session_id || '');
          } else if (payload && payload.type === 'question') {
            const q = payload.data?.question_text || '';
            setQuestionText(q);
            setQuestionIndex(payload.index || 0);
          }
        } catch {}
      };
      es.onerror = () => {
        es.close();
      };
    } catch {}
    return () => {
      esRef.current?.close();
      esRef.current = null;
    };
  }, []);
  const mm = String(Math.floor(elapsed / 60)).padStart(2, '0');
  const ss = String(elapsed % 60).padStart(2, '0');
  const percent = Math.min(100, answeredCount > 0 ? Math.round((answeredCount / Math.max(answeredCount, 1)) * 100) : 0);
  const onSubmit = async (act?: 'next' | 'quit') => {
    if (!sessionId) return;
    const action = act || 'next';
    setSubmitting(true);
    try {
      await fetch('http://localhost:8888/api/interview/submit/answer', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, answer, action })
      });
      if (action === 'next') {
        setAnsweredCount(prev => prev + 1);
        setAnswer('');
      } else {
        esRef.current?.close();
      }
    } catch {
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
            <AntCard className="rounded-2xl">
              <Space align="start" style={{ width: '100%' }}>
                <Avatar src="https://api.dicebear.com/7.x/avataaars/svg?seed=interviewer2" size={48} />
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
                  <Tag style={{ fontSize: 14 }}>字数建议 60-250</Tag>
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