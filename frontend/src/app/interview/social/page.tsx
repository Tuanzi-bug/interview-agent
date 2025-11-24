'use client';

import { Typography, Row, Col, Card as AntCard, Form, Select, Input, Button, Tag, Upload, message, Modal } from 'antd';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import type { UploadFile } from 'antd/es/upload/interface';
import { CheckCircleOutlined, VideoCameraOutlined, ToolOutlined } from '@ant-design/icons';
import BackendHealthCheck from '@/components/BackendHealthCheck';

const { Title, Paragraph } = Typography;

export default function SocialInterviewPage() {
  const [form] = Form.useForm();
  const [resumeFile, setResumeFile] = useState<UploadFile | null>(null);
  const [starting, setStarting] = useState(false);
  const [diagnosisVisible, setDiagnosisVisible] = useState(false);
  const router = useRouter();

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
                label="上传简历"
                name="resume"
                valuePropName="fileList"
                getValueFromEvent={(e) => (e?.fileList || [])}
                rules={[
                  { required: true, message: '请上传简历文件' },
                  {
                    validator: async (_, value) => {
                      const f: UploadFile | undefined = resumeFile || undefined;
                      if (!f) throw new Error('请上传简历文件');
                      const name = String(f.name || '').toLowerCase();
                      const okType = name.endsWith('.pdf') || name.endsWith('.doc') || name.endsWith('.docx');
                      if (!okType) throw new Error('仅支持 PDF/DOC/DOCX 格式');
                      const size = (f.size || 0);
                      if (size <= 0 || size > 2 * 1024 * 1024) throw new Error('文件大小需小于2MB');
                    },
                  },
                ]}
              >
                <Upload.Dragger
                  name="resume"
                  multiple={false}
                  maxCount={1}
                  accept=".pdf,.doc,.docx"
                  beforeUpload={() => false}
                  onChange={({ file, fileList }) => {
                    const f = fileList?.[0] || file;
                    setResumeFile(f || null);
                    const name = String((f?.name) || '').toLowerCase();
                    if (name && !(name.endsWith('.pdf'))) {
                      message.warning('当前后端仅支持PDF文件，DOC/DOCX将无法启动面试');
                    }
                  }}
                  onRemove={() => setResumeFile(null)}
                  disabled={starting}
                >
                  <p className="ant-upload-drag-icon">📄</p>
                  <p className="ant-upload-text">点击或拖拽上传简历（PDF/DOC/DOCX，不超过2MB）</p>
                  <p className="ant-upload-hint">用于生成更贴合你的面试问题</p>
                </Upload.Dragger>
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
                <Button
                  type="primary"
                  className="bg-green-500 w-full h-12 text-base"
                  loading={starting}
                  onClick={async () => {
                    try {
                      await form.validateFields();
                    } catch (e) {
                      message.error('请完善表单后再开始面试');
                      return;
                    }
                    const values = form.getFieldsValue();
                    const params = {
                      type: '综合面试',
                      domain: '社招',
                      difficulty: values.level,
                      position_name: values.job || '',
                      company_name: String(values.company_name || ''),
                    };
                    (window as any).__interviewParams = { ...params };
                    (window as any).__interviewResume = resumeFile;
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

