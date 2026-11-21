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
  { to: '/settings/quickbooks', label: 'QuickBooks', icon: 'pi-wallet' },
] as const
</script>

<style scoped>
.side-nav {
  width: var(--sidebar-width);
  background:
    linear-gradient(165deg, #2f2a26 0%, #241f1c 55%, #1c1917 100%);
  color: #f7f2ea;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  display: flex;
  flex-direction: column;
  z-index: 100;
  border-right: 1px solid rgba(255, 253, 249, 0.06);
}

.nav-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1.35rem 1rem;
  border-bottom: 1px solid rgba(255, 253, 249, 0.08);
}

.brand-icon {
  font-size: 1.25rem;
  color: var(--p-primary-300, #66c0ba);
}

.nav-header h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: 1.35rem;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.nav-links {
  flex: 1;
  padding: 0.75rem 0.5rem;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.8rem 1rem;
  min-height: 2.75rem;
  margin: 0.15rem 0;
  border-radius: 0.55rem;
  color: rgba(247, 242, 234, 0.72);
  text-decoration: none;
  border-left: none;
}

.nav-link:hover {
  background: rgba(255, 253, 249, 0.08);
  color: #fffdf9;
  transform: translateY(-1px);
}

.nav-link.router-link-active {
  background: color-mix(in srgb, var(--accent, #0f766e) 28%, transparent);
  color: #fffdf9;
  box-shadow: inset 3px 0 0 var(--p-primary-400, #3aa39c);
}

.nav-link .pi {
  font-size: 1.1rem;
}

.nav-link span:last-child {
  font-weight: 600;
}

@media (max-width: 767px) {
  .side-nav {
    display: none;
  }
}
</style>
