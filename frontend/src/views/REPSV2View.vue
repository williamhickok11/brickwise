<template>
  <div class="reps-v2-view">
    <div class="header">
      <h1>REPS V2 - Free Form</h1>
      <p class="subtitle">Just type or speak naturally. We'll parse it.</p>
    </div>

    <div v-if="store.error" class="error">{{ store.error }}</div>

    <div class="entry-form">
      <div class="form-group">
        <label>Property</label>
        <select v-model="selectedPropertyId" class="input">
          <option :value="null">General</option>
          <option
            v-for="property in properties"
            :key="property.id"
            :value="property.id"
          >
            {{ property.name }}
          </option>
        </select>
      </div>

      <div class="form-group">
        <label>Free Form Entry</label>
        <div class="text-input-wrapper">
          <textarea
            v-model="freeFormText"
            @input="handleTextChange"
            class="input textarea"
            placeholder="Example: '2.5 hours of property management fixing the sink at Main St' or 'Spent 3 hours on maintenance and repairs, drove 22 miles'"
            rows="6"
          />
          <button
            v-if="isVoiceSupported"
            type="button"
            @click="toggleVoiceInput"
            :class="['btn', 'btn-voice', { recording: isRecording }]"
          >
            {{ isRecording ? '🎤 Recording...' : '🎤 Voice' }}
          </button>
        </div>
        <p class="hint">
          Try: "2 hours maintenance fixing sink" or "Yesterday: 1.5h contractor oversight"
        </p>
      </div>

      <!-- Parsed Preview -->
      <div v-if="parsedData && freeFormText.trim()" class="parsed-preview">
        <h3>Parsed Information</h3>
        <div class="preview-grid">
          <div class="preview-item">
            <span class="preview-label">Date:</span>
            <span class="preview-value">{{ formatDisplayDate(parsedData.date) }}</span>
            <span :class="['confidence-badge', parsedData.confidence]">
              {{ parsedData.confidence }}
            </span>
          </div>
          <div class="preview-item">
            <span class="preview-label">Hours:</span>
            <span class="preview-value">{{ parsedData.hours || 'Not found' }}</span>
          </div>
          <div class="preview-item">
            <span class="preview-label">Category:</span>
            <select v-model="parsedData.category" class="input input-small">
              <option v-for="cat in CATEGORIES" :key="cat" :value="cat">
                {{ cat }}
              </option>
            </select>
          </div>
          <div class="preview-item">
            <span class="preview-label">Description:</span>
            <textarea
              v-model="parsedData.description"
              class="input textarea-small"
              rows="2"
            />
          </div>
          <div v-if="parsedData.mileage > 0" class="preview-item">
            <span class="preview-label">Mileage:</span>
            <span class="preview-value">{{ parsedData.mileage }} miles</span>
          </div>
        </div>
      </div>

      <div class="form-actions">
        <button
          @click="handleSave"
          class="btn btn-primary"
          :disabled="!canSave || store.loading"
        >
          {{ store.loading ? 'Saving...' : 'Save Entry' }}
        </button>
        <button
          @click="handleClear"
          class="btn btn-secondary"
          :disabled="store.loading"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Recent Entries -->
    <div v-if="recentEntries.length > 0" class="recent-entries">
      <h2>Recent Entries</h2>
      <div class="entries-list">
        <div
          v-for="entry in recentEntries"
          :key="entry.id"
          class="entry-card"
        >
          <div class="entry-header">
            <span class="entry-date">{{ formatDate(entry.date) }}</span>
            <span class="entry-hours">{{ entry.hours }}h</span>
          </div>
          <div class="entry-property">{{ getPropertyName(entry.property_id) }}</div>
          <div class="entry-category">{{ entry.category }}</div>
          <div class="entry-description">{{ entry.description }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useREPSStore } from '@/stores/reps'
import { usePropertyStore } from '@/stores/property'
import { parseREPSText, type ParsedREPSEntry } from '@/utils/repsParser'
import { CATEGORIES } from '@/types/time_entry'
import type { CreateTimeEntryRequest } from '@/types/time_entry'

