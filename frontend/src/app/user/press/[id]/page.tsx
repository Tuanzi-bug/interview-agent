'use client';

import { useMemo, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, List, Tag, Space, Divider } from 'antd';
import { BookOutlined, BulbOutlined, CheckCircleOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

const staticData = {
  questions: [
    {
      title: '我看你简历上提到了RAG驱动AI旅行助手项目，你主要负责的部分吗？',
      question:
        '我看你简历上提到了RAG驱动AI旅行助手项目，能先跟我简单介绍一下这个项目的背景和你主要负责的部分吗？',
      idea: [
        {
          label: '背景',
          value: '项目背景兼顾的真实性、个人职责的清晰度、技术方案的合理性、项目成果的可信度',
        },
        { label: '职责', value: '负责AI助手的研发落地，提升客户体验与员工效率' },
        { label: '挑战', value: '解决检索准确性与响应稳定性问题，构建可解释的AI响应' },
        { label: '成果', value: '整体满意度提升，业务指标优化，成功率提升' },
      ],
      reference:
        '该项目为RAG驱动的旅行助手，整合检索与生成，在稳定性与响应速度上做了优化，通过流式SSE提升交互体验，核心链路可用性达到99%以上。',
      followups: [
        '这个AI助手为何选择RAG而不是纯生成式？',
        '你在项目中的技术决策有哪些？如何权衡准确性与性能？',
        '项目上线后，用户反馈和指标变化如何？下一步优化方向？',
      ],
    },
    {
      title: '在这个AI旅行助手中如何做多源数据检索与融合？',
      question: '在这个AI旅行助手中，如何实现多源数据检索与融合以保证答复的可靠性？',
      idea: [
        { label: '背景', value: '多源数据，包括百科、攻略库、商家信息与用户生成内容' },
        { label: '职责', value: '负责检索管道设计与结果融合策略落地' },
        { label: '挑战', value: '异构数据的质量与时效性问题' },
        { label: '成果', value: '答复准确率与一致性提升' },
      ],
      reference:
        '采用分层检索与重排策略，BM25+向量检索结合，针对问句类别使用不同的融合权重与投票机制。',
      followups: ['你如何评估融合策略的效果？', '数据时效性问题如何处理？'],
    },
    {
      title: '在RAG系统中如何处理并发与延迟问题？',
      question: '在RAG系统中如何处理并发与延迟问题，确保用户体验？',
      idea: [
        { label: '背景', value: '高并发场景下检索与生成的协同' },
        { label: '职责', value: '优化请求调度与缓存策略' },
        { label: '挑战', value: '检索延迟与生成阻塞' },
        { label: '成果', value: '端到端延迟稳定在 100-200ms 量级（流式首包更快）' },
      ],
      reference: '采用异步管道与消息队列、向量缓存与热点文档预取，首包用SSE推送提升感知速度。',
      followups: ['为什么选择SSE而非WebSocket？', '缓存失效策略如何设计？'],
    },
  ],
};

export default function PressDetailPage() {
  const [selected, setSelected] = useState(0);
  const current = useMemo(
    () => staticData.questions[selected] || staticData.questions[0],
    [selected]
  );

  return (
    <div className="min-h-screen relative font-sans">
      {/* Decorative Background */}
      <div className="fixed top-0 right-0 w-[600px] h-[600px] bg-green-50/60 rounded-full blur-[120px] -translate-y-1/2 translate-x-1/3 pointer-events-none z-0" />
      <div className="fixed bottom-0 left-0 w-[600px] h-[600px] bg-blue-50/60 rounded-full blur-[120px] translate-y-1/2 -translate-x-1/3 pointer-events-none z-0" />

      <div className="container mx-auto px-4 relative z-10 pb-12">
        <div className="mb-8 animate-fade-in-up">
          <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight flex items-center gap-3">
            <BookOutlined className="text-blue-600" />
            押题详情
          </h1>
          <p className="text-slate-500 mt-2 ml-11">查看为您生成的精准面试题目与详细解析</p>
        </div>

        <Row
          gutter={[24, 24]}
          className="mt-2 animate-fade-in-up"
          style={{ animationDelay: '0.1s' }}
        >
          <Col xs={24} md={7} lg={6}>
            <AntCard
              className="rounded-2xl border-slate-100 shadow-lg shadow-slate-200/50 h-full"
              styles={{ body: { padding: 0 } }}
            >
              <div className="p-5 border-b border-slate-100 bg-slate-50/50 rounded-t-2xl">
                <Space align="center" className="font-bold text-slate-700">
                  <BookOutlined className="text-blue-500" />
                  <span>题目目录</span>
                </Space>
              </div>
              <div className="max-h-[calc(100vh-300px)] overflow-y-auto custom-scrollbar">
                <List
                  itemLayout="horizontal"
                  dataSource={staticData.questions}
                  split={false}
                  renderItem={(item, index) => (
                    <List.Item
                      onClick={() => setSelected(index)}
                      className={`transition-colors duration-200 cursor-pointer border-l-4 px-4 py-3 hover:bg-blue-50/50 ${
                        index === selected ? 'bg-blue-50 border-blue-500' : 'border-transparent'
                      }`}
                    >
                      <div className="w-full">
                        <div
                          className={`font-medium mb-1 line-clamp-2 ${index === selected ? 'text-blue-700' : 'text-slate-700'}`}
                        >
                          <span className="mr-2 text-slate-400">0{index + 1}.</span>
                          {item.title}
                        </div>
                      </div>
                    </List.Item>
                  )}
                />
              </div>
            </AntCard>
          </Col>
          <Col xs={24} md={17} lg={18}>
            <AntCard
              className="rounded-2xl border-slate-100 shadow-lg shadow-slate-200/50 min-h-[600px]"
              styles={{ body: { padding: 32 } }}
            >
              <div className="flex flex-col gap-6">
                <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
                  <h2 className="text-xl font-bold text-slate-800 leading-relaxed flex-1">
                    <span className="text-blue-600 mr-2 text-2xl">Q{selected + 1}.</span>
                    {current.question}
                  </h2>
                  <Space wrap>
                    <Tag
                      color="blue"
                      className="px-3 py-1 rounded-full border-0 bg-blue-50 text-blue-600 font-medium"
                    >
                      整体思路
                    </Tag>
                    <Tag
                      color="green"
                      className="px-3 py-1 rounded-full border-0 bg-green-50 text-green-600 font-medium"
                    >
                      参考答案
                    </Tag>
                    <Tag className="px-3 py-1 rounded-full border-slate-200 text-slate-500 hover:text-blue-600 cursor-pointer transition-colors">
                      收藏单题
                    </Tag>
                  </Space>
                </div>

                <div className="w-full h-px bg-slate-100" />

                <div className="space-y-8">
                  <section>
                    <div className="flex items-center gap-2 mb-4 text-lg font-bold text-slate-800">
                      <div className="w-8 h-8 rounded-lg bg-green-100 flex items-center justify-center text-green-600">
                        <CheckCircleOutlined />
                      </div>
                      答案思路
                    </div>
                    <Row gutter={[16, 16]}>
                      {current.idea.map((it, idx) => (
                        <Col key={idx} xs={24} md={12}>
                          <div className="bg-slate-50 rounded-xl p-4 border border-slate-100 hover:border-blue-200 hover:bg-blue-50/30 transition-colors h-full">
                            <div className="flex items-center gap-2 mb-2 font-bold text-slate-700">
                              <BulbOutlined className="text-yellow-500" />
                              {it.label}
                            </div>
                            <p className="text-slate-600 text-sm leading-relaxed">{it.value}</p>
                          </div>
                        </Col>
                      ))}
                    </Row>
                  </section>

                  <section>
                    <div className="flex items-center gap-2 mb-4 text-lg font-bold text-slate-800">
                      <div className="w-8 h-8 rounded-lg bg-indigo-100 flex items-center justify-center text-indigo-600">
                        <BookOutlined />
                      </div>
                      参考答案
                    </div>
                    <div className="bg-gradient-to-br from-slate-50 to-white rounded-xl p-6 border border-slate-100 text-slate-700 leading-loose">
                      {current.reference}
                    </div>
                  </section>

                  <section>
                    <div className="flex items-center gap-2 mb-4 text-lg font-bold text-slate-800">
                      <div className="w-8 h-8 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600">
                        <BulbOutlined />
                      </div>
                      可能追问
                    </div>
                    <div className="bg-white rounded-xl border border-slate-100 overflow-hidden">
                      <List
                        dataSource={current.followups}
                        renderItem={(t, i) => (
                          <List.Item className="px-6 py-4 hover:bg-slate-50 transition-colors border-slate-50">
                            <div className="flex gap-3">
                              <span className="text-orange-500 font-bold font-mono">0{i + 1}</span>
                              <span className="text-slate-700">{t}</span>
                            </div>
                          </List.Item>
                        )}
                      />
                    </div>
                  </section>
                </div>
              </div>
            </AntCard>
          </Col>
        </Row>
      </div>
    </div>
  );
}
