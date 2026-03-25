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

    <div class="chart-section">
      <h3>FPS Overview (Last 60 seconds)</h3>
      <LineChart :data="fpsHistory" :height="200" />
    </div>

    <div class="stations-section">
      <div class="section-header">
        <h2>Stations</h2>
        <div class="header-actions">
          <span class="count">{{ filteredStations.length }} stations</span>
          <button class="btn-add" @click="openAddModal">+ Add Station</button>
        </div>
      </div>

      <div class="search-bar">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="Search by ID or name..."
          class="search-input"
        />
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
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="station in filteredStations" :key="station.id" @click="openDetailModal(station)" class="clickable-row">
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
            <td class="actions-cell">
              <button class="btn-action" @click="openEditModal(station)">Edit</button>
              <button class="btn-action danger" @click="confirmDelete(station)">Delete</button>
            </td>
          </tr>
          <tr v-if="filteredStations.length === 0">
            <td colspan="8" class="no-data">
              {{ searchQuery ? 'No stations match your search' : 'No stations connected' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="last-update">
      Last updated: {{ lastUpdateFormatted }}
    </div>

    <StationModal
      v-if="showStationModal"
      :station="editingStation"
      @close="closeModal"
      @save="saveStation"
    />

    <ConfirmDialog
      v-if="showDeleteConfirm"
      title="Delete Station"
      :message="`Are you sure you want to delete station ${deletingStation?.id}?`"
      variant="danger"
      @confirm="deleteStation"
      @cancel="showDeleteConfirm = false"
    />

    <div v-if="showDetailModal" class="modal-overlay" @click.self="closeDetailModal">
      <div class="detail-modal">
        <div class="modal-header">
          <h3>Station {{ detailStation?.id }} Details</h3>
          <button class="close-btn" @click="closeDetailModal">&times;</button>
        </div>
        
        <div class="detail-stats">
          <div class="detail-stat">
            <span class="label">Status</span>
            <span :class="['status-badge', detailStation?.connected ? 'connected' : 'disconnected']">
              {{ detailStation?.connected ? 'Connected' : 'Disconnected' }}
            </span>
          </div>
          <div class="detail-stat">
            <span class="label">FPS</span>
            <span class="value">{{ detailStation?.fps?.toFixed(1) || '0.0' }}</span>
          </div>
          <div class="detail-stat">
            <span class="label">Frames</span>
            <span class="value">{{ formatNumber(detailStation?.frames_total || 0) }}</span>
          </div>
          <div class="detail-stat">
            <span class="label">Drops</span>
            <span class="value danger">{{ formatNumber(detailStation?.frames_dropped || 0) }}</span>
          </div>
        </div>

        <div class="outputs-section">
          <div class="outputs-header">
            <h4>Outputs</h4>
            <button class="btn-add-small" @click="openAddOutputModal">+ Add</button>
          </div>
          
          <table class="outputs-table">
            <thead>
              <tr>
                <th>Caster</th>
                <th>Mountpoint</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="output in (detailStation?.outputs || [])" :key="output.caster_id">
                <td>{{ getCasterName(output.caster_id) }}</td>
                <td>{{ output.mountpoint }}</td>
                <td>
                  <button 
                    class="toggle-btn" 
                    :class="output.enabled ? 'enabled' : 'disabled'"
                    @click="toggleOutput(output)"
                  >
                    {{ output.enabled ? 'Enabled' : 'Disabled' }}
                  </button>
                </td>
                <td>
                  <button class="btn-action danger" @click="removeOutput(output)">Remove</button>
                </td>
              </tr>
              <tr v-if="!detailStation?.outputs?.length">
                <td colspan="4" class="no-outputs">No outputs configured</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <OutputModal
      v-if="showOutputModal"
      :station-id="detailStation?.id"
      :output="editingOutput"
      @close="closeOutputModal"
      @save="saveOutput"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useStreamStore } from '../stores/stream'
import { useStationsStore } from '../stores/stations'
import { useCastersStore } from '../stores/casters'
import StationModal from '../components/StationModal.vue'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import OutputModal from '../components/OutputModal.vue'
import LineChart from '../components/LineChart.vue'

const streamStore = useStreamStore()
const stationsStore = useStationsStore()
const castersStore = useCastersStore()

const searchQuery = ref('')
const showStationModal = ref(false)
const editingStation = ref(null)
const showDeleteConfirm = ref(false)
const deletingStation = ref(null)
const showDetailModal = ref(false)
const detailStation = ref(null)
const showOutputModal = ref(false)
const editingOutput = ref(null)

const fpsHistory = ref([])
const MAX_FPS_HISTORY = 60

const filteredStations = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return streamStore.stations
  
  return streamStore.stations.filter(s => {
    const idMatch = String(s.id).includes(query)
    const nameMatch = (s.name || '').toLowerCase().includes(query)
    return idMatch || nameMatch
  })
})

watch(() => streamStore.lastData, (newData) => {
  if (newData && newData.stations) {
    const now = Date.now()
    newData.stations.forEach(s => {
      fpsHistory.value.push({
        x: now,
        y: s.fps || 0
      })
    })
    if (fpsHistory.value.length > MAX_FPS_HISTORY) {
      fpsHistory.value = fpsHistory.value.slice(-MAX_FPS_HISTORY)
    }
  }
})

const lastUpdateFormatted = computed(() => {
  if (!streamStore.lastUpdate) return '-'
  return new Date(streamStore.lastUpdate).toLocaleTimeString()
})

