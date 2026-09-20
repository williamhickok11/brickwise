import { defineStore } from 'pinia'
import { ref } from 'vue'
import { qboApi } from '@/api/client'
import type { PropertyQBOMappingInput, QBOClass, QBOStatus, QBOSyncResult } from '@/types/qbo'

export const useQBOStore = defineStore('qbo', () => {
  const status = ref<QBOStatus | null>(null)
  const classes = ref<QBOClass[]>([])
  const loading = ref(false)
  const syncing = ref(false)
  const error = ref<string | null>(null)

  async function fetchStatus() {
    loading.value = true
    error.value = null
    try {
      status.value = await qboApi.status()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load QuickBooks status'
    } finally {
      loading.value = false
    }
  }

  async function connect() {
    error.value = null
    try {
      const { authorize_url } = await qboApi.connect()
      window.location.href = authorize_url
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to start QuickBooks connect'
    }
  }

  async function disconnect() {
    error.value = null
    try {
      await qboApi.disconnect()
      classes.value = []
      await fetchStatus()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to disconnect QuickBooks'
    }
  }

  async function fetchClasses() {
    error.value = null
    try {
      classes.value = await qboApi.classes()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load QuickBooks classes'
    }
  }

  async function saveMappings(mappings: PropertyQBOMappingInput[]) {
    error.value = null
    loading.value = true
    try {
      status.value = await qboApi.updateMappings(mappings)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to save mappings'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function sync(): Promise<QBOSyncResult | null> {
    error.value = null
    syncing.value = true
    try {
      const result = await qboApi.sync()
      await fetchStatus()
      return result
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to sync from QuickBooks'
      return null
    } finally {
      syncing.value = false
    }
  }

  return {
    status,
    classes,
    loading,
    syncing,
    error,
    fetchStatus,
    connect,
    disconnect,
    fetchClasses,
    saveMappings,
    sync,
  }
})
