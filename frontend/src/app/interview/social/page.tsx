'use client';

import { Typography, Row, Col, Card as AntCard, Form, Select, Input, Button, Tag } from 'antd';
import { CheckCircleOutlined, VideoCameraOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export default function SocialInterviewPage() {
  const [form] = Form.useForm();

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">综合面试 · 社招简历面试</Title>
      <Paragraph className="text-gray-600 max-w-3xl">
        在综合面试模式中，系统会围绕你的简历、项目经历与岗位胜任力，从技术基础、项目落地、设计能力到沟通协作，构建环环追问的真实面试场景，帮助你快速查漏补缺与提升应对能力。
      </Paragraph>

      <Row gutter={[24, 24]} className="mt-2">
        <Col xs={24} md={16}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              '深挖技术本质逻辑，构建环环追问交叉',
              '聚焦架构设计能力，真实场景还原',
              '针对业务问题推演，技术方案落地',
              '逻辑体系梳理完整，洞察核心关键点',
            ].map((t, i) => (
              <div key={i} className="flex items-center gap-2 text-green-700">
                <CheckCircleOutlined />
                <span>{t}</span>
              </div>
            ))}
          </div>

          <AntCard className="rounded-2xl">
            <Form form={form} layout="vertical" initialValues={{ resume: '我的简历_v1.pdf', job: 'Java后端开发', level: '入门', city: '北京' }}>
              <Form.Item label="选择面试的简历" name="resume">
                <Select options={[{ value: '我的简历_v1.pdf', label: '我的简历_v1.pdf' }, { value: '加强版.pdf', label: '加强版.pdf' }]} popupMatchSelectWidth={false} />
              </Form.Item>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="岗位意向" name="job">
                    <Input placeholder="如：Java后端开发" />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="难度等级" name="level">
                    <Select options={[{ value: '入门', label: '入门' }, { value: '中级', label: '中级' }, { value: '进阶', label: '进阶' }]} />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="求职城市" name="city">
                    <Select options={[{ value: '北京', label: '北京' }, { value: '上海', label: '上海' }, { value: '深圳', label: '深圳' }]} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="面试时长" name="duration">
                    <Select options={[{ value: '30min', label: '30分钟' }, { value: '60min', label: '60分钟' }]} />
                  </Form.Item>
                </Col>
              </Row>

              <div className="mt-4">
                <Button type="primary" className="bg-green-500 w-full h-12 text-base">开始面试</Button>
                <div className="text-center text-gray-500 text-sm mt-2">1次体验价约等于20次AI陪练，单次2小时题目自动续集</div>
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

