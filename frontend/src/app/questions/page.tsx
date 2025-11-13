'use client';

import { useState, useMemo } from 'react';
import { Typography, Row, Col, Card as AntCard, Tag, Avatar, Button } from 'antd';
import { DatabaseOutlined, CodeOutlined, CloudOutlined, RobotOutlined, DeploymentUnitOutlined, ClusterOutlined, ApiOutlined, InboxOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

type Item = {
  key: string;
  title: string;
  desc: string;
  icon: React.ReactNode;
  tags: string[];
  popular?: boolean;
};

const TAGS = [
  '最受欢迎',
  'Java题库',
  'Golang题库',
  'C++题库',
  '计算机基础',
  '数据库',
  '编程语言',
  '前端题库',
  '后端组件',
  '后端工具',
  '场景设计',
  '云原生',
  'AI题库',
  'AI理论',
  'AI编程',
  '全部',
];

const DATA: Item[] = [
  { key: 'java-base', title: 'Java基础', desc: '聚焦Java 核心入门知识，涵盖变量与数据类型、面向对象等', icon: <CodeOutlined className="text-red-500 text-2xl" />, tags: ['Java题库', '编程语言'], popular: true },
  { key: 'java-collection', title: 'Java集合', desc: '聚焦 Java 集合框架（List、Map、Set 等）的实现与应用', icon: <CodeOutlined className="text-red-500 text-2xl" />, tags: ['Java题库', '编程语言'] },
  { key: 'java-concurrency', title: 'Java并发', desc: '剖析 Java 多线程、锁机制、并发工具类与实战', icon: <CodeOutlined className="text-red-500 text-2xl" />, tags: ['Java题库'] },
  { key: 'jvm', title: 'Java虚拟机', desc: '深入 JVM 内存模型、类加载、垃圾回收等核心机制', icon: <CodeOutlined className="text-red-500 text-2xl" />, tags: ['Java题库'] },
  { key: 'go-base', title: 'Golang基础', desc: '涵盖变量类型、函数、goroutine、通道与错误处理', icon: <CodeOutlined className="text-blue-500 text-2xl" />, tags: ['Golang题库', '编程语言'], popular: true },
  { key: 'go-container', title: 'Golang容器', desc: '探索 Go 语言容器（数组、切片、映射等）的实现与性能', icon: <CodeOutlined className="text-blue-500 text-2xl" />, tags: ['Golang题库'] },
  { key: 'go-concurrency', title: 'Golang并发', desc: '围绕 Goroutine、Channel、同步原语等并发模型', icon: <CodeOutlined className="text-blue-500 text-2xl" />, tags: ['Golang题库'] },
  { key: 'go-gc', title: 'GolangGC', desc: '剖析 Go 垃圾回收机制、GC 算法与调优', icon: <CodeOutlined className="text-blue-500 text-2xl" />, tags: ['Golang题库'] },
  { key: 'cpp', title: 'C++', desc: '涵盖面向对象、STL 容器、内存管理等核心知识', icon: <CodeOutlined className="text-indigo-500 text-2xl" />, tags: ['C++题库', '编程语言'] },
  { key: 'redis', title: 'Redis', desc: '覆盖 Redis 数据结构、缓存策略、集群架构与持久化', icon: <DatabaseOutlined className="text-red-500 text-2xl" />, tags: ['数据库', '后端组件'], popular: true },
  { key: 'mysql', title: 'MySQL', desc: '涵盖 SQL 优化、索引设计、事务机制、锁与隔离级别', icon: <DatabaseOutlined className="text-green-600 text-2xl" />, tags: ['数据库'] },
  { key: 'mq', title: '消息队列', desc: '包含 RabbitMQ、Kafka 等主流组件知识', icon: <ClusterOutlined className="text-blue-600 text-2xl" />, tags: ['后端组件'] },
  { key: 'ssm', title: 'SSM全家桶', desc: '涵盖 Spring IoC AOP、Spring MVC 请求流程与实战', icon: <DeploymentUnitOutlined className="text-green-600 text-2xl" />, tags: ['后端组件'] },
  { key: 'os', title: '操作系统', desc: '解析进程线程、内存调度、文件系统、IO 等', icon: <InboxOutlined className="text-amber-600 text-2xl" />, tags: ['计算机基础'] },
  { key: 'network', title: '计算机网络', desc: '梳理 TCP/IP 协议栈、网络分层、链路与安全', icon: <ApiOutlined className="text-pink-600 text-2xl" />, tags: ['计算机基础'] },
  { key: 'scenario', title: '场景题', desc: '聚焦高并发、数据库存储、缓存策略等核心场景题', icon: <CloudOutlined className="text-teal-600 text-2xl" />, tags: ['场景设计'] },
  { key: 'cloud-native', title: '云原生', desc: '容器编排、服务治理、可观测性与CI/CD流水线', icon: <CloudOutlined className="text-cyan-600 text-2xl" />, tags: ['云原生'] },
  { key: 'ai-coding', title: 'AI编程', desc: '模型调用、提示工程、工具函数与自动化工作流', icon: <RobotOutlined className="text-purple-600 text-2xl" />, tags: ['AI编程', 'AI题库'] },
  { key: 'ai-theory', title: 'AI理论', desc: '机器学习基础、优化算法、深度学习核心概念', icon: <RobotOutlined className="text-purple-600 text-2xl" />, tags: ['AI理论', 'AI题库'] },
];

export default function QuestionsPage() {
  const [active, setActive] = useState<string>('最受欢迎');

  const filtered = useMemo(() => {
    if (active === '全部') return DATA;
    if (active === '最受欢迎') return DATA.filter(i => i.popular);
    if (active === '编程语言') return DATA.filter(i => i.tags.includes('Java题库') || i.tags.includes('Golang题库') || i.tags.includes('C++题库'));
    return DATA.filter(i => i.tags.includes(active));
  }, [active]);

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">专项面试试题库</Title>
      <div className="bg-green-50 rounded-2xl p-6 mt-4">
        <div className="flex items-center gap-6">
          <div className="relative w-20 h-20">
            <Avatar className="w-20 h-20" style={{ backgroundColor: '#fff' }} src="https://api.dicebear.com/7.x/miniavs/svg?seed=cow" />
            <div className="absolute -top-2 -left-2 bg-green-500 text-white text-xs px-3 py-1 rounded-full shadow">快来选择你的题库</div>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 flex-1">
            {TAGS.map(t => (
              <button
                key={t}
                onClick={() => setActive(t)}
                className={`rounded-full px-4 py-2 border transition text-sm ${active === t ? 'bg-green-500 text-white border-green-500' : 'bg-white text-gray-700 hover:border-green-400'}`}
              >
                {t}
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className="mt-6">
        <Row gutter={[24, 24]}>
          {filtered.map(item => (
            <Col xs={24} md={12} key={item.key}>
              <AntCard className="rounded-2xl">
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 rounded-lg bg-gray-100 flex items-center justify-center">
                    {item.icon}
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <Title level={4} className="m-0">{item.title}</Title>
                      {item.popular && <Tag color="green">热门</Tag>}
                    </div>
                    <Paragraph className="text-gray-600">{item.desc}</Paragraph>
                  </div>
                </div>
              </AntCard>
            </Col>
          ))}
        </Row>
      </div>
    </div>
  );
}

