<template>
  <div class="dashboard">
    <header>
      <h1>RTCMv2 Relay Dashboard</h1>
      <div class="system-stats">
        <div class="stat">
          <span class="label">Active Stations</span>
          <span class="value">{{ systemStats.active_stations }}</span>
        </div>
        <div class="stat">
          <span class="label">Total Frames</span>
          <span class="value">{{ formatNumber(systemStats.total_frames) }}</span>
        </div>
        <div class="stat">
          <span class="label">Dropped</span>
          <span class="value">{{ formatNumber(systemStats.total_drops) }}</span>
        </div>
        <div class="stat">
          <span class="label">Uptime</span>
          <span class="value">{{ formatUptime(systemStats.uptime_seconds) }}</span>
        </div>
      </div>
    </header>

    <main>
      <div class="stations-header">
        <h2>Stations</h2>
        <button @click="refresh" class="refresh-btn">Refresh</button>
      </div>

      <table class="stations-table">
        <thead>
          <tr>
            <th>Station ID</th>
            <th>Name</th>
            <th>Status</th>
            <th>FPS</th>
            <th>Frames</th>
            <th>Dropped</th>
            <th>Last Seen</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="station in stations" :key="station.station_id">
            <td>{{ station.station_id }}</td>
            <td>{{ station.name }}</td>
            <td>
              <span :class="['status', station.connected ? 'connected' : 'disconnected']">
                {{ station.connected ? 'Connected' : 'Disconnected' }}
              </span>
            </td>
            <td>{{ station.fps.toFixed(1) }}</td>
            <td>{{ formatNumber(station.frames_total) }}</td>
            <td>{{ formatNumber(station.frames_dropped) }}</td>
            <td>{{ formatTime(station.last_seen) }}</td>
          </tr>
          <tr v-if="stations.length === 0">
            <td colspan="7" class="no-data">No stations connected</td>
          </tr>
        </tbody>
      </table>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useDashboardStore } from './stores/dashboard'

const store = useDashboardStore()
const stations = ref([])
const systemStats = ref({
  active_stations: 0,
  total_frames: 0,
  total_drops: 0,
  uptime_seconds: 0
})

let interval = null

onMounted(() => {
  refresh()
  interval = setInterval(refresh, 2000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})

async function refresh() {
  try {
    const [stationsRes, systemRes] = await Promise.all([
      store.fetchStations(),
      store.fetchSystem()
    ])
    stations.value = stationsRes.stations || []
    systemStats.value = systemRes
  } catch (err) {
    console.error('Failed to fetch data:', err)
  }
}

function formatNumber(n) {
  return new Intl.NumberFormat().format(n)
}

function formatUptime(seconds) {
  if (!seconds) return '0s'
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  if (h > 0) return `${h}h ${m}m`
  if (m > 0) return `${m}m ${s}s`
  return `${s}s`
}

function formatTime(nanos) {
  if (!nanos) return '-'
  const date = new Date(nanos / 1000000)
  return date.toLocaleTimeString()
}
</script>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f5f5f5;
  color: #333;
}

.dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

header {
  background: #2c3e50;
  color: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

header h1 {
  font-size: 1.5rem;
  margin-bottom: 15px;
}

.system-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 15px;
}

.stat {
  background: rgba(255,255,255,0.1);
  padding: 12px;
  border-radius: 4px;
}

.stat .label {
  display: block;
  font-size: 0.8rem;
  opacity: 0.8;
  margin-bottom: 4px;
}

.stat .value {
  font-size: 1.5rem;
  font-weight: bold;
}

.stations-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.stations-header h2 {
  font-size: 1.2rem;
}

.refresh-btn {
  background: #3498db;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.refresh-btn:hover { background: #2980b9; }

.stations-table {
  width: 100%;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.stations-table th,
.stations-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.stations-table th {
  background: #ecf0f1;
  font-weight: 600;
}

.status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8rem;
}

.status.connected { background: #27ae60; color: white; }
.status.disconnected { background: #e74c3c; color: white; }

.no-data {
  text-align: center;
  color: #999;
  padding: 30px !important;
}
</style>
