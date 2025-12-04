'use client';

import { Typography, Row, Col, Card as AntCard, Avatar, Tag, Button, Space, Table, Select, Upload, message, Spin, Popconfirm, Alert } from 'antd';
import { UploadOutlined, FileOutlined, DeleteOutlined, StarOutlined, StarFilled, InboxOutlined } from '@ant-design/icons';
import { useEffect, useState, useCallback } from 'react';
import Link from 'next/link';
import type { UploadProps } from 'antd';
import apiClient from '@/services/api/client';

const { Title, Paragraph, Text } = Typography;
const { Dragger } = Upload;

// 简历信息类型
interface ResumeInfo {
  id: number;
  user_id: number;
  file_name: string;
  file_size: number;
  file_type: string;
  is_default: number;
  created_at: number;
  updated_at: number;
}

const columns = [
  { title: '项目', dataIndex: 'project' },
  { title: '金币变动', dataIndex: 'coin' },
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
  const [resumes, setResumes] = useState<ResumeInfo[]>([]);
  const [uploading, setUploading] = useState(false);
  const [loadingResumes, setLoadingResumes] = useState(false);
  const [modelConfigured, setModelConfigured] = useState<boolean | null>(null);
  const [checkingConfig, setCheckingConfig] = useState<boolean>(false);

  // 获取简历列表
  const fetchResumes = useCallback(async () => {
    setLoadingResumes(true);
    try {
      const data: any = await apiClient.get('/resume/list');
      setResumes(data?.resumes || []);
    } catch (err) {
      console.error('获取简历列表失败:', err);
    } finally {
      setLoadingResumes(false);
    }
  }, []);

  // 上传简历
  const handleUpload = async (file: File) => {
    if (modelConfigured === false) {
      message.error('请先配置模型，否则无法上传简历');
      return false;
    }

    if (resumes.length >= 3) {
      message.warning('最多只能上传 3 份简历');
      return false;
    }

    const formData = new FormData();
    formData.append('resume', file);

    setUploading(true);
    try {
      const res: any = await apiClient.post('/resume/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      message.success('简历上传成功');
      fetchResumes();
    } catch (err: any) {
      message.error(err?.message || '简历上传失败');
    } finally {
      setUploading(false);
    }
    return false; // 阻止默认上传行为
  };

  // 删除简历
  const handleDelete = async (resumeId: number) => {
    try {
      await apiClient.delete(`/resume/${resumeId}`);
      message.success('简历已删除');
      fetchResumes();
    } catch (err: any) {
      message.error(err?.message || '删除失败');
    }
  };

  // 设为默认简历
  const handleSetDefault = async (resumeId: number) => {
    try {
      await apiClient.post('/resume/set-default', { resume_id: resumeId });
      message.success('已设为默认简历');
      fetchResumes();
    } catch (err: any) {
      message.error(err?.message || '设置默认简历失败');
    }
  };

  // 格式化文件大小
  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  };

  // 格式化时间
  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
  };

  // 上传配置
  const uploadProps: UploadProps = {
    name: 'resume',
    accept: '.pdf',
    showUploadList: false,
    beforeUpload: handleUpload,
  };

  useEffect(() => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    setCheckingConfig(true);
    fetch('http://localhost:8888/api/user/model/check', {
      method: 'GET',
      headers: {
        Authorization: token ? `Bearer ${token}` : '',
        'X-Auth-Token': token || '',
      },
    })
      .then(async (res) => {
        const data = await res.json().catch(() => null);
        const configured = !!(data && data.data && data.data.configured);
        setModelConfigured(configured);
      })
      .catch(() => {
        setModelConfigured(false);
      })
      .finally(() => {
        setCheckingConfig(false);
      });
  }, []);

  useEffect(() => {
    (async () => {
      try {
        const data: any = await apiClient.get('/user/profile');
        setProfile(data || null);
      } catch {}
    })();
    fetchResumes();
  }, [fetchResumes]);
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
                  <Tag color="gold">面试吧学员</Tag>
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
                  <div className="text-sm opacity-90">剩余金币</div>
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
                <div className="flex items-center justify-between mb-4">
                  <div className="font-medium">我的简历 ({resumes.length}/3)</div>
                </div>

                {/* 简历列表 */}
                <Spin spinning={loadingResumes}>
                  {resumes.length > 0 && (
                    <div className="space-y-2 mb-4">
                      {resumes.map((resume) => (
                        <div
                          key={resume.id}
                          className="flex items-center justify-between p-2 bg-gray-50 rounded-lg"
                        >
                          <div className="flex items-center gap-2">
                            <FileOutlined className="text-red-500" />
                            <span className="text-gray-800">{resume.file_name}</span>
                          </div>
                          <Popconfirm
                            title="确认删除"
                            description="删除后无法恢复，确定删除吗？"
                            onConfirm={() => handleDelete(resume.id)}
                            okText="确定"
                            cancelText="取消"
                          >
                            <Button
                              type="text"
                              size="small"
                              danger
                              icon={<DeleteOutlined />}
                            />
                          </Popconfirm>
                        </div>
                      ))}
                    </div>
                  )}

                  {/* 上传区域 */}
                  {resumes.length < 3 && (
                    <>
                      {!checkingConfig && modelConfigured === false && (
                        <Alert
                          message="模型未配置"
                          description={
                            <span>
                              无法上传简历，请先去 <Link href="/user/models" className="text-blue-500 underline">用户模型页面</Link> 配置模型
                            </span>
                          }
                          type="warning"
                          showIcon
                          className="mb-4"
                        />
                      )}
                      <Dragger {...uploadProps} disabled={uploading || !modelConfigured || checkingConfig}>
                        <p className="ant-upload-drag-icon">
                          {uploading ? <Spin /> : <InboxOutlined />}
                        </p>
                        <p className="ant-upload-text">
                          {uploading ? '上传中...' : (modelConfigured === false ? '请先配置模型' : '点击或拖拽文件到此区域上传')}
                        </p>
                        <p className="ant-upload-hint text-gray-500">
                          仅支持 PDF 格式，文件大小不超过 10MB
                        </p>
                      </Dragger>
                    </>
                  )}

                  {resumes.length >= 3 && (
                    <div className="text-center text-gray-500 py-4">
                      已达到简历数量上限，如需上传新简历请先删除旧简历
                    </div>
                  )}
                </Spin>
              </AntCard>
            </Col>

            <Col span={24}>
              <AntCard className="rounded-2xl" title="金币记录">
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
