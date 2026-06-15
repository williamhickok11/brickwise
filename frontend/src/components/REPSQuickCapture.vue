<template>
  <div class="quick-capture">
    <Message v-if="store.error" severity="error" class="mb-3" :closable="false">
      {{ store.error }}
    </Message>

    <Card class="mb-3">
      <template #title>1) Capture</template>
      <template #content>
        <label class="field-label" for="free-form">Free form entry</label>
        <div class="text-input-wrapper">
          <Textarea
            id="free-form"
            v-model="freeFormText"
            rows="5"
            fluid
            class="w-full"
            placeholder="Example: Yesterday 2.5 hours maintenance fixed sink"
            @input="handleTextChange"
          />
          <Button
            v-if="isVoiceSupported"
            type="button"
            :icon="isRecording ? 'pi pi-stop-circle' : 'pi pi-microphone'"
            :label="isRecording ? 'Stop' : 'Voice'"
            :severity="isRecording ? 'danger' : 'secondary'"
            class="voice-btn"
            size="large"
            @click="toggleVoiceInput"
          />
        </div>
        <small class="hint">Tip: Include hours explicitly (e.g. 1.5h, 2 hours).</small>
      </template>
    </Card>

    <Card class="mb-3">
      <template #title>2) Review & edit</template>
      <template #content>
        <FloatLabel class="mb-3">
          <Select
            id="capture-property"
            v-model="selectedPropertyId"
            :options="propertyOptions"
            optionLabel="label"
            optionValue="value"
            fluid
            class="w-full"
          />
          <label for="capture-property">Property</label>
        </FloatLabel>

        <div v-if="parsedData && freeFormText.trim()" class="parsed-preview">
          <Message v-if="parsedData.confidence === 'low'" severity="warn" class="mb-3" :closable="false">
            Hours were not confidently parsed. Enter hours manually before saving.
          </Message>
          <div class="preview-grid">
            <FloatLabel>
              <DatePicker
                id="parsed-date"
                v-model="parsedDateModel"
                dateFormat="yy-mm-dd"
                showIcon
                fluid
              />
              <label for="parsed-date">Date</label>
            </FloatLabel>
            <Tag :value="parsedData.confidence" :severity="confidenceSeverity" />
            <FloatLabel>
              <InputNumber id="parsed-hours" v-model="parsedData.hours" :min="0" :max="24" :step="0.25" fluid />
              <label for="parsed-hours">Hours</label>
            </FloatLabel>
            <FloatLabel>
              <Select
                id="parsed-category"
                v-model="parsedData.category"
                :options="categoryOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Category"
                fluid
              />
              <label for="parsed-category">Category</label>
            </FloatLabel>
            <FloatLabel class="span-full">
              <Textarea id="parsed-desc" v-model="parsedData.description" rows="2" fluid />
              <label for="parsed-desc">Description</label>
            </FloatLabel>
            <FloatLabel>
              <InputNumber id="parsed-mileage" v-model="parsedData.mileage" :min="0" :maxFractionDigits="1" fluid />
              <label for="parsed-mileage">Mileage</label>
            </FloatLabel>
            <div class="checkbox-row">
              <Checkbox v-model="parsedData.full_drive" inputId="parsed-drive" binary />
              <label for="parsed-drive">Full drive</label>
            </div>
          </div>
        </div>
        <p v-else class="empty-mini">Add free-form text above to preview parsed fields.</p>
      </template>
    </Card>

    <Card class="mb-3">
      <template #title>3) Copy & save</template>
      <template #content>
        <div v-if="parsedData" class="copy-panel">
          <div class="copy-item">
            <label class="field-label">Summary</label>
            <Textarea :modelValue="summaryOutput" rows="2" readonly fluid class="w-full" />
            <Button
              label="Copy summary"
              icon="pi pi-copy"
              severity="secondary"
              outlined
              size="large"
              :disabled="!canCopy"
              @click="copySummary"
            />
          </div>
          <div class="copy-item">
            <label class="field-label">CSV row</label>
            <Textarea :modelValue="csvOutput" rows="2" readonly fluid class="w-full" />
            <Button
              label="Copy CSV row"
              icon="pi pi-copy"
              severity="secondary"
              outlined
              size="large"
              :disabled="!canCopy"
              @click="copyCsv"
            />
          </div>
        </div>
        <p v-else class="empty-mini">Parse an entry before copying or saving.</p>

        <div class="form-actions mt-3">
          <Button
            label="Save entry"
            icon="pi pi-check"
            size="large"
            class="action-btn"
            :disabled="!canSave"
            :loading="store.loading"
            @click="handleSave"
          />
          <Button
            label="Clear"
            icon="pi pi-times"
            severity="secondary"
            outlined
            size="large"
            class="action-btn"
            :disabled="store.loading"
            @click="handleClear"
          />
        </div>
      </template>
    </Card>

    <div v-if="recentEntries.length > 0">
      <h3 class="section-heading">Recent entries</h3>
      <div class="recent-list">
        <Card v-for="entry in recentEntries" :key="entry.id" class="recent-card">
          <template #content>
            <div class="recent-header">
              <span class="recent-date">{{ formatDate(entry.date) }}</span>
              <Tag :value="`${entry.hours}h`" />
            </div>
            <p class="recent-meta">{{ getPropertyName(entry.property_id) }} · {{ entry.category }}</p>
            <p class="recent-desc">{{ entry.description }}</p>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Checkbox from 'primevue/checkbox'
