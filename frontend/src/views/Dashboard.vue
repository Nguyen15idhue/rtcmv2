<template>
  <div class="dashboard">
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">Active Stations</div>
        <div class="stat-value">{{ streamStore.system.active_stations }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Total Frames</div>
        <div class="stat-value">{{ formatNumber(streamStore.system.total_frames) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Dropped</div>
        <div class="stat-value">{{ formatNumber(streamStore.system.total_drops) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">Uptime</div>
        <div class="stat-value">{{ formatUptime(streamStore.system.uptime_seconds) }}</div>
      </div>
    </div>

    <div class="stations-section">
      <div class="section-header">
        <h2>Stations</h2>
        <span class="count">{{ streamStore.stations.length }} stations</span>
      </div>

      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Status</th>
            <th>FPS</th>
            <th>Frames</th>
            <th>Drops</th>
            <th>Last Seen</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="station in streamStore.stations" :key="station.id">
            <td class="id-cell">{{ station.id }}</td>
            <td>{{ station.name || '-' }}</td>
            <td>
              <span :class="['status-badge', station.connected ? 'connected' : 'disconnected']">
                {{ station.connected ? 'Connected' : 'Disconnected' }}
              </span>
            </td>
            <td class="fps-cell">{{ station.fps?.toFixed(1) || '0.0' }}</td>
            <td>{{ formatNumber(station.frames_total || 0) }}</td>
            <td class="drops-cell">{{ formatNumber(station.frames_dropped || 0) }}</td>
            <td class="time-cell">{{ formatTime(station.last_seen) }}</td>
          </tr>
          <tr v-if="streamStore.stations.length === 0">
            <td colspan="7" class="no-data">No stations connected</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="last-update">
      Last updated: {{ lastUpdateFormatted }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useStreamStore } from '../stores/stream'

const streamStore = useStreamStore()

const lastUpdateFormatted = computed(() => {
  if (!streamStore.lastUpdate) return '-'
  return new Date(streamStore.lastUpdate).toLocaleTimeString()
})

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

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 15px;
  margin-bottom: 25px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.stat-label {
  font-size: 0.85rem;
  color: #666;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 600;
  color: #2c3e50;
}

.stations-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
}

.section-header h2 {
  font-size: 1.1rem;
  color: #2c3e50;
}

.count {
  font-size: 0.85rem;
  color: #999;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.data-table th {
  background: #f8f9fa;
  font-weight: 600;
  font-size: 0.8rem;
  text-transform: uppercase;
  color: #666;
}

.data-table tr:hover {
  background: #f8f9fa;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-badge.connected {
  background: #d4edda;
  color: #155724;
}

.status-badge.disconnected {
  background: #f8d7da;
  color: #721c24;
}

.id-cell {
  font-family: monospace;
  font-weight: 600;
}

.fps-cell {
  color: #2980b9;
  font-weight: 600;
}

.drops-cell {
  color: #e74c3c;
}

.time-cell {
  font-size: 0.85rem;
  color: #666;
}

.no-data {
  text-align: center;
  color: #999;
  padding: 40px !important;
}

.last-update {
  text-align: right;
  font-size: 0.8rem;
  color: #999;
  margin-top: 10px;
}
</style>
