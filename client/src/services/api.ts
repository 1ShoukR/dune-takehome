import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user_data');
      window.location.href = '/auth';
    }
    return Promise.reject(error);
  }
);

// Auth API calls
export const authAPI = {
  login: async (email: string, password: string) => {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },

  register: async (name: string, email: string, password: string) => {
    const response = await api.post('/auth/register', { name, email, password });
    return response.data;
  },

  getProfile: async () => {
    const response = await api.get('/auth/profile');
    return response.data;
  },
};

export interface FormField {
  id: string;
  type: 'text' | 'textarea' | 'email' | 'number' | 'select' | 'radio' | 'checkbox' | 'rating';
  label: string;
  placeholder?: string;
  required: boolean;
  options?: string[];
  validation?: Record<string, string>;
  order: number;
}

export interface CreateFormRequest {
  title: string;
  description?: string;
  fields: FormField[];
  status?: 'draft' | 'published';
}

export interface SubmitFormResponse {
  [fieldId: string]: unknown; 
}


export const formsAPI = {
  getAllForms: async (status?: string) => {
    const params = status ? `?status=${status}` : '';
    const response = await api.get(`/forms${params}`);
    return response.data;
  },

  createForm: async (formData: CreateFormRequest) => {
    const response = await api.post('/forms', formData);
    return response.data;
  },

  getFormById: async (id: string) => {
    const response = await api.get(`/forms/${id}`);
    return response.data;
  },

  submitFormResponse: async (formId: string, responseData: SubmitFormResponse) => {
    const response = await api.post(`/forms/${formId}/responses`, responseData);
    return response.data;
  },
};

export default api;