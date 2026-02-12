import apiClient from './client';
import { ListPredictionResponse, GetPredictionDetailResponse } from '../../types/prediction';

export const predictionService = {
  getPredictionList: async (
    page: number = 1, 
    size: number = 10,
    status?: string,
    companyName?: string
  ) => {
    const params: Record<string, any> = { page, size };
    
    // Only add optional params if they are provided and not empty
    if (status && status !== '全部状态') {
      params.status = status;
    }
    if (companyName && companyName.trim() !== '') {
      params.company_name = companyName;
    }
    
    return apiClient.get<any, ListPredictionResponse>('/prediction/list', { 
      params 
    });
  },

  getPredictionDetail: async (id: number) => {
    // The user gave `http://localhost:8888/api/prediction/2` which implies GET /prediction/:id
    return apiClient.get<any, GetPredictionDetailResponse>(`/prediction/${id}`);
  }
};
