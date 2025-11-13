'use client';

import { useState } from 'react';
import { Typography, Row, Col, Card as AntCard, Slider, DatePicker, Select, Empty, Button } from 'antd';
import { CheckCircleOutlined } from '@ant-design/icons';
import Link from 'next/link';

const { Title } = Typography;
const { RangePicker } = DatePicker;

export default function InterviewRecordsPage() {
  const [scoreRange, setScoreRange] = useState<[number, number]>([0, 100]);
  const [filter, setFilter] = useState('全部');

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">面试记录</Title>

      <AntCard className="rounded-2xl mt-2" styles={{ body: { padding: 20 } }}>
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} md={12}>
            <div className="grid grid-cols-2 gap-6 items-center">
              <div className="text-center">
                <div className="text-3xl font-semibold">0</div>
                <div className="text-gray-500">已完成面试(次)</div>
              </div>
              <div className="text-center">
                <div className="text-3xl font-semibold">0</div>
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
        </div>
      </div>
    </div>
  );
}

