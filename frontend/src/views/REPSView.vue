<template>
  <div class="page-view reps-view">
    <div class="page-header">
      <div>
        <h1 class="page-title">REPS time tracking</h1>
        <p class="page-subtitle">Log hours for tax compliance</p>
      </div>
      <div v-if="activeTab === '0'" class="header-actions">
        <Button
          label="Export"
          icon="pi pi-download"
          severity="secondary"
          outlined
          size="large"
          @click="showExportModal = true"
        />
        <Button label="Add entry" icon="pi pi-plus" size="large" @click="openAddForm" />
      </div>
    </div>

    <Tabs v-model:value="activeTab">
      <TabList>
        <Tab value="0">Time log</Tab>
        <Tab value="1">Quick capture</Tab>
      </TabList>
      <TabPanels>
        <TabPanel value="0">
          <Card class="mb-3">
            <template #content>
              <div class="summary-row">
                <span class="summary-label">Total hours ({{ currentYear }})</span>
                <span class="summary-value">{{ currentYearTotal.toFixed(2) }}</span>
              </div>
            </template>
          </Card>

          <div class="filters-grid mb-3">
            <Select
              v-model="selectedPropertyId"
              :options="propertyFilterOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="All properties"
              fluid
              @change="applyFilters"
            />
            <Select
              v-model="selectedCategory"
              :options="categoryFilterOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="All categories"
              fluid
              @change="applyFilters"
            />
            <DatePicker
              v-model="startDateModel"
              placeholder="Start date"
              showIcon
              fluid
              @update:modelValue="applyFilters"
            />
            <DatePicker
              v-model="endDateModel"
              placeholder="End date"
              showIcon
              fluid
              @update:modelValue="applyFilters"
            />
          </div>

          <Message v-if="store.error" severity="error" class="mb-3" :closable="false">
            {{ store.error }}
          </Message>

          <div v-if="store.loading && entries.length === 0" class="state-block">
            <ProgressSpinner />
            <p>Loading entries…</p>
          </div>

          <Card v-else-if="entries.length === 0">
            <template #content>
              <div class="empty-state">
                <i class="pi pi-clock empty-icon" aria-hidden="true" />
                <p>No time entries yet. Add one to get started.</p>
                <Button label="Add entry" icon="pi pi-plus" size="large" @click="openAddForm" />
              </div>
            </template>
          </Card>

          <DataTable
            v-else
            :value="entries"
            stripedRows
            scrollable
            scrollHeight="flex"
            class="entries-table"
            :pt="{
              table: { style: 'min-width: 40rem' },
            }"
            @row-click="onRowClick"
          >
            <Column field="date" header="Date" style="min-width: 7rem">
              <template #body="{ data }">
                {{ formatDate(data.date) }}
              </template>
            </Column>
            <Column header="Property" style="min-width: 8rem">
              <template #body="{ data }">
                {{ getPropertyName(data.property_id) }}
              </template>
            </Column>
            <Column field="category" header="Category" style="min-width: 9rem" />
            <Column field="hours" header="Hours" style="min-width: 4rem">
              <template #body="{ data }">
                <Tag :value="`${data.hours}h`" />
              </template>
            </Column>
            <Column field="description" header="Description" style="min-width: 12rem" />
            <Column header="" style="width: 4rem">
              <template #body="{ data }">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  aria-label="Delete entry"
                  :loading="store.loading"
                  @click.stop="confirmDelete(data.id)"
                />
              </template>
            </Column>
          </DataTable>
        </TabPanel>

        <TabPanel value="1">
          <REPSQuickCapture />
        </TabPanel>
      </TabPanels>
    </Tabs>

    <Dialog
      v-model:visible="showForm"
      modal
      :header="editingEntry ? 'Edit entry' : 'Add time entry'"
      :style="{ width: 'min(100%, 36rem)' }"
      :breakpoints="{ '768px': '95vw' }"
    >
      <REPSEntryForm
        :entry="editingEntry"
        :initial-property-id="filterPropertyForForm"
        :initial-date="endDate || undefined"
        @save="handleSave"
        @update="handleUpdate"
        @save-and-add-another="handleSaveAndAddAnother"
        @cancel="closeForm"
      />
    </Dialog>

    <Dialog
      v-model:visible="showExportModal"
      modal
      header="Export time entries"
      :style="{ width: 'min(100%, 28rem)' }"
    >
      <p class="export-text">Export filtered entries to CSV for tax filing.</p>
      <div class="dialog-actions">
        <Button label="Export CSV" icon="pi pi-download" size="large" @click="handleExport" />
        <Button
          label="Cancel"
          severity="secondary"
          outlined
          size="large"
          @click="showExportModal = false"
        />
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import DatePicker from 'primevue/datepicker'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import Tab from 'primevue/tab'
import TabList from 'primevue/tablist'
import TabPanel from 'primevue/tabpanel'
import TabPanels from 'primevue/tabpanels'
import Tabs from 'primevue/tabs'
import Tag from 'primevue/tag'
import REPSEntryForm from '@/components/REPSEntryForm.vue'
import REPSQuickCapture from '@/components/REPSQuickCapture.vue'
import { useREPSStore } from '@/stores/reps'
import { usePropertyStore } from '@/stores/property'
import { CATEGORIES } from '@/types/time_entry'
import type {
  TimeEntry,
  CreateTimeEntryRequest,
  UpdateTimeEntryRequest,
} from '@/types/time_entry'

