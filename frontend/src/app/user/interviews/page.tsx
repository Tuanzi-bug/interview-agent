'use client';

import { useEffect, useMemo, useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Select, Empty, Button, Tag, message, Pagination, Spin } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';
import Link from 'next/link';
import apiClient from '@/services/api/client';

const { Title } = Typography;

export default function InterviewRecordsPage() {
  const [filter, setFilter] = useState('全部');
  const [list, setList] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1); // 前端分页当前页
  const pageSize = 6; // 前端分页：每页显示6条

  const fetchList = async () => {
    setLoading(true);
    try {
      // 一次性获取所有数据，然后在前端分页
      const res: any = await apiClient.get('http://localhost:8888/api/interview/records', { params: { page: 1, page_size: 1000 } });
      const data = res?.data || res;
      const items = (data?.records || []).map((it: any) => ({
        id: it.id,
        userId: it.user_id,
        title: it.title,
        type: it.type,
        difficulty: it.difficulty,
        domain: it.domain,
        companyName: it.company_name,
        status: it.status,
        createdAt: it.created_at,
        updatedAt: it.updated_at,
      }));
      setList(items);
    } catch (e: any) {
      message.error(e?.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  // 移除 mock 接口调用，改用真实接口

  useEffect(() => {
    fetchList();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const completedCount = useMemo(() => list.filter(it => it.status === 'completed').length, [list]);
  const totalCount = useMemo(() => list.length, [list]);
  
  const filteredList = useMemo(() => {
    let filtered = list;
    
    // 根据类型筛选
    if (filter !== '全部') {
      if (filter === '综合面试') {
        filtered = filtered.filter(it => it.type === '综合面试');
      } else if (filter === '社招' || filter === '校招') {
        filtered = filtered.filter(it => it.domain === filter);
      }
    }
    
    return filtered;
  }, [list, filter]);

  // 前端分页展示的数据
  const paginatedList = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return filteredList.slice(startIndex, endIndex);
  }, [filteredList, currentPage, pageSize]);

  // 处理筛选变化时重置页码
  useEffect(() => {
    setCurrentPage(1);
  }, [filter]);

  // 处理分页变化
  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">面试记录</Title>

      <AntCard className="rounded-2xl mt-2" styles={{ body: { padding: 20 } }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={12}>
            <div className="grid grid-cols-2 gap-6 items-center">
              <div className="text-center">
                <div className="text-3xl font-semibold">{totalCount}</div>
                <div className="text-gray-500">面试总数(次)</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold">{completedCount}</div>
                <div className="text-gray-500">已完成面试(次)</div>
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
          <span>类型筛选：</span>
          <Select value={filter} onChange={setFilter} style={{ width: 200 }} options={[{ value: '全部', label: '全部' }, { value: '综合面试', label: '综合面试' }, { value: '社招', label: '社招' }, { value: '校招', label: '校招' }]} />
        </div>

        <div className="mt-8">
          <Spin spinning={loading}>
          {filteredList.length === 0 ? (
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
            <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {paginatedList.map((it: any) => {
                const statusMap: Record<string, { text: string; color: string }> = {
                  pending: { text: '待面试', color: 'blue' },
                  in_progress: { text: '进行中', color: 'orange' },
                  completed: { text: '已完成', color: 'green' },
                };
                const statusInfo = statusMap[it.status] || { text: it.status, color: 'default' };
                const createdTime = it.createdAt ? new Date(it.createdAt).toLocaleString('zh-CN') : '-';
                
                return (
                  <AntCard key={it.id} className="rounded-2xl" styles={{ body: { padding: 16 } }} style={{ minWidth: 300 }}>
                    <div className="space-y-2 text-sm">
                      <div>面试类型：{it.type || '-'}</div>
                      <div>公司名称：{it.companyName || '-'}</div>
                      <div>难度等级：{it.difficulty || '-'}</div>
                      <div>领域：{it.domain || '-'}</div>
                      <div>完成状态：<Tag color={statusInfo.color}>{statusInfo.text}</Tag></div>
                      <div>创建时间：{createdTime}</div>
                    </div>
                    <div className="mt-4">
                      <Link href={`/user/interviews/results/${it.id}`} className="inline-block">
                        <Button type="primary">查看面试详情</Button>
                      </Link>
                    </div>
                  </AntCard>
                );
              })}
            </div>
            
            {/* 分页组件 */}
            {filteredList.length > pageSize && (
              <div className="flex justify-center mt-8">
                <Pagination
                  current={currentPage}
                  total={filteredList.length}
                  pageSize={pageSize}
                  onChange={handlePageChange}
                  showSizeChanger={false}
                  showTotal={(total) => `共 ${total} 条记录`}
                />
              </div>
            )}
            </>
          )}
          </Spin>
        </div>
      </div>
    </div>
  );
}

