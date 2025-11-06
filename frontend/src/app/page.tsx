'use client';

import { Typography, Row, Col } from 'antd';
import { BulbOutlined, FileTextOutlined, CodeOutlined } from '@ant-design/icons';
import Banner from '@/components/home/Banner';
import type { FC } from 'react';

const { Title, Paragraph } = Typography;

export default function Home() {
  // 功能卡片数据
  const features = [
    {
      title: '综合面试',
      description: '全方位评估专业能力、沟通技巧和解决问题的能力',
      icon: <BulbOutlined className="text-4xl text-primary mb-4" />,
    },
    {
      title: '简历押题',
      description: '深度解析简历，定制化预测面试考题，面试如开卷',
      icon: <FileTextOutlined className="text-4xl text-secondary mb-4" />,
    },
    {
      title: '专项面试',
      description: '针对特定技术领域进行深度评估，帮助查漏补缺',
      icon: <CodeOutlined className="text-4xl text-success mb-4" />,
    },
  ];

  return (
    <div className="container mx-auto px-4">
      {/* Banner区 */}
      <Banner />
      
      {/* 功能入口区 */}
      <section className="py-16">
        <div className="text-center mb-12">
          <Title level={2}>专业AI面试服务</Title>
          <Paragraph className="text-gray-600 max-w-2xl mx-auto">
            基于10万+大厂真题深度训练，95%命中率精准预测，让你的面试更有信心
          </Paragraph>
        </div>
        
        <Row gutter={[24, 24]}>
          {features.map((feature, index) => (
            <Col xs={24} md={8} key={index}>
              <div className="bg-white rounded-lg shadow-sm p-6 hover:shadow-md transition-all h-full text-center hover:translate-y-[-5px]">
                <div className="flex flex-col items-center">
                  {feature.icon}
                  <Title level={4} className="mb-2">{feature.title}</Title>
                  <Paragraph className="text-gray-600">{feature.description}</Paragraph>
                </div>
              </div>
            </Col>
          ))}
        </Row>
      </section>
    </div>
  );
}