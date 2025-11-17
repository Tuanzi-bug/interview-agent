'use client';

import { Typography, Card as AntCard, Table, Button, Space, Tag, Modal, Form, Input, Select, InputNumber, message } from 'antd';
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
  const [openEdit, setOpenEdit] = useState(false);
  const [editForm] = Form.useForm();
  const [editingId, setEditingId] = useState<number | null>(null);

  const fetchList = async (p = page, s = pageSize) => {
    setLoading(true);
    try {
      const res: any = await apiClient.get('/user/model/list', { params: { page: p, size: s } });
      const data = res?.data || res;
      const items: ModelItem[] = (data?.list || []).map((it: any) => ({
        id: it.id ?? it.ID,
        name: it.name ?? it.Name,
        modelKey: it.modelKey ?? it.model_key ?? it.ModelKey,
        protocol: it.protocol ?? it.Protocol,
        baseURL: it.baseURL ?? it.base_url ?? it.BaseURL,
        providerName: it.providerName ?? it.provider_name ?? it.ProviderName,
        status: it.status ?? it.Status,
        createdAt: it.createdAt ?? it.created_at ?? it.CreatedAt,
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
        icon_uri: v.iconURI,
        temperature: v.temperature,
        max_tokens: v.maxTokens,
        top_p: v.topP,
        top_k: v.topK,
        timeout: v.timeout,
        capability: {
          function_call: v.functionCall === true,
          json_mode: v.jsonMode === true,
          input_tokens: v.inputTokenLimit,
          max_tokens: v.outputTokenLimit,
        },
      };
      const payload = {
        name: v.name,
        model_key: v.modelKey,
        protocol: v.protocol,
        base_url: v.baseURL,
        api_key: v.apiSecret,
        provider_name: v.providerName,
        default_params: v.defaultParams || "{}",
        meta_id: v.metaId !== undefined && v.metaId !== null && v.metaId !== '' ? Number(v.metaId) : undefined,
        config_json: JSON.stringify(config),
        scope: 7,
        status: v.status !== undefined && v.status !== null ? Number(v.status) : 1,
      };
      await apiClient.post('/create/user-models', payload);
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
      { title: '创建时间', dataIndex: 'createdAt', render: (ts?: any) => (ts ? (typeof ts === 'number' ? new Date(ts).toLocaleString() : String(ts)) : '-') },
      {
        title: '操作',
        render: (_: any, row: ModelItem) => (
          <Space>
            <Button
              type="link"
              onClick={async () => {
                const t = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
                if (!t) {
                  message.warning('请先登录后再编辑');
                  return;
                }
                try {
                  const res: any = await apiClient.get(`/user/model/details/${row.id}`);
                  const detail = res?.data || res;
                  let cfg: any = {};
                  try {
                    cfg = detail?.config_json ? JSON.parse(detail.config_json) : {};
                  } catch (_e) {
                    cfg = {};
                  }
                  editForm.setFieldsValue({
                    name: detail?.name ?? row.name,
                    apiSecret: '',
                    modelKey: detail?.model_key ?? detail?.modelKey ?? row.modelKey,
                    providerName: detail?.provider_name ?? detail?.providerName ?? row.providerName,
                    protocol: detail?.protocol ?? row.protocol,
                    metaId: detail?.meta_id,
                    status: Number(detail?.status ?? row.status ?? 1),
                    baseURL: detail?.base_url ?? detail?.baseURL ?? row.baseURL,
                    defaultParams: detail?.default_params ?? '',
                    iconURI: cfg?.icon_uri,
                    temperature: cfg?.temperature,
                    maxTokens: cfg?.max_tokens,
                    topP: cfg?.top_p,
                    topK: cfg?.top_k,
                    timeout: cfg?.timeout,
                    functionCall: cfg?.capability?.function_call,
                    jsonMode: cfg?.capability?.json_mode,
                    inputTokenLimit: cfg?.capability?.input_tokens,
                    outputTokenLimit: cfg?.capability?.max_tokens,
                  });
                  setEditingId(row.id);
                  setOpenEdit(true);
                } catch (e: any) {
                  message.error(e?.response?.data?.message || '加载详情失败');
                }
              }}
            >
              编辑
            </Button>
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

      <AntCard className="rounded-2xl mt-2" extra={<Button type="primary" onClick={() => { const t = typeof window !== 'undefined' ? localStorage.getItem('token') : null; if (!t) { message.warning('请先登录后再创建'); return; } setOpenCreate(true); }}>创建模型</Button>}>
        <Table
          rowKey="id"
          loading={loading}
          columns={columns as any}
          dataSource={list}
          pagination={{ current: page, pageSize, total, onChange: setPage, showSizeChanger: true, onShowSizeChange: (_c, s) => setPageSize(s) }}
        />
      </AntCard>

      <Modal open={openCreate} title="创建模型" onCancel={() => setOpenCreate(false)} onOk={onCreate} okText="创建" width={800} styles={{ body: { maxHeight: '70vh', overflowY: 'auto' } }} destroyOnClose>
        <Form form={form} layout="vertical" initialValues={{ protocol: 'ark', providerName: 'OpenAI', status: 1, temperature: 0.7, maxTokens: 2048, topP: 0.9, topK: 40, timeout: 30, functionCall: true, jsonMode: true, inputTokenLimit: 128000, outputTokenLimit: 128000 }}>
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
            <Form.Item label="Meta ID" name="metaId">
              <InputNumber style={{ width: '100%' }} placeholder="如：1001" />
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

      <Modal
        open={openEdit}
        title="编辑模型"
        onCancel={() => setOpenEdit(false)}
        onOk={async () => {
          try {
            const v = await editForm.validateFields();
            const config = {
              icon_uri: v.iconURI,
              temperature: v.temperature,
              max_tokens: v.maxTokens,
              top_p: v.topP,
              top_k: v.topK,
              timeout: v.timeout,
              capability: {
                function_call: v.functionCall === true,
                json_mode: v.jsonMode === true,
                input_tokens: v.inputTokenLimit,
                max_tokens: v.outputTokenLimit,
              },
            };
            const payload: any = {
              name: v.name,
              model_key: v.modelKey,
              protocol: v.protocol,
              base_url: v.baseURL,
              provider_name: v.providerName,
              config_json: JSON.stringify(config),
              scope: 7,
            };
            if (v.defaultParams) payload.default_params = v.defaultParams;
            if (v.metaId !== undefined && v.metaId !== null && v.metaId !== '') payload.meta_id = Number(v.metaId);
            if (v.status !== undefined && v.status !== null) payload.status = Number(v.status);
            if (v.apiSecret) payload.api_key = v.apiSecret;
            if (!editingId) {
              message.error('未选择编辑的模型');
              return;
            }
            await apiClient.put(`/user/model/update/${editingId}`, payload);
            message.success('更新成功');
            setOpenEdit(false);
            editForm.resetFields();
            fetchList(page, pageSize);
          } catch (e: any) {
            if (e?.errorFields) return;
            message.error(e?.response?.data?.message || '更新失败');
          }
        }}
        okText="更新"
        width={800}
        styles={{ body: { maxHeight: '70vh', overflowY: 'auto' } }}
        destroyOnClose
      >
        <Form form={editForm} layout="vertical" initialValues={{ protocol: 'ark', providerName: 'OpenAI', status: 1, temperature: 0.7, maxTokens: 2048, topP: 0.9, topK: 40, timeout: 30, functionCall: true, jsonMode: true, inputTokenLimit: 128000, outputTokenLimit: 128000 }}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Form.Item label="模型名称" name="name" rules={[{ required: true, message: '请输入模型名称' }]}> 
              <Input placeholder="如：My GPT-4 Model" maxLength={100} />
            </Form.Item>
            <Form.Item label="API 秘钥" name="apiSecret">
              <Input.Password placeholder="留空则不更新密钥" maxLength={500} />
            </Form.Item>
            <Form.Item label="模型 Key" name="modelKey" rules={[{ required: true, message: '请输入模型 Key' }]}> 
              <Input placeholder="如：gpt-4" maxLength={100} />
            </Form.Item>
            <Form.Item label="提供商名称" name="providerName" rules={[{ required: true, message: '请输入提供商名称' }]}> 
              <Select options={[{ value: 'OpenAI', label: 'OpenAI' }, { value: 'Anthropic', label: 'Anthropic' }, { value: '火山', label: '火山' }, { value: 'Ollama', label: 'Ollama' }]} />
            </Form.Item>
            <Form.Item label="协议" name="protocol" rules={[{ required: true }]}> 
              <Select options={[{ value: 'ark', label: 'ark' }, { value: 'openai', label: 'openai' }, { value: 'ollama', label: 'ollama' }]} />
            </Form.Item>
            <Form.Item label="Meta ID" name="metaId">
              <InputNumber style={{ width: '100%' }} placeholder="如：1001" />
            </Form.Item>
            <Form.Item label="状态" name="status">
              <Select options={[{ value: 1, label: '启用' }, { value: 0, label: '停用' }]} />
            </Form.Item>
            <Form.Item label="基础 URI" name="baseURL" rules={[{ required: true, message: '请输入基础 URI' }]}> 
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