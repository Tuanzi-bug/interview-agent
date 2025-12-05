// API 配置文件
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8888/api';

// 面试相关接口
export const INTERVIEW_API = {
  // 启动面试流
  START_STREAM: `${API_BASE_URL}/mianshi/stream/start`,
  // 结束面试
  END_INTERVIEW: `${API_BASE_URL}/mianshi/interview/end`,
  // 提交答案
  SUBMIT_ANSWER: `${API_BASE_URL}/interview/submit/answer`,
  // 获取答题记录
  GET_ANSWER_RECORD: `${API_BASE_URL}/interview/answer-record`,
};

// 用户相关接口
export const USER_API = {
  LOGIN: `${API_BASE_URL}/user/login`,
  REGISTER: `${API_BASE_URL}/user/register`,
  GET_PROFILE: `${API_BASE_URL}/user/profile`,
};
