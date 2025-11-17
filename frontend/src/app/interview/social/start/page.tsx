'use client';

import { useEffect, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Space, Tag, Button, Input, Avatar, Progress } from 'antd';
import { AudioOutlined, CustomerServiceOutlined, QuestionCircleOutlined, SendOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

export default function SocialInterviewStartPage() {
  const [elapsed, setElapsed] = useState(0);
  useEffect(() => {
    const timer = setInterval(() => setElapsed(prev => prev + 1), 1000);
    return () => clearInterval(timer);
  }, []);
  const mm = String(Math.floor(elapsed / 60)).padStart(2, '0');
  const ss = String(elapsed % 60).padStart(2, '0');

  return (
    <div className="container mx-auto px-4">
      <Space align="center" className="mb-2" wrap>
        <Title level={2} style={{ margin: 0 }}>综合面试 · 社招简历面试</Title>
        <Space size={16} className="ml-2">
          <span style={{ fontSize: 16 }}>已回答 0 题</span>
          <Space align="center" size={8}>
            <span style={{ fontSize: 16 }}>进度</span>
            <Progress percent={0} size="small" style={{ width: 140 }} />
            <Tag color="green" style={{ fontSize: 16, padding: '4px 10px' }}>0%</Tag>
          </Space>
          <Tag style={{ fontSize: 16, padding: '4px 10px' }}>{mm}:{ss}</Tag>
          <Button className="bg-green-500" type="primary">结束面试</Button>
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
                    我看你简历上提到了RAG驱动AI旅行助手项目，能先跟我简单介绍一下这个项目的背景和你主要负责的部分吗？
                  </div>
                  <Text type="secondary">22:39:12</Text>
                </Space>
              </Space>
            </AntCard>

            <AntCard className="rounded-2xl">
              <Space direction="vertical" style={{ width: '100%' }}>
                <Input.TextArea rows={8} placeholder="请在此作答，建议结构化回答（背景/职责/挑战/成果/反思）" maxLength={500} showCount />
                <Space align="center" size={16}>
                  <Tag style={{ fontSize: 14 }}>已回答 0 题</Tag>
                  <Tag style={{ fontSize: 14 }}>字数建议 80-300</Tag>
                </Space>
                <Space size={12}>
                  <Button>求助</Button>
                  <Button type="primary" icon={<SendOutlined />}>提交答案</Button>
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