'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Typography, Row, Col, Card as AntCard, Form, Select, Button, Tag, message, Alert } from 'antd';
import { CheckCircleOutlined, VideoCameraOutlined, CaretRightOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

const GROUPED_OPTIONS = [
  { label: '标准语言', options: [
    { value: 'Java', label: 'Java' },
    { value: 'Go', label: 'Go' },
    { value: 'C/C++', label: 'C/C++' },
    { value: 'Rust', label: 'Rust' },
    { value: 'PHP', label: 'PHP' },
    { value: 'Node.js', label: 'Node.js' },
  ]},
  { label: '后端组件', options: [
    { value: 'Redis', label: 'Redis' },
    { value: 'MySQL', label: 'MySQL' },
    { value: 'Kafka', label: 'Kafka' },
    { value: 'MongoDB', label: 'MongoDB' },
  ]},
  { label: '云原生与运维', options: [
    { value: 'Docker', label: 'Docker' },
    { value: 'Kubernetes', label: 'Kubernetes' },
    { value: 'Nginx', label: 'Nginx' },
  ]},
  { label: '计算机基础', options: [
    { value: '操作系统', label: '操作系统' },
    { value: '计算机网络', label: '计算机网络' },
    { value: '数据结构与算法', label: '数据结构与算法' },
  ]},
];

export default function SpecialInterviewPage() {
  const [stack, setStack] = useState<string>('Go');
  const [starting, setStarting] = useState(false);
  const [modelConfigured, setModelConfigured] = useState<boolean | null>(null);
  const [checkingConfig, setCheckingConfig] = useState<boolean>(false);
  const [form] = Form.useForm();
  const router = useRouter();

  useEffect(() => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    setCheckingConfig(true);
    fetch('http://localhost:8888/api/user/model/check', {
      method: 'GET',
      headers: {
        Authorization: token ? `Bearer ${token}` : '',
        'X-Auth-Token': token || '',
      },
    })
      .then(async (res) => {
        const data = await res.json().catch(() => null);
        const configured = !!(data && data.data && data.data.configured);
        setModelConfigured(configured);
      })
      .catch(() => {
        setModelConfigured(false);
      })
      .finally(() => {
        setCheckingConfig(false);
      });
  }, []);

  const handleStart = async () => {
    if (!modelConfigured) {
      message.error('未配置模型，无法开始面试');
      return;
    }
    try {
      const values = await form.validateFields();
      setStarting(true);
      
      const params = {
        type: '专项面试',
        domain: values.stack,
        difficulty: values.level
      };
      
      (window as any).__interviewParams = { ...params };
      try { sessionStorage.setItem('interviewParams', JSON.stringify(params)); } catch {}
      
      router.push('/interview/special/start');
    } catch (e) {
      message.error('请选择专项类别和难度等级');
    } finally {
      setStarting(false);
    }
  };

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">专项面试 · {stack}</Title>
      <Paragraph className="text-gray-600 max-w-3xl">
        选择专项方向后，系统会围绕该技术栈构建真实面试场景，聚焦高频问题与深度追问，结合行业通用标准输出结构化评估与改进建议。
      </Paragraph>

      <Row gutter={[24, 24]} className="mt-2">
        <Col xs={24} md={16}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              '精准拆解，直击岗位核心高频要点',
              '动静结合，链路梳理',
              '高密度追问，定位能力边界',
              '实战模拟，还原面试真实效果',
            ].map((t, i) => (
              <div key={i} className="flex items-center gap-2 text-green-700">
                <CheckCircleOutlined />
                <span>{t}</span>
              </div>
            ))}
          </div>

          <AntCard className="rounded-2xl">
            <Form form={form} layout="vertical" initialValues={{ stack: stack, level: '简单' }}>
              <Form.Item label="专项类别" name="stack">
                <Select
                  popupMatchSelectWidth={false}
                  options={GROUPED_OPTIONS}
                  value={stack}
                  onChange={(v) => setStack(v)}
                />
              </Form.Item>
              <Form.Item label="难度等级" name="level">
                <Select options={[{ value: '简单', label: '简单' }, { value: '中等', label: '中等' }, { value: '复杂', label: '复杂' }]} />
              </Form.Item>
              <div className="mt-2">
                {!checkingConfig && modelConfigured === false && (
                  <Alert
                    message="模型未配置"
                    description={
                      <span>
                        请去 <Link href="/user/models" className="text-blue-500 underline">用户模型页面</Link> 配置模型
                      </span>
                    }
                    type="warning"
                    showIcon
                    className="mb-4"
                  />
                )}
                {checkingConfig && (
                  <Tag color="default" className="mb-2">正在检查模型配置</Tag>
                )}
                <Button 
                  type="primary" 
                  className="bg-green-500 w-full h-12 text-base" 
                  onClick={handleStart}
                  loading={starting}
                  disabled={starting || checkingConfig || modelConfigured === false}
                >
                  首次专项面试免费
                </Button>
                <div className="text-center text-gray-500 text-sm mt-2">单次专项面试约30-60分钟，系统自动续集题目链路</div>
              </div>
            </Form>
          </AntCard>
        </Col>
        <Col xs={24} md={8}>
          <AntCard className="rounded-2xl">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2"><VideoCameraOutlined /><span>功能演示</span></div>
              <Tag color="green">推荐观看</Tag>
            </div>
            <div className="w-full h-48 md:h-60 bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 rounded-xl flex flex-col items-center justify-center text-slate-600 relative overflow-hidden group cursor-pointer transition-all hover:shadow-lg border border-slate-100">
              <div className="absolute inset-0 bg-[linear-gradient(45deg,#0000_25%,rgba(0,0,0,0.02)_0,rgba(0,0,0,0.02)_50%,#0000_0,#0000_75%,rgba(0,0,0,0.02)_0)] bg-[length:20px_20px] opacity-50" />
              
              <div className="w-16 h-16 bg-white rounded-full flex items-center justify-center text-green-500 shadow-md transform scale-95 group-hover:scale-110 transition-all duration-300 z-10 group-hover:text-green-600">
                <CaretRightOutlined style={{ fontSize: '32px', marginLeft: '4px' }} />
              </div>
              
              <div className="mt-4 font-medium z-10 group-hover:text-slate-800 transition-colors">功能演示视频</div>
              <div className="text-xs text-slate-400 mt-1 z-10">点击播放 (演示)</div>
            </div>
          </AntCard>
        </Col>
      </Row>
    </div>
  );
}

