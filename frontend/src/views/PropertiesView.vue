<template>
  <div class="page-view properties">
    <div class="page-header">
      <div>
        <h1 class="page-title">Properties</h1>
        <p class="page-subtitle">Manage your rental portfolio</p>
      </div>
      <Button
        v-if="!showAddForm"
        label="Add Property"
        icon="pi pi-plus"
        size="large"
        @click="openForm"
      />
      <Button
        v-else
        label="Cancel"
        icon="pi pi-times"
        severity="secondary"
        size="large"
        outlined
        @click="closeForm"
      />
    </div>

    <Message v-if="store.error" severity="error" class="mb-3" :closable="false">
      {{ store.error }}
    </Message>

    <Card v-if="showAddForm" class="mb-4">
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
            <div>
              <h3 class="property-name">{{ property.name }}</h3>
              <p class="property-address">{{ property.address }}</p>
              <Tag :value="property.property_type" severity="secondary" class="capitalize" />
            </div>
            <Button
              icon="pi pi-trash"
              severity="danger"
              outlined
              aria-label="Delete property"
              :loading="store.loading"
              @click="confirmDelete(property.id)"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
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
import { usePropertyStore } from '@/stores/property'
import type { CreatePropertyRequest } from '@/types/property'

const route = useRoute()
const router = useRouter()
const store = usePropertyStore()
const confirm = useConfirm()
const toast = useToast()

const propertyTypes = [
  { label: 'Residential', value: 'residential' },
  { label: 'Commercial', value: 'commercial' },
  { label: 'Mixed use', value: 'mixed' },
]

const form = ref<CreatePropertyRequest>({
  name: '',
  address: '',
  property_type: '',
  default_mileage: 0,
})

const properties = computed(() => store.properties ?? [])
const showAddForm = computed(() => route.query.add === '1')

onMounted(() => {
  store.fetchProperties()
})

function openForm() {
  router.push({ path: '/properties', query: { ...route.query, add: '1' } })
}

function closeForm() {
  const q = { ...route.query }
  delete q.add
  router.replace({ path: '/properties', query: q })
}

async function handleSubmit() {
  try {
    await store.createProperty(form.value)
    form.value = { name: '', address: '', property_type: '', default_mileage: 0 }
    closeForm()
    toast.add({ severity: 'success', summary: 'Property created', life: 3000 })
  } catch {
    // Error handled by store
  }
}

function confirmDelete(id: number) {
  confirm.require({
    message: 'Are you sure you want to delete this property?',
    header: 'Delete property',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    acceptClass: 'p-button-danger',
    accept: async () => {
      await store.deleteProperty(id)
      toast.add({ severity: 'success', summary: 'Property deleted', life: 3000 })
    },
  })
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
