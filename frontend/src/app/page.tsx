'use client';

import { Typography, Row, Col, Button, Card as AntCard, Rate, Collapse, Avatar, Tag } from 'antd';
import { BulbOutlined, FileTextOutlined, CodeOutlined, CompassOutlined, FlagOutlined, EyeOutlined, ThunderboltOutlined, SmileOutlined, SwitcherOutlined, SendOutlined, TeamOutlined, ExperimentOutlined, VideoCameraOutlined, UserOutlined } from '@ant-design/icons';
import Banner from '@/components/home/Banner';
import type { FC } from 'react';

const { Title, Paragraph } = Typography;

export default function Home() {
  const testimonials = [
    {
      text: '非科班出身，自学一年多总感觉基础不扎实。牛面的简历押题功能太神了，针对我的项目经历生成的题目命中率很高，90%都在实际面试中遇到过。特别是Spring框架深度问题、IOC到AOP，从题目的逻辑出发让我能由浅入深串联知识体系，这种从值到到深度的感觉真的很棒。',
      user: '35岁重启人生',
      title: '社招Java开发',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=a',
    },
    {
      text: '996的工作节奏根本没时间找人mock interview。牛面24小时随时随地陪练，效率极高。每天晚上坚持练习30分钟，岗位面试官会问到系统架构设计、链路梳理思维，牛面都能覆盖到位。最后还给了提升建议，面完真实的大厂，和正式面试时居然遇到70%相似问题！',
      user: 'smartbob',
      title: '在职提升-Go开发',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=b',
    },
    {
      text: '用了牛面模拟面试，追问技术细节的能力大幅提升！模拟面试会根据我的回答深入追问，连Redis AOF、优化这种偏细节都能展开，最后还给了提升建议。面完真实的大厂，和正式面试时居然遇到70%相似问题！',
      user: 'Lex',
      title: '架构师',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=c',
    },
    {
      text: '工作三年想冲击大厂，但系统设计这块一直是短板。牛面详细评估报告+改进建议与题目分析，高效锤炼，不断改进，亲身实战。现在在面试中能条理清晰地讲解方案，拿到满意的薪资。',
      user: '静以修身',
      title: '后端开发',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=d',
    },
    {
      text: '很棒！价格可能就是人工服务的1/10，效果至少是8、9成，性价比超高！',
      user: 'jackey',
      title: '后端开发',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=e',
    },
    {
      text: 'Django ORM优化和导出视图这些都能有涉及，面试准备很全面。',
      user: '默然',
      title: 'Python开发',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=f',
    },
  ];

  return (
    <div className="container mx-auto px-4">
      
      <div className="w-full bg-gradient-to-r from-purple-500 to-purple-600 text-white rounded-lg my-4 px-4 py-2 text-sm flex items-center justify-between">
        <div>
          秋招特惠：免费简历押题体验+综合面试单次免费！快来无限问答，简历押题&专项面试免费体验！
        </div>
        <div className="hidden sm:block opacity-90">
          今日进行中：题量 +224；出题速度 +3.2 题/小时；板块活跃度 +32.1%
        </div>
      </div>

      
      <section className="py-8">
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={14}>
            <Title level={2} className="leading-snug">
              <span className="text-primary">牛面 综合面试</span>
              <br />全维度拷问 真实场景还原
              <br />充分准备 <span className="text-primary">自信</span>迎战！
            </Title>
            <Paragraph className="text-gray-600 mt-4">
              同构单位：大企业标准综合考察 | 详细问题复盘与弱项分析 | 深度面试针对性训练持续提升
            </Paragraph>
            <div className="flex gap-3 mt-6">
              <Button type="primary" size="large">立即使用</Button>
              <Button size="large" icon={<VideoCameraOutlined />}>视频演示</Button>
            </div>
            <div className="mt-6 flex items-center text-gray-500 text-sm gap-6">
              <div className="flex items-center gap-2"><Tag color="green">注册用户</Tag><span>12,056+</span></div>
              <div className="flex items-center gap-2"><Tag color="blue">AI面试次数</Tag><span>5,000+</span></div>
              <div className="flex items-center gap-2"><Tag color="purple">斩获Offer</Tag><span>500+</span></div>
              <div className="flex items-center gap-2"><Tag color="red">涨薪幅度</Tag><span>40%</span></div>
            </div>
          </Col>
          <Col xs={24} md={10}>
            <div className="bg-orange-100 rounded-2xl h-64 md:h-80 flex items-center justify-center">
              <span className="text-orange-500 text-3xl">示例插画</span>
            </div>
          </Col>
        </Row>
      </section>

      
      <section className="py-10">
        <div className="text-center mb-10">
          <Title level={2}>通过牛面能得到什么服务？</Title>
        </div>
        <Row gutter={[24, 24]}>
          <Col xs={24} md={8}>
            <AntCard className="rounded-2xl">
              <div className="flex items-center gap-3 mb-4">
                <BulbOutlined className="text-3xl text-primary" />
                <Title level={4} className="m-0">综合面试</Title>
              </div>
              <Paragraph className="text-gray-600">高度还原真实面试场景，全面考察技术基础、架构设计、项目经验和团队协作。</Paragraph>
              <ul className="text-gray-600 space-y-2">
                <li>• 多维度深入问答，层层追问细节</li>
                <li>• 结构化评估报告与改进建议</li>
                <li>• 支持多岗位角色模拟</li>
              </ul>
              <Button type="primary" className="mt-6 w-full">选择综合面试</Button>
            </AntCard>
          </Col>
          <Col xs={24} md={8}>
            <AntCard className="rounded-2xl">
              <div className="flex items-center gap-3 mb-4">
                <FileTextOutlined className="text-3xl text-blue-500" />
                <Title level={4} className="m-0">简历押题</Title>
              </div>
              <Paragraph className="text-gray-600">深度分析你的简历，生成个性化的专属面试题目，让你面试如开卷。</Paragraph>
              <ul className="text-gray-600 space-y-2">
                <li>• 精准命中高频考点</li>
                <li>• 项目经历细节追问</li>
                <li>• 输出押题清单</li>
              </ul>
              <Button type="primary" className="mt-6 w-full">选择简历押题</Button>
            </AntCard>
          </Col>
          <Col xs={24} md={8}>
            <AntCard className="rounded-2xl">
              <div className="flex items-center gap-3 mb-4">
                <CodeOutlined className="text-3xl text-purple-500" />
                <Title level={4} className="m-0">专项面试</Title>
              </div>
              <Paragraph className="text-gray-600">针对特定技术领域进行深度评估，帮助查漏补缺，快速补齐短板。</Paragraph>
              <ul className="text-gray-600 space-y-2">
                <li>• 数据库/中间件/性能优化</li>
                <li>• 系统设计/并发/网络</li>
                <li>• 前后端/算法等</li>
              </ul>
              <Button type="primary" className="mt-6 w-full">选择专项面试</Button>
            </AntCard>
          </Col>
        </Row>
      </section>

      
      <section className="py-6">
        <div className="text-center mb-8">
          <Title level={2}>真实口碑，有据可依</Title>
          <Paragraph className="text-gray-600 max-w-3xl mx-auto">不搞虚的，来自用户的真实反馈，覆盖技术栈、框架、架构等多类用户群体的面试提升体验与成果。</Paragraph>
        </div>
        <Row gutter={[24, 24]}>
          {testimonials.map((t, i) => (
            <Col xs={24} md={8} key={i}>
              <AntCard className="rounded-2xl">
                <div className="flex justify-between mb-4">
                  <div className="text-3xl text-primary">“</div>
                  <Rate disabled defaultValue={5} />
                </div>
                <Paragraph className="text-gray-700">{t.text}</Paragraph>
                <div className="flex items-center gap-3 mt-6">
                  <Avatar src={t.avatar} />
                  <div>
                    <div className="font-medium">{t.user}</div>
                    <div className="text-gray-500 text-sm">{t.title}</div>
                  </div>
                </div>
              </AntCard>
            </Col>
          ))}
        </Row>
      </section>

      
      <section className="py-10">
        <div className="text-center mb-8">
          <Title level={2}>为什么选择牛面？</Title>
        </div>
        <Row gutter={[24, 24]}>
          {[
            { icon: <CompassOutlined />, title: '大厂真题轰炸', desc: '基于字节、阿里、腾讯真题深度训练，95%命中率精准押题，直击考点！', color: 'text-green-500' },
            { icon: <FlagOutlined />, title: '简历押题神器', desc: '深入分析你的简历，生成定制化的专属面试题目，让你面试如开卷', color: 'text-blue-500' },
            { icon: <EyeOutlined />, title: '还原真实标准', desc: '多名大厂面试官共同设计，还原真实评估体系：你是待定还是强烈推荐？', color: 'text-purple-500' },
            { icon: <ThunderboltOutlined />, title: '进步肉眼可见改', desc: '详细评估报告+改进建议与题目分析，高效锤炼，不断改进，亲身实战锤炼', color: 'text-red-500' },
            { icon: <SmileOutlined />, title: '成本碾压传统模式', desc: '1次私教费=20次AI特训，省下90%成本，随时开练无需预约，时间花在刀刃上', color: 'text-purple-500' },
            { icon: <SwitcherOutlined />, title: '难度自由切换', desc: '入门/进阶/挑战级难度一键切换，从小白到大神全阶级覆盖，自己掌控节奏', color: 'text-blue-500' },
            { icon: <SendOutlined />, title: '临场状态激活', desc: '模拟高压追问，实时1小时节奏练习，让你上场即巅峰，关键时刻绝不掉链子', color: 'text-green-500' },
            { icon: <TeamOutlined />, title: '领域洞察', desc: '针对不同领域岗位的真实问题库与场景化训练，全面提升', color: 'text-teal-500' },
          ].map((f, idx) => (
            <Col xs={24} md={6} key={idx}>
              <AntCard className="rounded-2xl h-full">
                <div className={`w-10 h-10 rounded-full bg-gray-100 flex items-center justify-center mb-3 ${f.color}`}>
                  <span className="text-xl">{f.icon}</span>
                </div>
                <Title level={4} className="m-0 mb-2">{f.title}</Title>
                <Paragraph className="text-gray-600">{f.desc}</Paragraph>
              </AntCard>
            </Col>
          ))}
        </Row>
      </section>

      
      <section className="py-10">
        <Title level={2} className="mb-6">常见问题</Title>
        <Collapse
          items={[
            { key: '1', label: '为什么要用牛面，而不是豆包、ChatGPT这些通用AI？', children: <Paragraph>牛面针对求职面试场景深度定制，真实面试标准、追问逻辑与评估体系都更贴近用人单位的要求。</Paragraph> },
            { key: '2', label: '牛面到底解决了什么问题', children: <Paragraph>帮助你在真实面试来临前发现薄弱点并针对性训练，输出结构化评估与改进建议。</Paragraph> },
            { key: '3', label: '牛面适合什么样的人使用？', children: <Paragraph>从校招到社招，从转岗到晋升，皆可使用；支持多岗位面试模拟。</Paragraph> },
            { key: '4', label: '收费标准是什么，性价比如何？', children: <Paragraph>单次体验低成本，会员价格更划算；与线下私教相比成本约为1/10。</Paragraph> },
            { key: '5', label: '牛面的AI是否真正理解并回答问题？', children: <Paragraph>基于大厂真题与结构化知识库训练，具备追问能力与场景还原，输出更专业。</Paragraph> },
            { key: '6', label: '牛面押题效果怎么样？', children: <Paragraph>押题命中率高达95%，覆盖核心技术栈，帮助你在面试中游刃有余。</Paragraph> },
          ]}
        />
      </section>

      
      <div className="fixed right-4 bottom-24 flex flex-col gap-3">
        <Button shape="round" className="w-12 h-12 shadow" icon={<UserOutlined />} />
        <Button shape="round" className="w-12 h-12 shadow" icon={<ExperimentOutlined />} />
        <Button shape="round" className="w-12 h-12 shadow" icon={<SendOutlined />} />
        <Button shape="round" className="w-12 h-12 shadow" icon={<VideoCameraOutlined />} />
      </div>

      
      <div className="fixed right-6 bottom-6 flex items-center gap-3">
        <Button type="primary" shape="circle">↑</Button>
        <Button shape="circle">✉️</Button>
      </div>
    </div>
  );
}
