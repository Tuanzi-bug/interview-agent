import { useState, useEffect, useCallback, useRef } from 'react';
import { API_BASE_URL } from '@/config/api';

interface ProgressUpdate {
  upload_id: string;
  status: 'pending' | 'extracting' | 'extracted' | 'analyzing' | 'completed' | 'failed';
  progress: number;
  stage: string;
  error_msg?: string;
  resume_id?: number;
}

interface UseResumeUploadProgressReturn {
  progress: number;
  status: string;
  stage: string;
  error: string | null;
  resumeId: number | null;
  isConnected: boolean;
  connect: (uploadId: string) => void;
  disconnect: () => void;
  reset: () => void;
}

export const useResumeUploadProgress = (): UseResumeUploadProgressReturn => {
  const [progress, setProgress] = useState<number>(0);
  const [status, setStatus] = useState<string>('pending');
  const [stage, setStage] = useState<string>('准备中...');
  const [error, setError] = useState<string | null>(null);
  const [resumeId, setResumeId] = useState<number | null>(null);
  const [isConnected, setIsConnected] = useState<boolean>(false);
  
  const eventSourceRef = useRef<EventSource | null>(null);

  const disconnect = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
      setIsConnected(false);
    }
  }, []);

  const reset = useCallback(() => {
    disconnect();
    setProgress(0);
    setStatus('pending');
    setStage('准备中...');
    setError(null);
    setResumeId(null);
  }, [disconnect]);

  const connect = useCallback((uploadId: string) => {
    // Prevent multiple connections
    if (eventSourceRef.current) {
      disconnect();
    }

    const token = localStorage.getItem('token');
    if (!token) {
      setError('未登录，无法建立连接');
      return;
    }

    // Construct URL with token query param since EventSource doesn't support headers
    // Use the API_BASE_URL which might already contain /api suffix
    // If API_BASE_URL ends with /, remove it to avoid double slashes
    const baseUrl = (API_BASE_URL || '').replace(/\/$/, '');
    const url = `${baseUrl}/resume/upload/progress/${uploadId}?token=${token}`;
    
    try {
      const eventSource = new EventSource(url);
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        setIsConnected(true);
        setError(null);
      };

      eventSource.onerror = (e) => {
        // EventSource error handling is limited, usually means connection failed
        console.error('SSE connection error:', e);
        // Don't immediately fail on error as it might be a temporary network issue
        // But if readyState is CLOSED (2), then we should probably fail
        if (eventSource.readyState === 2) {
          setIsConnected(false);
          setError('连接中断，请重试');
          disconnect();
        }
      };

      // Listen for specific events
      eventSource.addEventListener('connected', (event: MessageEvent) => {
        try {
          const data = JSON.parse(event.data);
          // console.log('SSE Connected:', data);
        } catch (e) {
          console.error('Error parsing connected event:', e);
        }
      });

      eventSource.addEventListener('progress', (event: MessageEvent) => {
        try {
          const data: ProgressUpdate = JSON.parse(event.data);
          
          setProgress(data.progress);
          setStatus(data.status);
          setStage(data.stage);
          
          if (data.status === 'failed') {
            setError(data.error_msg || '上传失败');
            disconnect();
          } else if (data.status === 'completed') {
            if (data.resume_id) {
              setResumeId(data.resume_id);
            }
            // Don't disconnect immediately, let the UI handle the completion state
            // or wait for the component to unmount/reset
            disconnect(); 
          }
        } catch (e) {
          console.error('Error parsing progress event:', e);
        }
      });

    } catch (e: any) {
      setError(e.message || '无法建立连接');
      setIsConnected(false);
    }
  }, [disconnect]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  return {
    progress,
    status,
    stage,
    error,
    resumeId,
    isConnected,
    connect,
    disconnect,
    reset
  };
};

export default useResumeUploadProgress;
