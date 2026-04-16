import { defineStore } from 'pinia'
import { ref } from 'vue'
import { residenceApi } from '@/api/client'
import type { CreateResidenceRequest, Residence, UpdateResidenceRequest } from '@/types/residence'

export const useResidenceStore = defineStore('residences', () => {
  const currentResidents = ref<Residence[]>([])
  const formerResidents = ref<Residence[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchCurrentResidents(propertyId: number) {
    loading.value = true
    error.value = null
    try {
      currentResidents.value = await residenceApi.list({
        property_id: propertyId,
        is_active: true,
      })
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch current residents'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchFormerResidents(propertyId: number) {
    loading.value = true
    error.value = null
    try {
      formerResidents.value = await residenceApi.list({
        property_id: propertyId,
        is_active: false,
      })
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch former residents'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createResident(data: CreateResidenceRequest) {
    loading.value = true
    error.value = null
    try {
      const residence = await residenceApi.create(data)
      // Optimistically update current/ former list based on `is_active`.
      if (residence.is_active) currentResidents.value.unshift(residence)
      else formerResidents.value.unshift(residence)
      return residence
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create resident'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateResident(id: number, data: UpdateResidenceRequest) {
    loading.value = true
    error.value = null
    try {
      const residence = await residenceApi.update(id, data)

      // Remove from both arrays first.
      currentResidents.value = currentResidents.value.filter((r) => r.id !== id)
      formerResidents.value = formerResidents.value.filter((r) => r.id !== id)

      // Add back to the right array.
      if (residence.is_active) currentResidents.value.unshift(residence)
      else formerResidents.value.unshift(residence)

      return residence
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update resident'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    currentResidents,
    formerResidents,
    loading,
    error,
    fetchCurrentResidents,
    fetchFormerResidents,
    createResident,
    updateResident,
  }
})

