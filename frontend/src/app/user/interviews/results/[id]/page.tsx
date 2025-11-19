'use client';

import { useMemo } from 'react';
import { Typography, Card as AntCard, Row, Col, List, Tag, Button, Avatar } from 'antd';
import { useParams } from 'next/navigation';

const { Title, Paragraph, Text } = Typography;

const mockBasic = {
  candidate: 'LittleBear',
  resume: 'golang.pdf',
  type: '综合面试',
  score: 65,
  difficulty: '挑战',
  company: '腾讯',
  position: '软件开发-后台开发方向',
  duration: '13分钟12秒',
  time: '2025-11-13 22:37:53',
};

const mockPerformance = {
  comment:
    '候选人整体表现良好，技术基础扎实，沟通清晰。建议加强在系统设计和架构方面的深度思考。',
  dimensions: [
    { dimension_name: '技术基础与实践能力', evaluation: '掌握基础数据结构与算法，实现规范，逻辑清晰。', score: 68 },
    { dimension_name: '项目经历真实性验证', evaluation: '能清晰描述项目背景与职责，举例充分。', score: 62 },
    { dimension_name: '系统架构设计思维', evaluation: '对架构理解到位，需加强扩展性与容错思考。', score: 55 },
    { dimension_name: '技术深度与前瞻视野', evaluation: '对新技术保持关注，有学习与反思能力。', score: 58 },
    { dimension_name: '团队协作与沟通能力', evaluation: '表达清晰，主动沟通与互动，解释思路清楚。', score: 66 },
  ],
};

const mockRecords = [
  {
    order: 1,
    content: '项目考察',
    comment: {
      score: 85,
      key_points: '项目架构、技术选型、团队协作',
      difficulty: '中等',
      strengths: '表达清晰，逻辑完整，技术细节讲解充分',
      weaknesses: '缺少对项目性能优化的讨论',
      suggestion: '建议补充项目的性能指标和优化方案',
      know_points: '系统设计、技术栈选择、项目管理',
      thinking: '从项目背景→技术选型→实现细节→性能优化的思路讲解',
      reference: '应包含核心架构、技术栈、难点与解决方案',
    },
    message: [
      { order: 1, question: '这个项目的核心难点是什么？', answer: '高并发场景的数据一致性，使用消息队列与分布式锁。' },
      { order: 2, question: '你们是如何处理数据一致性的？', answer: 'Redis分布式锁 + RabbitMQ，保证操作原子性。' },
      { order: 3, question: '性能指标如何？', answer: '支持10万QPS，平均响应 50ms。' },
    ],
  },
  {
    order: 2,
    content: '技术考察',
    comment: {
      score: 78,
      key_points: '服务拆分、通信方式、服务治理',
      difficulty: '中高',
      strengths: '基础扎实，能举例说明',
      weaknesses: '服务治理实践不足',
      suggestion: '建议学习 Service Mesh 相关技术',
      know_points: '微服务、RPC、服务发现、负载均衡',
      thinking: '从单体问题→微服务优势→实现挑战的讲解',
      reference: '应包含拆分原则与通信方式等',
    },
    message: [
      { order: 1, question: '微服务和单体架构的主要区别是什么？', answer: '微服务可独立部署与扩展，单体为整体应用。' },
      { order: 2, question: '微服务之间如何通信？', answer: '同步RPC（如gRPC）与异步消息队列。' },
    ],
  },
];

function RadarChart({ items, size = 520 }: { items: { dimension_name: string; score: number }[]; size?: number }) {
  const radius = 140;
  const cx = size / 2;
  const cy = size / 2;
  const points = items.map((it, i) => {
    const angle = (2 * Math.PI * i) / items.length - Math.PI / 2;
    const r = (Math.max(0, Math.min(100, Number(it.score))) / 100) * radius;
    return [cx + r * Math.cos(angle), cy + r * Math.sin(angle)];
  });
  const axis = items.map((_, i) => {
    const angle = (2 * Math.PI * i) / items.length - Math.PI / 2;
    return [cx + radius * Math.cos(angle), cy + radius * Math.sin(angle)];
  });
  const poly = points.map(p => p.join(',')).join(' ');
  return (
    <div style={{ width: size, height: size }}>
      <svg width={size} height={size}>
        <circle cx={cx} cy={cy} r={radius} fill="#f6ffed" stroke="#b7eb8f" />
        {axis.map((p, i) => (
          <line key={i} x1={cx} y1={cy} x2={p[0]} y2={p[1]} stroke="#d9d9d9" />
        ))}
        <polygon points={poly} fill="rgba(82,196,26,0.25)" stroke="#52c41a" />
        {axis.map((p, i) => (
          <text key={i} x={p[0]} y={p[1]} dx={p[0] > cx ? 8 : -8} dy={p[1] > cy ? 16 : -8} textAnchor={p[0] > cx ? 'start' : 'end'} fontSize={14} fontWeight={600} fill="#434343">
            {items[i].dimension_name}
          </text>
        ))}
      </svg>
    </div>
  );
}

