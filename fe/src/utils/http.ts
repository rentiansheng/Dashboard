import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { message } from 'antd';

// 定义通用的API响应接口
export interface ApiResponse<T = any> {
  retcode: number;
  message: string;
  data?: T;
}


// 创建axios实例
const httpClient: AxiosInstance = axios.create({
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
httpClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加token等认证信息
    // const token = localStorage.getItem('token');
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
httpClient.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response;
    
    // 检查retcode是否不等于0
    if (data.retcode !== 0) {
      // 在页面顶部显示错误消息
      message.error(data.message || '请求失败');
      // 可以根据需要决定是否reject这个promise
      return Promise.reject(new Error(data.message || '请求失败'));
    }
    
    return response;
  },
  (error) => {
    // 处理网络错误或其他HTTP错误
    let errorMessage = '网络错误，请稍后重试';
    
    if (error.response) {
      // 服务器返回了错误状态码
      const { status, data } = error.response;
      errorMessage = data?.message || `请求失败 (${status})`;
    } else if (error.request) {
      // 请求已发出但没有收到响应
      errorMessage = '网络连接失败，请检查网络连接';
    } else {
      // 其他错误
      errorMessage = error.message || '未知错误';
    }
    
    message.error(errorMessage);
    return Promise.reject(error);
  }
);

// 封装常用的HTTP方法
export const http = {
  get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return httpClient.get<ApiResponse<T>>(url, config).then(res => res.data);
  },
  
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return httpClient.post<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
  
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return httpClient.put<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
  
  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return httpClient.delete<ApiResponse<T>>(url, config).then(res => res.data);
  },
  
  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return httpClient.patch<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
};

export default http; 