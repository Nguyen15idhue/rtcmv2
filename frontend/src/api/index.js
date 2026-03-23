import axios from 'axios'

const api = axios.create({
  baseURL: '',
  timeout: 10000
})

let apiKey = localStorage.getItem('api_key') || ''

api.interceptors.request.use(config => {
  if (apiKey) {
    config.headers['X-API-Key'] = apiKey
  }
  return config
})

api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      console.error('API Key required or invalid')
    }
    return Promise.reject(error)
  }
)

export const setApiKey = (key) => {
  apiKey = key
  localStorage.setItem('api_key', key)
}

export const getApiKey = () => apiKey

export const stationsApi = {
  getAll: () => api.get('/api/stations'),
  create: (data) => api.post('/api/station', data),
  update: (id, data) => api.put(`/api/station/${id}`, data),
  delete: (id) => api.delete(`/api/station/${id}`),
  addOutput: (id, data) => api.post(`/api/station/${id}/output`, data),
  removeOutput: (id, casterId) => api.delete(`/api/station/${id}/output/${casterId}`),
  toggleOutput: (id, casterId) => api.put(`/api/station/${id}/output/${casterId}/toggle`)
}

export const castersApi = {
  getAll: () => api.get('/api/casters'),
  create: (data) => api.post('/api/caster', data),
  delete: (id) => api.delete(`/api/caster/${id}`)
}

export const systemApi = {
  getSystem: () => api.get('/api/system'),
  getConfig: () => api.get('/api/config'),
  reload: () => api.post('/api/reload'),
  health: () => api.get('/api/health')
}

export default api
