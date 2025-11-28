'use client';

import { Layout, Typography, Button, Badge, Dropdown, Modal, Tabs, Form, Input, message } from 'antd';
import Link from 'next/link';
import { BellOutlined, UserOutlined, DownOutlined } from '@ant-design/icons';
import type { FC } from 'react';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import apiClient from '@/services/api/client';

const { Header } = Layout;
const { Title } = Typography;

const Navbar: FC = () => {
  const router = useRouter();
  const [openAuth, setOpenAuth] = useState(false);
  const [activeKey, setActiveKey] = useState<'login' | 'register'>('login');
  const [authed, setAuthed] = useState(false);
  const [user, setUser] = useState<{ username?: string; email?: string } | null>(null);
  const [loginForm] = Form.useForm();
  const [registerForm] = Form.useForm();

  useEffect(() => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const u = typeof window !== 'undefined' ? localStorage.getItem('user') : null;
    setAuthed(!!token);
    setUser(u ? JSON.parse(u) : null);
  }, []);

  const doLogin = async (values: { email: string; password: string }) => {
    try {
      const res: any = await apiClient.post('/user/login', values);
      const data = res?.data || res;
      const token = data?.token || data?.accessToken;
      if (!token) {
        message.error('登录失败：缺少令牌');
        return;
      }
      localStorage.setItem('token', token);
      try { document.cookie = `token=${token};path=/;max-age=${60 * 60 * 24}`; } catch {}
      if (data?.user) {
        localStorage.setItem('user', JSON.stringify(data.user));
        setUser(data.user);
      } else {
        localStorage.setItem('user', JSON.stringify({ email: values.email }));
        setUser({ email: values.email });
      }
      setAuthed(true);
      setOpenAuth(false);
      message.success('登录成功');
    } catch (e: any) {
      message.error(e?.response?.data?.message || '登录失败');
    }
  };

  const doRegister = async (values: { username: string; email: string; password: string }) => {
    try {
      const data: any = await apiClient.post('/user/register', values);
      const token = data?.token;
      const userData = data?.user;
      if (!token || !userData) {
        message.error('注册失败：返回数据缺失');
        return;
      }
      localStorage.setItem('token', token);
      try { document.cookie = `token=${token};path=/;max-age=${60 * 60 * 24}`; } catch {}
      localStorage.setItem('user', JSON.stringify(userData));
      setUser(userData);
      setAuthed(true);
      setOpenAuth(false);
      message.success('注册并登录成功');
    } catch (e: any) {
      message.error(e?.response?.data?.message || '注册失败');
    }
  };

  const logout = async () => {
    try {
      await apiClient.post('/user/logout', {});
    } catch (e) {
    }
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setAuthed(false);
    setUser(null);
    message.success('已退出登录');
    router.push('/');
  };

  return (
    <Header className="sticky top-0 z-50 bg-white shadow-sm border-b">
      <div className="container mx-auto px-4 flex items-center justify-between h-full">
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
            <span className="text-white text-lg">面</span>
          </div>
          <Title level={3} className="m-0">面试吧</Title>
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
          {authed ? (
            <Dropdown
              trigger={["hover"]}
              menu={{
                items: [
                  { key: 'center', label: <Link href="/user/center">个人中心</Link> },
                  { key: 'interviews', label: <Link href="/user/interviews">面试记录</Link> },
                  { key: 'press', label: <Link href="/user/press">押题记录</Link> },
                  { key: 'notes', label: <Link href="/user/notes">笔记列表</Link> },
                  { key: 'models', label: <Link href="/user/models">用户模型</Link> },
                  { key: 'logout', label: <a onClick={logout}>退出登录</a> },
                ],
              }}
            >
              <Button icon={<UserOutlined />}>{user?.username || user?.email || '用户'}</Button>
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => { setActiveKey('login'); setOpenAuth(true); }}>登录 / 注册</Button>
          )}
        </div>
      </div>
      <Modal
        open={openAuth}
        onCancel={() => setOpenAuth(false)}
        footer={null}
        title="账号登录 / 注册"
        destroyOnClose
      >
        <Tabs activeKey={activeKey} onChange={(k) => setActiveKey(k as 'login' | 'register')} items={[
          {
            key: 'login',
            label: '登录',
            children: (
              <Form form={loginForm} layout="vertical" onFinish={doLogin} initialValues={{ email: '', password: '' }}>
                <Form.Item label="邮箱" name="email" rules={[{ required: true, message: '请输入邮箱' }]}>
                  <Input placeholder="请输入邮箱" />
                </Form.Item>
                <Form.Item label="密码" name="password" rules={[{ required: true, message: '请输入密码' }]}>
                  <Input.Password placeholder="请输入密码" />
                </Form.Item>
                <Button type="primary" htmlType="submit" className="w-full">登录</Button>
              </Form>
            ),
          },
          {
            key: 'register',
            label: '注册',
            children: (
              <Form form={registerForm} layout="vertical" onFinish={doRegister} initialValues={{ username: '', email: '', password: '' }}>
                <Form.Item label="用户名" name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                  <Input placeholder="请输入用户名" />
                </Form.Item>
                <Form.Item label="邮箱" name="email" rules={[{ required: true, message: '请输入邮箱' }]}>
                  <Input placeholder="请输入邮箱" />
                </Form.Item>
                <Form.Item label="密码" name="password" rules={[{ required: true, message: '请输入密码' }]}>
                  <Input.Password placeholder="请输入密码" />
                </Form.Item>
                <Button type="primary" htmlType="submit" className="w-full">注册并登录</Button>
              </Form>
            ),
          },
        ]} />
      </Modal>
    </Header>
  );
};

export default Navbar;
