<template>
  <div class="page-view property-edit">
    <div class="page-header">
      <div>
        <Button
          label="Back to properties"
          icon="pi pi-arrow-left"
          text
          class="back-link"
          @click="goBack"
        />
        <h1 class="page-title">{{ pageTitle }}</h1>
        <p class="page-subtitle">Update details, photos, financials, or remove this property</p>
      </div>
    </div>

    <Message v-if="loadError" severity="error" class="mb-3" :closable="false">
      {{ loadError }}
    </Message>
    <Message v-else-if="store.error" severity="error" class="mb-3" :closable="false">
      {{ store.error }}
    </Message>

    <div v-if="loadingProperty" class="state-block">
      <ProgressSpinner />
      <p>Loading property…</p>
    </div>

    <template v-else-if="property">
      <Card class="mb-4">
        <template #title>Property details</template>
        <template #content>
          <form class="property-form" @submit.prevent="saveDetails">
            <div class="form-stack">
              <FloatLabel>
                <InputText id="edit-name" v-model="form.name" class="w-full" required fluid />
                <label for="edit-name">Property name</label>
              </FloatLabel>
              <FloatLabel>
                <InputText id="edit-address" v-model="form.address" class="w-full" required fluid />
                <label for="edit-address">Address</label>
              </FloatLabel>
              <FloatLabel>
                <Select
                  id="edit-type"
                  v-model="form.property_type"
                  :options="propertyTypes"
                  optionLabel="label"
                  optionValue="value"
                  placeholder="Select type"
                  class="w-full"
                  fluid
                />
                <label for="edit-type">Property type</label>
              </FloatLabel>
              <FloatLabel>
                <InputNumber
                  id="edit-mileage"
                  v-model="form.default_mileage"
                  :min="0"
                  :maxFractionDigits="1"
                  class="w-full"
                  fluid
                />
                <label for="edit-mileage">Default mileage (round trip)</label>
              </FloatLabel>
            </div>
            <Button
              type="submit"
              label="Save property details"
              icon="pi pi-check"
              size="large"
              class="mt-3 w-full"
              :loading="store.loading"
            />
          </form>
        </template>
      </Card>

      <Card class="mb-4">
        <template #title>Photos</template>
        <template #content>
          <p class="section-help">
            Upload multiple photos. The main photo is used as the thumbnail on the Properties list.
          </p>
          <div class="photo-toolbar">
            <input
              ref="imageInputRef"
              type="file"
              accept="image/jpeg,image/png,image/webp,image/gif,.jpg,.jpeg,.png,.webp,.gif"
              multiple
              class="file-input"
              @change="onImagesSelected"
            />
            <Button
              label="Add photos"
              icon="pi pi-images"
              size="large"
              :loading="uploadingImages"
              @click="imageInputRef?.click()"
            />
          </div>

          <div v-if="loadingImages" class="state-block compact">
            <ProgressSpinner style="width: 2rem; height: 2rem" />
            <p>Loading photos…</p>
          </div>
          <p v-else-if="images.length === 0" class="empty-mini">No photos yet.</p>
          <div v-else class="photo-grid">
            <div
              v-for="image in images"
              :key="image.id"
              class="photo-tile"
              :class="{ 'is-main': image.is_main }"
            >
              <img
                :src="propertyImageFileUrl(property.id, image.id)"
                :alt="image.file_name"
                class="photo-thumb"
              />
              <div class="photo-meta">
                <Tag v-if="image.is_main" value="Main" severity="info" />
                <span class="photo-name" :title="image.file_name">{{ image.file_name }}</span>
              </div>
              <div class="photo-actions">
                <Button
                  v-if="!image.is_main"
                  label="Set as main"
                  icon="pi pi-star"
                  size="small"
                  outlined
                  :loading="imageActionId === image.id"
                  @click="setMainImage(image.id)"
                />
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  aria-label="Delete photo"
                  :loading="imageActionId === image.id"
                  @click="confirmDeleteImage(image.id)"
                />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="mb-4">
        <template #title>Financials</template>
        <template #content>
          <form class="financials-form" @submit.prevent="saveFinancials">
            <h4 class="form-section-title">Property value</h4>
            <div class="form-stack">
              <FloatLabel>
                <InputNumber
                  id="edit-current-value"
                  v-model="financialsForm.current_value"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  class="w-full"
                  fluid
                />
                <label for="edit-current-value">Current value</label>
              </FloatLabel>
              <FloatLabel>
                <DatePicker
                  id="edit-value-updated"
                  v-model="valueUpdatedModel"
                  dateFormat="yy-mm-dd"
                  showIcon
                  fluid
                />
                <label for="edit-value-updated">Value as of</label>
              </FloatLabel>
            </div>

            <h4 class="form-section-title">Mortgage</h4>
            <div class="form-stack">
              <FloatLabel>
                <InputNumber
                  id="edit-current-balance"
                  v-model="financialsForm.mortgage.current_balance"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  class="w-full"
                  fluid
                />
                <label for="edit-current-balance">Current balance</label>
              </FloatLabel>
              <FloatLabel>
                <InputNumber
                  id="edit-original-loan"
                  v-model="financialsForm.mortgage.original_loan"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  class="w-full"
                  fluid
                />
                <label for="edit-original-loan">Original loan</label>
              </FloatLabel>
              <FloatLabel>
                <InputNumber
                  id="edit-interest-rate"
                  v-model="financialsForm.mortgage.interest_rate"
                  suffix="%"
                  :min="0"
                  :max="100"
                  :maxFractionDigits="3"
                  class="w-full"
                  fluid
                />
                <label for="edit-interest-rate">Interest rate</label>
              </FloatLabel>
              <FloatLabel>
                <InputNumber
                  id="edit-monthly-pi"
                  v-model="financialsForm.mortgage.monthly_pi"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  class="w-full"
                  fluid
                />
                <label for="edit-monthly-pi">Monthly P&amp;I</label>
              </FloatLabel>
              <FloatLabel>
                <InputNumber
                  id="edit-monthly-payment"
                  v-model="financialsForm.mortgage.monthly_payment"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  class="w-full"
                  fluid
                />
                <label for="edit-monthly-payment">Monthly payment (incl. escrow)</label>
              </FloatLabel>
              <FloatLabel>
                <DatePicker
                  id="edit-balance-as-of"
                  v-model="balanceAsOfModel"
                  dateFormat="yy-mm-dd"
                  showIcon
                  fluid
                />
                <label for="edit-balance-as-of">Balance as of</label>
              </FloatLabel>
              <FloatLabel>
                <InputText
                  id="edit-lender"
                  v-model="financialsForm.mortgage.lender"
                  class="w-full"
                  fluid
                />
                <label for="edit-lender">Lender</label>
              </FloatLabel>
            </div>

            <div class="dialog-actions">
              <Button
                type="submit"
                label="Save financials"
                icon="pi pi-check"
                size="large"
                :loading="savingFinancials"
              />
              <Button
                v-if="hasMortgage"
                type="button"
                label="Remove mortgage"
                icon="pi pi-trash"
                outlined
                size="large"
                :loading="savingFinancials"
                @click="removeMortgage"
              />
            </div>
          </form>
        </template>
      </Card>

      <Card class="danger-card">
        <template #title>Delete property</template>
        <template #content>
          <p class="danger-copy">
            Permanently remove this property from your portfolio. This cannot be undone.
          </p>
          <Button
            label="Delete property"
            icon="pi pi-trash"
            severity="danger"
            outlined
            size="large"
            :loading="store.loading"
            @click="confirmDelete"
          />
        </template>
      </Card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DatePicker from 'primevue/datepicker'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import FloatLabel from 'primevue/floatlabel'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import { propertyApi, mortgageApi, propertyImageApi, propertyImageFileUrl } from '@/api/client'
