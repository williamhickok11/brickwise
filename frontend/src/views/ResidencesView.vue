<template>
  <div class="page-view residences">
    <div class="page-header">
      <div>
        <h1 class="page-title">Residences</h1>
        <p class="page-subtitle">Track current and former residents per property</p>
      </div>
    </div>

    <Message v-if="propertyStore.error" severity="error" class="mb-3" :closable="false">
      {{ propertyStore.error }}
    </Message>

    <Card class="mb-3">
      <template #title>Select a property</template>
      <template #content>
        <div v-if="propertyStore.loading && properties.length === 0" class="state-block">
          <ProgressSpinner />
          <p>Loading properties…</p>
        </div>
        <div v-else-if="properties.length === 0" class="empty-state">
          <p>No properties yet. Add one in Properties.</p>
          <Button label="Go to properties" icon="pi pi-building" as="router-link" to="/properties" />
        </div>
        <div v-else class="property-picker">
          <Button
            v-for="property in properties"
            :key="property.id"
            :label="property.name"
            :severity="selectedPropertyId === property.id ? undefined : 'secondary'"
            :outlined="selectedPropertyId !== property.id"
            size="large"
            class="property-chip"
            @click="selectProperty(property.id)"
          />
        </div>
      </template>
    </Card>

    <template v-if="selectedPropertyId !== null">
      <div class="detail-header">
        <div>
          <h2 class="section-title">Residents</h2>
          <p class="detail-subtitle">
            For: <strong>{{ selectedPropertyName }}</strong>
          </p>
        </div>
        <Button
          v-if="!showAddForm"
          label="Add resident"
          icon="pi pi-user-plus"
          size="large"
          @click="showAddForm = true"
        />
        <Button
          v-else
          label="Cancel"
          icon="pi pi-times"
          severity="secondary"
          outlined
          size="large"
          @click="cancelAdd"
        />
      </div>

      <Message v-if="residenceStore.error" severity="error" class="mb-3" :closable="false">
        {{ residenceStore.error }}
      </Message>

      <Card v-if="showAddForm" class="mb-3">
        <template #title>New resident</template>
        <template #content>
          <form class="form-grid" @submit.prevent="handleAddResident">
            <FloatLabel>
              <InputText id="r-name" v-model="form.name" required fluid />
              <label for="r-name">Name</label>
            </FloatLabel>
            <FloatLabel>
              <InputText id="r-phone" v-model="form.phone" fluid />
              <label for="r-phone">Phone (optional)</label>
            </FloatLabel>
            <FloatLabel>
              <InputText id="r-email" v-model="form.email" type="email" fluid />
              <label for="r-email">Email (optional)</label>
            </FloatLabel>
            <FloatLabel>
              <DatePicker id="r-start" v-model="startDateModel" dateFormat="yy-mm-dd" showIcon fluid />
              <label for="r-start">Move-in date</label>
            </FloatLabel>
            <FloatLabel>
              <DatePicker id="r-end" v-model="endDateModel" dateFormat="yy-mm-dd" showIcon fluid />
              <label for="r-end">Move-out date (optional)</label>
            </FloatLabel>
            <div class="checkbox-row">
              <Checkbox v-model="form.is_active" inputId="r-active" binary />
              <label for="r-active">Active resident</label>
            </div>
            <Button
              type="submit"
              label="Add resident"
              icon="pi pi-check"
              size="large"
              class="submit-btn"
              :loading="residenceStore.loading"
            />
          </form>
        </template>
      </Card>

      <Card class="mb-3">
        <template #title>
          Current residents ({{ residenceStore.currentResidents.length }})
        </template>
        <template #content>
          <div v-if="residenceStore.loading && residenceStore.currentResidents.length === 0" class="state-block">
            <ProgressSpinner />
          </div>
          <p v-else-if="residenceStore.currentResidents.length === 0" class="empty-mini">
            No current residents yet.
          </p>
          <div v-else class="resident-list">
            <div v-for="r in residenceStore.currentResidents" :key="r.id" class="resident-row">
              <div>
                <div class="resident-name">{{ r.name }}</div>
                <div class="resident-meta">
                  <span>Move-in: {{ r.start_date }}</span>
                  <span v-if="r.phone"> · {{ r.phone }}</span>
                  <span v-else-if="r.email"> · {{ r.email }}</span>
                </div>
              </div>
              <Button
                label="Move out"
                severity="secondary"
                outlined
                size="large"
                @click="startMoveOut(r)"
              />
            </div>
          </div>

          <Card v-if="moveOutTarget" class="moveout-card mt-3">
            <template #title>Move out: {{ moveOutTarget.name }}</template>
            <template #content>
              <FloatLabel class="mb-3">
                <DatePicker v-model="moveOutDateModel" dateFormat="yy-mm-dd" showIcon fluid />
                <label>End date</label>
              </FloatLabel>
              <div class="dialog-actions">
                <Button
                  label="Confirm move out"
                  icon="pi pi-check"
                  size="large"
                  :loading="residenceStore.loading"
                  :disabled="!moveOutEndDate"
                  @click="confirmMoveOut"
                />
                <Button
                  label="Cancel"
                  severity="secondary"
                  outlined
                  size="large"
                  @click="cancelMoveOut"
                />
              </div>
            </template>
          </Card>
        </template>
      </Card>

      <Card>
        <template #title>
          Former residents ({{ residenceStore.formerResidents.length }})
        </template>
        <template #content>
          <Button
            :label="showFormer ? 'Hide former residents' : 'Show former residents'"
            :icon="showFormer ? 'pi pi-chevron-up' : 'pi pi-chevron-down'"
            severity="secondary"
            text
            size="large"
            class="mb-3"
            @click="showFormer = !showFormer"
          />
          <div v-if="showFormer">
            <p v-if="residenceStore.formerResidents.length === 0" class="empty-mini">
              No former residents yet.
            </p>
            <div v-else class="resident-list">
              <div
                v-for="r in residenceStore.formerResidents"
                :key="r.id"
                class="resident-row former"
              >
                <div>
                  <div class="resident-name">{{ r.name }}</div>
                  <div class="resident-meta">
                    <span>Move-in: {{ r.start_date }} · End: {{ r.end_date ?? '—' }}</span>
                    <span v-if="r.phone"> · {{ r.phone }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import DatePicker from 'primevue/datepicker'
import FloatLabel from 'primevue/floatlabel'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import { useToast } from 'primevue/usetoast'
import { usePropertyStore } from '@/stores/property'
import { useResidenceStore } from '@/stores/residence'
import type { CreateResidenceRequest, UpdateResidenceRequest } from '@/types/residence'

const propertyStore = usePropertyStore()
const residenceStore = useResidenceStore()
const toast = useToast()

const properties = computed(() => propertyStore.properties ?? [])

const selectedPropertyId = ref<number | null>(null)
const showAddForm = ref(false)
const showFormer = ref(false)

const form = ref<CreateResidenceRequest>({
  property_id: 0,
  name: '',
  phone: '',
  email: '',
  start_date: todayISO(),
  end_date: undefined,
  is_active: true,
})

const startDateModel = computed({
  get: () => parseDate(form.value.start_date),
  set: (val: Date | null) => {
    if (val) form.value.start_date = toISO(val)
  },
})

const endDateModel = computed({
  get: () => (form.value.end_date ? parseDate(form.value.end_date) : null),
  set: (val: Date | null) => {
    form.value.end_date = val ? toISO(val) : undefined
  },
})

function todayISO(): string {
  return new Date().toISOString().split('T')[0]
}

function parseDate(iso: string): Date {
  const [y, m, d] = iso.split('-').map(Number)
  return new Date(y, m - 1, d)
}

function toISO(date: Date): string {
  return date.toISOString().split('T')[0]
}

const selectedProperty = computed(() => {
  if (selectedPropertyId.value === null) return null
  return properties.value.find((p) => p.id === selectedPropertyId.value) ?? null
})

const selectedPropertyName = computed(() => selectedProperty.value?.name ?? 'Unknown')

onMounted(() => {
  propertyStore.fetchProperties()
})

watch(
  () => properties.value,
  (props) => {
    if (selectedPropertyId.value === null && props.length > 0) {
      selectedPropertyId.value = props[0].id
    }
  },
)

watch(
  () => selectedPropertyId.value,
  async (id) => {
    if (id === null) return
    showAddForm.value = false
    showFormer.value = false
    await residenceStore.fetchCurrentResidents(id)
    await residenceStore.fetchFormerResidents(id)
  },
)

function selectProperty(propertyId: number) {
  selectedPropertyId.value = propertyId
}

function resetForm() {
  form.value = {
    property_id: selectedPropertyId.value ?? 0,
    name: '',
    phone: '',
    email: '',
    start_date: todayISO(),
    end_date: undefined,
    is_active: true,
  }
}

function cancelAdd() {
  showAddForm.value = false
  resetForm()
}

async function handleAddResident() {
  if (selectedPropertyId.value === null) return

  const data: CreateResidenceRequest = {
    property_id: selectedPropertyId.value,
    name: form.value.name.trim(),
    phone: form.value.phone.trim(),
    email: form.value.email.trim(),
    start_date: form.value.start_date,
    end_date: form.value.end_date || undefined,
    is_active: form.value.is_active,
  }

  await residenceStore.createResident(data)
  await residenceStore.fetchCurrentResidents(selectedPropertyId.value)
  await residenceStore.fetchFormerResidents(selectedPropertyId.value)
  showAddForm.value = false
  resetForm()
  toast.add({ severity: 'success', summary: 'Resident added', life: 3000 })
}

const moveOutTargetId = ref<number | null>(null)
const moveOutEndDate = ref<string>(todayISO())

const moveOutDateModel = computed({
  get: () => parseDate(moveOutEndDate.value),
  set: (val: Date | null) => {
    if (val) moveOutEndDate.value = toISO(val)
  },
})

const moveOutTarget = computed(() => {
  if (moveOutTargetId.value === null) return null
  return (
    residenceStore.currentResidents.find((r) => r.id === moveOutTargetId.value) ??
    residenceStore.formerResidents.find((r) => r.id === moveOutTargetId.value) ??
    null
  )
})

function startMoveOut(r: { id: number }) {
  moveOutTargetId.value = r.id
  moveOutEndDate.value = todayISO()
}

function cancelMoveOut() {
  moveOutTargetId.value = null
  moveOutEndDate.value = todayISO()
}

async function confirmMoveOut() {
  if (selectedPropertyId.value === null) return
  if (!moveOutTarget.value) return
  if (!moveOutEndDate.value) return

  const r = moveOutTarget.value
  const data: UpdateResidenceRequest = {
    property_id: r.property_id,
    name: r.name,
    phone: r.phone,
    email: r.email,
    start_date: r.start_date,
    end_date: moveOutEndDate.value,
    is_active: false,
  }

  await residenceStore.updateResident(r.id, data)
  await residenceStore.fetchCurrentResidents(selectedPropertyId.value)
  await residenceStore.fetchFormerResidents(selectedPropertyId.value)
  cancelMoveOut()
  toast.add({ severity: 'success', summary: 'Resident moved out', life: 3000 })
}
</script>

<style scoped>
.page-header {
  margin-bottom: 1.5rem;
}

.property-picker {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.property-chip {
  flex: 1 1 auto;
  min-width: 8rem;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.section-title {
  margin: 0;
  font-size: 1.2rem;
}

.detail-subtitle {
  margin: 0.25rem 0 0;
  color: var(--text-secondary);
}

.form-grid {
  display: grid;
  gap: 1.25rem;
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.submit-btn {
  width: 100%;
}

.resident-list {
  display: grid;
  gap: 0.75rem;
}

.resident-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  padding: 0.75rem;
  border-radius: var(--p-border-radius-md);
  background: var(--p-surface-50);
  border: 1px solid var(--p-surface-200);
}

.resident-row.former {
  opacity: 0.9;
}

.resident-name {
  font-weight: 700;
  margin-bottom: 0.2rem;
}

.resident-meta {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.state-block,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  text-align: center;
  color: var(--text-secondary);
}

.empty-mini {
  margin: 0;
  color: var(--text-secondary);
}

.dialog-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

@media (min-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr 1fr;
  }

  .submit-btn {
    grid-column: 1 / -1;
    max-width: 16rem;
  }

  .dialog-actions {
    flex-direction: row;
  }
}
</style>