import DatePicker from 'primevue/datepicker'
import FloatLabel from 'primevue/floatlabel'
import InputNumber from 'primevue/inputnumber'
import Message from 'primevue/message'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import { useREPSStore } from '@/stores/reps'
import { usePropertyStore } from '@/stores/property'
import { parseREPSText, type ParsedREPSEntry } from '@/utils/repsParser'
import { CATEGORIES } from '@/types/time_entry'
import type { CreateTimeEntryRequest } from '@/types/time_entry'

const store = useREPSStore()
const propertyStore = usePropertyStore()
const toast = useToast()

const freeFormText = ref('')
const selectedPropertyId = ref<number | null>(null)
const parsedData = ref<ParsedREPSEntry | null>(null)
const isRecording = ref(false)
const recognition = ref<SpeechRecognition | null>(null)

const properties = computed(() => propertyStore.properties)
const recentEntries = computed(() => store.entries.slice(0, 5))

const propertyOptions = computed(() => [
  { label: 'General', value: null },
  ...properties.value.map((p) => ({ label: p.name, value: p.id })),
])

const categoryOptions = computed(() =>
  CATEGORIES.map((c) => ({ label: c, value: c })),
)

const parsedDateModel = computed({
  get: () => {
    if (!parsedData.value?.date) return null
    const [y, m, d] = parsedData.value.date.split('-').map(Number)
    return new Date(y, m - 1, d)
  },
  set: (val: Date | null) => {
    if (!parsedData.value || !val) return
    parsedData.value.date = val.toISOString().split('T')[0]
  },
})

const confidenceSeverity = computed(() => {
  if (!parsedData.value) return 'secondary'
  const map = { high: 'success', medium: 'warn', low: 'danger' } as const
  return map[parsedData.value.confidence] ?? 'secondary'
})

const isVoiceSupported = computed(() => {
  return 'webkitSpeechRecognition' in window || 'SpeechRecognition' in window
})

const canSave = computed(() => {
  if (!parsedData.value) return false
  return (
    parsedData.value.hours > 0 &&
    parsedData.value.description.trim().length > 0 &&
    parsedData.value.category.trim().length > 0
  )
})

const canCopy = computed(() => canSave.value)

const selectedPropertyName = computed(() => {
  if (selectedPropertyId.value === null) return 'General'
  const property = properties.value.find((p) => p.id === selectedPropertyId.value)
  return property?.name || 'Unknown'
})

const summaryOutput = computed(() => {
  if (!parsedData.value) return ''
  return `${parsedData.value.date} | ${selectedPropertyName.value} | ${parsedData.value.hours.toFixed(2)}h | ${parsedData.value.category} | ${parsedData.value.description.trim()}`
})

