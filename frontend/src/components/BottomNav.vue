<template>
  <nav class="bottom-nav" aria-label="Main navigation">
    <router-link
      v-for="item in navItems"
      :key="item.to"
      :to="item.to"
      class="nav-item"
      active-class="active"
    >
      <i :class="['pi', item.icon]" aria-hidden="true" />
      <span class="nav-label">{{ item.label }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: 'pi-chart-bar' },
  { to: '/properties', label: 'Properties', icon: 'pi-building' },
  { to: '/residences', label: 'Residences', icon: 'pi-home' },
  { to: '/reps', label: 'REPS', icon: 'pi-clock' },
] as const
</script>

<style scoped>
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-around;
  align-items: stretch;
  background: var(--p-surface-0, #fff);
  border-top: 1px solid var(--p-surface-200, #e2e8f0);
  padding: 0.35rem 0 calc(0.35rem + env(safe-area-inset-bottom, 0px));
  z-index: 1000;
  box-shadow: 0 -2px 12px rgba(15, 23, 42, 0.08);
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.2rem;
  padding: 0.5rem 0.75rem;
  min-height: 3.25rem;
  text-decoration: none;
  color: var(--text-secondary);
  transition: color 0.2s;
  flex: 1;
}

.nav-item .pi {
  font-size: 1.25rem;
}

.nav-item.active {
  color: var(--p-primary-500, #1a6fb5);
}

.nav-label {
  font-size: 0.7rem;
  font-weight: 600;
}

@media (min-width: 768px) {
  .bottom-nav {
    display: none;
  }
}
</style>
