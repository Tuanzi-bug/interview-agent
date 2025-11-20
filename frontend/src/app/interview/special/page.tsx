'use client';

import { useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Form, Select, Button, Tag } from 'antd';
import { CheckCircleOutlined, VideoCameraOutlined } from '@ant-design/icons';

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
  const [form] = Form.useForm();

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
            <Form form={form} layout="vertical" initialValues={{ stack: stack, level: '入门' }}>
              <Form.Item label="专项类别" name="stack">
                <Select
                  popupMatchSelectWidth={false}
                  options={GROUPED_OPTIONS}
                  value={stack}
                  onChange={(v) => setStack(v)}
                />
              </Form.Item>
              <Form.Item label="难度等级" name="level">
                <Select options={[{ value: '入门', label: '入门' }, { value: '中级', label: '中级' }, { value: '进阶', label: '进阶' }]} />
              </Form.Item>
              <div className="mt-2">
                <Button type="primary" className="bg-green-500 w-full h-12 text-base">首次专项面试免费</Button>
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
            <div className="w-full h-48 md:h-60 bg-gray-100 rounded-xl flex items-center justify-center text-gray-600">
              <VideoCameraOutlined className="text-3xl mr-2" />示例视频
            </div>
          </AntCard>
        </Col>
      </Row>
    </div>
  );
}

