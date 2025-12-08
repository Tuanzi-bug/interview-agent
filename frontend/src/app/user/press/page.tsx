'use client';

import { useState } from 'react';
import {
  Typography,
  Table,
  DatePicker,
  Select,
  Input,
  Space,
  Tag,
  Button,
  Card as AntCard,
} from 'antd';

const { Title } = Typography;
const { RangePicker } = DatePicker;

type RecordItem = {
  key: number;
  resume: string;
  status: '已出题' | '进行中' | '失败';
  type: '精准岗位押题' | '让机器面试题';
  level: '入门' | '中级' | '进阶';
  company: string;
  job: string;
  time: string;
};

const DATA: RecordItem[] = [
  {
    key: 1,
    resume: '我的简历_v1.pdf',
    status: '已出题',
    type: '精准岗位押题',
    level: '入门',
    company: '字节跳动',
    job: 'Java后端开发',
    time: '2024-11-01 19:30',
  },
  {
    key: 2,
    resume: '校招版.pdf',
    status: '进行中',
    type: '让机器面试题',
    level: '中级',
    company: '美团',
    job: 'Golang开发',
    time: '2024-11-05 09:10',
  },
  {
    key: 3,
    resume: '我的简历_v1.pdf',
    status: '已出题',
    type: '精准岗位押题',
    level: '进阶',
    company: '阿里巴巴',
    job: '后端架构',
    time: '2024-11-07 14:22',
  },
];

export default function PressRecordsPage() {
  const [selectedKeys, setSelectedKeys] = useState<number[]>([]);
  const [status, setStatus] = useState<string>('全部状态');
  const [company, setCompany] = useState<string>('');

  const filtered = DATA.filter(
    (r) =>
      (status === '全部状态' || r.status === status) &&
      (company === '' || r.company.includes(company))
  );

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

        <div
          className="bg-white rounded-3xl p-8 border border-slate-100 shadow-xl shadow-slate-200/50 animate-fade-in-up"
          style={{ animationDelay: '0.1s' }}
        >
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
                { value: '失败', label: '失败' },
              ]}
              className="min-w-[140px] h-10"
              size="large"
              variant="filled"
            />
            <Input
              placeholder="搜索公司名称..."
              value={company}
              onChange={(e) => setCompany(e.target.value)}
              className="w-[240px] h-10 rounded-lg border-slate-200 hover:border-blue-400 focus:border-blue-500"
              variant="filled"
            />
            <div className="flex-1" />
            <Button
              type="primary"
              className="bg-blue-600 hover:bg-blue-500 h-10 px-6 rounded-lg shadow-blue-200"
            >
              查询
            </Button>
          </div>

          <Table
            rowSelection={{
              selectedRowKeys: selectedKeys,
              onChange: (keys) => setSelectedKeys(keys as number[]),
            }}
            columns={[
              {
                title: '使用的简历',
                dataIndex: 'resume',
                render: (text) => <span className="font-medium text-slate-700">{text}</span>,
              },
              {
                title: '状态',
                dataIndex: 'status',
                render: (v: RecordItem['status']) => {
                  const colorMap = {
                    已出题: {
                      color: 'success',
                      bg: 'bg-green-50',
                      text: 'text-green-600',
                      border: 'border-green-100',
                    },
                    进行中: {
                      color: 'processing',
                      bg: 'bg-blue-50',
                      text: 'text-blue-600',
                      border: 'border-blue-100',
                    },
                    失败: {
                      color: 'error',
                      bg: 'bg-red-50',
                      text: 'text-red-600',
                      border: 'border-red-100',
                    },
                  };
                  const style = colorMap[v] || colorMap['进行中'];
                  return (
                    <div
                      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${style.bg} ${style.text} border ${style.border}`}
                    >
                      <span
                        className={`w-1.5 h-1.5 rounded-full mr-1.5 ${v === '已出题' ? 'bg-green-500' : v === '进行中' ? 'bg-blue-500' : 'bg-red-500'}`}
                      ></span>
                      {v}
                    </div>
                  );
                },
              },
              {
                title: '押题类型',
                dataIndex: 'type',
                render: (text) => <span className="text-slate-600">{text}</span>,
              },
              {
                title: '难度等级',
                dataIndex: 'level',
                render: (text) => (
                  <span
                    className={`font-medium ${text === '进阶' ? 'text-purple-600' : text === '中级' ? 'text-blue-600' : 'text-slate-600'}`}
                  >
                    {text}
                  </span>
                ),
              },
              {
                title: '公司名称',
                dataIndex: 'company',
                render: (text) => <span className="font-bold text-slate-800">{text}</span>,
              },
              {
                title: '岗位名称',
                dataIndex: 'job',
                render: (text) => <span className="text-slate-600">{text}</span>,
              },
              {
                title: '押题时间',
                dataIndex: 'time',
                render: (text) => <span className="text-slate-500 text-sm font-mono">{text}</span>,
              },
              {
                title: '操作',
                render: (_: any, row: any) => (
                  <Space size="small">
                    <a
                      href={`/user/press/${row.key}`}
                      className="text-blue-600 hover:text-blue-500 font-medium"
                    >
                      查看详情
                    </a>
                    <span className="text-slate-300">|</span>
                    <Button
                      type="link"
                      className="text-slate-500 hover:text-blue-600 p-0 h-auto font-normal"
                    >
                      继续押题
                    </Button>
                  </Space>
                ),
              },
            ]}
            dataSource={filtered}
            pagination={{
              pageSize: 20,
              className: 'mt-6',
              showTotal: (total) => <span className="text-slate-500">共 {total} 条记录</span>,
            }}
            className="modern-table"
            rowClassName="hover:bg-slate-50 transition-colors"
          />

          {selectedKeys.length > 0 && (
            <div className="mt-4 p-3 bg-blue-50 border border-blue-100 rounded-xl flex items-center gap-2 text-blue-700 text-sm animate-fade-in-up">
              <span className="font-medium">已选择 {selectedKeys.length} 项</span>
              <div className="h-4 w-px bg-blue-200 mx-2" />
              <Button
                size="small"
                type="text"
                className="text-blue-600 hover:text-blue-800 hover:bg-blue-100/50"
              >
                批量导出
              </Button>
              <Button size="small" type="text" danger className="hover:bg-red-50">
                批量删除
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
