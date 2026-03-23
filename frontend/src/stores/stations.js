import { defineStore } from 'pinia'
import { ref } from 'vue'
import { stationsApi } from '../api'

export const useStationsStore = defineStore('stations', () => {
  const stations = ref([])
  const unassigned = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchStations = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await stationsApi.getAll()
      stations.value = res.data.stations || []
      unassigned.value = res.data.unassigned || []
    } catch (e) {
      error.value = e.message
      console.error('Failed to fetch stations:', e)
    } finally {
      loading.value = false
    }
  }

  const createStation = async (data) => {
    try {
      const res = await stationsApi.create(data)
      stations.value.push(res.data)
      return res.data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const updateStation = async (id, data) => {
    try {
      const res = await stationsApi.update(id, data)
      const index = stations.value.findIndex(s => s.id === id)
      if (index !== -1) {
        stations.value[index] = res.data
      }
      return res.data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const deleteStation = async (id) => {
    try {
      await stationsApi.delete(id)
      stations.value = stations.value.filter(s => s.id !== id)
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const addOutput = async (stationId, data) => {
    try {
      const res = await stationsApi.addOutput(stationId, data)
      const station = stations.value.find(s => s.id === stationId)
      if (station) {
        if (!station.outputs) station.outputs = []
        station.outputs.push(res.data)
      }
      return res.data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const removeOutput = async (stationId, casterId) => {
    try {
      await stationsApi.removeOutput(stationId, casterId)
      const station = stations.value.find(s => s.id === stationId)
      if (station && station.outputs) {
        station.outputs = station.outputs.filter(o => o.caster_id !== casterId)
      }
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  const toggleOutput = async (stationId, casterId) => {
    try {
      const res = await stationsApi.toggleOutput(stationId, casterId)
      const station = stations.value.find(s => s.id === stationId)
      if (station && station.outputs) {
        const output = station.outputs.find(o => o.caster_id === casterId)
        if (output) {
          output.enabled = res.data.enabled
        }
      }
      return res.data
    } catch (e) {
      error.value = e.message
      throw e
    }
  }

  return {
    stations,
    unassigned,
    loading,
    error,
    fetchStations,
    createStation,
    updateStation,
    deleteStation,
    addOutput,
    removeOutput,
    toggleOutput
  }
})
