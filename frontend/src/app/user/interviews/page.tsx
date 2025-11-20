'use client';

import { useEffect, useMemo, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Slider, DatePicker, Select, Empty, Button, Tag, Space, message, Avatar } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';
import Link from 'next/link';
import apiClient from '@/services/api/client';

const { Title } = Typography;
const { RangePicker } = DatePicker;

export default function InterviewRecordsPage() {
  const [scoreRange, setScoreRange] = useState<[number, number]>([0, 100]);
  const [filter, setFilter] = useState('全部');
  const [list, setList] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [mockList, setMockList] = useState<any[]>([]);

  const fetchList = async (p = page, s = pageSize) => {
    setLoading(true);
    try {
      const res: any = await apiClient.get('/interview/records', { params: { page: p, page_size: s } });
      const data = res?.data || res;
      const items = (data?.records || data?.list || []).map((it: any) => ({
        id: it.id ?? it.ID,
        title: it.title ?? it.Title,
        status: it.status ?? it.Status,
        score: it.score ?? it.Score,
        duration: it.duration ?? it.Duration,
        createdAt: it.created_at ?? it.createdAt ?? it.CreatedAt,
      }));
      setList(items);
      setTotal(data?.total ?? data?.Total ?? items.length);
    } catch (e: any) {
      message.error(e?.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchMockList = async () => {
    setLoading(true);
    try {
      const res: any = await apiClient.get('/interview/records/mock');
      const data = res?.data || res;
      const arr = Array.isArray(data) ? data : (data?.data || []);
      setMockList(arr);
    } catch (e: any) {
      message.error(e?.response?.data?.message || '加载面试记录失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchList(page, pageSize);
    fetchMockList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const completedCount = useMemo(() => mockList.length, [mockList]);
  const averageScore = useMemo(() => {
    const scores = mockList.map(i => Number(i.score)).filter(s => !isNaN(s));
    if (scores.length === 0) return 0;
    return Math.round((scores.reduce((a, b) => a + b, 0) / scores.length) * 10) / 10;
  }, [mockList]);
  const filteredMockList = useMemo(() => {
    const min = scoreRange?.[0] ?? 0;
    const max = scoreRange?.[1] ?? 100;
    return mockList.filter(it => {
      const s = Number(it?.score);
      if (isNaN(s)) return true;
      return s >= min && s <= max;
    });
  }, [mockList, scoreRange]);

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">面试记录</Title>

      <AntCard className="rounded-2xl mt-2" styles={{ body: { padding: 20 } }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={12}>
            <div className="grid grid-cols-2 gap-6 items-center">
              <div className="text-center">
                <div className="text-3xl font-semibold">{completedCount}</div>
                <div className="text-gray-500">已完成面试(次)</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold">{averageScore}</div>
                <div className="text-gray-500">平均得分(分)</div>
              </div>
            </div>
          </Col>
          <Col xs={24} md={12}>
            <div className="space-y-2">
              {[
                '你的面试次数在全站用户中位于靠前 0%',
                '你的近7天没有进行面试',
                '你的面试均分在全站用户中位于靠前 0%',
                '你的提高分数为 0分',
              ].map((t, i) => (
                <div key={i} className="flex items-center gap-2 text-green-700">
                  <CheckCircleOutlined />
                  <span>{t}</span>
                </div>
              ))}
            </div>
          </Col>
        </Row>
      </AntCard>

      <div className="mt-6">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <span>分数筛选</span>
            <div className="w-64">
              <Slider range value={scoreRange} onChange={(v) => setScoreRange(v as [number, number])} />
            </div>
          </div>
          <RangePicker />
          <Select value={filter} onChange={setFilter} options={[{ value: '全部', label: '全部' }, { value: '综合面试', label: '综合面试' }, { value: '社招', label: '社招' }, { value: '校招', label: '校招' }]} />
        </div>

        <div className="mt-8">
          {filteredMockList.length === 0 ? (
            <Empty
              imageStyle={{ height: 120 }}
              description={
                <div>
                  <div>暂时无面试记录</div>
                  <div className="mt-2">
                    可以进行 <Link href="/interview/social">社招简历面试</Link> 或 <Link href="/interview/campus">校招简历面试</Link>
                  </div>
                </div>
              }
            />
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredMockList.map((it: any, idx: number) => {
                const avatar = String(it?.avatar_url || '').replace(/[`\s]/g, '');
                return (
                  <AntCard key={idx} className="rounded-2xl" styles={{ body: { padding: 16 } }} style={{ minWidth: 300, minHeight: 300 }}>
                    <div className="flex items-start gap-3">
                      <Avatar src={avatar} size={56} />
                      <div className="flex-1">
                        <div className="text-xl font-semibold">{it?.interview_type || '-'}</div>
                      </div>
                      <div className="text-base font-medium text-green-700">{`测试分数：${it?.score ?? '-'}`}</div>
                    </div>
                    <div className="mt-4 space-y-2 text-sm">
                      <div>公司名称：{it?.company_name || '-'}</div>
                      <div>岗位名称：{it?.position_name || '-'}</div>
                      <div>难度等级：{it?.difficulty || '-'}</div>
                      <div>简历名称：{it?.resume_name || '-'}</div>
                      <div>测试时间：{it?.interview_time || '-'}</div>
                    </div>
                    <div className="mt-4">
                      <Link href={`/user/interviews/results/${idx + 1}`} className="inline-block">
                        <Button type="primary">查看面试结果</Button>
                      </Link>
                    </div>
                  </AntCard>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

