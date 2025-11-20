'use client';

import { useMemo, useState } from 'react';
import { Typography, Input, Select, Space, Button, Card as AntCard, Table, Tag, Empty, message } from 'antd';

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

const ANSWERS: Record<number, string> = {
  1: '在Redis持久化方面，AOF用于记录写操作日志，RDB用于快照。结合两者可以在性能与可靠性之间取得平衡：AOF建议使用everysec策略保证数据安全，RDB用于低频全量备份。对于热点与大对象需谨慎持久化，避免阻塞与膨胀。',
  2: 'Go并发最佳实践包括合理使用goroutine与channel、避免共享可变状态、通过context控制生命周期、使用sync包（如WaitGroup、Mutex）做并发协调，并配合worker池与限流实现稳定的吞吐。',
  3: 'MySQL索引设计要点：优先考虑查询条件的选择性；前缀索引优化长字符串；合理使用覆盖索引减少回表；联合索引遵循最左前缀原则；避免在高频更新的低选择性字段上建立索引。',
};

export default function NotesPage() {
  const [active, setActive] = useState<'押题笔记' | '面试笔记'>('押题笔记');
  const [keyword, setKeyword] = useState('');
  const [tag, setTag] = useState<string | undefined>();
  const [notes, setNotes] = useState<Note[]>(NOTES);
  const [expandedKeys, setExpandedKeys] = useState<number[]>([]);

  const filtered = useMemo(() => {
    return notes.filter(n => n.type === active)
      .filter(n => (keyword ? n.title.includes(keyword) || n.source.includes(keyword) : true))
      .filter(n => (tag ? n.tags.includes(tag) : true));
  }, [active, keyword, tag]);

  const toggleExpand = (key: number) => {
    setExpandedKeys(prev => (prev.includes(key) ? prev.filter(k => k !== key) : [...prev, key]));
  };

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
            rowKey="key"
            expandable={{
              expandedRowRender: (row: Note) => (
                <AntCard className="bg-green-50" styles={{ body: { padding: 16 } }}>
                  <Space direction="vertical" style={{ width: '100%' }}>
                    <Space align="center" className="text-green-700">
                      <Tag color="green">答案思路</Tag>
                    </Space>
                    <Space wrap>
                      {row.tags.map(t => (
                        <Tag key={t}>{t}</Tag>
                      ))}
                    </Space>
                    <div className="mt-2" />
                    <Tag>参考答案</Tag>
                    <div>{ANSWERS[row.key] || '暂无答案'}</div>
                  </Space>
                </AntCard>
              ),
              expandedRowKeys: expandedKeys,
              onExpandedRowsChange: keys => setExpandedKeys(keys as number[]),
              expandIcon: () => null,
            }}
            onRow={(row: Note) => ({
              onClick: () => toggleExpand(row.key),
            })}
            pagination={{ pageSize: 10 }}
            dataSource={filtered}
            columns={[
              { title: '标题', dataIndex: 'title', render: (_: any, row: Note) => <span>{row.title}</span> },
              { title: '来源', dataIndex: 'source' },
              { title: '创建时间', dataIndex: 'time' },
              { title: '标签', dataIndex: 'tags', render: (tags: string[]) => <Space>{tags.map(t => <Tag key={t}>{t}</Tag>)}</Space> },
              { title: '操作', render: (_: any, row: Note) => (
                <Space>
                  <Button type="link" danger onClick={(e) => { e.stopPropagation(); setNotes(prev => prev.filter(n => n.key !== row.key)); }}>删除</Button>
                  <Button type="link" onClick={async (e) => { e.stopPropagation(); try { await navigator.clipboard.writeText(ANSWERS[row.key] || ''); message.success('答案已复制'); } catch { message.error('复制失败'); } }}>复制</Button>
                </Space>
              ) },
            ]}
          />
        )}
      </AntCard>
    </div>
  );
}

