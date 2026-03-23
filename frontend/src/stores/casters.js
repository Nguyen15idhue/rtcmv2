import { defineStore } from 'pinia'
import { ref } from 'vue'
import { castersApi } from '../api'

export const useCastersStore = defineStore('casters', () => {
  const casters = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchCasters = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await castersApi.getAll()
      casters.value = res.data.casters || []
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch casters:', e)
    } finally {
      loading.value = false
    }
  }

  const createCaster = async (data) => {
    try {
      const res = await castersApi.create(data)
      casters.value.push(res.data)
      return res.data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const deleteCaster = async (id) => {
    try {
      await castersApi.delete(id)
      casters.value = casters.value.filter(c => c.id !== id)
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const getCasterById = (id) => {
    return casters.value.find(c => c.id === id)
  }

  return {
    casters,
    loading,
    error,
    fetchCasters,
    createCaster,
    deleteCaster,
    getCasterById
  }
})
