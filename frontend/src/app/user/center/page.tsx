'use client';

import { Typography, Row, Col, Card as AntCard, Avatar, Tag, Button, Space, Table, Select } from 'antd';
import { useEffect, useState } from 'react';
import apiClient from '@/services/api/client';

const { Title, Paragraph, Text } = Typography;

const columns = [
  { title: '项目', dataIndex: 'project' },
  { title: '牛币变动', dataIndex: 'coin' },
  { title: '支付金额', dataIndex: 'amount' },
  { title: '交易渠道', dataIndex: 'channel' },
  { title: '交易时间', dataIndex: 'time' },
  { title: '交易单号', dataIndex: 'orderId' },
];

const data = [
  { key: 1, project: '专项面试-Redis', coin: '+20', amount: '¥0.00', channel: '免费体验', time: '2024-10-01 20:12', orderId: 'FREE-001' },
  { key: 2, project: '简历押题', coin: '-10', amount: '¥9.90', channel: '微信支付', time: '2024-11-02 12:45', orderId: 'WX-20241102-123456' },
];

export default function UserCenterPage() {
  const [profile, setProfile] = useState<{ id?: number; username?: string; email?: string } | null>(null);
  useEffect(() => {
    (async () => {
      try {
        const data: any = await apiClient.get('/user/profile');
        setProfile(data || null);
      } catch {}
    })();
  }, []);
  return (
    <div className="container mx-auto px-4">
      <Row gutter={[24, 24]}>
        <Col xs={24} md={8}>
          <AntCard className="rounded-2xl">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Avatar size={64} src="https://api.dicebear.com/7.x/adventurer/svg?seed=LB" />
                <div>
                  <div className="font-medium text-lg">{profile?.username || '未登录'}</div>
                  <Tag color="gold">牛面学员</Tag>
                </div>
              </div>
              
            </div>

            <div className="mt-6 space-y-2 text-sm text-gray-700">
              <div>用户名：{profile?.username ?? '-'}</div>
              <div>邮箱：{profile?.email ?? '-'}</div>
            </div>

            <AntCard className="rounded-2xl mt-6 bg-green-500 text-white" variant="outlined">
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-sm opacity-90">剩余牛币</div>
                  <div className="text-3xl font-semibold">0</div>
                </div>
                <Button>充值</Button>
              </div>
              <div className="mt-4 text-sm space-y-1 opacity-90">
                <div>专项面试：0次</div>
                <div>综合面试：0次</div>
                <div>简历押题：0次</div>
              </div>
            </AntCard>
          </AntCard>
        </Col>

        <Col xs={24} md={16}>
          <Row gutter={[16, 16]}>
            <Col span={24}>
              <AntCard className="rounded-2xl">
                <div className="flex items-center justify-between">
                  <div className="font-medium">我的简历 (0/3)</div>
                  <Select size="small" value="全部" options={[{ value: '全部', label: '全部' }, { value: '校招', label: '校招' }, { value: '社招', label: '社招' }]} />
                </div>
                <div className="mt-3 h-24 border border-dashed rounded-lg flex items-center justify-center text-gray-500">
                  支持上传 pdf、doc、docx 格式
                </div>
              </AntCard>
            </Col>

            <Col span={24}>
              <AntCard className="rounded-2xl" title="牛币记录">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex gap-2">
                    <Select size="small" value="全部类型" options={[{ value: '全部类型', label: '全部类型' }, { value: '收入', label: '收入' }, { value: '支出', label: '支出' }]} />
                    <Select size="small" value="最近30天" options={[{ value: '最近30天', label: '最近30天' }, { value: '最近90天', label: '最近90天' }]} />
                  </div>
                </div>
                <Table size="small" columns={columns} dataSource={data} pagination={{ pageSize: 20 }} />
              </AntCard>
            </Col>
          </Row>
        </Col>
      </Row>
    </div>
  );
}
