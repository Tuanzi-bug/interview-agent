'use client';

import { useEffect } from 'react';
import { Button, Spin } from 'antd';
import { useUploadProgressStore } from '@/store/uploadProgressStore';

const ResumeUploadFloatingCard = () => {
  const {
    stage,
    status,
    isBackground,
    showModal,
    isModalVisible,
    isConnected,
    progress,
    error,
    syncFromStorage,
  } = useUploadProgressStore();

  const isProcessing = status === 'pending' || status === 'extracting' || status === 'analyzing';

  useEffect(() => {
    syncFromStorage();
  }, [syncFromStorage]);

  if (!isBackground || !isProcessing) {
    return null;
  }

  return (
    <div className="fixed bottom-6 right-6 z-[9999]">
      <div className="bg-white/95 backdrop-blur border border-slate-200 shadow-lg rounded-2xl px-4 py-3 flex items-center gap-3">
        <div className="w-9 h-9 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
          <Spin size="small" />
        </div>
        <div>
          <div className="text-sm font-medium text-slate-800">后台处理中</div>
          <div className="text-xs text-slate-500">{stage || '处理中...'}</div>
          <div className="text-[11px] text-slate-400">
            {progress}% · {isConnected ? '已连接' : '连接中'}
            {error ? ` · ${error}` : ''}
          </div>
        </div>
        <Button size="small" type="primary" onClick={showModal}>
          查看进度
        </Button>
      </div>
    </div>
  );
};

export default ResumeUploadFloatingCard;
