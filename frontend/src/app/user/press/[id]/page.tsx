'use client';

import { useEffect, useState, useMemo } from 'react';
import { Typography, Row, Col, Card as AntCard, List, Tag, Space, Divider, Spin, message, Collapse, Button, Badge } from 'antd';
import { 
  BookOutlined, 
  ArrowLeftOutlined, 
  BulbOutlined, 
  AimOutlined, 
  CheckCircleOutlined,
  MenuOutlined,
  RocketOutlined
} from '@ant-design/icons';
import { useParams, useRouter } from 'next/navigation';
import { predictionService } from '@/services/api/prediction';
import { GetPredictionDetailResponse, PredictionQuestion } from '@/types/prediction';

const { Title, Paragraph, Text } = Typography;
const { Panel } = Collapse;

export default function PressDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [detail, setDetail] = useState<GetPredictionDetailResponse | null>(null);
  const [selected, setSelected] = useState(0);

  useEffect(() => {
    if (params.id) {
      fetchDetail(Number(params.id));
    }
  }, [params.id]);

  const fetchDetail = async (id: number) => {
    setLoading(true);
    try {
      const res = await predictionService.getPredictionDetail(id);
      setDetail(res);
    } catch (error) {
      console.error('Failed to fetch prediction detail:', error);
      message.error('获取押题详情失败');
    } finally {
      setLoading(false);
    }
  };

  const parseFollowUp = (followUp: string): string[] => {
    try {
      return JSON.parse(followUp);
    } catch (e) {
      return [];
    }
  };

  const sortedQuestions = useMemo(() => {
    if (!detail?.questions) return [];
    return [...detail.questions].sort((a, b) => a.sort - b.sort);
  }, [detail]);

  const currentQuestion = useMemo(() => {
    return sortedQuestions[selected] || null;
  }, [sortedQuestions, selected]);

  if (loading) {
    return (
      <div className="min-h-screen flex justify-center items-center bg-slate-50">
        <Spin size="large" tip="正在生成详情..." />
      </div>
    );
  }

  if (!detail && !loading) {
     return (
      <div className="min-h-screen flex flex-col justify-center items-center text-slate-500 bg-slate-50">
        <div className="text-6xl mb-4">😕</div>
        <Text className="text-lg mb-4">未找到相关押题数据</Text>
        <Button type="primary" onClick={() => router.back()} className="bg-blue-600">返回列表</Button>
      </div>
    );
  }

  return (
    <div className="min-h-screen relative font-sans bg-slate-50/30">
      {/* Tech Background Elements */}
      <div className="fixed top-0 right-0 w-[800px] h-[800px] bg-indigo-100/40 rounded-full blur-[120px] -translate-y-1/2 translate-x-1/3 pointer-events-none z-0" />
      <div className="fixed bottom-0 left-0 w-[800px] h-[800px] bg-cyan-100/40 rounded-full blur-[120px] translate-y-1/2 -translate-x-1/3 pointer-events-none z-0" />
      
      <div className="container mx-auto px-4 py-8 relative z-10">
        {/* Header Navigation */}
        <div className="mb-8 animate-fade-in-up">
          <Button 
            type="text" 
            icon={<ArrowLeftOutlined />} 
            onClick={() => router.back()}
            className="text-slate-500 hover:text-blue-600 hover:bg-blue-50 mb-4 px-0"
          >
            返回列表
          </Button>
          
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-600 rounded-2xl shadow-lg shadow-blue-200">
              <BookOutlined className="text-2xl text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-slate-900 tracking-tight m-0">
                押题详情
              </h1>
              <p className="text-slate-500 mt-1 text-sm">
                ID: #{params.id} · 共 {sortedQuestions.length} 道预测题目
              </p>
            </div>
          </div>
        </div>

        <Row gutter={[24, 24]} className="animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
          {/* Left Sidebar: Question Directory */}
          <Col xs={24} md={8} lg={6}>
            <div className="sticky top-6">
              <AntCard 
                className="rounded-2xl border-slate-200/60 shadow-xl shadow-slate-200/40 backdrop-blur-xl bg-white/80 overflow-hidden"
                styles={{ body: { padding: 0 } }}
              >
                <div className="p-5 border-b border-slate-100 bg-slate-50/50">
                  <Space align="center" className="font-bold text-slate-700">
                    <MenuOutlined className="text-blue-500" />
                    <span>题目目录</span>
                  </Space>
                </div>
                
                <div className="max-h-[calc(100vh-250px)] overflow-y-auto custom-scrollbar p-3">
                  <List
                    dataSource={sortedQuestions}
                    split={false}
                    renderItem={(item, index) => {
                      const isSelected = index === selected;
                      return (
                        <div
                          onClick={() => setSelected(index)}
                          className={`
                            group relative mb-2 rounded-xl p-4 cursor-pointer transition-all duration-300
                            ${isSelected 
                              ? 'bg-gradient-to-r from-blue-50 to-white border border-blue-100 shadow-sm' 
                              : 'hover:bg-slate-50 border border-transparent hover:border-slate-100'
                            }
                          `}
                        >
                          {isSelected && (
                            <div className="absolute left-0 top-1/2 -translate-y-1/2 h-8 w-1 bg-blue-500 rounded-r-full" />
                          )}
                          
                          <div className="flex items-start gap-3">
                            <span className={`
                              shrink-0 flex items-center justify-center w-6 h-6 rounded-lg text-xs font-mono font-bold mt-0.5
                              ${isSelected 
                                ? 'bg-blue-600 text-white shadow-blue-200' 
                                : 'bg-slate-100 text-slate-400 group-hover:bg-slate-200'
                              }
                            `}>
                              {item.sort}
                            </span>
                            <div className={`text-sm leading-snug line-clamp-2 ${isSelected ? 'text-slate-800 font-medium' : 'text-slate-500'}`}>
                              {item.question}
                            </div>
                          </div>
                        </div>
                      );
                    }}
                  />
                </div>
              </AntCard>
            </div>
          </Col>

          {/* Right Content: Detail View */}
          <Col xs={24} md={16} lg={18}>
            {currentQuestion ? (
              <div className="space-y-6">
                {/* Main Question Card */}
                <AntCard 
                  className="rounded-2xl border-slate-200/60 shadow-xl shadow-slate-200/40 backdrop-blur-xl bg-white/90"
                  styles={{ body: { padding: '2rem' } }}
                >
                  <div className="flex items-start gap-4 mb-8">
                    <div className="shrink-0 flex flex-col items-center gap-1">
                      <span className="flex items-center justify-center w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-mono text-xl font-bold shadow-lg shadow-blue-200">
                        Q{currentQuestion.sort}
                      </span>
                    </div>
                    <div className="flex-1">
                      <Title level={3} className="!text-slate-800 !mb-2 !mt-0 leading-tight">
                        {currentQuestion.question}
                      </Title>
                      {currentQuestion.content && (
                        <div className="mt-4 bg-amber-50/60 p-4 rounded-xl border border-amber-100/50 text-amber-900/80 text-sm leading-relaxed flex gap-3">
                          <div className="shrink-0 mt-0.5 text-amber-500">💡</div>
                          <div>{currentQuestion.content}</div>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                    {/* Focus Section */}
                    <div className="bg-slate-50/80 rounded-2xl p-5 border border-slate-100 hover:border-blue-200 transition-colors group">
                      <div className="flex items-center gap-2 mb-3 text-blue-600">
                        <AimOutlined className="text-lg" />
                        <span className="font-bold text-sm uppercase tracking-wider">考察重点</span>
                      </div>
                      <Paragraph className="text-slate-600 mb-0 leading-relaxed group-hover:text-slate-800 transition-colors">
                        {currentQuestion.focus}
                      </Paragraph>
                    </div>

                    {/* Thinking Path Section */}
                    <div className="bg-slate-50/80 rounded-2xl p-5 border border-slate-100 hover:border-purple-200 transition-colors group">
                      <div className="flex items-center gap-2 mb-3 text-purple-600">
                        <BulbOutlined className="text-lg" />
                        <span className="font-bold text-sm uppercase tracking-wider">答题思路</span>
                      </div>
                      <Paragraph 
                        ellipsis={{ rows: 4, expandable: true, symbol: '展开更多' }}
                        className="text-slate-600 mb-0 leading-relaxed group-hover:text-slate-800 transition-colors"
                      >
                        {currentQuestion.thinking_path}
                      </Paragraph>
                    </div>
                  </div>

                  {/* Reference Answer */}
                  <div className="relative overflow-hidden rounded-2xl bg-slate-900 text-slate-300 shadow-2xl shadow-slate-900/10">
                    <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500" />
                    <div className="p-6">
                      <div className="flex items-center gap-2 mb-4 text-white/90">
                        <CheckCircleOutlined className="text-green-400" />
                        <span className="font-bold">参考答案</span>
                      </div>
                      <div className="prose prose-invert max-w-none prose-p:leading-loose prose-p:text-slate-300">
                        <div className="whitespace-pre-wrap font-sans">
                          {currentQuestion.reference_answer}
                        </div>
                      </div>
                    </div>
                  </div>
                </AntCard>

                {/* Follow Up Section */}
                {parseFollowUp(currentQuestion.follow_up).length > 0 && (
                  <AntCard 
                    className="rounded-2xl border-slate-200/60 shadow-lg shadow-slate-200/30 backdrop-blur-xl bg-white/80"
                    title={
                      <div className="flex items-center gap-2 text-slate-800">
                        <RocketOutlined className="text-indigo-500" />
                        <span>追问建议</span>
                      </div>
                    }
                  >
                    <List
                      dataSource={parseFollowUp(currentQuestion.follow_up)}
                      split={false}
                      renderItem={(item, i) => (
                        <List.Item className="px-0 py-3 hover:bg-slate-50/50 rounded-lg transition-colors px-4 -mx-4">
                          <div className="flex gap-4 items-start w-full">
                            <span className="shrink-0 flex items-center justify-center w-6 h-6 rounded-full bg-indigo-50 text-indigo-600 text-xs font-bold mt-0.5">
                              {i + 1}
                            </span>
                            <span className="text-slate-600 leading-relaxed">{item}</span>
                          </div>
                        </List.Item>
                      )}
                    />
                  </AntCard>
                )}
              </div>
            ) : (
              <div className="h-full min-h-[400px] flex flex-col items-center justify-center text-slate-400 bg-white/50 rounded-3xl border border-slate-100 border-dashed">
                <BookOutlined className="text-4xl mb-4 opacity-50" />
                <p>请从左侧选择题目查看详情</p>
              </div>
            )}
          </Col>
        </Row>
      </div>
    </div>
  );
}
