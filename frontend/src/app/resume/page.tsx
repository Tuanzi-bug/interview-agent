'use client';

import { useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Form, Select, Input, Radio, Button, Collapse, Tag } from 'antd';
import { FileTextOutlined, VideoCameraOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

export default function ResumePressPage() {
  const [form] = Form.useForm();

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">简历押题</Title>

      <Row gutter={[24, 24]} className="mt-4">
        <Col xs={24} md={16}>
          <AntCard className="rounded-2xl bg-green-50" styles={{ body: { padding: 16 } }}>
            <div className="text-green-700">
              <div className="font-medium mb-2">注意事项</div>
              <ul className="space-y-2 text-sm">
                <li>押题会根据你的简历内容生成，一次约提供至少20道题。</li>
                <li>添加/切换简历内容会影响押题的范围，系统会为你保留历史题目。</li>
              </ul>
            </div>
          </AntCard>

          <AntCard className="rounded-2xl mt-4">
            <Form form={form} layout="vertical" initialValues={{ resume: '我的简历_v1.pdf', type: '精准岗位押题', language: 'Java', job: '后端开发', city: '北京', level: '入门' }}>
              <Form.Item label="选择押题的简历" name="resume">
                <Select options={[{ value: '我的简历_v1.pdf', label: '我的简历_v1.pdf' }, { value: '校招版.pdf', label: '校招版.pdf' }]} popupMatchSelectWidth={false} />
              </Form.Item>

              <Form.Item label="押题类型" name="type">
                <Radio.Group>
                  <Radio value="精准岗位押题">精准岗位押题</Radio>
                  <Radio value="让机器面试题">让机器面试题</Radio>
                </Radio.Group>
              </Form.Item>

              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="语言" name="language">
                    <Select options={[{ value: 'Java', label: 'Java' }, { value: 'Golang', label: 'Golang' }, { value: 'Python', label: 'Python' }]} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="岗位意向" name="job">
                    <Input placeholder="如：Java后端开发" />
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
                  <Form.Item label="难度等级" name="level">
                    <Select options={[{ value: '入门', label: '入门' }, { value: '中级', label: '中级' }, { value: '进阶', label: '进阶' }]} />
                  </Form.Item>
                </Col>
              </Row>

              <Collapse className="mt-2" items={[{ key: 'adv', label: '高级选项', children: (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="text-sm text-gray-600">偏好题型</div>
                    <Select mode="multiple" placeholder="选择偏好题型" options={[{ value: '项目深挖', label: '项目深挖' }, { value: '系统设计', label: '系统设计' }, { value: '性能优化', label: '性能优化' }]} />
                  </div>
                  <div className="space-y-2">
                    <div className="text-sm text-gray-600">重点领域</div>
                    <Select mode="multiple" placeholder="选择重点领域" options={[{ value: '数据库', label: '数据库' }, { value: '中间件', label: '中间件' }, { value: '并发网络', label: '并发网络' }]} />
                  </div>
                </div>
              ) }]} />

              <div className="mt-6">
                <Button type="primary" className="bg-green-500 w-full h-12 text-base">开始简历押题</Button>
                <div className="text-center text-gray-500 text-sm mt-2">首次免费押题 20 题题量</div>
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
            <div className="rounded-xl overflow-hidden bg-black/5">
              <div className="w-full h-48 md:h-60 bg-gray-100 rounded-xl flex items-center justify-center text-gray-600">
                <VideoCameraOutlined className="text-3xl mr-2" />示例视频
              </div>
            </div>
          </AntCard>
        </Col>
      </Row>

      <div className="mt-10 grid grid-cols-1 md:grid-cols-3 gap-6">
        <AntCard className="rounded-2xl">
          <div className="flex items-center gap-3">
            <FileTextOutlined className="text-green-500 text-2xl" />
            <div>
              <div className="font-medium">快速定位</div>
              <Paragraph className="text-gray-600 m-0">深入解析简历条目，生成对应问答清单，直击考点。</Paragraph>
            </div>
          </div>
        </AntCard>
        <AntCard className="rounded-2xl">
          <div className="flex items-center gap-3">
            <FileTextOutlined className="text-green-500 text-2xl" />
            <div>
              <div className="font-medium">快速剖析</div>
              <Paragraph className="text-gray-600 m-0">结合岗位要求与项目经历，输出结构化追问路径。</Paragraph>
            </div>
          </div>
        </AntCard>
        <AntCard className="rounded-2xl">
          <div className="flex items-center gap-3">
            <FileTextOutlined className="text-green-500 text-2xl" />
            <div>
              <div className="font-medium">直接学习</div>
              <Paragraph className="text-gray-600 m-0">对题清单搭配参考答案与延伸阅读，立即提升。</Paragraph>
            </div>
          </div>
        </AntCard>
      </div>
    </div>
  );
}
