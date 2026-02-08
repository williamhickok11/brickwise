<template>
  <div class="reps-view">
    <div class="header">
      <h1>REPS Time Tracking</h1>
      <div class="header-actions">
        <button @click="showExportModal = true" class="btn btn-secondary">
          Export
        </button>
        <button @click="showForm = true" class="btn btn-primary">
          Add Entry
        </button>
      </div>
    </div>

    <div class="summary">
      <div class="summary-item">
        <span class="summary-label">Total Hours ({{ currentYear }})</span>
        <span class="summary-value">{{ currentYearTotal.toFixed(2) }}</span>
      </div>
    </div>

    <div class="filters">
      <select v-model="selectedPropertyId" @change="applyFilters" class="input">
        <option :value="null">All Properties</option>
        <option :value="0">General</option>
        <option
          v-for="property in properties"
          :key="property.id"
          :value="property.id"
        >
          {{ property.name }}
        </option>
      </select>

      <select v-model="selectedCategory" @change="applyFilters" class="input">
        <option value="">All Categories</option>
        <option v-for="category in CATEGORIES" :key="category" :value="category">
          {{ category }}
        </option>
      </select>

      <input
        v-model="startDate"
        type="date"
        @change="applyFilters"
        class="input"
        placeholder="Start Date"
      />
      <input
        v-model="endDate"
        type="date"
        @change="applyFilters"
        class="input"
        placeholder="End Date"
      />
    </div>

    <div v-if="store.error" class="error">{{ store.error }}</div>

    <div v-if="store.loading && entries.length === 0" class="loading">
      Loading...
    </div>

    <div v-else-if="entries.length === 0" class="empty">
      No time entries yet. Add one above!
    </div>

    <div v-else class="entries-list">
      <div
        v-for="entry in entries"
        :key="entry.id"
        class="entry-card"
        @click="editEntry(entry)"
      >
        <div class="entry-header">
          <span class="entry-date">{{ formatDate(entry.date) }}</span>
          <span class="entry-hours">{{ entry.hours }}h</span>
        </div>
        <div class="entry-property">
          {{ getPropertyName(entry.property_id) }}
        </div>
        <div class="entry-category">{{ entry.category }}</div>
        <div class="entry-description">{{ entry.description }}</div>
        <div v-if="entry.notes" class="entry-notes">{{ entry.notes }}</div>
        <div v-if="entry.mileage > 0" class="entry-mileage">
          {{ entry.mileage }} miles
        </div>
        <div class="entry-actions">
          <button
            @click.stop="deleteEntry(entry.id)"
            class="btn btn-danger btn-small"
            :disabled="store.loading"
          >
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Entry Form Modal -->
    <div v-if="showForm" class="modal-overlay" @click="showForm = false">
      <div class="modal-content" @click.stop>
        <REPSEntryForm
          :entry="editingEntry"
          :initial-property-id="selectedPropertyId"
          :initial-date="endDate || undefined"
          @save="handleSave"
          @update="handleUpdate"
          @save-and-add-another="handleSaveAndAddAnother"
          @cancel="showForm = false"
        />
      </div>
    </div>

    <!-- Export Modal -->
    <div v-if="showExportModal" class="modal-overlay" @click="showExportModal = false">
      <div class="modal-content" @click.stop>
        <h2>Export Time Entries</h2>
        <p>Export filtered entries to CSV for tax filing.</p>
        <div class="form-actions">
          <button @click="handleExport" class="btn btn-primary">Export</button>
          <button @click="showExportModal = false" class="btn btn-secondary">
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useREPSStore } from '@/stores/reps'
import { usePropertyStore } from '@/stores/property'
import REPSEntryForm from '@/components/REPSEntryForm.vue'
import { CATEGORIES } from '@/types/time_entry'
import type {
  TimeEntry,
  CreateTimeEntryRequest,
  UpdateTimeEntryRequest,
} from '@/types/time_entry'

const store = useREPSStore()
const propertyStore = usePropertyStore()

