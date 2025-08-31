import axios from 'axios';
import Cookies from 'js-cookie';
import toast from 'react-hot-toast';

// Create axios instance
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = Cookies.get('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    console.error('API Error:', error);
    
    if (error.response?.status === 401) {
      // Unauthorized - remove token and redirect to login
      Cookies.remove('auth_token');
      window.location.href = '/login';
      return Promise.reject(error);
    }
    
    if (error.response?.status === 429) {
      toast.error('Rate limit exceeded. Please try again later.');
    }
    
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (email, password) => {
    const response = await api.post('/auth/login', { email, password });
    return response;
  },

  register: async (email, password) => {
    const response = await api.post('/auth/register', { email, password });
    return response;
  },

  verifyToken: async (token) => {
    // This would be implemented on backend
    const response = await api.get('/auth/verify', {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.user;
  },

  logout: () => {
    Cookies.remove('auth_token');
  },
};

// SMS API
export const smsAPI = {
  sendSMS: async (phoneNumber, message, deviceId = null) => {
    const response = await api.post('/api/send-sms', {
      phone_number: phoneNumber,
      message: message,
      device_id: deviceId,
    });
    return response;
  },

  sendBulkSMS: async (phoneNumbers, message, deviceId = null) => {
    const response = await api.post('/api/send-bulk-sms', {
      phone_numbers: phoneNumbers,
      message: message,
      device_id: deviceId,
    });
    return response;
  },

  getHistory: async (page = 1, limit = 10, filters = {}) => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...filters,
    });
    
    const response = await api.get(`/api/sms-history?${params}`);
    return response;
  },

  updateStatus: async (messageId, status) => {
    const response = await api.put(`/api/sms/${messageId}/status`, { status });
    return response;
  },
};

// Call API
export const callAPI = {
  makeCall: async (phoneNumber, deviceId = null) => {
    const response = await api.post('/api/make-call', {
      phone_number: phoneNumber,
      device_id: deviceId,
    });
    return response;
  },

  getHistory: async (page = 1, limit = 10, filters = {}) => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
      ...filters,
    });
    
    const response = await api.get(`/api/call-history?${params}`);
    return response;
  },

  updateStatus: async (callId, status, duration = null) => {
    const response = await api.put(`/api/calls/${callId}/status`, { 
      status, 
      duration 
    });
    return response;
  },
};

// Device API
export const deviceAPI = {
  getDevices: async () => {
    const response = await api.get('/api/devices');
    return response.devices;
  },

  registerDevice: async (deviceId, name, phoneNumber = null) => {
    const response = await api.post('/api/devices', {
      device_id: deviceId,
      name: name,
      phone_number: phoneNumber,
    });
    return response;
  },

  updateDevice: async (deviceId, updates) => {
    const response = await api.put(`/api/devices/${deviceId}`, updates);
    return response;
  },

  deleteDevice: async (deviceId) => {
    const response = await api.delete(`/api/devices/${deviceId}`);
    return response;
  },

  getDeviceStatus: async (deviceId) => {
    const response = await api.get(`/api/devices/${deviceId}/status`);
    return response;
  },
};

// Dashboard API
export const dashboardAPI = {
  getStats: async () => {
    const response = await api.get('/api/dashboard/stats');
    return response;
  },

  getRecentActivity: async (limit = 10) => {
    const response = await api.get(`/api/dashboard/activity?limit=${limit}`);
    return response;
  },

  getAnalytics: async (period = '7d') => {
    const response = await api.get(`/api/dashboard/analytics?period=${period}`);
    return response;
  },
};

// User API (for admin)
export const userAPI = {
  getUsers: async (page = 1, limit = 10) => {
    const response = await api.get(`/api/admin/users?page=${page}&limit=${limit}`);
    return response;
  },

  createUser: async (userData) => {
    const response = await api.post('/api/admin/users', userData);
    return response;
  },

  updateUser: async (userId, updates) => {
    const response = await api.put(`/api/admin/users/${userId}`, updates);
    return response;
  },

  deleteUser: async (userId) => {
    const response = await api.delete(`/api/admin/users/${userId}`);
    return response;
  },

  toggleUserStatus: async (userId) => {
    const response = await api.post(`/api/admin/users/${userId}/toggle-status`);
    return response;
  },
};

// Export default api instance for custom requests
export default api;