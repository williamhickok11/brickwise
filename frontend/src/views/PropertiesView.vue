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
      <template #title>{{ editingId ? 'Edit property' : 'New property' }}</template>
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
            :label="editingId ? 'Save changes' : 'Create property'"
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
              <p v-if="property.current_value > 0" class="property-value">
                Value: {{ formatCurrency(property.current_value) }}
                <span v-if="property.value_updated_at" class="value-date">
                  (as of {{ property.value_updated_at.split('T')[0] }})
                </span>
              </p>
            </div>
            <div class="property-actions">
              <Button
                icon="pi pi-pencil"
                outlined
                aria-label="Edit property"
                @click="openEdit(property)"
              />
              <Button
                icon="pi pi-dollar"
                outlined
                aria-label="Edit financials"
                @click="openFinancials(property)"
              />
              <Button
                icon="pi pi-trash"
                outlined
                aria-label="Delete property"
                :loading="store.loading"
                @click="confirmDelete(property.id)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <Dialog
      v-model:visible="financialsVisible"
      :header="financialsProperty ? `Financials: ${financialsProperty.name}` : 'Financials'"
      modal
      class="financials-dialog"
      :style="{ width: 'min(36rem, 95vw)' }"
    >
      <form class="financials-form" @submit.prevent="saveFinancials">
        <h4 class="form-section-title">Property value</h4>
        <div class="form-stack">
          <FloatLabel>
            <InputNumber
              id="current-value"
              v-model="financialsForm.current_value"
              mode="currency"
              currency="USD"
              locale="en-US"
              :min="0"
              class="w-full"
              fluid
            />
            <label for="current-value">Current value</label>
          </FloatLabel>
          <FloatLabel>
            <DatePicker
              id="value-updated"
              v-model="valueUpdatedModel"
              dateFormat="yy-mm-dd"
              showIcon
              fluid
            />
            <label for="value-updated">Value as of</label>
          </FloatLabel>
        </div>

        <h4 class="form-section-title">Mortgage</h4>
        <div class="form-stack">
          <FloatLabel>
            <InputNumber
              id="current-balance"
              v-model="financialsForm.mortgage.current_balance"
              mode="currency"
              currency="USD"
              locale="en-US"
              :min="0"
              class="w-full"
              fluid
            />
            <label for="current-balance">Current balance</label>
          </FloatLabel>
          <FloatLabel>
            <InputNumber
              id="original-loan"
              v-model="financialsForm.mortgage.original_loan"
              mode="currency"
              currency="USD"
              locale="en-US"
              :min="0"
              class="w-full"
              fluid
            />
            <label for="original-loan">Original loan</label>
          </FloatLabel>
          <FloatLabel>
            <InputNumber
              id="interest-rate"
              v-model="financialsForm.mortgage.interest_rate"
              suffix="%"
              :min="0"
              :max="100"
              :maxFractionDigits="3"
              class="w-full"
              fluid
            />
            <label for="interest-rate">Interest rate</label>
          </FloatLabel>
          <FloatLabel>
            <InputNumber
              id="monthly-pi"
              v-model="financialsForm.mortgage.monthly_pi"
              mode="currency"
              currency="USD"
              locale="en-US"
              :min="0"
              class="w-full"
              fluid
            />
            <label for="monthly-pi">Monthly P&amp;I</label>
          </FloatLabel>
          <FloatLabel>
            <InputNumber
              id="monthly-payment"
              v-model="financialsForm.mortgage.monthly_payment"
              mode="currency"
              currency="USD"
              locale="en-US"
              :min="0"
              class="w-full"
              fluid
            />
            <label for="monthly-payment">Monthly payment (incl. escrow)</label>
          </FloatLabel>
          <FloatLabel>
            <DatePicker
              id="balance-as-of"
              v-model="balanceAsOfModel"
              dateFormat="yy-mm-dd"
              showIcon
              fluid
            />
            <label for="balance-as-of">Balance as of</label>
          </FloatLabel>
          <FloatLabel>
            <InputText id="lender" v-model="financialsForm.mortgage.lender" class="w-full" fluid />
            <label for="lender">Lender</label>
          </FloatLabel>
        </div>

        <div class="dialog-actions">
          <Button
            type="submit"
            label="Save financials"
            icon="pi pi-check"
            :loading="savingFinancials"
          />
          <Button
            v-if="hasMortgage"
            type="button"
            label="Remove mortgage"
            icon="pi pi-trash"
            outlined
            :loading="savingFinancials"
            @click="removeMortgage"
          />
        </div>
      </form>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DatePicker from 'primevue/datepicker'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import FloatLabel from 'primevue/floatlabel'
import Message from 'primevue/message'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import { propertyApi, mortgageApi } from '@/api/client'
import { usePropertyStore } from '@/stores/property'
import { PROPERTY_TYPES, type CreatePropertyRequest, type Property } from '@/types/property'
import type { UpsertMortgageRequest } from '@/types/mortgage'

const route = useRoute()
const router = useRouter()
const store = usePropertyStore()
const confirm = useConfirm()
const toast = useToast()

const propertyTypes = PROPERTY_TYPES

const emptyForm = (): CreatePropertyRequest => ({
  name: '',
  address: '',
  property_type: '',
  default_mileage: 0,
})

const form = ref<CreatePropertyRequest>(emptyForm())
const editingId = ref<number | null>(null)

