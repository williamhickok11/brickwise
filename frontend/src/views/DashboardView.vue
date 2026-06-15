<template>
  <div class="page-view dashboard">
    <div class="welcome-section">
      <h1 class="page-title">Welcome to Brickwise</h1>
      <p class="page-subtitle">Manage your property portfolio with ease</p>
    </div>

    <div class="stats-grid">
      <Card v-for="stat in stats" :key="stat.label" class="stat-card">
        <template #content>
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-label">{{ stat.label }}</div>
        </template>
      </Card>
    </div>

    <section class="quick-actions">
      <h2 class="section-title">Quick actions</h2>
      <div class="actions-grid">
        <Card
          v-for="action in actions"
          :key="action.title"
          :class="['action-card', { disabled: action.disabled }]"
        >
          <template #content>
            <component
              :is="action.to ? 'router-link' : 'div'"
              :to="action.to"
              class="action-link"
            >
              <i :class="['pi', action.icon, 'action-icon']" aria-hidden="true" />
              <h3>{{ action.title }}</h3>
              <p>{{ action.description }}</p>
            </component>
          </template>
        </Card>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import Card from 'primevue/card'
import { usePropertyStore } from '@/stores/property'

const store = usePropertyStore()
const propertyCount = computed(() => store.properties.length)

const stats = computed(() => [
  { value: propertyCount.value, label: 'Properties' },
  { value: '$0', label: 'Monthly revenue' },
  { value: '0', label: 'Active units' },
])

const actions = [
  {
    title: 'Add property',
    description: 'Register a new property in your portfolio',
    icon: 'pi-plus-circle',
    to: '/properties',
    disabled: false,
  },
  {
    title: 'Log REPS time',
    description: 'Track hours for tax compliance',
    icon: 'pi-clock',
    to: '/reps',
    disabled: false,
  },
  {
    title: 'View reports',
    description: 'Coming soon',
    icon: 'pi-chart-line',
    to: null,
    disabled: true,
  },
]

onMounted(() => {
  store.fetchProperties()
})
</script>

<style scoped>
.welcome-section {
  margin-bottom: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
  gap: 1rem;
  margin-bottom: 2.5rem;
}

.stat-card :deep(.p-card-content) {
  text-align: center;
  padding: 1.5rem 1rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--p-primary-500);
  margin-bottom: 0.35rem;
}

.stat-label {
  color: var(--text-secondary);
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-weight: 600;
}

.section-title {
  margin: 0 0 1.25rem;
  font-size: 1.25rem;
  color: var(--text-primary);
}

.actions-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

.action-link {
  display: block;
  text-decoration: none;
  color: inherit;
}

.action-icon {
  font-size: 2rem;
  color: var(--p-primary-500);
  display: block;
  margin-bottom: 0.75rem;
}

.action-card h3 {
  margin: 0 0 0.35rem;
  font-size: 1.05rem;
}

.action-card p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.action-card.disabled {
  opacity: 0.55;
  pointer-events: none;
}

.action-card:not(.disabled) :deep(.p-card) {
  transition: box-shadow 0.2s, transform 0.2s;
}

.action-card:not(.disabled):hover :deep(.p-card) {
  box-shadow: 0 4px 16px rgba(15, 23, 42, 0.1);
  transform: translateY(-2px);
}

@media (min-width: 768px) {
  .actions-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
