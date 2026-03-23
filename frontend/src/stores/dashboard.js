import { defineStore } from 'pinia'
import axios from 'axios'

const api = axios.create({
  baseURL: '/api'
})

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    stations: [],
    system: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchStations() {
      const res = await api.get('/stations')
      return res.data
    },

    async fetchSystem() {
      const res = await api.get('/system')
      return res.data
    },

    async checkHealth() {
      const res = await api.get('/health')
      return res.data
    }
  }
})
