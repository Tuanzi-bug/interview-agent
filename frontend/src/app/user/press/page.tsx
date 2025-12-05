'use client';

import { useEffect, useState } from 'react';
import { Typography, Table, DatePicker, Select, Input, Space, Tag, Button, Card as AntCard, message } from 'antd';
import { useRouter } from 'next/navigation';
import { predictionService } from '@/services/api/prediction';
import { PredictionRecordItem } from '@/types/prediction';
import type { ColumnsType } from 'antd/es/table';

const { Title } = Typography;
const { RangePicker } = DatePicker;

export default function PressRecordsPage() {
  const router = useRouter();
  const [selectedKeys, setSelectedKeys] = useState<React.Key[]>([]);
  const [status, setStatus] = useState<string>('全部状态');
  const [company, setCompany] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<PredictionRecordItem[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const fetchData = async (page: number, size: number) => {
    setLoading(true);
    try {
      const res = await predictionService.getPredictionList(page, size);
      setData(res.list || []);
      setPagination({
        ...pagination,
        current: res.page,
        pageSize: res.size,
        total: Number(res.total),
      });
    } catch (error) {
      console.error('Failed to fetch prediction list:', error);
      message.error('获取押题记录失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(pagination.current, pagination.pageSize);
  }, []);

  const handleTableChange = (newPagination: any) => {
    fetchData(newPagination.current, newPagination.pageSize);
  };

  const columns: ColumnsType<PredictionRecordItem> = [
    {
        title: 'ID',
        dataIndex: 'id',
        width: 80,
        render: (text) => <span className="text-slate-500">#{text}</span>
    },
    {
        title: '岗位',
        dataIndex: 'job_title',
        render: (text) => <span className="font-medium text-slate-700">{text}</span>
    },
    {
        title: '公司',
        dataIndex: 'company',
        render: (text) => text || <span className="text-slate-400">-</span>
    },
    {
        title: '难度',
        dataIndex: 'difficulty',
        render: (text) => {
            let color = 'blue';
            if (text === 'Hard' || text === '困难' || text === '进阶') color = 'red';
            if (text === 'Medium' || text === '中等' || text === '中级') color = 'orange';
            if (text === 'Easy' || text === '简单' || text === '入门') color = 'green';
            return <Tag color={color} bordered={false}>{text}</Tag>;
        }
    },
    {
        title: '类型',
        dataIndex: 'prediction_type',
        render: (text) => <Tag color="cyan" bordered={false}>{text}</Tag>
    },
    {
        title: '语言',
        dataIndex: 'language',
        render: (text) => <Tag bordered={false}>{text}</Tag>
    },
    {
        title: '创建时间',
        dataIndex: 'created_at',
        render: (text) => <span className="text-slate-500 text-sm">{text}</span>
    },
    {
        title: '操作',
        key: 'action',
        render: (_, record) => (
            <Button 
                type="link" 
                className="px-0 text-blue-600 hover:text-blue-500"
                onClick={() => router.push(`/user/press/${record.id}`)}
            >
                查看详情
            </Button>
        )
    }
  ];

  return (
    <div className="min-h-screen relative font-sans">
      {/* Decorative Background */}
      <div className="fixed top-0 right-0 w-[600px] h-[600px] bg-purple-50/60 rounded-full blur-[120px] -translate-y-1/2 translate-x-1/3 pointer-events-none z-0" />
      <div className="fixed bottom-0 left-0 w-[600px] h-[600px] bg-pink-50/60 rounded-full blur-[120px] translate-y-1/2 -translate-x-1/3 pointer-events-none z-0" />

      <div className="container mx-auto px-4 relative z-10">
        <div className="mb-8 animate-fade-in-up">
          <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight">押题记录</h1>
          <p className="text-slate-500 mt-2">回顾你的历史押题，追踪面试预测准确度</p>
        </div>

        <div className="bg-white rounded-3xl p-8 border border-slate-100 shadow-xl shadow-slate-200/50 animate-fade-in-up" style={{ animationDelay: '0.1s' }}>
          <div className="flex flex-wrap items-center gap-4 mb-8 bg-slate-50/50 p-4 rounded-2xl border border-slate-100">
            <RangePicker 
              className="border-slate-200 hover:border-blue-400 rounded-lg h-10" 
              variant="filled"
            />
            <Select 
              value={status} 
              onChange={setStatus} 
              options={[
                { value: '全部状态', label: '全部状态' }, 
                { value: '已出题', label: '已出题' }, 
                { value: '进行中', label: '进行中' }, 
                { value: '失败', label: '失败' }
              ]} 
              className="min-w-[140px] h-10"
              size="large"
              variant="filled"
            />
            <Input 
              placeholder="搜索公司名称..." 
              value={company} 
              onChange={e => setCompany(e.target.value)} 
              className="w-[240px] h-10 rounded-lg border-slate-200 hover:border-blue-400 focus:border-blue-500"
              variant="filled"
            />
            <div className="flex-1" />
            <Button 
                type="primary" 
                className="bg-blue-600 hover:bg-blue-500 h-10 px-6 rounded-lg shadow-blue-200"
                onClick={() => fetchData(1, pagination.pageSize)}
            >
              刷新
            </Button>
          </div>

          <Table
            rowSelection={{
              selectedRowKeys: selectedKeys,
              onChange: keys => setSelectedKeys(keys),
            }}
            columns={columns}
            dataSource={data}
            rowKey="id"
            pagination={{
                current: pagination.current,
                pageSize: pagination.pageSize,
                total: pagination.total,
                showSizeChanger: true,
                showTotal: (total) => `共 ${total} 条记录`,
            }}
            loading={loading}
            onChange={handleTableChange}
          />
        </div>
      </div>
    </div>
  );
}