import { usePropertyStore } from '@/stores/property'
import { PROPERTY_TYPES, type CreatePropertyRequest, type Property } from '@/types/property'
import type { PropertyImage } from '@/types/property_image'
import type { UpsertMortgageRequest } from '@/types/mortgage'

const route = useRoute()
const router = useRouter()
const store = usePropertyStore()
const confirm = useConfirm()
const toast = useToast()

const propertyTypes = PROPERTY_TYPES
const property = ref<Property | null>(null)
const loadingProperty = ref(true)
const loadError = ref<string | null>(null)
const savingFinancials = ref(false)
const hasMortgage = ref(false)
const images = ref<PropertyImage[]>([])
const loadingImages = ref(false)
const uploadingImages = ref(false)
const imageActionId = ref<number | null>(null)
const imageInputRef = ref<HTMLInputElement | null>(null)

const emptyForm = (): CreatePropertyRequest => ({
  name: '',
  address: '',
  property_type: '',
  default_mileage: 0,
})

const form = ref<CreatePropertyRequest>(emptyForm())

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

const propertyId = computed(() => {
  const raw = route.params.id
  const id = Number(Array.isArray(raw) ? raw[0] : raw)
  return Number.isFinite(id) ? id : NaN
})

const pageTitle = computed(() => (property.value ? `Edit ${property.value.name}` : 'Edit property'))

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

