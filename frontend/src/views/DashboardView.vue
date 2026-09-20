<template>
  <div class="page-view dashboard">
    <div class="welcome-section">
      <h1 class="page-title">Dashboard</h1>
      <p class="page-subtitle">Portfolio overview at a glance</p>
    </div>

    <Message v-if="dashboardStore.error" severity="error" class="mb-3" :closable="false">
      {{ dashboardStore.error }}
    </Message>

    <div v-if="dashboardStore.loading && !dashboardStore.summary" class="state-block">
      <ProgressSpinner />
      <p>Loading dashboard…</p>
    </div>

    <template v-else-if="dashboardStore.summary">
      <div class="stats-grid">
        <Card v-for="stat in stats" :key="stat.label" class="stat-card">
          <template #content>
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
            <div v-if="stat.hint" class="stat-hint">{{ stat.hint }}</div>
          </template>
        </Card>
      </div>

      <section class="portfolio-section">
        <h2 class="section-title">Portfolio breakdown</h2>
        <Card>
          <template #content>
            <DataTable
              :value="dashboardStore.summary.properties"
              stripedRows
              responsiveLayout="scroll"
              class="portfolio-table"
            >
              <Column field="name" header="Property" />
              <Column field="current_value" header="Value">
                <template #body="{ data }">{{ formatCurrency(data.current_value) }}</template>
              </Column>
              <Column field="mortgage_balance" header="Mortgage">
                <template #body="{ data }">{{ formatCurrency(data.mortgage_balance) }}</template>
              </Column>
              <Column field="equity" header="Equity">
                <template #body="{ data }">{{ formatCurrency(data.equity) }}</template>
              </Column>
              <Column field="monthly_rent" header="Rent/mo">
                <template #body="{ data }">{{ formatCurrency(data.monthly_rent) }}</template>
              </Column>
              <Column field="monthly_pi" header="P&amp;I/mo">
                <template #body="{ data }">{{ formatCurrency(data.monthly_pi) }}</template>
              </Column>
              <Column field="net_cash_flow" header="Net/mo">
                <template #body="{ data }">
                  <span :class="data.net_cash_flow >= 0 ? 'positive' : 'negative'">
                    {{ formatCurrency(data.net_cash_flow) }}
                  </span>
                </template>
              </Column>
            </DataTable>
            <p v-if="dashboardStore.summary.properties.length === 0" class="empty-mini">
              No properties yet. Add one to see portfolio breakdown.
            </p>
          </template>
        </Card>
      </section>
    </template>

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
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import { useDashboardStore } from '@/stores/dashboard'

const dashboardStore = useDashboardStore()

const stats = computed(() => {
  const s = dashboardStore.summary
  if (!s) {
    return [
      { value: '$0', label: 'Net worth', hint: null },
      { value: '$0', label: 'Portfolio value', hint: null },
      { value: '$0', label: 'Monthly revenue', hint: null },
      { value: '$0', label: 'Est. cash flow', hint: 'pre-expenses' },
    ]
  }
  return [
    { value: formatCurrency(s.total_equity), label: 'Net worth', hint: null },
    { value: formatCurrency(s.portfolio_value), label: 'Portfolio value', hint: null },
    { value: formatCurrency(s.monthly_revenue), label: 'Monthly revenue', hint: null },
    {
      value: formatCurrency(s.estimated_cash_flow),
      label: 'Est. cash flow',
      hint: 'pre-expenses',
    },
  ]
})

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

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(amount)
}

onMounted(() => {
  dashboardStore.fetchSummary()
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

.stat-hint {
  margin-top: 0.35rem;
  font-size: 0.75rem;
  color: var(--text-secondary);
  font-style: italic;
}

.portfolio-section {
  margin-bottom: 2.5rem;
}

.section-title {
  margin: 0 0 1.25rem;
  font-size: 1.25rem;
  color: var(--text-primary);
}

.portfolio-table :deep(.p-datatable) {
  font-size: 0.9rem;
}

.positive {
  color: var(--p-green-600);
  font-weight: 600;
}

.negative {
  color: var(--p-red-500);
  font-weight: 600;
}

.empty-mini {
  margin: 0;
  padding: 1rem;
  text-align: center;
  color: var(--text-secondary);
}

.state-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
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
