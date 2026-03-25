<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal">
      <div class="modal-header">
        <h3>{{ isEdit ? 'Edit Output' : 'Add Output' }}</h3>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="modal-body">
        <div class="form-group">
          <label>Caster *</label>
          <select v-model.number="formData.caster_id" required>
            <option value="">Select a caster</option>
            <option v-for="caster in casters" :key="caster.id" :value="caster.id">
              {{ caster.name }} ({{ caster.host }}:{{ caster.port }})
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label>Mountpoint *</label>
          <input 
            v-model="formData.mountpoint" 
            type="text" 
            placeholder="e.g., /BASE1"
            required
          />
        </div>

        <div class="form-group checkbox-group">
          <label>
            <input v-model="formData.enabled" type="checkbox" />
            Enabled
          </label>
        </div>
        
        <div class="modal-footer">
          <button type="button" class="btn-cancel" @click="$emit('close')">
            Cancel
          </button>
          <button type="submit" class="btn-save">
            {{ isEdit ? 'Update' : 'Add' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useCastersStore } from '../stores/casters'

const props = defineProps({
  stationId: {
    type: Number,
    required: true
  },
  output: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'save'])

const castersStore = useCastersStore()
const casters = computed(() => castersStore.casters)

const isEdit = computed(() => !!props.output)

const formData = ref({
  caster_id: '',
  mountpoint: '',
  enabled: true
})

if (props.output) {
  formData.value = {
    caster_id: props.output.caster_id,
    mountpoint: props.output.mountpoint,
    enabled: props.output.enabled
  }
}

function handleSubmit() {
  if (!formData.value.caster_id) {
    alert('Please select a caster')
    return
  }
  if (!formData.value.mountpoint) {
    alert('Mountpoint is required')
    return
  }
  emit('save', { ...formData.value })
}
</script>

<style scoped>
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

.modal {
  background: white;
  border-radius: 8px;
  width: 450px;
  max-width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.1rem;
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

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 0.85rem;
  color: #666;
  margin-bottom: 6px;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.95rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3498db;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-group input {
  width: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 16px;
  border-top: 1px solid #eee;
  margin-top: 16px;
}

.btn-cancel,
.btn-save {
  padding: 10px 20px;
  border-radius: 4px;
  font-size: 0.9rem;
  cursor: pointer;
  border: none;
}

.btn-cancel {
  background: #f0f0f0;
  color: #666;
}

.btn-cancel:hover {
  background: #e0e0e0;
}

.btn-save {
  background: #3498db;
  color: white;
}

.btn-save:hover {
  background: #2980b9;
}
</style>