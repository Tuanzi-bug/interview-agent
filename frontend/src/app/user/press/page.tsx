'use client';

import { useState } from 'react';
import { Typography, Table, DatePicker, Select, Input, Space, Tag, Button, Card as AntCard } from 'antd';

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
  { key: 1, resume: '我的简历_v1.pdf', status: '已出题', type: '精准岗位押题', level: '入门', company: '字节跳动', job: 'Java后端开发', time: '2024-11-01 19:30' },
  { key: 2, resume: '校招版.pdf', status: '进行中', type: '让机器面试题', level: '中级', company: '美团', job: 'Golang开发', time: '2024-11-05 09:10' },
  { key: 3, resume: '我的简历_v1.pdf', status: '已出题', type: '精准岗位押题', level: '进阶', company: '阿里巴巴', job: '后端架构', time: '2024-11-07 14:22' },
];

export default function PressRecordsPage() {
  const [selectedKeys, setSelectedKeys] = useState<number[]>([]);
  const [status, setStatus] = useState<string>('全部状态');
  const [company, setCompany] = useState<string>('');

  const filtered = DATA.filter(r => (status === '全部状态' || r.status === status) && (company === '' || r.company.includes(company)));

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">押题记录</Title>

      <AntCard className="rounded-2xl mt-2">
        <Space size="middle" wrap>
          <RangePicker />
          <Select value={status} onChange={setStatus} options={[{ value: '全部状态', label: '全部状态' }, { value: '已出题', label: '已出题' }, { value: '进行中', label: '进行中' }, { value: '失败', label: '失败' }]} />
          <Input placeholder="公司名称" value={company} onChange={e => setCompany(e.target.value)} style={{ width: 200 }} />
        </Space>

        <Table
          className="mt-4"
          rowSelection={{
            selectedRowKeys: selectedKeys,
            onChange: keys => setSelectedKeys(keys as number[]),
          }}
          columns={[
            { title: '使用的简历', dataIndex: 'resume' },
            { title: '状态', dataIndex: 'status', render: (v: RecordItem['status']) => <Tag color={v === '已出题' ? 'green' : v === '进行中' ? 'blue' : 'red'}>{v}</Tag> },
            { title: '押题类型', dataIndex: 'type' },
            { title: '难度等级', dataIndex: 'level' },
            { title: '公司名称', dataIndex: 'company' },
            { title: '岗位名称', dataIndex: 'job' },
            { title: '押题时间', dataIndex: 'time' },
            { title: '操作', render: (_: any, row: any) => <Space><a href={`/user/press/${row.key}`}>查看详情</a><Button type="link">继续押题</Button></Space> },
          ]}
          dataSource={filtered}
          pagination={{ pageSize: 20 }}
        />

        <div className="text-gray-500 text-sm mt-2">已选 {selectedKeys.length} 条记录</div>
      </AntCard>
    </div>
  );
}

