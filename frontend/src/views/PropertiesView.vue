<template>
  <div class="properties">
    <div class="header">
      <h2>Properties</h2>
      <button @click="showForm = !showForm" class="btn btn-primary">
        {{ showForm ? 'Cancel' : 'Add Property' }}
      </button>
    </div>

    <form v-if="showForm" @submit.prevent="handleSubmit" class="property-form">
      <input
        v-model="form.name"
        type="text"
        placeholder="Property Name"
        required
        class="input"
      />
      <input
        v-model="form.address"
        type="text"
        placeholder="Address"
        required
        class="input"
      />
      <select v-model="form.property_type" required class="input">
        <option value="">Select Type</option>
        <option value="residential">Residential</option>
        <option value="commercial">Commercial</option>
        <option value="mixed">Mixed Use</option>
      </select>
      <input
        v-model.number="form.default_mileage"
        type="number"
        step="0.1"
        min="0"
        placeholder="Default Mileage (round trip)"
        class="input"
      />
      <button type="submit" class="btn btn-primary" :disabled="store.loading">
        {{ store.loading ? 'Creating...' : 'Create' }}
      </button>
    </form>

    <div v-if="store.error" class="error">{{ store.error }}</div>

    <div v-if="store.loading && properties.length === 0" class="loading">Loading...</div>

    <div v-else-if="properties.length === 0" class="empty">No properties yet. Add one above!</div>

    <div v-else class="property-list">
      <div v-for="property in properties" :key="property.id" class="property-card">
        <div class="property-info">
          <h3>{{ property.name }}</h3>
          <p class="address">{{ property.address }}</p>
          <span class="badge">{{ property.property_type }}</span>
        </div>
        <button
          @click="handleDelete(property.id)"
          class="btn btn-danger"
          :disabled="store.loading"
        >
          Delete
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { usePropertyStore } from '@/stores/property'
import type { CreatePropertyRequest } from '@/types/property'

const store = usePropertyStore()
const showForm = ref(false)
const form = ref<CreatePropertyRequest>({
  name: '',
  address: '',
  property_type: '',
  default_mileage: 0,
})

const properties = computed(() => store.properties)

onMounted(() => {
  store.fetchProperties()
})

async function handleSubmit() {
  try {
    await store.createProperty(form.value)
    form.value = { name: '', address: '', property_type: '', default_mileage: 0 }
    showForm.value = false
  } catch (err) {
    // Error handled by store
  }
}

async function handleDelete(id: number) {
  if (confirm('Are you sure you want to delete this property?')) {
    await store.deleteProperty(id)
  }
}
</script>

<style scoped>
.properties {
  padding: 1rem;
  padding-bottom: 80px; /* Space for bottom nav */
}

@media (min-width: 768px) {
  .properties {
    padding: 2rem;
    padding-bottom: 2rem;
  }
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header h2 {
  margin: 0;
}

.property-form {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 2rem;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.input {
  flex: 1;
  min-width: 200px;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2980b9;
}

.btn-danger {
  background: #e74c3c;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #c0392b;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  background: #fee;
  color: #c33;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.loading,
.empty {
  text-align: center;
  padding: 3rem;
  color: #666;
}

.property-list {
  display: grid;
  gap: 1rem;
}

.property-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.property-info h3 {
  margin: 0 0 0.5rem 0;
}

.address {
  color: #666;
  margin-bottom: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #ecf0f1;
  border-radius: 12px;
  font-size: 0.875rem;
  text-transform: capitalize;
}
</style>
