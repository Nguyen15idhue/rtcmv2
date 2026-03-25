<template>
  <div class="chart-container">
    <canvas ref="canvasRef" :height="height"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 200
  }
})

const canvasRef = ref(null)
let animationId = null

function drawChart() {
  if (!canvasRef.value || props.data.length === 0) return
  
  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  const width = canvas.parentElement?.clientWidth || 600
  
  canvas.width = width
  canvas.height = props.height
  
  ctx.clearRect(0, 0, width, props.height)
  
  if (props.data.length < 2) return
  
  const padding = 30
  const chartWidth = width - padding * 2
  const chartHeight = props.height - padding * 2
  
  const maxY = Math.max(...props.data.map(d => d.y), 10)
  const minY = 0
  
  ctx.strokeStyle = '#3498db'
  ctx.lineWidth = 2
  ctx.beginPath()
  
  props.data.forEach((point, i) => {
    const x = padding + (i / (props.data.length - 1)) * chartWidth
    const y = padding + chartHeight - ((point.y - minY) / (maxY - minY)) * chartHeight
    
    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  
  ctx.stroke()
  
  ctx.fillStyle = '#666'
  ctx.font = '10px sans-serif'
  ctx.fillText('0', 5, props.height - 5)
  ctx.fillText(maxY.toFixed(0), 5, 15)
}

watch(() => props.data, () => {
  drawChart()
}, { deep: true })

onMounted(() => {
  drawChart()
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }
})
</script>

<style scoped>
.chart-container {
  width: 100%;
  background: white;
  border-radius: 4px;
}

canvas {
  width: 100%;
}
</style>