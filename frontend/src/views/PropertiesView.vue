<template>
  <div class="page-view properties">
    <div class="page-header">
      <div>
        <h1 class="page-title">Properties</h1>
        <p class="page-subtitle">Manage your rental portfolio</p>
      </div>
      <Button
        v-if="!showPropertyForm"
        label="Add Property"
        icon="pi pi-plus"
        size="large"
        @click="openForm"
      />
      <Button
        v-else
        label="Cancel"
        icon="pi pi-times"
        size="large"
        outlined
        @click="closeForm"
      />
    </div>

    <Message v-if="store.error" severity="error" class="mb-3" :closable="false">
      {{ store.error }}
    </Message>

    <Card v-if="showPropertyForm" class="mb-4">
      <template #title>New property</template>
      <template #content>
        <form class="property-form" @submit.prevent="handleSubmit">
          <div class="form-stack">
            <FloatLabel>
              <InputText id="name" v-model="form.name" class="w-full" required fluid />
              <label for="name">Property name</label>
            </FloatLabel>
            <FloatLabel>
              <InputText id="address" v-model="form.address" class="w-full" required fluid />
              <label for="address">Address</label>
            </FloatLabel>
            <FloatLabel>
              <Select
                id="type"
                v-model="form.property_type"
                :options="propertyTypes"
                optionLabel="label"
                optionValue="value"
                placeholder="Select type"
                class="w-full"
                fluid
              />
              <label for="type">Property type</label>
            </FloatLabel>
            <FloatLabel>
              <InputNumber
                id="mileage"
                v-model="form.default_mileage"
                :min="0"
                :maxFractionDigits="1"
                class="w-full"
                fluid
              />
              <label for="mileage">Default mileage (round trip)</label>
            </FloatLabel>
          </div>
          <Button
            type="submit"
            label="Create property"
            icon="pi pi-check"
            size="large"
            class="mt-3 w-full"
            :loading="store.loading"
          />
        </form>
      </template>
    </Card>

    <div v-if="store.loading && properties.length === 0" class="state-block">
      <ProgressSpinner />
      <p>Loading properties…</p>
    </div>

    <Card v-else-if="properties.length === 0">
      <template #content>
        <div class="empty-state">
          <i class="pi pi-building empty-icon" aria-hidden="true" />
          <p>No properties yet. Add your first property to get started.</p>
        </div>
      </template>
    </Card>

    <div v-else class="property-list">
      <Card v-for="property in properties" :key="property.id" class="property-card">
        <template #content>
          <div class="property-row">
            <div class="property-main">
              <img
                v-if="property.main_image_id"
                :src="propertyImageFileUrl(property.id, property.main_image_id)"
                :alt="`${property.name} thumbnail`"
                class="property-thumb"
              />
              <div
                v-else
                class="property-thumb placeholder"
                aria-hidden="true"
              >
                <i class="pi pi-building" />
              </div>
              <div>
                <h3 class="property-name">{{ property.name }}</h3>
                <p class="property-address">{{ property.address }}</p>
                <Tag :value="property.property_type" severity="secondary" class="capitalize" />
                <p v-if="property.current_value > 0" class="property-value">
                  Value: {{ formatCurrency(property.current_value) }}
                  <span v-if="property.value_updated_at" class="value-date">
                    (as of {{ property.value_updated_at.split('T')[0] }})
                  </span>
                </p>
              </div>
            </div>
            <div class="property-actions">
              <Button
                icon="pi pi-pencil"
                label="Edit"
                outlined
                aria-label="Edit property"
                @click="openEdit(property)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import FloatLabel from 'primevue/floatlabel'
import Message from 'primevue/message'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import { propertyImageFileUrl } from '@/api/client'
import { usePropertyStore } from '@/stores/property'
import { PROPERTY_TYPES, type CreatePropertyRequest, type Property } from '@/types/property'

const route = useRoute()
const router = useRouter()
const store = usePropertyStore()
const toast = useToast()

const propertyTypes = PROPERTY_TYPES

const emptyForm = (): CreatePropertyRequest => ({
  name: '',
  address: '',
  property_type: '',
  default_mileage: 0,
})

const form = ref<CreatePropertyRequest>(emptyForm())

const properties = computed(() => store.properties ?? [])
const showPropertyForm = computed(() => route.query.add === '1')

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(amount)
}

onMounted(() => {
  store.fetchProperties()
})

function openForm() {
  form.value = emptyForm()
  router.push({ path: '/properties', query: { ...route.query, add: '1' } })
}

function openEdit(property: Property) {
  router.push({ name: 'property-edit', params: { id: String(property.id) } })
}

function closeForm() {
  form.value = emptyForm()
  const q = { ...route.query }
  delete q.add
  router.replace({ path: '/properties', query: q })
}

async function handleSubmit() {
  try {
    await store.createProperty(form.value)
    closeForm()
    toast.add({ severity: 'success', summary: 'Property created', life: 3000 })
  } catch {
    // Error handled by store
  }
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.property-list {
  display: grid;
  gap: 1rem;
}

.property-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.property-main {
  display: flex;
  align-items: flex-start;
  gap: 0.85rem;
  min-width: 0;
}

.property-thumb {
  width: 4.5rem;
  height: 4.5rem;
  object-fit: cover;
  border-radius: 0.4rem;
  flex-shrink: 0;
  background: var(--p-surface-100, #f1f5f9);
}

.property-thumb.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--p-primary-300, #66c0ba);
  font-size: 1.35rem;
}

.property-actions {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.property-value {
  margin: 0.5rem 0 0;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.value-date {
  font-size: 0.8rem;
}

.property-name {
  margin: 0 0 0.35rem;
  font-size: 1.1rem;
}

.property-address {
  margin: 0 0 0.5rem;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.state-block,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 2.5rem;
  color: var(--p-primary-300);
}

.capitalize {
  text-transform: capitalize;
}

@media (min-width: 768px) {
  .property-list {
    grid-template-columns: repeat(2, 1fr);
  }

  .form-stack {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
  }

  .form-stack > :first-child,
  .form-stack > :nth-child(2) {
    grid-column: span 1;
  }
}
</style>
