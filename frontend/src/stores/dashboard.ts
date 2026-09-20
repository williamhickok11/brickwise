import { defineStore } from 'pinia'
import { ref } from 'vue'
import { dashboardApi } from '@/api/client'
import type { DashboardSummary } from '@/types/dashboard'

export const useDashboardStore = defineStore('dashboard', () => {
  const summary = ref<DashboardSummary | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSummary() {
    loading.value = true
    error.value = null
    try {
      summary.value = await dashboardApi.summary()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch dashboard summary'
    } finally {
      loading.value = false
    }
  }

  return {
    summary,
    loading,
    error,
    fetchSummary,
  }
})
