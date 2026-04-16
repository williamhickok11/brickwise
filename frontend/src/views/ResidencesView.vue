<template>
  <div class="residences">
    <div class="header">
      <h2>Residences</h2>
      <div class="subtext">Track current residents per property, then see former residents later.</div>
    </div>

    <div v-if="propertyStore.error" class="error">{{ propertyStore.error }}</div>

    <div class="section">
      <h3>Select a property</h3>

      <div v-if="propertyStore.loading && properties.length === 0" class="loading">
        Loading properties...
      </div>
      <div v-else-if="properties.length === 0" class="empty">
        No properties yet. Add one in `Properties`.
      </div>

      <div v-else class="property-list">
        <div
          v-for="property in properties"
          :key="property.id"
          class="property-card"
          :class="{ selected: selectedPropertyId === property.id }"
          role="button"
          tabindex="0"
          @click="selectProperty(property.id)"
          @keydown.enter="selectProperty(property.id)"
        >
          <div class="property-info">
            <h4>{{ property.name }}</h4>
            <p class="address">{{ property.address }}</p>
            <span class="badge">{{ property.property_type }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="selectedPropertyId !== null" class="section">
      <div class="detail-header">
        <div>
          <h3>Residents</h3>
          <p class="detail-subtitle">
            For: <strong>{{ selectedPropertyName }}</strong>
          </p>
        </div>

        <button v-if="!showAddForm" class="btn btn-primary" type="button" @click="showAddForm = true">
          Add Resident
        </button>
        <button v-else class="btn btn-primary" type="button" @click="cancelAdd">Cancel</button>
      </div>

      <div v-if="residenceStore.error" class="error">{{ residenceStore.error }}</div>

      <form v-if="showAddForm" class="resident-form" @submit.prevent="handleAddResident">
        <div class="form-grid">
          <label class="field">
            <span>Name</span>
            <input v-model="form.name" type="text" class="input" required placeholder="Resident name" />
          </label>

          <label class="field">
            <span>Phone</span>
            <input v-model="form.phone" type="text" class="input" placeholder="Optional" />
          </label>

          <label class="field">
            <span>Email</span>
            <input v-model="form.email" type="email" class="input" placeholder="Optional" />
          </label>

          <label class="field">
            <span>Move-in date</span>
            <input v-model="form.start_date" type="date" class="input" required />
          </label>

          <label class="field">
            <span>Move-out / end date</span>
            <input v-model="form.end_date" type="date" class="input" />
            <span class="hint">Optional (blank = unknown)</span>
          </label>

          <label class="field checkbox">
            <input v-model="form.is_active" type="checkbox" />
            <span>Active resident</span>
          </label>
        </div>

        <div class="form-actions">
          <button class="btn btn-primary" type="submit" :disabled="residenceStore.loading">
            {{ residenceStore.loading ? 'Saving...' : 'Add Resident' }}
          </button>
        </div>
      </form>

      <div class="current-section">
        <div class="section-title">
          <h4>Current residents</h4>
          <span class="muted">({{ residenceStore.currentResidents.length }})</span>
        </div>

        <div v-if="residenceStore.loading && residenceStore.currentResidents.length === 0" class="loading">
          Loading...
        </div>
        <div v-else-if="residenceStore.currentResidents.length === 0" class="empty">
          No current residents yet.
        </div>

        <div v-else class="resident-list">
          <div v-for="r in residenceStore.currentResidents" :key="r.id" class="resident-card">
            <div class="resident-info">
              <div class="resident-name">{{ r.name }}</div>
              <div class="resident-meta">
                <span class="meta-item">Move-in: {{ r.start_date }}</span>
                <span class="meta-item" v-if="r.phone">Phone: {{ r.phone }}</span>
                <span class="meta-item" v-else-if="r.email">Email: {{ r.email }}</span>
              </div>
            </div>
            <button class="btn btn-secondary" type="button" @click="startMoveOut(r)">Move out</button>
          </div>
        </div>

        <div v-if="moveOutTarget" class="moveout">
          <h4>Move out</h4>
          <p class="moveout-subtitle">
            Resident: <strong>{{ moveOutTarget.name }}</strong>
          </p>

          <label class="field">
            <span>End date</span>
            <input v-model="moveOutEndDate" type="date" class="input" required />
          </label>

          <div class="form-actions">
            <button
              class="btn btn-primary"
              type="button"
              :disabled="residenceStore.loading || !moveOutEndDate"
              @click="confirmMoveOut"
            >
              {{ residenceStore.loading ? 'Updating...' : 'Confirm move out' }}
            </button>
            <button class="btn btn-secondary" type="button" :disabled="residenceStore.loading" @click="cancelMoveOut">
              Cancel
            </button>
          </div>
        </div>
      </div>

      <div class="former-section">
        <div class="section-title">
          <h4>Former residents</h4>
          <span class="muted">({{ residenceStore.formerResidents.length }})</span>
        </div>

        <button class="btn btn-ghost" type="button" @click="showFormer = !showFormer">
          {{ showFormer ? 'Hide former residents' : 'Show former residents' }}
        </button>

        <div v-if="showFormer" class="resident-list">
          <div v-if="residenceStore.formerResidents.length === 0" class="empty">
            No former residents yet.
          </div>

          <div v-else>
            <div
              v-for="r in residenceStore.formerResidents"
              :key="r.id"
              class="resident-card resident-card-former"
            >
              <div class="resident-info">
                <div class="resident-name">{{ r.name }}</div>
                <div class="resident-meta">
                  <span class="meta-item">Move-in: {{ r.start_date }}</span>
                  <span class="meta-item">End: {{ r.end_date ?? '—' }}</span>
                  <span class="meta-item" v-if="r.phone">Phone: {{ r.phone }}</span>
                  <span class="meta-item" v-else-if="r.email">Email: {{ r.email }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { usePropertyStore } from '@/stores/property'
import { useResidenceStore } from '@/stores/residence'
import type { CreateResidenceRequest, UpdateResidenceRequest } from '@/types/residence'

const propertyStore = usePropertyStore()
const residenceStore = useResidenceStore()

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

function todayISO(): string {
  return new Date().toISOString().split('T')[0]
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
}

// Move-out workflow: updates a current resident to become former.
const moveOutTargetId = ref<number | null>(null)
const moveOutEndDate = ref<string>(todayISO())

const moveOutTarget = computed(() => {
  if (moveOutTargetId.value === null) return null
  return (
    residenceStore.currentResidents.find((r) => r.id === moveOutTargetId.value) ??
    residenceStore.formerResidents.find((r) => r.id === moveOutTargetId.value) ??
    null
  )
})

function startMoveOut(r: { id: number; name: string }) {
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
}
</script>

<style scoped>
.residences {
  padding: 1rem;
  padding-bottom: 80px; /* Space for bottom nav */
}

@media (min-width: 768px) {
  .residences {
    padding: 2rem;
    padding-bottom: 2rem;
  }
}

.header {
  margin-bottom: 1.5rem;
}

.header h2 {
  margin: 0 0 0.5rem 0;
}

.subtext {
  color: #666;
  font-size: 0.9rem;
}

.section {
  margin-bottom: 1.5rem;
}

.section h3 {
  margin: 0 0 1rem 0;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1rem;
}

.detail-subtitle {
  margin: 0.25rem 0 0 0;
  color: #666;
}

.property-list {
  display: grid;
  gap: 1rem;
}

.property-card {
  background: white;
  padding: 1.25rem;
  border-radius: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
}

.property-card.selected {
  border: 2px solid #3498db;
}

.property-info h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.05rem;
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

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s;
  cursor: pointer;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2980b9;
}

.btn-secondary {
  background: #95a5a6;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #7f8c8d;
}

.btn-ghost {
  background: transparent;
  border: 1px solid #ddd;
  color: #333;
  padding: 0.6rem 1rem;
  border-radius: 6px;
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
  padding: 2rem 1rem;
  color: #666;
}

.current-section,
.former-section {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  margin-bottom: 1.25rem;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.section-title h4 {
  margin: 0;
  font-size: 1.05rem;
}

.muted {
  color: #666;
  font-size: 0.95rem;
}

.resident-list {
  display: grid;
  gap: 0.75rem;
}

.resident-card {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: flex-start;
  background: #fafafa;
  padding: 0.9rem;
  border-radius: 8px;
  border: 1px solid #eee;
}

.resident-card-former {
  background: #f8f9fa;
}

.resident-name {
  font-weight: 700;
  margin-bottom: 0.25rem;
}

.resident-meta {
  display: grid;
  gap: 0.25rem;
}

.meta-item {
  color: #666;
  font-size: 0.9rem;
}

.resident-form {
  background: #f8f9fa;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.form-grid {
  display: grid;
  gap: 1rem;
}

.field {
  display: grid;
  gap: 0.4rem;
}

.field span {
  font-weight: 500;
  color: #333;
}

.input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  font-family: inherit;
}

.hint {
  color: #666;
  font-size: 0.85rem;
  margin-top: 0.1rem;
}

.checkbox {
  flex-direction: row;
  align-items: center;
  grid-template-columns: auto 1fr;
}

.form-actions {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
}

.moveout {
  margin-top: 1rem;
  padding: 1rem;
  border-radius: 8px;
  border: 2px solid #3498db;
  background: #f7fbff;
}

.moveout h4 {
  margin: 0 0 0.5rem 0;
}

.moveout-subtitle {
  margin: 0 0 1rem 0;
  color: #666;
}

@media (min-width: 768px) {
  .property-list {
    grid-template-columns: 1fr 1fr;
  }

  .form-grid {
    grid-template-columns: 1fr 1fr;
  }

  .resident-card {
    align-items: center;
  }
}
</style>
