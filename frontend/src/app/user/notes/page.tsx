'use client';

import { useMemo, useState } from 'react';
import { Typography, Input, Select, Space, Button, Card as AntCard, Table, Tag, Empty } from 'antd';

const { Title } = Typography;

type Note = {
  key: number;
  type: '押题笔记' | '面试笔记';
  title: string;
  source: string;
  time: string;
  tags: string[];
};

const NOTES: Note[] = [
  { key: 1, type: '押题笔记', title: 'Redis持久化要点', source: '简历押题-Redis', time: '2024-11-02 12:40', tags: ['Redis', 'AOF'] },
  { key: 2, type: '面试笔记', title: 'Go并发最佳实践', source: '综合面试-社招', time: '2024-11-06 19:20', tags: ['Go', '并发'] },
  { key: 3, type: '押题笔记', title: 'MySQL索引设计', source: '简历押题-MySQL', time: '2024-11-03 21:05', tags: ['MySQL', '索引'] },
];

export default function NotesPage() {
  const [active, setActive] = useState<'押题笔记' | '面试笔记'>('押题笔记');
  const [keyword, setKeyword] = useState('');
  const [tag, setTag] = useState<string | undefined>();

  const filtered = useMemo(() => {
    return NOTES.filter(n => n.type === active)
      .filter(n => (keyword ? n.title.includes(keyword) || n.source.includes(keyword) : true))
      .filter(n => (tag ? n.tags.includes(tag) : true));
  }, [active, keyword, tag]);

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">笔记记录</Title>

      <div className="mt-2 flex items-center justify-between">
        <Space>
          <Button type={active === '押题笔记' ? 'primary' : 'default'} onClick={() => setActive('押题笔记')}>押题笔记</Button>
          <Button type={active === '面试笔记' ? 'primary' : 'default'} onClick={() => setActive('面试笔记')}>面试笔记</Button>
        </Space>
        <Space wrap>
          <Input placeholder="输入关键词" value={keyword} onChange={e => setKeyword(e.target.value)} style={{ width: 220 }} />
          <Select placeholder="选择标签" allowClear value={tag} onChange={setTag} style={{ width: 180 }} options={[{ value: 'Redis', label: 'Redis' }, { value: 'AOF', label: 'AOF' }, { value: 'Go', label: 'Go' }, { value: '并发', label: '并发' }, { value: 'MySQL', label: 'MySQL' }, { value: '索引', label: '索引' }]} />
          <Button type="primary">搜索</Button>
          <Button onClick={() => { setKeyword(''); setTag(undefined); }}>重置</Button>
        </Space>
      </div>

      <AntCard className="rounded-2xl mt-4">
        {filtered.length === 0 ? (
          <Empty imageStyle={{ height: 120 }} description={`${active}为空`} />
        ) : (
          <Table
            pagination={{ pageSize: 10 }}
            dataSource={filtered}
            columns={[
              { title: '标题', dataIndex: 'title' },
              { title: '来源', dataIndex: 'source' },
              { title: '创建时间', dataIndex: 'time' },
              { title: '标签', dataIndex: 'tags', render: (tags: string[]) => <Space>{tags.map(t => <Tag key={t}>{t}</Tag>)}</Space> },
              { title: '操作', render: () => <Space><Button type="link">查看</Button><Button type="link">编辑</Button></Space> },
            ]}
          />
        )}
      </AntCard>
    </div>
  );
}