function ScoreGauge({ score }: { score: number }) {
  const size = 260;
  const cx = size / 2;
  const cy = size / 2;
  const r = 100;
  const start = Math.PI;
  const end = 0;
  const angle = (Math.max(0, Math.min(100, score)) / 100) * Math.PI;
  const arcPath = (ang: number, color: string) => {
    const sx = cx + r * Math.cos(start);
    const sy = cy + r * Math.sin(start);
    const ex = cx + r * Math.cos(start + ang);
    const ey = cy + r * Math.sin(start + ang);
    return (
      <path d={`M ${sx} ${sy} A ${r} ${r} 0 0 1 ${ex} ${ey}`} stroke={color} strokeWidth={14} fill="none" />
    );
  };
  return (
    <div style={{ width: size, height: size / 1.4 }}>
      <svg width={size} height={size / 1.4}>
        <path d={`M ${cx - r} ${cy} A ${r} ${r} 0 0 1 ${cx + r} ${cy}`} stroke="#e6f4ff" strokeWidth={14} fill="none" />
        {arcPath(angle, '#52c41a')}
        <text x={cx} y={cy - 10} textAnchor="middle" fontSize={12} fill="#8c8c8c">本次面试评分</text>
        <text x={cx} y={cy + 24} textAnchor="middle" fontSize={36} fontWeight={600} fill="#262626">{score}</text>
      </svg>
    </div>
  );
}

export default function InterviewResultDetailPage() {
  const params = useParams() as any;
  const id = Number(params?.id || 1);
  const perf = mockPerformance;
  const basic = mockBasic;
  const records = mockRecords;

  const avgScore = useMemo(() => {
    const s = perf.dimensions.map(d => d.score);
    return Math.round((s.reduce((a, b) => a + b, 0) / s.length) * 10) / 10;
  }, [perf]);

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">面试结果</Title>
      
      <AntCard className="rounded-2xl mt-2" styles={{ body: { padding: 20 } }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={14}>
            <div className="flex items-start gap-3">
              <Avatar size={48}>LB</Avatar>
              <div>
                <div className="text-lg font-semibold">{basic.candidate}</div>
                <div className="mt-3 space-y-2 text-sm">
                  <div>面试类型：{basic.type}</div>
                  <div>使用简历：{basic.resume}</div>
                  <div>面试难度：{basic.difficulty}</div>
                  <div>公司名称：{basic.company}</div>
                  <div>岗位名称：{basic.position}</div>
                  <div>面试时间：{basic.time}</div>
                </div>
              </div>
            </div>
          </Col>
          <Col xs={24} md={10}>
            <div className="flex items-center justify-between gap-4">
              <ScoreGauge score={basic.score} />
              <div className="flex flex-col gap-2">
                <Button>下载面试报告</Button>
                <Button type="primary">AI 技能补充建议</Button>
              </div>
            </div>
          </Col>
        </Row>
      </AntCard>

      
      <AntCard className="rounded-2xl mt-4" styles={{ body: { padding: 20 } }}>
        <div className="text-xl font-bold mb-2">面试官点评</div>
        <Paragraph className="text-base leading-relaxed">{perf.comment}</Paragraph>
      </AntCard>

      <div className="mt-4">
        <div className="text-xl font-bold mb-2">面试表现</div>
        <div className="flex justify-center">
          <RadarChart items={perf.dimensions} size={520} />
        </div>
        <div className="mt-6">
          <Row gutter={[24, 24]}>
            {perf.dimensions.map((d, i) => (
              <Col key={i} xs={24} md={8}>
                <AntCard className="rounded-2xl" styles={{ body: { padding: 20 } }}>
                  <div className="flex items-center gap-2 text-green-700">
                    <Tag color="green">{d.dimension_name}</Tag>
                  </div>
                  <Paragraph className="text-base leading-relaxed" style={{ marginTop: 8 }}>{d.evaluation}</Paragraph>
                  <div className="text-base text-gray-700 font-medium">{d.score} 分</div>
                </AntCard>
              </Col>
            ))}
          </Row>
        </div>
      </div>

      

      <AntCard className="rounded-2xl mt-6" styles={{ body: { padding: 20 } }}>
        <div className="text-lg font-semibold mb-2">本次答题记录</div>
        <List
          dataSource={records}
          renderItem={(rec) => (
            <List.Item>
              <div style={{ width: '100%' }}>
                <div className="flex items-center gap-2">
                  <Tag color="green">{rec.order}/{records.length}</Tag>
                  <Text strong>{rec.content}</Text>
                </div>
                <AntCard size="small" className="mt-2" title="问答">
                  <List
                    dataSource={rec.message}
                    renderItem={(m) => (
                      <List.Item>
                        <div>
                          <Text strong className="text-base">{m.order}. 问题：</Text>
                          <Text className="text-base">{m.question}</Text>
                          <div className="mt-1">
                            <Text strong className="text-base">回答：</Text>
                            <Text className="text-base leading-relaxed">{m.answer}</Text>
                          </div>
                        </div>
                      </List.Item>
                    )}
                  />
                </AntCard>
                <AntCard className="mt-3" styles={{ body: { padding: 20, background: '#f5f5f5' } }}>
                  <div className="space-y-3">
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">本题得分</div><div className="flex-1 text-lg font-semibold"><Tag color="blue">{rec.comment.score}分</Tag></div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">难度等级</div><div className="flex-1 text-lg font-semibold"><Tag>{rec.comment.difficulty}</Tag></div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">关键点</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.key_points}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">优势</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.strengths}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">不足</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.weaknesses}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">建议</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.suggestion}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">考点</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.know_points}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">思路</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.thinking}</div></div>
                    <div className="flex gap-6"><div className="w-32 text-gray-800 text-lg font-bold">参考答案</div><div className="flex-1 text-lg leading-relaxed">{rec.comment.reference}</div></div>
                  </div>
                </AntCard>
              </div>
            </List.Item>
          )}
        />
      </AntCard>
    </div>
  );
}