function openAddModal() {
  editingStation.value = null
  showStationModal.value = true
}

function openEditModal(station) {
  editingStation.value = station
  showStationModal.value = true
}

function closeModal() {
  showStationModal.value = false
  editingStation.value = null
}

async function saveStation(data) {
  try {
    if (editingStation.value) {
      await stationsStore.updateStation(data.id, data)
    } else {
      await stationsStore.createStation(data)
    }
    await stationsStore.fetchStations()
    closeModal()
  } catch (e) {
    alert('Error: ' + e.message)
  }
}

function confirmDelete(station) {
  deletingStation.value = station
  showDeleteConfirm.value = true
}

async function deleteStation() {
  try {
    await stationsStore.deleteStation(deletingStation.value.id)
    await stationsStore.fetchStations()
    showDeleteConfirm.value = false
    deletingStation.value = null
  } catch (e) {
    alert('Error: ' + e.message)
  }
}

function openDetailModal(station) {
  detailStation.value = station
  showDetailModal.value = true
}

function closeDetailModal() {
  showDetailModal.value = false
  detailStation.value = null
}

function openAddOutputModal() {
  editingOutput.value = null
  showOutputModal.value = true
}

function openEditOutputModal(output) {
  editingOutput.value = output
  showOutputModal.value = true
}

function closeOutputModal() {
  showOutputModal.value = false
  editingOutput.value = null
}

async function saveOutput(data) {
  try {
    if (editingOutput.value) {
      await stationsStore.removeOutput(detailStation.value.id, editingOutput.value.caster_id)
    }
    await stationsStore.addOutput(detailStation.value.id, data)
    await stationsStore.fetchStations()
    closeOutputModal()
    const updated = streamStore.stations.find(s => s.id === detailStation.value.id)
    if (updated) {
      detailStation.value = updated
    }
  } catch (e) {
    alert('Error: ' + e.message)
  }
}

async function toggleOutput(output) {
  try {
    await stationsStore.toggleOutput(detailStation.value.id, output.caster_id)
    await stationsStore.fetchStations()
    const updated = streamStore.stations.find(s => s.id === detailStation.value.id)
    if (updated) {
      detailStation.value = updated
    }
  } catch (e) {
    alert('Error: ' + e.message)
  }
}

async function removeOutput(output) {
  if (!confirm('Remove this output?')) return
  try {
    await stationsStore.removeOutput(detailStation.value.id, output.caster_id)
    await stationsStore.fetchStations()
    const updated = streamStore.stations.find(s => s.id === detailStation.value.id)
    if (updated) {
      detailStation.value = updated
    }
  } catch (e) {
    alert('Error: ' + e.message)
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

function getCasterName(casterId) {
  const caster = castersStore.casters.find(c => c.id === casterId)
  return caster ? caster.name : `Caster ${casterId}`
}

castersStore.fetchCasters()
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.count {
  font-size: 0.85rem;
  color: #999;
}

.btn-add {
  background: #27ae60;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
}

.btn-add:hover {
  background: #219a52;
}

.search-bar {
  padding: 12px 20px;
  border-bottom: 1px solid #eee;
}

.search-input {
  width: 100%;
  max-width: 300px;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
}

.search-input:focus {
  outline: none;
  border-color: #3498db;
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

.actions-cell {
  display: flex;
  gap: 8px;
}

.btn-action {
  padding: 4px 8px;
  border: none;
  border-radius: 3px;
  font-size: 0.75rem;
  cursor: pointer;
  background: #e0e0e0;
  color: #333;
}

.btn-action:hover {
  background: #d0d0d0;
}

.btn-action.danger {
  background: #fdecea;
  color: #e74c3c;
}

.btn-action.danger:hover {
  background: #fbdcd7;
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

.clickable-row {
  cursor: pointer;
}

.clickable-row:hover {
  background: #f0f8ff !important;
}

.chart-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 20px;
  margin-bottom: 20px;
}

.chart-section h3 {
  font-size: 1rem;
  color: #333;
  margin-bottom: 15px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.detail-modal {
  background: white;
  border-radius: 8px;
  width: 600px;
  max-width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.detail-modal .modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.detail-modal .modal-header h3 {
  margin: 0;
  font-size: 1.2rem;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #999;
  padding: 0;
  line-height: 1;
}

.detail-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.detail-stat {
  text-align: center;
}

.detail-stat .label {
  display: block;
  font-size: 0.8rem;
  color: #666;
  margin-bottom: 5px;
}

.detail-stat .value {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

.detail-stat .value.danger {
  color: #e74c3c;
}

.outputs-section {
  padding: 20px;
}

.outputs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.outputs-header h4 {
  margin: 0;
  font-size: 1rem;
  color: #333;
}

.btn-add-small {
  background: #3498db;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.btn-add-small:hover {
  background: #2980b9;
}

.outputs-table {
  width: 100%;
  border-collapse: collapse;
}

.outputs-table th,
.outputs-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.outputs-table th {
  background: #f8f9fa;
  font-size: 0.8rem;
  font-weight: 600;
  color: #666;
}

.toggle-btn {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 0.75rem;
  border: none;
  cursor: pointer;
}

.toggle-btn.enabled {
  background: #d4edda;
  color: #155724;
}

.toggle-btn.disabled {
  background: #f8d7da;
  color: #721c24;
}

.no-outputs {
  text-align: center;
  color: #999;
  padding: 20px !important;
}
</style>