const route = useRoute()
const router = useRouter()
const store = useREPSStore()
const propertyStore = usePropertyStore()
const confirm = useConfirm()
const toast = useToast()

const activeTab = ref('0')
const showForm = ref(false)
const showExportModal = ref(false)
const editingEntry = ref<TimeEntry | null>(null)
const selectedPropertyId = ref<number | null | 'all'>('all')
const selectedCategory = ref<string>('')
const startDate = ref<string>('')
const endDate = ref<string>('')

const entries = computed(() => store.entries)
const properties = computed(() => propertyStore.properties)
const currentYearTotal = computed(() => store.currentYearTotal)
const currentYear = computed(() => new Date().getFullYear())

const propertyFilterOptions = computed(() => [
  { label: 'All properties', value: 'all' },
  { label: 'General', value: 0 },
  ...properties.value.map((p) => ({ label: p.name, value: p.id })),
])

const categoryFilterOptions = computed(() => [
  { label: 'All categories', value: '' },
  ...CATEGORIES.map((c) => ({ label: c, value: c })),
])

const filterPropertyForForm = computed(() => {
  if (selectedPropertyId.value === 'all' || selectedPropertyId.value === 0) {
    return selectedPropertyId.value === 0 ? null : undefined
  }
  return selectedPropertyId.value as number
})

const startDateModel = computed({
  get: () => (startDate.value ? parseDate(startDate.value) : null),
  set: (val: Date | null) => {
    startDate.value = val ? toISO(val) : ''
  },
})

const endDateModel = computed({
  get: () => (endDate.value ? parseDate(endDate.value) : null),
  set: (val: Date | null) => {
    endDate.value = val ? toISO(val) : ''
  },
})

watch(
  () => route.query.tab,
  (tab) => {
    if (tab === 'capture') {
      activeTab.value = '1'
    }
  },
  { immediate: true },
)

watch(activeTab, (tab) => {
  const query = { ...route.query }
  if (tab === '1') {
    query.tab = 'capture'
  } else {
    delete query.tab
  }
  router.replace({ query })
})

onMounted(() => {
  store.fetchEntries()
  propertyStore.fetchProperties()
})

function parseDate(iso: string): Date {
  const [y, m, d] = iso.split('-').map(Number)
  return new Date(y, m - 1, d)
}

function toISO(date: Date): string {
  return date.toISOString().split('T')[0]
}

function applyFilters() {
  const filter: Record<string, unknown> = {}
  if (selectedPropertyId.value !== 'all') {
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

function openAddForm() {
  editingEntry.value = null
  showForm.value = true
}

function closeForm() {
  showForm.value = false
  editingEntry.value = null
}

function onRowClick(event: { data: TimeEntry }) {
  editingEntry.value = event.data
  showForm.value = true
}

function handleSave(data: CreateTimeEntryRequest) {
  store
    .createEntry(data)
    .then(() => {
      closeForm()
      applyFilters()
      toast.add({ severity: 'success', summary: 'Entry saved', life: 3000 })
    })
    .catch(() => {})
}

function handleUpdate(id: number, data: UpdateTimeEntryRequest) {
  store
    .updateEntry(id, data)
    .then(() => {
      closeForm()
      applyFilters()
      toast.add({ severity: 'success', summary: 'Entry updated', life: 3000 })
    })
    .catch(() => {})
}

function handleSaveAndAddAnother(data: CreateTimeEntryRequest) {
  store
    .createEntry(data)
    .then(() => {
      editingEntry.value = null
      applyFilters()
      toast.add({ severity: 'success', summary: 'Entry saved', life: 3000 })
    })
    .catch(() => {})
}

function confirmDelete(id: number) {
  confirm.require({
    message: 'Are you sure you want to delete this entry?',
    header: 'Delete entry',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    acceptClass: 'p-button-danger',
    accept: () => {
      store.deleteEntry(id).then(() => {
        applyFilters()
        toast.add({ severity: 'success', summary: 'Entry deleted', life: 3000 })
      })
    },
  })
}

function handleExport() {
  const filter: Record<string, unknown> = {}
  if (selectedPropertyId.value !== 'all') {
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
    toast.add({ severity: 'success', summary: 'Export started', life: 3000 })
  })
}
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.25rem;
  flex-wrap: wrap;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-weight: 500;
  color: var(--text-secondary);
}

.summary-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--p-primary-500);
}

.filters-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.entries-table {
  margin-top: 0.5rem;
}

.entries-table :deep(.p-datatable-tbody > tr) {
  cursor: pointer;
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

.export-text {
  margin: 0 0 1.25rem;
  color: var(--text-secondary);
}

.dialog-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

@media (min-width: 768px) {
  .filters-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .dialog-actions {
    flex-direction: row;
  }
}

@media (min-width: 1024px) {
  .filters-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}
</style>