const properties = computed(() => store.properties ?? [])
const showAddForm = computed(() => route.query.add === '1')
const showPropertyForm = computed(() => showAddForm.value || editingId.value !== null)

const financialsVisible = ref(false)
const financialsProperty = ref<Property | null>(null)
const savingFinancials = ref(false)
const hasMortgage = ref(false)

const emptyMortgage = (): UpsertMortgageRequest => ({
  original_loan: 0,
  current_balance: 0,
  interest_rate: 0,
  monthly_payment: 0,
  monthly_pi: 0,
  balance_as_of: undefined,
  lender: '',
  notes: '',
})

const financialsForm = ref({
  current_value: 0,
  value_updated_at: todayISO(),
  mortgage: emptyMortgage(),
})

const valueUpdatedModel = computed({
  get: () => parseDate(financialsForm.value.value_updated_at),
  set: (val: Date | null) => {
    if (val) financialsForm.value.value_updated_at = toISO(val)
  },
})

const balanceAsOfModel = computed({
  get: () =>
    financialsForm.value.mortgage.balance_as_of
      ? parseDate(financialsForm.value.mortgage.balance_as_of)
      : null,
  set: (val: Date | null) => {
    financialsForm.value.mortgage.balance_as_of = val ? toISO(val) : undefined
  },
})

function todayISO(): string {
  return new Date().toISOString().split('T')[0]
}

function parseDate(iso: string): Date {
  const [y, m, d] = iso.split('T')[0].split('-').map(Number)
  return new Date(y, m - 1, d)
}

function toISO(date: Date): string {
  return date.toISOString().split('T')[0]
}

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(amount)
}

function openFinancials(property: Property) {
  financialsProperty.value = property
  financialsForm.value = {
    current_value: property.current_value ?? 0,
    value_updated_at: property.value_updated_at
      ? property.value_updated_at.split('T')[0]
      : todayISO(),
    mortgage: emptyMortgage(),
  }
  hasMortgage.value = false
  financialsVisible.value = true
  loadMortgage(property.id)
}

async function loadMortgage(propertyId: number) {
  try {
    const m = await mortgageApi.get(propertyId)
    financialsForm.value.mortgage = {
      original_loan: m.original_loan,
      current_balance: m.current_balance,
      interest_rate: m.interest_rate,
      monthly_payment: m.monthly_payment,
      monthly_pi: m.monthly_pi,
      balance_as_of: m.balance_as_of ?? undefined,
      lender: m.lender,
      notes: m.notes,
    }
    hasMortgage.value = true
  } catch {
    // No mortgage on file
  }
}

async function saveFinancials() {
  if (!financialsProperty.value) return
  savingFinancials.value = true
  try {
    const updated = await propertyApi.updateValue(financialsProperty.value.id, {
      current_value: financialsForm.value.current_value,
      value_updated_at: financialsForm.value.value_updated_at,
    })
    const idx = store.properties.findIndex((p) => p.id === updated.id)
    if (idx >= 0) store.properties[idx] = updated

    const m = financialsForm.value.mortgage
    const hasMortgageData =
      m.current_balance > 0 ||
      m.original_loan > 0 ||
      m.monthly_pi > 0 ||
      m.monthly_payment > 0

    if (hasMortgageData) {
      await mortgageApi.upsert(financialsProperty.value.id, m)
      hasMortgage.value = true
    }

    financialsVisible.value = false
    toast.add({ severity: 'success', summary: 'Financials saved', life: 3000 })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to save financials',
      detail: err instanceof Error ? err.message : undefined,
      life: 5000,
    })
  } finally {
    savingFinancials.value = false
  }
}

async function removeMortgage() {
  if (!financialsProperty.value) return
  savingFinancials.value = true
  try {
    await mortgageApi.delete(financialsProperty.value.id)
    financialsForm.value.mortgage = emptyMortgage()
    hasMortgage.value = false
    toast.add({ severity: 'success', summary: 'Mortgage removed', life: 3000 })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to remove mortgage',
      detail: err instanceof Error ? err.message : undefined,
      life: 5000,
    })
  } finally {
    savingFinancials.value = false
  }
}

onMounted(() => {
  store.fetchProperties()
})

function resetForm() {
  form.value = emptyForm()
  editingId.value = null
}

function openForm() {
  editingId.value = null
  form.value = emptyForm()
  router.push({ path: '/properties', query: { ...route.query, add: '1' } })
}

function openEdit(property: Property) {
  const q = { ...route.query }
  delete q.add
  router.replace({ path: '/properties', query: q })
  editingId.value = property.id
  form.value = {
    name: property.name,
    address: property.address,
    property_type: (property.property_type as CreatePropertyRequest['property_type']) || '',
    default_mileage: property.default_mileage ?? 0,
  }
}

function closeForm() {
  resetForm()
  const q = { ...route.query }
  delete q.add
  router.replace({ path: '/properties', query: q })
}

async function handleSubmit() {
  try {
    if (editingId.value !== null) {
      await store.updateProperty(editingId.value, form.value)
      closeForm()
      toast.add({ severity: 'success', summary: 'Property updated', life: 3000 })
      return
    }
    await store.createProperty(form.value)
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

.form-section-title {
  margin: 1.25rem 0 0.75rem;
  font-size: 0.95rem;
  color: var(--text-primary);
}

.form-section-title:first-child {
  margin-top: 0;
}

.financials-form .form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.dialog-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.5rem;
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

  .financials-form .form-stack {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
  }
}
</style>