function populateFromProperty(p: Property) {
  property.value = p
  form.value = {
    name: p.name,
    address: p.address,
    property_type: (p.property_type as CreatePropertyRequest['property_type']) || '',
    default_mileage: p.default_mileage ?? 0,
  }
  financialsForm.value = {
    current_value: p.current_value ?? 0,
    value_updated_at: p.value_updated_at ? p.value_updated_at.split('T')[0] : todayISO(),
    mortgage: emptyMortgage(),
  }
}

async function loadMortgage(id: number) {
  hasMortgage.value = false
  try {
    const m = await mortgageApi.get(id)
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

async function loadImages(id: number) {
  loadingImages.value = true
  try {
    images.value = await propertyImageApi.list(id)
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to load photos',
      detail: err instanceof Error ? err.message : undefined,
      life: 5000,
    })
    images.value = []
  } finally {
    loadingImages.value = false
  }
}

function syncMainImageId() {
  if (!property.value) return
  const main = images.value.find((img) => img.is_main)
  property.value.main_image_id = main ? main.id : null
  const idx = store.properties.findIndex((p) => p.id === property.value!.id)
  if (idx >= 0) {
    store.properties[idx] = {
      ...store.properties[idx],
      main_image_id: property.value.main_image_id,
    }
  }
}

async function onImagesSelected(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  input.value = ''
  if (!property.value || files.length === 0) return

  uploadingImages.value = true
  try {
    for (const file of files) {
      if (!file.type.startsWith('image/')) {
        toast.add({
          severity: 'warn',
          summary: 'Skipped file',
          detail: `${file.name} is not an image`,
          life: 4000,
        })
        continue
      }
      await propertyImageApi.upload(property.value.id, file)
    }
    await loadImages(property.value.id)
    syncMainImageId()
    toast.add({ severity: 'success', summary: 'Photos uploaded', life: 3000 })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to upload photos',
      detail: err instanceof Error ? err.message : undefined,
      life: 5000,
    })
  } finally {
    uploadingImages.value = false
  }
}

async function setMainImage(imageId: number) {
  if (!property.value) return
  imageActionId.value = imageId
  try {
    await propertyImageApi.setMain(property.value.id, imageId)
    images.value = images.value.map((img) => ({
      ...img,
      is_main: img.id === imageId,
    }))
    syncMainImageId()
    toast.add({ severity: 'success', summary: 'Main photo updated', life: 3000 })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: 'Failed to set main photo',
      detail: err instanceof Error ? err.message : undefined,
      life: 5000,
    })
  } finally {
    imageActionId.value = null
  }
}

function confirmDeleteImage(imageId: number) {
  if (!property.value) return
  confirm.require({
    message: 'Delete this photo?',
    header: 'Delete photo',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    accept: async () => {
      if (!property.value) return
      imageActionId.value = imageId
      try {
        await propertyImageApi.delete(property.value.id, imageId)
        await loadImages(property.value.id)
        syncMainImageId()
        toast.add({ severity: 'success', summary: 'Photo deleted', life: 3000 })
      } catch (err) {
        toast.add({
          severity: 'error',
          summary: 'Failed to delete photo',
          detail: err instanceof Error ? err.message : undefined,
          life: 5000,
        })
      } finally {
        imageActionId.value = null
      }
    },
  })
}

