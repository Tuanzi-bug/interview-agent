'use client';

import { Layout, Typography, Button } from 'antd';
import type { FC } from 'react';

const { Header } = Layout;
const { Title } = Typography;

const Navbar: FC = () => {
  return (
    <Header className="sticky top-0 z-50 bg-white shadow-sm border-b">
      <div className="container mx-auto px-4 flex items-center justify-between h-full">
        <div className="flex items-center">
          <Title level={3} className="m-0 text-primary">
            牛面AI面试
          </Title>
        </div>
        
        <nav className="hidden md:flex items-center space-x-6">
          <a href="/" className="text-gray-800 hover:text-primary font-medium">
            首页
          </a>
          <a href="/" className="text-gray-600 hover:text-primary">
            面试服务
          </a>
          <a href="/" className="text-gray-600 hover:text-primary">
            真题题库
          </a>
          <a href="/" className="text-gray-600 hover:text-primary">
            关于我们
          </a>
        </nav>
        
        <div className="flex items-center space-x-2">
          <Button>登录</Button>
          <Button type="primary">注册</Button>
        </div>
      </div>
    </Header>
  );
};

export default Navbar;