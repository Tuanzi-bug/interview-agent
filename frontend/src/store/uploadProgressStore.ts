import { create } from 'zustand';
import { API_BASE_URL } from '@/config/api';

type UploadStatus = 'pending' | 'extracting' | 'extracted' | 'analyzing' | 'completed' | 'failed';

interface ProgressUpdate {
  upload_id: string;
  status: UploadStatus;
  progress: number;
  stage: string;
  error_msg?: string;
  error?: string;
  resume_id?: number;
}

interface UploadProgressState {
  uploadId: string | null;
  progress: number;
  status: UploadStatus;
  stage: string;
  error: string | null;
  resumeId: number | null;
  isConnected: boolean;
  isModalVisible: boolean;
  isBackground: boolean;
  lastMessageAt: number | null;
  connect: (uploadId: string, options?: { showModal?: boolean; asBackground?: boolean }) => void;
  disconnect: () => void;
  reset: () => void;
  showModal: () => void;
  hideToBackground: () => void;
  syncFromStorage: () => void;
  markCompleted: () => void;
  markFailed: (error?: string) => void;
}

const STORAGE_KEY = 'resume_upload_pending';

const getStoredUploadId = () => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(STORAGE_KEY);
};

const storeUploadId = (uploadId: string) => {
  if (typeof window === 'undefined') return;
  localStorage.setItem(STORAGE_KEY, uploadId);
};

const clearStoredUploadId = () => {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(STORAGE_KEY);
};

let eventSource: EventSource | null = null;
let statusRef: UploadStatus = 'pending';
let resetTimer: ReturnType<typeof setTimeout> | null = null;

export const useUploadProgressStore = create<UploadProgressState>((set, get) => ({
  uploadId: null,
  progress: 0,
  status: 'pending',
  stage: '准备中...',
  error: null,
  resumeId: null,
  isConnected: false,
  isModalVisible: false,
  isBackground: false,
  lastMessageAt: null,
  connect: (uploadId, options) => {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }

    if (resetTimer) {
      clearTimeout(resetTimer);
      resetTimer = null;
    }
    const token = localStorage.getItem('token');
    if (!token) {
      set({ error: '未登录，无法建立连接', isConnected: false });
      return;
    }

    const baseUrl = (API_BASE_URL || '/api').replace(/\/$/, '');
    const url = `${baseUrl}/resume/upload/progress/${uploadId}?token=${token}`;

    const showModal = options?.showModal ?? true;
    const asBackground = options?.asBackground ?? false;

    statusRef = 'pending';
    storeUploadId(uploadId);

    set({
      uploadId,
      progress: 0,
      status: 'pending',
      stage: '准备中...',
      error: null,
      resumeId: null,
      isModalVisible: showModal,
      isBackground: asBackground,
      lastMessageAt: Date.now(),
    });

    eventSource = new EventSource(url);

    eventSource.onopen = () => {
      set({ isConnected: true, error: null });
    };

    eventSource.onerror = () => {
      if (eventSource && eventSource.readyState === 2) {
        if (statusRef === 'completed' || statusRef === 'failed') {
          get().disconnect();
          return;
        }
        set({ isConnected: false, error: '连接中断，请重试' });
        get().disconnect();
      }
    };

    eventSource.addEventListener('complete', (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        const status = data?.status as UploadStatus | undefined;
        if (status === 'failed') {
          get().markFailed(data?.error || data?.error_msg || '上传失败');
          return;
        }
        get().markCompleted();
      } catch {
        get().markCompleted();
      }
    });

    eventSource.addEventListener('timeout', () => {
      set({ error: '进度连接超时，请刷新重试', isConnected: false });
    });

    eventSource.addEventListener('progress', (event: MessageEvent) => {
      try {
        const data: ProgressUpdate = JSON.parse(event.data);
        const safeProgress = Math.min(100, Math.max(0, Math.round(Number(data.progress) || 0)));

        statusRef = data.status;

        set({
          progress: safeProgress,
          status: data.status,
          stage: data.stage,
          resumeId: data.resume_id ?? null,
          lastMessageAt: Date.now(),
        });

        if (data.status === 'failed') {
          get().markFailed(data.error_msg || data.error || '上传失败');
        }

        if (data.status === 'completed') {
          get().markCompleted();
        }
      } catch (error) {
        set({ error: '进度解析失败', lastMessageAt: Date.now() });
      }
    });
  },
  disconnect: () => {
    if (eventSource) {
      eventSource.close();
      eventSource = null;
    }
    set({ isConnected: false });
  },
  reset: () => {
    if (resetTimer) {
      clearTimeout(resetTimer);
      resetTimer = null;
    }
    get().disconnect();
    clearStoredUploadId();
    statusRef = 'pending';
    set({
      uploadId: null,
      progress: 0,
      status: 'pending',
      stage: '准备中...',
      error: null,
      resumeId: null,
      isModalVisible: false,
      isBackground: false,
      lastMessageAt: null,
    });
  },
  showModal: () => set({ isModalVisible: true, isBackground: false }),
  hideToBackground: () => set({ isModalVisible: false, isBackground: true }),
  syncFromStorage: () => {
    const pendingId = getStoredUploadId();
    if (pendingId && !get().uploadId) {
      get().connect(pendingId, { showModal: false, asBackground: true });
    }
  },
  markCompleted: () => {
    if (resetTimer) {
      clearTimeout(resetTimer);
      resetTimer = null;
    }
    set({ progress: 100, status: 'completed', stage: '完成', isBackground: false });
    resetTimer = setTimeout(() => {
      get().reset();
    }, 1500);
  },
  markFailed: (errorMessage) => {
    if (resetTimer) {
      clearTimeout(resetTimer);
      resetTimer = null;
    }
    set({
      status: 'failed',
      error: errorMessage || '上传失败',
      isBackground: false,
    });
    resetTimer = setTimeout(() => {
      get().reset();
    }, 1500);
  },
}));

export const getPendingUploadId = getStoredUploadId;
