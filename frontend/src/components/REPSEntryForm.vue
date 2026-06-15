<template>
  <div class="entry-form">
    <h2 class="form-title">{{ editing ? 'Edit entry' : 'Add time entry' }}</h2>

    <div class="form-stack">
      <div class="date-row">
        <FloatLabel class="flex-1">
          <DatePicker
            id="entry-date"
            v-model="dateModel"
            dateFormat="yy-mm-dd"
            showIcon
            fluid
            class="w-full"
          />
          <label for="entry-date">Date</label>
        </FloatLabel>
        <Button
          type="button"
          label="Yesterday"
          severity="secondary"
          outlined
          size="large"
          @click="setYesterday"
        />
      </div>

      <FloatLabel>
        <Select
          id="entry-property"
          v-model="form.property_id"
          :options="propertyOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Property"
          fluid
          class="w-full"
        />
        <label for="entry-property">Property</label>
      </FloatLabel>

      <div>
        <label class="field-label">Category</label>
        <div class="category-grid">
          <Button
            v-for="category in CATEGORIES"
            :key="category"
            type="button"
            :label="category"
            :severity="form.category === category ? undefined : 'secondary'"
            :outlined="form.category !== category"
            size="large"
            class="category-btn"
            @click="form.category = category"
          />
        </div>
      </div>

      <FloatLabel>
        <InputNumber
          id="entry-hours"
          v-model="form.hours"
          :min="0"
          :max="24"
          :minFractionDigits="0"
          :maxFractionDigits="2"
          :step="0.25"
          fluid
          class="w-full"
        />
        <label for="entry-hours">Hours</label>
      </FloatLabel>

      <div class="description-wrap">
        <FloatLabel>
          <Textarea
            id="entry-description"
            v-model="form.description"
            rows="3"
            fluid
            class="w-full"
            required
          />
          <label for="entry-description">Description</label>
        </FloatLabel>
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

      <FloatLabel>
        <InputText id="entry-notes" v-model="form.notes" fluid class="w-full" />
        <label for="entry-notes">Notes (optional)</label>
      </FloatLabel>

      <div class="mileage-row">
        <div class="checkbox-row">
          <Checkbox v-model="form.full_drive" inputId="full-drive" binary @change="handleFullDriveChange" />
          <label for="full-drive">Full drive</label>
        </div>
        <FloatLabel v-if="form.full_drive" class="flex-1">
          <InputNumber
            id="entry-mileage"
            v-model="form.mileage"
            :min="0"
            :maxFractionDigits="1"
            fluid
            class="w-full"
          />
          <label for="entry-mileage">Mileage</label>
        </FloatLabel>
      </div>
    </div>

    <Message v-if="error" severity="error" class="mt-3" :closable="false">{{ error }}</Message>

    <div class="form-actions">
      <Button
        v-if="!editing"
        type="button"
        label="Save & add another"
        icon="pi pi-plus"
        size="large"
        class="action-btn"
        :loading="loading"
        @click="handleSaveAndAddAnother"
      />
      <Button
        type="button"
        :label="loading ? 'Saving…' : editing ? 'Update' : 'Save & done'"
        icon="pi pi-check"
        size="large"
        class="action-btn"
        :loading="loading"
        @click="handleSave"
      />
      <Button
        type="button"
        label="Cancel"
        severity="secondary"
        outlined
        size="large"
        class="action-btn"
        :disabled="loading"
        @click="$emit('cancel')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import DatePicker from 'primevue/datepicker'
import FloatLabel from 'primevue/floatlabel'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
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

const dateModel = computed({
  get: () => {
    const [y, m, d] = form.value.date.split('-').map(Number)
    return new Date(y, m - 1, d)
  },
  set: (val: Date | null) => {
    if (!val) return
    form.value.date = val.toISOString().split('T')[0]
  },
})

const properties = computed(() => propertyStore.properties)

const propertyOptions = computed(() => [
  { label: 'General', value: null },
  ...properties.value.map((p) => ({ label: p.name, value: p.id })),
])

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
    const SpeechRecognitionCtor = window.SpeechRecognition ?? window.webkitSpeechRecognition
    if (!SpeechRecognitionCtor) return
    recognition.value = new SpeechRecognitionCtor()
    recognition.value.continuous = false
    recognition.value.interimResults = false
    recognition.value.lang = 'en-US'

    recognition.value.onresult = (event: SpeechRecognitionEvent) => {
      form.value.description = event.results[0][0].transcript
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
  },
)

watch(
  () => props.initialDate,
  (newDate) => {
    if (newDate) {
      form.value.date = newDate
    }
  },
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
  padding: 0.25rem;
}

.form-title {
  margin: 0 0 1.25rem;
  font-size: 1.25rem;
}

.form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.field-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
}

.date-row {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
}

.category-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.5rem;
}

.category-btn {
  justify-content: flex-start;
  text-align: left;
}

.description-wrap {
  position: relative;
}

.voice-btn {
  position: absolute;
  right: 0.5rem;
  bottom: 0.5rem;
}

.mileage-row {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.checkbox-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.action-btn {
  width: 100%;
}

@media (min-width: 768px) {
  .category-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .mileage-row {
    flex-direction: row;
    align-items: flex-end;
  }
}
</style>
