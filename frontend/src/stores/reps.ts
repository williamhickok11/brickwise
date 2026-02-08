import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { timeEntryApi } from '@/api/client'
import type {
  TimeEntry,
  CreateTimeEntryRequest,
  UpdateTimeEntryRequest,
  TimeEntryFilter,
} from '@/types/time_entry'

export const useREPSStore = defineStore('reps', () => {
  const entries = ref<TimeEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const currentYearTotal = computed(() => {
    const currentYear = new Date().getFullYear()
    return entries.value
      .filter((e) => new Date(e.date).getFullYear() === currentYear)
      .reduce((sum, e) => sum + e.hours, 0)
  })

  async function fetchEntries(filter?: TimeEntryFilter) {
    loading.value = true
    error.value = null
    try {
      entries.value = await timeEntryApi.list(filter)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch time entries'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createEntry(data: CreateTimeEntryRequest) {
    loading.value = true
    error.value = null
    try {
      const entry = await timeEntryApi.create(data)
      entries.value.unshift(entry) // Add to beginning
      return entry
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create time entry'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateEntry(id: number, data: UpdateTimeEntryRequest) {
    loading.value = true
    error.value = null
    try {
      const entry = await timeEntryApi.update(id, data)
      const index = entries.value.findIndex((e) => e.id === id)
      if (index !== -1) {
        entries.value[index] = entry
      }
      return entry
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update time entry'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteEntry(id: number) {
    loading.value = true
    error.value = null
    try {
      await timeEntryApi.delete(id)
      entries.value = entries.value.filter((e) => e.id !== id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete time entry'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function exportEntries(filter?: TimeEntryFilter) {
    loading.value = true
    error.value = null
    try {
      const blob = await timeEntryApi.export(filter)
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `reps_time_entries_${new Date().toISOString().split('T')[0]}.csv`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to export time entries'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    entries,
    loading,
    error,
    currentYearTotal,
    fetchEntries,
    createEntry,
    updateEntry,
    deleteEntry,
    exportEntries,
  }
})