const store = useREPSStore()
const propertyStore = usePropertyStore()

const freeFormText = ref('')
const selectedPropertyId = ref<number | null>(null)
const parsedData = ref<ParsedREPSEntry | null>(null)
const isRecording = ref(false)
const recognition = ref<SpeechRecognition | null>(null)

const properties = computed(() => propertyStore.properties)
const recentEntries = computed(() => store.entries.slice(0, 5))

const isVoiceSupported = computed(() => {
  return 'webkitSpeechRecognition' in window || 'SpeechRecognition' in window
})

const canSave = computed(() => {
  if (!parsedData.value) return false
  return parsedData.value.hours > 0 && parsedData.value.description.trim().length > 0
})

onMounted(() => {
  store.fetchEntries()
  propertyStore.fetchProperties()

  if (isVoiceSupported.value) {
    const SpeechRecognition =
      (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition
    recognition.value = new SpeechRecognition()
    recognition.value.continuous = false
    recognition.value.interimResults = false
    recognition.value.lang = 'en-US'

    recognition.value.onresult = (event: any) => {
      const transcript = event.results[0][0].transcript
      freeFormText.value = transcript
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
  } catch (err) {
    // Error handled by store
  }
}

function handleClear() {
  freeFormText.value = ''
  parsedData.value = null
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function formatDisplayDate(dateString: string): string {
  const date = new Date(dateString)
  const today = new Date()
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)

  if (date.toDateString() === today.toDateString()) {
    return 'Today'
  } else if (date.toDateString() === yesterday.toDateString()) {
    return 'Yesterday'
  } else {
    return formatDate(dateString)
  }
}

function getPropertyName(propertyId: number | null): string {
  if (propertyId === null) return 'General'
  const property = properties.value.find((p) => p.id === propertyId)
  return property?.name || 'Unknown'
}
</script>

<style scoped>
.reps-v2-view {
  padding: 1rem;
  padding-bottom: 80px;
  max-width: 800px;
  margin: 0 auto;
}

.header {
  margin-bottom: 2rem;
}

.header h1 {
  margin: 0 0 0.5rem 0;
  font-size: 1.75rem;
  color: #2c3e50;
}

.subtitle {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.entry-form {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

.textarea {
  resize: vertical;
  min-height: 120px;
}

.text-input-wrapper {
  position: relative;
}

.btn-voice {
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  background: #95a5a6;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-voice.recording {
  background: #e74c3c;
}

.hint {
  margin: 0.5rem 0 0 0;
  font-size: 0.875rem;
  color: #666;
  font-style: italic;
}

.parsed-preview {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  margin-top: 1.5rem;
  border: 2px solid #3498db;
}

.parsed-preview h3 {
  margin: 0 0 1rem 0;
  font-size: 1.1rem;
  color: #2c3e50;
}

.preview-grid {
  display: grid;
  gap: 1rem;
}

.preview-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.preview-label {
  font-weight: 500;
  color: #666;
  min-width: 80px;
}

.preview-value {
  flex: 1;
  color: #333;
}

.input-small {
  flex: 1;
  min-width: 200px;
}

.textarea-small {
  flex: 1;
  min-width: 200px;
}

.confidence-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
}

.confidence-badge.high {
  background: #d4edda;
  color: #155724;
}

.confidence-badge.medium {
  background: #fff3cd;
  color: #856404;
}

.confidence-badge.low {
  background: #f8d7da;
  color: #721c24;
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  flex: 1;
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

.recent-entries {
  margin-top: 2rem;
}

.recent-entries h2 {
  margin: 0 0 1rem 0;
  font-size: 1.25rem;
  color: #2c3e50;
}

.entries-list {
  display: grid;
  gap: 1rem;
}

.entry-card {
  background: white;
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
  color: #333;
}

@media (max-width: 768px) {
  .reps-v2-view {
    padding: 0.5rem;
  }

  .entry-form {
    padding: 1rem;
  }

  .preview-item {
    flex-direction: column;
    align-items: stretch;
  }

  .preview-label {
    min-width: auto;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
