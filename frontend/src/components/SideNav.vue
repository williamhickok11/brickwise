<template>
  <aside class="side-nav">
    <div class="nav-header">
      <h2>Brickwise</h2>
    </div>
    <nav class="nav-links">
      <router-link to="/dashboard" class="nav-link">
        <span class="nav-icon">📊</span>
        <span>Dashboard</span>
      </router-link>
      <router-link to="/properties" class="nav-link">
        <span class="nav-icon">🏢</span>
        <span>Properties</span>
      </router-link>
      <router-link to="/residences" class="nav-link">
        <span class="nav-icon">🏠</span>
        <span>Residences</span>
      </router-link>
      <router-link to="/reps" class="nav-link">
        <span class="nav-icon">⏱️</span>
        <span>REPS</span>
      </router-link>
      <router-link to="/reps-v2" class="nav-link">
        <span class="nav-icon">🎤</span>
        <span>REPS V2</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup lang="ts">
// Side navigation component
</script>

<style scoped>
.side-nav {
  width: 240px;
  background: #2c3e50;
  color: white;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  display: flex;
  flex-direction: column;
}

.nav-header {
  padding: 1.5rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.nav-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.nav-links {
  flex: 1;
  padding: 1rem 0;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.5rem;
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.15);
  color: white;
  border-left-color: #3498db;
}

.nav-icon {
  font-size: 1.25rem;
}

.nav-link span:last-child {
  font-weight: 500;
}

@media (max-width: 767px) {
  .side-nav {
    display: none;
  }
}
</style>
