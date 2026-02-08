<template>
  <div class="landing">
    <div class="welcome-section">
      <h1>Welcome to Brickwise</h1>
      <p class="subtitle">Manage your property portfolio with ease</p>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-value">{{ propertyCount }}</div>
        <div class="stat-label">Properties</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">$0</div>
        <div class="stat-label">Monthly Revenue</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">0</div>
        <div class="stat-label">Active Units</div>
      </div>
    </div>

    <div class="quick-actions">
      <h2>Quick Actions</h2>
      <div class="actions-grid">
        <router-link to="/properties" class="action-card">
          <span class="action-icon">➕</span>
          <h3>Add Property</h3>
          <p>Register a new property in your portfolio</p>
        </router-link>
        <div class="action-card disabled">
          <span class="action-icon">📊</span>
          <h3>View Reports</h3>
          <p>Coming soon</p>
        </div>
        <div class="action-card disabled">
          <span class="action-icon">💰</span>
          <h3>Financial Overview</h3>
          <p>Coming soon</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { usePropertyStore } from '@/stores/property'

const store = usePropertyStore()
const propertyCount = computed(() => store.properties.length)

onMounted(() => {
  store.fetchProperties()
})
</script>

<style scoped>
.landing {
  padding: 2rem 0;
}

.welcome-section {
  margin-bottom: 3rem;
}

.welcome-section h1 {
  font-size: 2.5rem;
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.subtitle {
  font-size: 1.25rem;
  color: #666;
  margin: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: #3498db;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #666;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.quick-actions h2 {
  margin: 0 0 1.5rem 0;
  color: #2c3e50;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.action-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-decoration: none;
  color: inherit;
  transition: transform 0.2s, box-shadow 0.2s;
  border: 2px solid transparent;
}

.action-card:hover:not(.disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  border-color: #3498db;
}

.action-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.action-icon {
  font-size: 2.5rem;
  display: block;
  margin-bottom: 1rem;
}

.action-card h3 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.action-card p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}
</style>