const csvOutput = computed(() => {
  if (!parsedData.value) return ''
  const row = [
    parsedData.value.date,
    selectedPropertyName.value,
    parsedData.value.category,
    parsedData.value.hours.toFixed(2),
    parsedData.value.description.trim(),
    parsedData.value.mileage.toString(),
    parsedData.value.full_drive ? 'true' : 'false',
  ]
  return row.map(csvEscape).join(',')
})

onMounted(() => {
  store.fetchEntries()
  propertyStore.fetchProperties()

  if (isVoiceSupported.value) {
    const SpeechRecognitionCtor = window.SpeechRecognition ?? window.webkitSpeechRecognition
    if (!SpeechRecognitionCtor) return
    recognition.value = new SpeechRecognitionCtor()
    recognition.value.continuous = false
    recognition.value.interimResults = false
    recognition.value.lang = 'en-US'

    recognition.value.onresult = (event: SpeechRecognitionEvent) => {
      const transcript = event.results[0][0].transcript
      freeFormText.value = freeFormText.value.trim()
        ? `${freeFormText.value.trim()} ${transcript}`.trim()
        : transcript
      handleTextChange()
      isRecording.value = false
    }

    recognition.value.onerror = () => {
      isRecording.value = false
    }

    recognition.value.onend = () => {
      isRecording.value = false
    }
  }
})

onUnmounted(() => {
  if (recognition.value && isRecording.value) {
    recognition.value.stop()
  }
})

function handleTextChange() {
  if (freeFormText.value.trim()) {
    parsedData.value = parseREPSText(freeFormText.value)
  } else {
    parsedData.value = null
  }
}

function toggleVoiceInput() {
  if (!recognition.value) return
  if (isRecording.value) {
    recognition.value.stop()
  } else {
    isRecording.value = true
    recognition.value.start()
  }
}

async function handleSave() {
  if (!parsedData.value || !canSave.value) return

  const entry: CreateTimeEntryRequest = {
    property_id: selectedPropertyId.value,
    date: parsedData.value.date,
    category: parsedData.value.category,
    description: parsedData.value.description,
    hours: parsedData.value.hours,
    notes: '',
    mileage: parsedData.value.mileage,
    full_drive: parsedData.value.full_drive,
  }

  try {
    await store.createEntry(entry)
    handleClear()
    await store.fetchEntries()
    toast.add({ severity: 'success', summary: 'Entry saved', life: 3000 })
  } catch {
    // Error handled by store
  }
}

function handleClear() {
  freeFormText.value = ''
  parsedData.value = null
}

async function copySummary() {
  await copyText(summaryOutput.value)
  toast.add({ severity: 'info', summary: 'Summary copied', life: 2000 })
}

async function copyCsv() {
  await copyText(csvOutput.value)
  toast.add({ severity: 'info', summary: 'CSV row copied', life: 2000 })
}

async function copyText(text: string) {
  if (!text || !canCopy.value) return
  await navigator.clipboard.writeText(text)
}

function csvEscape(value: string): string {
  const escaped = value.replace(/"/g, '""')
  return `"${escaped}"`
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
</script>

<style scoped>
.field-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
}

.text-input-wrapper {
  position: relative;
}

.voice-btn {
  position: absolute;
  right: 0.5rem;
  bottom: 0.5rem;
}

.hint {
  display: block;
  margin-top: 0.5rem;
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.parsed-preview {
  border: 1px solid var(--p-primary-200);
  border-radius: var(--p-border-radius-md);
  padding: 1rem;
  background: var(--p-primary-50);
}

.preview-grid {
  display: grid;
  gap: 1rem;
}

.span-full {
  grid-column: 1 / -1;
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.copy-panel {
  display: grid;
  gap: 1rem;
}

.copy-item {
  display: grid;
  gap: 0.5rem;
}

.empty-mini {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.action-btn {
  width: 100%;
}

.section-heading {
  margin: 0 0 1rem;
  font-size: 1.1rem;
}

.recent-list {
  display: grid;
  gap: 0.75rem;
}

.recent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.35rem;
}

.recent-date {
  font-weight: 600;
}

.recent-meta {
  margin: 0 0 0.25rem;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.recent-desc {
  margin: 0;
  font-size: 0.95rem;
}

@media (min-width: 768px) {
  .preview-grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
