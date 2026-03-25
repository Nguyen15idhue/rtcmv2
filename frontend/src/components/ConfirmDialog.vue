<template>
  <div class="confirm-dialog-overlay" @click.self="$emit('cancel')">
    <div class="confirm-dialog">
      <div class="confirm-header">
        <h3>{{ title }}</h3>
      </div>
      <div class="confirm-body">
        <p>{{ message }}</p>
      </div>
      <div class="confirm-footer">
        <button class="btn-cancel" @click="$emit('cancel')">Cancel</button>
        <button class="btn-confirm" :class="variant" @click="$emit('confirm')">Confirm</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: 'Confirm'
  },
  message: {
    type: String,
    default: 'Are you sure?'
  },
  variant: {
    type: String,
    default: 'danger'
  }
})

defineEmits(['confirm', 'cancel'])
</script>

<style scoped>
.confirm-dialog-overlay {
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

.confirm-dialog {
  background: white;
  border-radius: 8px;
  width: 400px;
  max-width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.confirm-header {
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.confirm-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #333;
}

.confirm-body {
  padding: 20px;
}

.confirm-body p {
  margin: 0;
  color: #666;
}

.confirm-footer {
  padding: 16px 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn-cancel,
.btn-confirm {
  padding: 8px 16px;
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

.btn-confirm {
  color: white;
}

.btn-confirm.danger {
  background: #e74c3c;
}

.btn-confirm.danger:hover {
  background: #c0392b;
}

.btn-confirm.primary {
  background: #3498db;
}

.btn-confirm.primary:hover {
  background: #2980b9;
}
</style>