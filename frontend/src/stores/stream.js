import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useStreamStore = defineStore('stream', () => {
  const connected = ref(false)
  const lastData = ref(null)
  const lastUpdate = ref(null)
  const stations = ref([])
  const system = ref({
    active_stations: 0,
    total_frames: 0,
    total_drops: 0,
    uptime_seconds: 0
  })

  let eventSource = null

  const connect = () => {
    if (eventSource) {
      eventSource.close()
    }

    eventSource = new EventSource('/api/stream')

    eventSource.onopen = () => {
      connected.value = true
      console.log('SSE connected')
    }

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        lastData.value = data
        lastUpdate.value = Date.now()
        
        if (data.stations) {
          stations.value = data.stations
        }
        if (data.system) {
          system.value = data.system
        }
      } catch (e) {
        console.error('Failed to parse SSE data:', e)
      }
    }

    eventSource.onerror = (error) => {
      console.error('SSE error:', error)
      connected.value = false
      eventSource.close()
      
      setTimeout(() => {
        if (!connected.value) {
          connect()
        }
      }, 3000)
    }
  }

  const disconnect = () => {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
    connected.value = false
  }

  return {
    connected,
    lastData,
    lastUpdate,
    stations,
    system,
    connect,
    disconnect
  }
})
