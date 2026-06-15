<template>
  <aside class="side-nav" aria-label="Sidebar navigation">
    <div class="nav-header">
      <span class="brand-icon pi pi-th-large" aria-hidden="true" />
      <h2>Brickwise</h2>
    </div>
    <nav class="nav-links">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="nav-link"
        active-class="router-link-active"
      >
        <i :class="['pi', item.icon]" aria-hidden="true" />
        <span>{{ item.label }}</span>
      </router-link>
    </nav>
  </aside>
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
.side-nav {
  width: var(--sidebar-width);
  background: #1e293b;
  color: #f8fafc;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  display: flex;
  flex-direction: column;
  z-index: 100;
}

.nav-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1.25rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.brand-icon {
  font-size: 1.25rem;
  color: var(--p-primary-400, #3d9ed9);
}

.nav-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
}

.nav-links {
  flex: 1;
  padding: 0.75rem 0;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.25rem;
  min-height: 2.75rem;
  color: rgba(248, 250, 252, 0.75);
  text-decoration: none;
  transition: background 0.2s, color 0.2s;
  border-left: 3px solid transparent;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
  border-left-color: var(--p-primary-400, #3d9ed9);
}

.nav-link .pi {
  font-size: 1.1rem;
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
