'use client';

import { Layout, Typography, Button, Badge, Dropdown } from 'antd';
import Link from 'next/link';
import { BellOutlined, UserOutlined, DownOutlined } from '@ant-design/icons';
import type { FC } from 'react';

const { Header } = Layout;
const { Title } = Typography;

const Navbar: FC = () => {
  return (
    <Header className="sticky top-0 z-50 bg-white shadow-sm border-b">
      <div className="container mx-auto px-4 flex items-center justify-between h-full">
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-white text-lg">牛</span>
          </div>
          <Title level={3} className="m-0">牛面</Title>
        </div>

        <nav className="hidden md:flex items-center space-x-6">
          <Link href="/" className="text-gray-800 hover:text-primary font-medium">首页</Link>
          <Link href="/questions" className="text-gray-700 hover:text-primary">面试题库
            <Badge count={"free"} color="#52c41a" className="ml-2" />
          </Link>
          <Link href="/resume" className="text-gray-700 hover:text-primary">简历押题</Link>
          <Dropdown
            menu={{
              items: [
                { key: 'social', label: <Link href="/interview/social">社招简历面试</Link> },
                { key: 'campus', label: <Link href="/interview/campus">校招简历面试</Link> },
              ],
            }}
          >
            <a className="text-gray-700 hover:text-primary">
              综合面试 <DownOutlined className="ml-1" />
              <Badge count={"hot"} color="#fa541c" className="ml-2" />
            </a>
          </Dropdown>
          <Link href="/interview/special" className="text-gray-700 hover:text-primary">专项面试</Link>
          <Link href="/" className="text-gray-700 hover:text-primary">邀请有礼</Link>
          <Link href="/" className="text-gray-700 hover:text-primary">使用手册</Link>
        </nav>

        <div className="flex items-center space-x-3">
          <Button className="bg-yellow-300 hover:bg-yellow-400 border-none">充值中心</Button>
          <Button icon={<BellOutlined />} />
          <Dropdown
            trigger={["hover"]}
            menu={{
              items: [
                { key: 'center', label: <Link href="/user/center">个人中心</Link> },
                { key: 'interviews', label: <Link href="/user/interviews">面试记录</Link> },
                { key: 'press', label: <Link href="/user/press">押题记录</Link> },
                { key: 'notes', label: <Link href="/user/notes">笔记列表</Link> },
              ],
            }}
          >
            <Button icon={<UserOutlined />} />
          </Dropdown>
        </div>
      </div>
    </Header>
  );
};

export default Navbar;