async function loadProperty(id: number) {
  loadingProperty.value = true
  loadError.value = null
  property.value = null
  images.value = []
  try {
    if (!store.properties.length) {
      await store.fetchProperties()
    }
    let found = store.properties.find((p) => p.id === id) ?? null
    if (!found) {
      found = await propertyApi.get(id)
    }
    populateFromProperty(found)
    await Promise.all([loadMortgage(id), loadImages(id)])
  } catch (err) {
    loadError.value = err instanceof Error ? err.message : 'Failed to load property'
  } finally {
    loadingProperty.value = false
  }
}

function goBack() {
  router.push({ name: 'properties' })
}

async function saveDetails() {
  if (!property.value) return
  try {
    const updated = await store.updateProperty(property.value.id, form.value)
    property.value = updated
    toast.add({ severity: 'success', summary: 'Property updated', life: 3000 })
  } catch {
    // Error handled by store
  }
}

async function saveFinancials() {
  if (!property.value) return
  savingFinancials.value = true
  try {
    const updated = await propertyApi.updateValue(property.value.id, {
      current_value: financialsForm.value.current_value,
      value_updated_at: financialsForm.value.value_updated_at,
    })
    const idx = store.properties.findIndex((p) => p.id === updated.id)
    if (idx >= 0) store.properties[idx] = updated
    property.value = updated

    const m = financialsForm.value.mortgage
    const hasMortgageData =
      m.current_balance > 0 ||
      m.original_loan > 0 ||
      m.monthly_pi > 0 ||
      m.monthly_payment > 0

    if (hasMortgageData) {
      await mortgageApi.upsert(property.value.id, m)
      hasMortgage.value = true
    }

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
  if (!property.value) return
  savingFinancials.value = true
  try {
    await mortgageApi.delete(property.value.id)
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

function confirmDelete() {
  if (!property.value) return
  const id = property.value.id
  confirm.require({
    message: 'Are you sure you want to delete this property?',
    header: 'Delete property',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancel',
    acceptLabel: 'Delete',
    accept: async () => {
      try {
        await store.deleteProperty(id)
        toast.add({ severity: 'success', summary: 'Property deleted', life: 3000 })
        goBack()
      } catch {
        // Error handled by store
      }
    },
  })
}

watch(
  propertyId,
  (id) => {
    if (!Number.isFinite(id)) {
      loadError.value = 'Invalid property id'
      loadingProperty.value = false
      return
    }
    loadProperty(id)
  },
  { immediate: true },
)
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

.back-link {
  margin: 0 0 0.35rem -0.5rem;
}

.section-help {
  margin: 0 0 1rem;
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.photo-toolbar {
  margin-bottom: 1rem;
}

.file-input {
  display: none;
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(10.5rem, 1fr));
  gap: 1rem;
}

.photo-tile {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.5rem;
  border: 1px solid var(--p-surface-200, #e2e8f0);
  border-radius: 0.5rem;
  background: var(--p-surface-0, #fff);
}

.photo-tile.is-main {
  border-color: var(--p-primary-400, #3aa39c);
  box-shadow: 0 0 0 1px var(--p-primary-200, #a8d4ef);
}

.photo-thumb {
  width: 100%;
  aspect-ratio: 4 / 3;
  object-fit: cover;
  border-radius: 0.35rem;
  background: var(--p-surface-100, #f1f5f9);
}

.photo-meta {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  min-height: 1.5rem;
}

.photo-name {
  font-size: 0.8rem;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.photo-actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  flex-wrap: wrap;
}

.empty-mini {
  margin: 0;
  color: var(--text-secondary);
}

.state-block.compact {
  padding: 1rem;
}

.form-stack {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-section-title {
  margin: 1.25rem 0 0.75rem;
  font-size: 0.95rem;
  color: var(--text-primary);
}

.form-section-title:first-child {
  margin-top: 0;
}

.dialog-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1.5rem;
}

.danger-card {
  border: 1px solid var(--p-red-200, #fecaca);
}

.danger-copy {
  margin: 0 0 1rem;
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.state-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
}

@media (min-width: 768px) {
  .form-stack {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
  }
}
</style>
