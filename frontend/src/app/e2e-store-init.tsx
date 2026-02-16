'use client';

import { useEffect } from 'react';
import { useUploadProgressStore } from '@/store/uploadProgressStore';

const E2EStoreInit = () => {
  useEffect(() => {
    if (process.env.NODE_ENV !== 'production') {
      (
        window as typeof window & {
          __e2eUploadStore?: ReturnType<typeof useUploadProgressStore.getState>;
        }
      ).__e2eUploadStore = useUploadProgressStore.getState();
    }
  }, []);

  return null;
};

export default E2EStoreInit;
