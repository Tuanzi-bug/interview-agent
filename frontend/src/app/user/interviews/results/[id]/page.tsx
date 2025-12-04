'use client';

import { useEffect, useMemo, useState } from 'react';
import { Typography, Card as AntCard, Row, Col, List, Tag, Button, Avatar, Spin, message } from 'antd';
import { useParams } from 'next/navigation';
import apiClient from '@/services/api/client';
import { useAuth } from '@/hooks/useAuth';

const { Title, Paragraph, Text } = Typography;

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
  const id = Number(params?.id || 0);
  const { user, login } = useAuth();

  const [loading, setLoading] = useState(true);
  const [interviewInfo, setInterviewInfo] = useState<any>(null);
  const [evaluation, setEvaluation] = useState<any>(null);
  const [answerRecords, setAnswerRecords] = useState<any[]>([]);

  const fetchData = async () => {
    if (!id) return;
    setLoading(true);
    try {
      // 0. Fetch User Profile if missing
      if (!user) {
        try {
          const userRes: any = await apiClient.get('http://localhost:8888/api/user/profile');
          if (userRes && userRes.username) {
            login({
              id: String(userRes.id),
              name: userRes.username,
              email: userRes.email,
              avatar: userRes.avatar
            });
          }
        } catch (e) {
          console.error("Failed to fetch user profile", e);
        }
      }

      // 1. Fetch Interview Info (from list)
      const listRes: any = await apiClient.get('http://localhost:8888/api/interview/records', { params: { page: 1, page_size: 1000 } });
      const listData = listRes?.records || [];
      const info = listData.find((item: any) => item.id === id);
      setInterviewInfo(info);

      // 2. Fetch Evaluation Report
      const evalRes: any = await apiClient.get('http://localhost:8888/api/mianshi/evaluation', { params: { report_id: id } });
      setEvaluation(evalRes);

      // 3. Fetch Answer Records
      const recordRes: any = await apiClient.get('http://localhost:8888/api/mianshi/answer-record', { params: { report_id: id } });
      if (recordRes && recordRes.records) {
        setAnswerRecords(recordRes.records);
      } else if (Array.isArray(recordRes)) {
        setAnswerRecords(recordRes);
      } else {
        setAnswerRecords(recordRes?.records || []);
      }
    } catch (e: any) {
      console.error(e);
      message.error('获取面试详情失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <Spin size="large" tip="正在加载面试结果..." />
      </div>
    );
  }

  if (!evaluation) {
    return (
      <div className="container mx-auto px-4 mt-8">
        <AntCard>
          <div className="text-center py-8 text-gray-500">
            暂无面试报告数据，请稍后重试或确认面试是否已完成。
          </div>
        </AntCard>
      </div>
    );
  }

  const basic = {
    candidate: user?.name || '未知用户',
    resume: '暂无', // 接口暂未返回简历名称
    type: interviewInfo?.type || '未知',
    score: evaluation.score,
    difficulty: interviewInfo?.difficulty || '未知',
    company: interviewInfo?.company_name || '未指定',
    position: interviewInfo?.position_name || '未指定',
    duration: interviewInfo?.duration ? `${Math.floor(interviewInfo.duration / 60)}分钟${interviewInfo.duration % 60}秒` : '未知',
    time: interviewInfo?.created_at ? new Date(interviewInfo.created_at).toLocaleString('zh-CN') : '未知',
  };

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">面试结果</Title>
      
      <AntCard className="rounded-2xl mt-2" styles={{ body: { padding: 20 } }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={14}>
            <div className="flex items-start gap-3">
              <Avatar size={48}>{basic.candidate.substring(0, 2).toUpperCase()}</Avatar>
              <div>
                <div className="text-lg font-semibold">{basic.candidate}</div>
                <div className="mt-3 space-y-2 text-sm">
                  <div>面试类型：{basic.type}</div>
                  {/* <div>使用简历：{basic.resume}</div> */}
                  <div>面试难度：{basic.difficulty}</div>
                  <div>公司名称：{basic.company}</div>
                  <div>岗位名称：{basic.position}</div>
                  <div>面试时长：{basic.duration}</div>
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
        <Paragraph className="text-base leading-relaxed">{evaluation.comment}</Paragraph>
      </AntCard>

      <div className="mt-4">
        <div className="text-xl font-bold mb-2">面试表现</div>
        <div className="flex justify-center">
          {evaluation.dimensions && evaluation.dimensions.length > 0 && (
             <RadarChart items={evaluation.dimensions} size={520} />
          )}
        </div>
        <div className="mt-6">
          <Row gutter={[24, 24]}>
            {evaluation.dimensions?.map((d: any, i: number) => (
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
        {answerRecords && answerRecords.length > 0 ? (
        <List
          dataSource={answerRecords}
          renderItem={(rec: any) => (
            <List.Item>
              <div style={{ width: '100%' }}>
                <div className="flex items-center gap-2">
                  <Tag color="green">{rec.order}/{answerRecords.length}</Tag>
                  <Text strong>{rec.content}</Text>
                </div>
                <AntCard size="small" className="mt-2" title="问答">
                  <List
                    dataSource={rec.message || []}
                    renderItem={(m: any) => (
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
                {rec.comment && (
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
                )}
              </div>
            </List.Item>
          )}
        />
        ) : (
          <div className="text-center py-4 text-gray-500">暂无答题记录</div>
        )}
      </AntCard>
    </div>
  );
}