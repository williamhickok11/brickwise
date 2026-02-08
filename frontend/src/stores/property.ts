import { defineStore } from 'pinia'
import { ref } from 'vue'
import { propertyApi } from '@/api/client'
import type { Property, CreatePropertyRequest } from '@/types/property'

export const usePropertyStore = defineStore('property', () => {
  const properties = ref<Property[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchProperties() {
    loading.value = true
    error.value = null
    try {
      properties.value = await propertyApi.list()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch properties'
    } finally {
      loading.value = false
    }
  }

  async function createProperty(data: CreatePropertyRequest) {
    loading.value = true
    error.value = null
    try {
      const property = await propertyApi.create(data)
      properties.value.push(property)
      return property
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create property'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteProperty(id: number) {
    loading.value = true
    error.value = null
    try {
      await propertyApi.delete(id)
      properties.value = properties.value.filter((p) => p.id !== id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete property'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    properties,
    loading,
    error,
    fetchProperties,
    createProperty,
    deleteProperty,
  }
})
