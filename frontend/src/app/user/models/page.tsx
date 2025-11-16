'use client';

import { Typography, Card as AntCard, Table, Button, Space, Tag, Modal, Form, Input, Select, message } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import apiClient from '@/services/api/client';

const { Title } = Typography;

type ModelItem = {
  id: number;
  name: string;
  modelKey: string;
  protocol: string;
  baseURL?: string;
  providerName?: string;
  status?: number;
  createdAt?: number;
};

export default function UserModelsPage() {
  const [list, setList] = useState<ModelItem[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [openCreate, setOpenCreate] = useState(false);
  const [form] = Form.useForm();

  const fetchList = async (p = page, s = pageSize) => {
    setLoading(true);
    try {
      const res: any = await apiClient.get('/user/model/list', { params: { page: p, size: s } });
      const data = res?.data || res;
      const items: ModelItem[] = (data?.list || []).map((it: any) => ({
        id: it.id,
        name: it.name,
        modelKey: it.modelKey,
        protocol: it.protocol,
        baseURL: it.baseURL,
        providerName: it.providerName,
        status: it.status,
        createdAt: it.createdAt,
      }));
      setList(items);
      setTotal(data?.total || items.length);
    } catch (e: any) {
      message.error(e?.response?.data?.message || '加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchList(page, pageSize);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const onCreate = async () => {
    try {
      const v = await form.validateFields();
      const config = {
        apiSecret: v.apiSecret,
        iconURI: v.iconURI,
        temperature: v.temperature,
        maxTokens: v.maxTokens,
        topP: v.topP,
        topK: v.topK,
        timeout: v.timeout,
        functionCall: v.functionCall,
        jsonMode: v.jsonMode,
        inputTokenLimit: v.inputTokenLimit,
        outputTokenLimit: v.outputTokenLimit,
        concurrency: v.concurrency,
      };
      const payload = {
        name: v.name,
        modelKey: v.modelKey,
        protocol: v.protocol,
        providerName: v.providerName,
        baseURL: v.baseURL,
        status: v.status,
        defaultParams: v.defaultParams,
        configJSON: JSON.stringify(config),
      };
      await apiClient.post('/user/create/model', payload);
      message.success('创建成功');
      setOpenCreate(false);
      form.resetFields();
      fetchList(1, pageSize);
      setPage(1);
    } catch (e: any) {
      if (e?.errorFields) return;
      message.error(e?.response?.data?.message || '创建失败');
    }
  };

  const onDelete = async (id: number) => {
    try {
      await apiClient.delete(`/user/model/delete/${id}`);
      message.success('删除成功');
      fetchList(page, pageSize);
    } catch (e: any) {
      message.error(e?.response?.data?.message || '删除失败');
    }
  };

  const columns = useMemo(
    () => [
      { title: '模型名称', dataIndex: 'name' },
      { title: '模型 Key', dataIndex: 'modelKey' },
      { title: '协议', dataIndex: 'protocol', render: (v: string) => <Tag color="blue">{v}</Tag> },
      { title: '提供商', dataIndex: 'providerName' },
      { title: '状态', dataIndex: 'status', render: (v: number) => <Tag color={v === 1 ? 'green' : 'red'}>{v === 1 ? '启用' : '停用'}</Tag> },
      { title: '创建时间', dataIndex: 'createdAt', render: (ts?: number) => (ts ? new Date(ts).toLocaleString() : '-') },
      {
        title: '操作',
        render: (_: any, row: ModelItem) => (
          <Space>
            <Button type="link">编辑</Button>
            <Button type="link" danger onClick={() => onDelete(row.id)}>删除</Button>
          </Space>
        ),
      },
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [page, pageSize]
  );

  return (
    <div className="container mx-auto px-4">
      <Title level={2} className="mt-2">用户模型管理</Title>

      <AntCard className="rounded-2xl mt-2" extra={<Button type="primary" onClick={() => setOpenCreate(true)}>创建模型</Button>}>
        <Table
          rowKey="id"
          loading={loading}
          columns={columns as any}
          dataSource={list}
          pagination={{ current: page, pageSize, total, onChange: setPage, showSizeChanger: true, onShowSizeChange: (_c, s) => setPageSize(s) }}
        />
      </AntCard>

      <Modal open={openCreate} title="创建模型" onCancel={() => setOpenCreate(false)} onOk={onCreate} okText="创建" width={800} styles={{ body: { maxHeight: '70vh', overflowY: 'auto' } }} destroyOnClose>
        <Form form={form} layout="vertical" initialValues={{ protocol: 'ark', providerName: 'OpenAI', status: 1, temperature: 0.7, maxTokens: 2048, topP: 0.9, topK: 40, timeout: 30, functionCall: false, jsonMode: false, inputTokenLimit: 128000, outputTokenLimit: 128000, concurrency: 7 }}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Form.Item label="模型名称" name="name" rules={[{ required: true, message: '请输入模型名称' }]}>
              <Input placeholder="如：My GPT-4 Model" maxLength={100} />
            </Form.Item>
            <Form.Item label="API 秘钥" name="apiSecret">
              <Input.Password placeholder="请输入平台 API Key" maxLength={500} />
            </Form.Item>
            <Form.Item label="模型 Key" name="modelKey" rules={[{ required: true, message: '请输入模型 Key' }]}>
              <Input placeholder="如：gpt-4" maxLength={100} />
            </Form.Item>
            <Form.Item label="提供商名称" name="providerName">
              <Select options={[{ value: 'OpenAI', label: 'OpenAI' }, { value: 'Anthropic', label: 'Anthropic' }, { value: '火山', label: '火山' }, { value: 'Ollama', label: 'Ollama' }]} />
            </Form.Item>
            <Form.Item label="协议" name="protocol" rules={[{ required: true }]}> 
              <Select options={[{ value: 'ark', label: 'ark' }, { value: 'openai', label: 'openai' }, { value: 'ollama', label: 'ollama' }]} />
            </Form.Item>
            <Form.Item label="并发限制" name="concurrency">
              <Input type="number" placeholder="如：7" />
            </Form.Item>
            <Form.Item label="状态" name="status">
              <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '停用' }]} />
            </Form.Item>
            <Form.Item label="基础 URI" name="baseURL">
              <Input placeholder="API 基础接口地址，如：https://api.xxx.com" maxLength={500} />
            </Form.Item>
          </div>
          <Form.Item label="请求参数" name="defaultParams">
            <Input.TextArea placeholder="JSON 形式的默认参数" rows={3} />
          </Form.Item>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Form.Item label="图标 URI" name="iconURI">
              <Input placeholder="如：https://xxx/icon.png" maxLength={200} />
            </Form.Item>
            <Form.Item label="温度 (Temperature)" name="temperature">
              <Input type="number" step="0.1" min={0} max={2} />
            </Form.Item>
            <Form.Item label="最大 Token 数" name="maxTokens">
              <Input type="number" />
            </Form.Item>
            <Form.Item label="Top P" name="topP">
              <Input type="number" step="0.1" min={0} max={1} />
            </Form.Item>
            <Form.Item label="Top K" name="topK">
              <Input type="number" />
            </Form.Item>
            <Form.Item label="超时时间(秒)" name="timeout">
              <Input type="number" />
            </Form.Item>
            <Form.Item label="函数调用" name="functionCall" valuePropName="checked">
              <Select options={[{ value: true, label: '开启' }, { value: false, label: '关闭' }]} />
            </Form.Item>
            <Form.Item label="JSON 模式" name="jsonMode" valuePropName="checked">
              <Select options={[{ value: true, label: '开启' }, { value: false, label: '关闭' }]} />
            </Form.Item>
            <Form.Item label="输入 Token 限制" name="inputTokenLimit">
              <Input type="number" />
            </Form.Item>
            <Form.Item label="输出 Token 限制" name="outputTokenLimit">
              <Input type="number" />
            </Form.Item>
          </div>
        </Form>
      </Modal>
    </div>
  );
}