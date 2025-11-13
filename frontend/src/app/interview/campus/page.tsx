'use client';

import { Typography, Row, Col, Card as AntCard, Form, Select, Input, Button, Tag } from 'antd';
import { CheckCircleOutlined, VideoCameraOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export default function CampusInterviewPage() {
  const [form] = Form.useForm();

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">综合面试 · 校招简历面试</Title>
      <Paragraph className="text-gray-600 max-w-3xl">
        立标校招面试要求，围绕简历条目与技术栈进行真实追问；从基础能力到应用能力、学习与思考、技术热情与潜力，构建符合校招的结构化评估。
      </Paragraph>

      <Row gutter={[24, 24]} className="mt-2">
        <Col xs={24} md={16}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              '重视系统基础能力测评，构建逻辑理解力',
              '校招岗位标准化题目，链路化表达能力',
              '构建学习与举例表达，形成对知识的掌握',
              '将解题过程结构化，判断潜力与成长力',
            ].map((t, i) => (
              <div key={i} className="flex items-center gap-2 text-green-700">
                <CheckCircleOutlined />
                <span>{t}</span>
              </div>
            ))}
          </div>

          <AntCard className="rounded-2xl">
            <Form form={form} layout="vertical" initialValues={{ resume: '校招版.pdf', language: 'Java', title: '算法岗', level: '入门' }}>
              <Form.Item label="选择面试的简历" name="resume">
                <Select options={[{ value: '校招版.pdf', label: '校招版.pdf' }, { value: '我的简历_v1.pdf', label: '我的简历_v1.pdf' }]} popupMatchSelectWidth={false} />
              </Form.Item>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="语言" name="language">
                    <Select options={[{ value: 'Java', label: 'Java' }, { value: 'Golang', label: 'Golang' }, { value: 'Python', label: 'Python' }]} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="岗位名称" name="title">
                    <Input placeholder="如：算法岗" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="难度等级" name="level">
                    <Select options={[{ value: '入门', label: '入门' }, { value: '中级', label: '中级' }, { value: '进阶', label: '进阶' }]} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="面试时长" name="duration">
                    <Select options={[{ value: '30min', label: '30分钟' }, { value: '60min', label: '60分钟' }]} />
                  </Form.Item>
                </Col>
              </Row>

              <div className="mt-4">
                <Button type="primary" className="bg-green-500 w-full h-12 text-base">开始校招面试</Button>
                <div className="text-center text-gray-500 text-sm mt-2">1次面试时长 20分钟-60分钟，自动串联题目与追问链路</div>
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

