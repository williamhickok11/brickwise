<template>
  <div class="page-view qbo-settings">
    <div class="page-header">
      <div>
        <h1 class="page-title">QuickBooks</h1>
        <p class="page-subtitle">Connect books for per-property cash flow</p>
      </div>
    </div>

    <Message v-if="banner" :severity="banner.severity" class="mb-3" :closable="true" @close="banner = null">
      {{ banner.text }}
    </Message>

    <Message v-if="qboStore.error" severity="error" class="mb-3" :closable="false">
      {{ qboStore.error }}
    </Message>

    <div v-if="qboStore.loading && !qboStore.status" class="state-block">
      <ProgressSpinner />
      <p>Loading QuickBooks status…</p>
    </div>

    <template v-else-if="qboStore.status">
      <Card class="mb-4">
        <template #title>Connection</template>
        <template #content>
          <p class="hint">{{ qboStore.status.message }}</p>

          <div class="status-row">
            <span class="status-label">App credentials</span>
            <Tag
              :severity="qboStore.status.configured ? 'success' : 'warn'"
              :value="qboStore.status.configured ? 'Configured' : 'Missing env vars'"
            />
          </div>
          <div class="status-row">
            <span class="status-label">Company link</span>
            <Tag
              :severity="qboStore.status.connected ? 'success' : 'secondary'"
              :value="qboStore.status.connected ? 'Connected' : 'Not connected'"
            />
          </div>
          <div v-if="qboStore.status.realm_id" class="status-row">
            <span class="status-label">Realm ID</span>
            <span class="mono">{{ qboStore.status.realm_id }}</span>
          </div>
          <div v-if="qboStore.status.last_synced_at" class="status-row">
            <span class="status-label">Last sync</span>
            <span>{{ formatDateTime(qboStore.status.last_synced_at) }}</span>
          </div>

          <div class="actions">
            <Button
              v-if="!qboStore.status.connected"
              label="Connect QuickBooks"
              icon="pi pi-link"
              size="large"
              :disabled="!qboStore.status.configured"
              @click="qboStore.connect()"
            />
            <template v-else>
              <Button
                label="Refresh from QuickBooks"
                icon="pi pi-refresh"
                size="large"
                :loading="qboStore.syncing"
                @click="handleSync"
              />
              <Button
                label="Disconnect"
                icon="pi pi-times"
                severity="secondary"
                outlined
                size="large"
                @click="handleDisconnect"
              />
            </template>
          </div>

          <Message v-if="!qboStore.status.configured" severity="warn" class="mt-3" :closable="false">
            Set <code>QBO_CLIENT_ID</code> and <code>QBO_CLIENT_SECRET</code> on the API server, then restart.
            Register an app at developer.intuit.com and add redirect URI
            <code>http://localhost:8080/api/v1/integrations/qbo/callback</code>.
          </Message>
        </template>
      </Card>

      <Card v-if="qboStore.status.connected">
        <template #title>Property → Class mapping</template>
        <template #content>
          <p class="hint">
            Your accounting team must use <strong>one QuickBooks Class per property</strong> and code
            income and expenses to that Class. Brickwise reads cash-basis net income by Class for the
            current month.
          </p>

          <div v-if="propertyStore.properties.length === 0" class="empty-mini">
            Add properties first, then map them to Classes.
          </div>

          <div v-else class="mapping-list">
            <div v-for="property in propertyStore.properties" :key="property.id" class="mapping-row">
              <div class="property-name">{{ property.name }}</div>
              <Select
                v-model="draft[property.id]"
                :options="classOptions"
                optionLabel="label"
                optionValue="value"
                placeholder="Select Class"
                class="w-full"
                showClear
                fluid
              />
            </div>
          </div>

          <Button
            class="mt-3"
            label="Save mappings"
            icon="pi pi-check"
            size="large"
            :loading="qboStore.loading"
            :disabled="propertyStore.properties.length === 0"
            @click="handleSaveMappings"
          />
        </template>
      </Card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import { usePropertyStore } from '@/stores/property'
import { useQBOStore } from '@/stores/qbo'
import type { PropertyQBOMappingInput } from '@/types/qbo'

const route = useRoute()
const router = useRouter()
const propertyStore = usePropertyStore()
const qboStore = useQBOStore()

const banner = ref<{ severity: 'success' | 'error' | 'info'; text: string } | null>(null)
const draft = reactive<Record<number, string | null>>({})

const classOptions = computed(() =>
  qboStore.classes.map((c) => ({
    label: c.name,
    value: c.id,
  })),
)

watch(
  () => qboStore.status,
  (status) => {
    if (!status) return
    for (const p of propertyStore.properties) {
      const link = status.mappings.find((m) => m.property_id === p.id)
      draft[p.id] = link?.qbo_class_id ?? null
    }
  },
  { immediate: true },
)

watch(
  () => propertyStore.properties,
  (props) => {
    const status = qboStore.status
    if (!status) return
    for (const p of props) {
      if (draft[p.id] === undefined) {
        const link = status.mappings.find((m) => m.property_id === p.id)
        draft[p.id] = link?.qbo_class_id ?? null
      }
    }
  },
  { deep: true },
)

async function handleSaveMappings() {
  const classById = new Map(qboStore.classes.map((c) => [c.id, c.name]))
  const mappings: PropertyQBOMappingInput[] = propertyStore.properties.map((p) => {
    const classId = draft[p.id] ?? ''
    return {
      property_id: p.id,
      qbo_class_id: classId,
      qbo_class_name: classId ? classById.get(classId) ?? '' : '',
    }
  })
  try {
    await qboStore.saveMappings(mappings)
    banner.value = { severity: 'success', text: 'Mappings saved.' }
  } catch {
    /* error set on store */
  }
}

async function handleSync() {
  const result = await qboStore.sync()
  if (result) {
    banner.value = {
      severity: 'success',
      text: `Synced ${result.updated} propert${result.updated === 1 ? 'y' : 'ies'} from QuickBooks.`,
    }
  }
}

async function handleDisconnect() {
  await qboStore.disconnect()
  banner.value = { severity: 'info', text: 'Disconnected from QuickBooks.' }
}

function formatDateTime(iso: string): string {
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(iso))
}

onMounted(async () => {
  if (route.query.connected === '1') {
    banner.value = { severity: 'success', text: 'QuickBooks connected successfully.' }
    router.replace({ path: route.path, query: {} })
  } else if (route.query.error) {
    banner.value = {
      severity: 'error',
      text: `QuickBooks connection failed: ${String(route.query.error)}`,
    }
    router.replace({ path: route.path, query: {} })
  }

  await Promise.all([propertyStore.fetchProperties(), qboStore.fetchStatus()])
  if (qboStore.status?.connected) {
    await qboStore.fetchClasses()
  }
})
</script>

<style scoped>
.page-header {
  margin-bottom: 1.5rem;
}

.hint {
  margin: 0 0 1.25rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

.status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.65rem 0;
  border-bottom: 1px solid var(--p-surface-200, #e2e8f0);
}

.status-label {
  font-weight: 600;
  font-size: 0.9rem;
}

.mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 0.85rem;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 1.25rem;
}

.mapping-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mapping-row {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.5rem;
}

.property-name {
  font-weight: 600;
}

.empty-mini {
  padding: 1rem 0;
  color: var(--text-secondary);
}

.state-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
  color: var(--text-secondary);
}

code {
  font-size: 0.85em;
}

@media (min-width: 640px) {
  .mapping-row {
    grid-template-columns: minmax(10rem, 14rem) 1fr;
    align-items: center;
  }
}
</style>
