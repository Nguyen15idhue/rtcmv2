<template>
  <div class="app">
    <nav class="navbar">
      <div class="nav-brand">
        <span class="brand-icon">📡</span>
        <span class="brand-text">RTCMv2 Relay</span>
      </div>
      <div class="nav-links">
        <router-link to="/">Dashboard</router-link>
        <router-link to="/stations">Stations</router-link>
        <router-link to="/casters">Casters</router-link>
        <router-link to="/settings">Settings</router-link>
      </div>
      <div class="nav-status">
        <span :class="['status-dot', streamStore.connected ? 'connected' : 'disconnected']"></span>
        {{ streamStore.connected ? 'Live' : 'Offline' }}
      </div>
    </nav>

    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useStreamStore } from './stores/stream'

const streamStore = useStreamStore()

onMounted(() => {
  streamStore.connect()
})
</script>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: #f5f5f5;
  color: #333;
}

.app {
  min-height: 100vh;
}

.navbar {
  background: #2c3e50;
  color: white;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 30px;
  height: 56px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 1.1rem;
}

.brand-icon {
  font-size: 1.3rem;
}

.nav-links {
  display: flex;
  gap: 5px;
  flex: 1;
}

.nav-links a {
  color: rgba(255,255,255,0.8);
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.nav-links a:hover {
  background: rgba(255,255,255,0.1);
  color: white;
}

.nav-links a.router-link-active {
  background: rgba(255,255,255,0.2);
  color: white;
}

.nav-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  color: rgba(255,255,255,0.8);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #e74c3c;
}

.status-dot.connected {
  background: #27ae60;
}

.main-content {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}
</style>
