'use client';

import { Typography, Row, Col, Card as AntCard, Form, Select, Input, Button, Tag, message, Modal, Spin } from 'antd';
import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { CheckCircleOutlined, VideoCameraOutlined, ToolOutlined, FileOutlined } from '@ant-design/icons';
import BackendHealthCheck from '@/components/BackendHealthCheck';
import apiClient from '@/services/api/client';

const { Title, Paragraph } = Typography;

// 简历信息类型
interface ResumeInfo {
  id: number;
  file_name: string;
}

export default function SocialInterviewPage() {
  const [form] = Form.useForm();
  const [selectedResumeId, setSelectedResumeId] = useState<number | null>(null);
  const [resumes, setResumes] = useState<ResumeInfo[]>([]);
  const [loadingResumes, setLoadingResumes] = useState(false);
  const [starting, setStarting] = useState(false);
  const [diagnosisVisible, setDiagnosisVisible] = useState(false);
  const [modelConfigured, setModelConfigured] = useState<boolean | null>(null);
  const [checkingConfig, setCheckingConfig] = useState<boolean>(false);
  const router = useRouter();

  // 获取用户简历列表
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

  useEffect(() => {
    fetchResumes();
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
  }, [fetchResumes]);

  return (
    <div className="container mx-auto px-4">
      <div className="flex justify-between items-center mt-2 mb-2">
        <Title level={2} style={{ margin: 0 }}>综合面试 · 社招简历面试</Title>
        <Button 
          icon={<ToolOutlined />} 
          onClick={() => setDiagnosisVisible(true)}
        >
          后端服务诊断
        </Button>
      </div>
      <Paragraph className="text-gray-600 max-w-3xl">
        在综合面试模式中，系统会围绕你的简历、项目经历与岗位胜任力，从技术基础、项目落地、设计能力到沟通协作，构建环环追问的真实面试场景，帮助你快速查漏补缺与提升应对能力。
      </Paragraph>
      
      <Modal
        title="后端服务诊断"
        open={diagnosisVisible}
        onCancel={() => setDiagnosisVisible(false)}
        footer={null}
        width={700}
      >
        <BackendHealthCheck />
      </Modal>

      <Row gutter={[24, 24]} className="mt-2">
        <Col xs={24} md={16}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            {[
              '深挖技术本质逻辑，构建环环追问交叉',
              '聚焦架构设计能力，真实场景还原',
              '针对业务问题推演，技术方案落地',
              '逻辑体系梳理完整，洞察核心关键点',
            ].map((t, i) => (
              <div key={i} className="flex items-center gap-2 text-green-700">
                <CheckCircleOutlined />
                <span>{t}</span>
              </div>
            ))}
          </div>

          <AntCard className="rounded-2xl">
            <Form
              form={form}
              layout="vertical"
              initialValues={{ job: 'Java后端开发', level: '入门' }}
            >
              <Form.Item
                label="选择简历"
                name="resume_id"
                rules={[{ required: true, message: '请选择简历' }]}
              >
                <Select
                  placeholder="请选择已上传的简历"
                  loading={loadingResumes}
                  disabled={starting}
                  onChange={(value) => setSelectedResumeId(value)}
                  notFoundContent={loadingResumes ? <Spin size="small" /> : '暂无简历，请先在个人中心上传'}
                  options={resumes.map((r) => ({
                    value: r.id,
                    label: (
                      <div className="flex items-center gap-2">
                        <FileOutlined className="text-red-500" />
                        <span>{r.file_name}</span>
                      </div>
                    ),
                  }))}
                />
              </Form.Item>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="岗位意向" name="job">
                    <Input placeholder="如：Java后端开发" />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item label="难度等级" name="level" rules={[{ required: true, message: '请选择难度等级' }]}> 
                    <Select options={[{ value: '入门', label: '入门' }, { value: '中级', label: '中级' }, { value: '进阶', label: '进阶' }]} />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  {/* <Form.Item label="面试时长" name="duration">
                    <Select options={[{ value: '30min', label: '30分钟' }, { value: '60min', label: '60分钟' }]} />
                  </Form.Item> */}
                </Col>
              </Row>
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item label="目标公司（可选）" name="company_name">
                    <Input placeholder="如：字节跳动" maxLength={100} />
                  </Form.Item>
                </Col>
              </Row>

              <div className="mt-4">
                {checkingConfig ? (
                  <Tag color="default" className="mb-2">正在检查模型配置</Tag>
                ) : modelConfigured ? (
                  <Tag color="green" className="mb-2">模型已配置</Tag>
                ) : (
                  <Tag color="red" className="mb-2">模型未配置</Tag>
                )}
                <Button
                  type="primary"
                  className="bg-green-500 w-full h-12 text-base"
                  loading={starting}
                  disabled={starting || checkingConfig || modelConfigured === false}
                  onClick={async () => {
                    try {
                      await form.validateFields();
                    } catch (e) {
                      message.error('请完善表单后再开始面试');
                      return;
                    }
                    if (!modelConfigured) {
                      message.error('未配置模型，无法开始面试');
                      return;
                    }
                    const values = form.getFieldsValue();
                    const params = {
                      type: '综合面试',
                      domain: '社招',
                      difficulty: values.level,
                      position_name: values.job || '',
                      company_name: String(values.company_name || ''),
                      resume_id: values.resume_id,
                    };
                    (window as any).__interviewParams = { ...params };
                    try { sessionStorage.setItem('interviewParams', JSON.stringify(params)); } catch {}
                    setStarting(true);
                    router.push('/interview/social/start');
                  }}
                >
                  开始面试
                </Button>
                <div className="text-center text-gray-500 text-sm mt-2">1次体验价约等于20次AI陪练，单次2小时题目自动续集</div>
              </div>
            </Form>
          </AntCard>
        </Col>
        <Col xs={24} md={8}>
          <AntCard className="rounded-2xl">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2"><VideoCameraOutlined /><span>功能演示</span></div>
              <Tag color="green">推荐观看</Tag>
            </div>
            <div className="w-full h-48 md:h-60 bg-gray-100 rounded-xl flex items-center justify-center text-gray-600">
              <VideoCameraOutlined className="text-3xl mr-2" />示例视频
            </div>
          </AntCard>
        </Col>
      </Row>
    </div>
  );
}

