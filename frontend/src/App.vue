<template>
  <div id="app">
    <SideNav />
    <main class="main-content">
      <router-view />
    </main>
    <BottomNav />
  </div>
</template>

<script setup lang="ts">
import SideNav from './components/SideNav.vue'
import BottomNav from './components/BottomNav.vue'
</script>

<style scoped>
#app {
  display: flex;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  margin-left: 0;
  padding: 0;
  background: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 60px; /* Space for bottom nav */
}

@media (min-width: 768px) {
  .main-content {
    margin-left: 240px;
    padding-bottom: 0;
  }
}

.main-content > :deep(*) {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}
</style>
