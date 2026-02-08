<template>
  <nav class="bottom-nav">
    <router-link to="/dashboard" class="nav-item" active-class="active">
      <span class="nav-icon">📊</span>
      <span class="nav-label">Dashboard</span>
    </router-link>
    <router-link to="/properties" class="nav-item" active-class="active">
      <span class="nav-icon">🏢</span>
      <span class="nav-label">Properties</span>
    </router-link>
    <router-link to="/residences" class="nav-item" active-class="active">
      <span class="nav-icon">🏠</span>
      <span class="nav-label">Residences</span>
    </router-link>
    <router-link to="/reps" class="nav-item" active-class="active">
      <span class="nav-icon">⏱️</span>
      <span class="nav-label">REPS</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
// Bottom navigation component for mobile-first design
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-around;
  align-items: center;
  background: white;
  border-top: 1px solid #e0e0e0;
  padding: 0.5rem 0;
  z-index: 1000;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.5rem 1rem;
  text-decoration: none;
  color: #666;
  transition: color 0.2s;
  min-width: 60px;
  flex: 1;
}

.nav-item.active {
  color: #3498db;
}

.nav-icon {
  font-size: 1.5rem;
}

.nav-label {
  font-size: 0.75rem;
  font-weight: 500;
}

@media (min-width: 768px) {
  .bottom-nav {
    display: none;
  }
}
</style>
