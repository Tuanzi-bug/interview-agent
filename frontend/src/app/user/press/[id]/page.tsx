'use client';

import { useMemo, useState, useEffect } from 'react';
import { Typography, Row, Col, Card as AntCard, List, Tag, Space, Divider, Skeleton, message } from 'antd';
import { BookOutlined, BulbOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { useParams } from 'next/navigation';
import { predictionService } from '@/services/api/prediction';
import type { PredictionQuestion } from '@/types/prediction';

const { Title, Paragraph, Text } = Typography;

// UI View Model
interface QuestionViewModel {
  title: string;
  question: string;
  idea: { label: string; value: string }[];
  reference: string;
  followups: string[];
}

export default function PressDetailPage() {
  const params = useParams();
  const id = Number(params.id);
  
  const [selected, setSelected] = useState(0);
  const [loading, setLoading] = useState(true);
  const [questions, setQuestions] = useState<QuestionViewModel[]>([]);
  const [recordMetadata, setRecordMetadata] = useState<{
    difficulty: string;
    prediction_type: string;
    job_title: string;
    company: string;
  } | null>(null);

  useEffect(() => {
    if (!id) return;

    const fetchDetail = async () => {
      setLoading(true);
      try {
        const res = await predictionService.getPredictionDetail(id);
        if (res && res.questions) {
          // Store record-level metadata
          setRecordMetadata({
            difficulty: res.difficulty || '中等',
            prediction_type: res.prediction_type || '校招',
            job_title: res.job_title || '后端开发',
            company: res.company || '',
          });

          const viewModels: QuestionViewModel[] = res.questions.map((q: PredictionQuestion) => {
            let followups: string[] = [];
            try {
              // Try parsing if it looks like a JSON array, otherwise treat as single string or empty
              if (q.follow_up && (q.follow_up.startsWith('[') || q.follow_up.startsWith('{'))) {
                 followups = JSON.parse(q.follow_up);
              } else if (q.follow_up) {
                 followups = [q.follow_up];
              }
            } catch (e) {
              console.warn('Failed to parse follow_up:', q.follow_up);
              followups = q.follow_up ? [q.follow_up] : [];
            }

            return {
              title: q.question, // Use question as title
              question: q.question,
              idea: [
                { label: '考察重点', value: q.focus },
                { label: '思考路径', value: q.thinking_path },
                { label: '问题详解', value: q.content },
              ].filter(item => item.value), // Filter out empty values
              reference: q.reference_answer,
              followups: Array.isArray(followups) ? followups : [],
            };
          });
          setQuestions(viewModels);
        }
      } catch (error) {
        console.error('Failed to fetch prediction detail:', error);
        message.error('获取详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchDetail();
  }, [id]);

  const current = useMemo(
    () => questions[selected] || questions[0],
    [selected, questions]
  );

  if (loading) {
     return (
        <div className="min-h-screen container mx-auto px-4 py-12">
            <Skeleton active paragraph={{ rows: 4 }} />
        </div>
     )
  }

  if (!current) {
      return (
        <div className="min-h-screen container mx-auto px-4 py-12 text-center text-slate-500">
            暂无数据
        </div>
      )
  }

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
              <div className="max-h-[calc(100vh-300px)] overflow-y-auto custom-scrollbar p-2">
                <List
                  itemLayout="horizontal"
                  dataSource={questions}
                  split={false}
                  renderItem={(item, index) => {
                    const isSelected = index === selected;
                    return (
                      <div
                        onClick={() => setSelected(index)}
                        className={`
                          group relative cursor-pointer transition-all duration-200 ease-out mb-2 rounded-lg
                          ${isSelected 
                            ? 'bg-blue-50/80 shadow-sm ring-1 ring-blue-100' 
                            : 'hover:bg-slate-50 hover:shadow-sm hover:ring-1 hover:ring-slate-100'}
                        `}
                      >
                        <div className={`
                          absolute left-0 top-2 bottom-2 w-1 rounded-r-full transition-colors duration-200
                          ${isSelected ? 'bg-blue-500' : 'bg-transparent group-hover:bg-slate-200'}
                        `} />

                        <div className="flex gap-3 py-3 pr-3 pl-4">
                          <div className={`
                            flex-shrink-0 font-mono text-sm pt-[2px] w-6 text-right transition-colors
                            ${isSelected ? 'text-blue-600 font-bold' : 'text-slate-400 font-medium group-hover:text-slate-500'}
                          `}>
                            {(index + 1).toString().padStart(2, '0')}
                          </div>

                          <div className="flex-1 min-w-0 space-y-1.5">
                            {/* Question Text */}
                            <div 
                              className={`
                                text-sm leading-relaxed transition-colors
                                ${isSelected ? 'text-slate-800 font-medium' : 'text-slate-600 group-hover:text-slate-900'}
                              `}
                              title={item.question}
                            >
                              <div className={`
                                line-clamp-2
                                ${!isSelected && '[mask-image:linear-gradient(180deg,black_70%,transparent)]'}
                              `}>
                                {item.title}
                              </div>
                            </div>

                            {/* Metadata Tags */}
                            <div className="flex items-center gap-2 flex-wrap">
                              <span className={`
                                inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium
                                transition-all duration-200
                                ${isSelected
                                  ? 'bg-blue-100 text-blue-700 ring-1 ring-blue-200'
                                  : 'bg-slate-100 text-slate-600'
                                }
                              `}>
                                {recordMetadata?.prediction_type || '技术问题'}
                              </span>
                              <span className={`
                                inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium
                                transition-all duration-200
                                ${recordMetadata?.difficulty === '简单' && 'bg-green-100 text-green-700'}
                                ${recordMetadata?.difficulty === '中等' && 'bg-yellow-100 text-yellow-700'}
                                ${recordMetadata?.difficulty === '困难' && 'bg-red-100 text-red-700'}
                                ${!recordMetadata?.difficulty && 'bg-slate-100 text-slate-600'}
                                ${isSelected && 'ring-1 ring-current/20'}
                              `}>
                                {recordMetadata?.difficulty || '中等'}
                              </span>
                            </div>
                          </div>
                        </div>
                      </div>
                    );
                  }}
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
                    <div className="bg-gradient-to-br from-slate-50 to-white rounded-xl p-6 border border-slate-100 text-slate-700 leading-loose whitespace-pre-wrap">
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