const showForm = ref(false)
const showExportModal = ref(false)
const editingEntry = ref<TimeEntry | null>(null)
const selectedPropertyId = ref<number | null>(null)
const selectedCategory = ref<string>('')
const startDate = ref<string>('')
const endDate = ref<string>('')

const entries = computed(() => store.entries)
const properties = computed(() => propertyStore.properties)
const currentYearTotal = computed(() => store.currentYearTotal)
const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  store.fetchEntries()
  propertyStore.fetchProperties()
})

function applyFilters() {
  const filter: any = {}
  if (selectedPropertyId.value !== null) {
    filter.property_id = selectedPropertyId.value === 0 ? null : selectedPropertyId.value
  }
  if (selectedCategory.value) {
    filter.category = selectedCategory.value
  }
  if (startDate.value) {
    filter.start_date = startDate.value
  }
  if (endDate.value) {
    filter.end_date = endDate.value
  }
  store.fetchEntries(filter)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function getPropertyName(propertyId: number | null): string {
  if (propertyId === null) return 'General'
  const property = properties.value.find((p) => p.id === propertyId)
  return property?.name || 'Unknown'
}

function editEntry(entry: TimeEntry) {
  editingEntry.value = entry
  showForm.value = true
}

function handleSave(data: CreateTimeEntryRequest) {
  store
    .createEntry(data)
    .then(() => {
      showForm.value = false
      editingEntry.value = null
      applyFilters()
    })
    .catch(() => {
      // Error handled by store
    })
}

function handleUpdate(id: number, data: UpdateTimeEntryRequest) {
  store
    .updateEntry(id, data)
    .then(() => {
      showForm.value = false
      editingEntry.value = null
      applyFilters()
    })
    .catch(() => {
      // Error handled by store
    })
}

function handleSaveAndAddAnother(data: CreateTimeEntryRequest) {
  store
    .createEntry(data)
    .then(() => {
      editingEntry.value = null
      applyFilters()
      // Form stays open for next entry
    })
    .catch(() => {
      // Error handled by store
    })
}

function deleteEntry(id: number) {
  if (confirm('Are you sure you want to delete this entry?')) {
    store.deleteEntry(id).then(() => {
      applyFilters()
    })
  }
}

function handleExport() {
  const filter: any = {}
  if (selectedPropertyId.value !== null) {
    filter.property_id = selectedPropertyId.value === 0 ? null : selectedPropertyId.value
  }
  if (selectedCategory.value) {
    filter.category = selectedCategory.value
  }
  if (startDate.value) {
    filter.start_date = startDate.value
  }
  if (endDate.value) {
    filter.end_date = endDate.value
  }
  store.exportEntries(filter).then(() => {
    showExportModal.value = false
  })
}
</script>

<style scoped>
.reps-view {
  padding: 1rem;
  padding-bottom: 80px; /* Space for bottom nav */
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.header h1 {
  margin: 0;
  font-size: 1.75rem;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.summary {
  background: white;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.summary-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-weight: 500;
  color: #666;
}

.summary-value {
  font-size: 1.5rem;
  font-weight: 600;
  color: #3498db;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.input {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.entries-list {
  display: grid;
  gap: 1rem;
}

.entry-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.entry-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.entry-date {
  font-weight: 600;
  color: #333;
}

.entry-hours {
  font-weight: 600;
  color: #3498db;
}

.entry-property {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.25rem;
}

.entry-category {
  font-size: 0.875rem;
  color: #3498db;
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.entry-description {
  margin-bottom: 0.5rem;
  color: #333;
}

.entry-notes {
  font-size: 0.875rem;
  color: #666;
  font-style: italic;
  margin-bottom: 0.25rem;
}

.entry-mileage {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.entry-actions {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid #eee;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
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

.btn-danger {
  background: #e74c3c;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #c0392b;
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
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

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .reps-view {
    padding: 0.5rem;
  }

  .header {
    flex-direction: column;
    align-items: stretch;
  }

  .header-actions {
    width: 100%;
  }

  .header-actions .btn {
    flex: 1;
  }

  .filters {
    grid-template-columns: 1fr;
  }

  .entry-card {
    padding: 1rem;
  }
}
</style>
