<template>
  <div id="app">
    <SideNav />
    <main class="main-content">
      <div class="page-container">
        <router-view />
      </div>
    </main>
    <BottomNav />
    <Toast position="top-center" />
    <ConfirmDialog />
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
  background: var(--page-bg);
  min-height: 100vh;
  padding-bottom: var(--nav-height-safe);
}

.page-container {
  max-width: 75rem;
  margin: 0 auto;
  padding: 0 1rem;
}

@media (min-width: 768px) {
  .main-content {
    margin-left: var(--sidebar-width);
    padding-bottom: 0;
  }

  .page-container {
    padding: 0 2rem;
  }
}
</style>
