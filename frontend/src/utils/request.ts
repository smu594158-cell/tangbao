import axios from 'axios';
import { message } from 'antd';

// 配置请求拦截器，自动带上 Token
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => {
  return Promise.reject(error);
});

// 配置响应拦截器，处理全局错误和自动登出
axios.interceptors.response.use((response) => {
  // 检查业务自定义的权限错误码 (比如 40004: 未登录或已过期, 40007: 无权限)
  if (response.data && response.data.code === 40004) {
    message.error('登录已过期，请重新登录');
    localStorage.removeItem('token');
    localStorage.removeItem('token_expiry');
    localStorage.removeItem('user');
    if (window.location.pathname !== '/login') {
      window.location.href = '/login';
    }
  } else if (response.data && response.data.code === 40007) {
    message.error('禁止访问：权限不足');
    if (window.location.pathname !== '/') {
      window.location.href = '/';
    }
  }
  return response;
}, (error) => {
  if (error.response) {
    switch (error.response.status) {
      case 401:
        message.error('未授权，请登录');
        localStorage.removeItem('token');
        localStorage.removeItem('token_expiry');
        localStorage.removeItem('user');
        if (window.location.pathname !== '/login') {
          window.location.href = '/login';
        }
        break;
      case 403:
        message.error('拒绝访问');
        break;
      case 429:
        message.error('请求过于频繁，请稍后再试');
        break;
      case 500:
        message.error('服务器内部错误');
        break;
      default:
        message.error(`请求错误: ${error.response.status}`);
    }
  } else {
    message.error('网络连接异常');
  }
  return Promise.reject(error);
});

export default axios;
