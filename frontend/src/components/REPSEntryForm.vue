<template>
  <div class="entry-form">
    <h2>{{ editing ? 'Edit Entry' : 'Add Time Entry' }}</h2>

    <div class="form-group">
      <label>Date</label>
      <div class="date-inputs">
        <input
          v-model="form.date"
          type="date"
          class="input"
          required
        />
        <button
          type="button"
          @click="setYesterday"
          class="btn btn-secondary"
        >
          Yesterday
        </button>
      </div>
    </div>

    <div class="form-group">
      <label>Property</label>
      <select v-model="form.property_id" class="input">
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
      <label>Category</label>
      <div class="category-buttons">
        <button
          v-for="category in CATEGORIES"
          :key="category"
          type="button"
          @click="form.category = category"
          :class="['btn', 'btn-category', { active: form.category === category }]"
        >
          {{ category }}
        </button>
      </div>
    </div>

    <div class="form-group">
      <label>Hours</label>
      <input
        v-model.number="form.hours"
        type="number"
        step="0.25"
        min="0"
        max="24"
        class="input"
        placeholder="0.0"
        required
      />
    </div>

    <div class="form-group">
      <label>Description</label>
      <div class="description-input">
        <textarea
          v-model="form.description"
          class="input textarea"
          placeholder="Describe the work performed..."
          rows="3"
          required
        />
        <button
          type="button"
          @click="toggleVoiceInput"
          :class="['btn', 'btn-voice', { recording: isRecording }]"
          :disabled="!isVoiceSupported"
        >
          {{ isRecording ? '🎤 Recording...' : '🎤 Voice' }}
        </button>
      </div>
    </div>

    <div class="form-group">
      <label>Notes (optional)</label>
      <input
        v-model="form.notes"
        type="text"
        class="input"
        placeholder="Additional notes..."
      />
    </div>

    <div class="form-group">
      <label class="checkbox-label">
        <input
          v-model="form.full_drive"
          type="checkbox"
          @change="handleFullDriveChange"
        />
        <span>Full Drive</span>
      </label>
      <input
        v-if="form.full_drive"
        v-model.number="form.mileage"
        type="number"
        step="0.1"
        min="0"
        class="input"
        placeholder="Mileage"
      />
    </div>

    <div v-if="error" class="error">{{ error }}</div>

    <div class="form-actions">
      <button
        v-if="!editing"
        type="button"
        @click="handleSaveAndAddAnother"
        class="btn btn-primary"
        :disabled="loading"
      >
        {{ loading ? 'Saving...' : 'Save & Add Another' }}
      </button>
      <button
        type="button"
        @click="handleSave"
        class="btn btn-primary"
        :disabled="loading"
      >
        {{ loading ? 'Saving...' : editing ? 'Update' : 'Save & Done' }}
      </button>
      <button
        type="button"
        @click="$emit('cancel')"
        class="btn btn-secondary"
        :disabled="loading"
      >
        Cancel
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { usePropertyStore } from '@/stores/property'
import { CATEGORIES } from '@/types/time_entry'
import type { CreateTimeEntryRequest, UpdateTimeEntryRequest, TimeEntry } from '@/types/time_entry'

const props = defineProps<{
  entry?: TimeEntry | null
  initialPropertyId?: number | null
  initialDate?: string
}>()

const emit = defineEmits<{
  save: [data: CreateTimeEntryRequest]
  update: [id: number, data: UpdateTimeEntryRequest]
  cancel: []
  saveAndAddAnother: [data: CreateTimeEntryRequest]
}>()

const propertyStore = usePropertyStore()
const loading = ref(false)
const error = ref<string | null>(null)
const isRecording = ref(false)
const recognition = ref<SpeechRecognition | null>(null)

const editing = computed(() => !!props.entry)

const form = ref<CreateTimeEntryRequest>({
  property_id: props.initialPropertyId ?? null,
  date: props.initialDate ?? new Date().toISOString().split('T')[0],
  category: '',
  description: '',
  hours: 0,
  notes: '',
  mileage: 0,
  full_drive: false,
})

const properties = computed(() => propertyStore.properties)

const isVoiceSupported = computed(() => {
  return 'webkitSpeechRecognition' in window || 'SpeechRecognition' in window
})

onMounted(() => {
  if (props.entry) {
    form.value = {
      property_id: props.entry.property_id,
      date: props.entry.date.split('T')[0],
      category: props.entry.category,
      description: props.entry.description,
      hours: props.entry.hours,
      notes: props.entry.notes,
      mileage: props.entry.mileage,
      full_drive: props.entry.full_drive,
    }
  }

  if (isVoiceSupported.value) {
    const SpeechRecognition = (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition
    recognition.value = new SpeechRecognition()
    recognition.value.continuous = false
    recognition.value.interimResults = false
    recognition.value.lang = 'en-US'

    recognition.value.onresult = (event: any) => {
      const transcript = event.results[0][0].transcript
      form.value.description = transcript
      isRecording.value = false
    }

    recognition.value.onerror = () => {
      isRecording.value = false
      error.value = 'Voice recognition failed. Please try again.'
    }

    recognition.value.onend = () => {
      isRecording.value = false
    }
  }

  propertyStore.fetchProperties()
})

onUnmounted(() => {
  if (recognition.value && isRecording.value) {
    recognition.value.stop()
  }
})

watch(
  () => props.initialPropertyId,
  (newId) => {
    if (newId !== undefined) {
      form.value.property_id = newId
    }
  }
)

watch(
  () => props.initialDate,
  (newDate) => {
    if (newDate) {
      form.value.date = newDate
    }
  }
)

function setYesterday() {
  const yesterday = new Date()
  yesterday.setDate(yesterday.getDate() - 1)
  form.value.date = yesterday.toISOString().split('T')[0]
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

function handleFullDriveChange() {
  if (form.value.full_drive && form.value.property_id) {
    const property = properties.value.find((p) => p.id === form.value.property_id)
    if (property && property.default_mileage > 0) {
      form.value.mileage = property.default_mileage
    }
  }
}

function handleSave() {
  error.value = null
  if (editing.value && props.entry) {
    emit('update', props.entry.id, form.value as UpdateTimeEntryRequest)
  } else {
    emit('save', form.value)
  }
}

function handleSaveAndAddAnother() {
  error.value = null
  emit('saveAndAddAnother', form.value)
  // Reset form but keep property and date
  const savedPropertyId = form.value.property_id
  const savedDate = form.value.date
  form.value = {
    property_id: savedPropertyId,
    date: savedDate,
    category: '',
    description: '',
    hours: 0,
    notes: '',
    mileage: 0,
    full_drive: false,
  }
}
</script>

<style scoped>
.entry-form {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  max-width: 600px;
  margin: 0 auto;
}

h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.5rem;
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

.date-inputs {
  display: flex;
  gap: 0.5rem;
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
  min-height: 80px;
}

.description-input {
  position: relative;
}

.btn-voice {
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  padding: 0.5rem;
  font-size: 0.875rem;
}

.btn-voice.recording {
  background: #e74c3c;
  color: white;
}

.category-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
}

.btn-category {
  padding: 0.75rem;
  text-align: left;
  font-size: 0.875rem;
  white-space: normal;
}

.btn-category.active {
  background: #3498db;
  color: white;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.checkbox-label input[type='checkbox'] {
  width: auto;
  cursor: pointer;
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 2rem;
  flex-wrap: wrap;
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
  min-width: 120px;
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

@media (max-width: 768px) {
  .entry-form {
    padding: 1rem;
  }

  .category-buttons {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }
}
</style>
