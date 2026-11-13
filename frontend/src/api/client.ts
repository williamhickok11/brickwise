import axios from 'axios'
import type { Property, CreatePropertyRequest, UpdatePropertyValueRequest } from '@/types/property'
import type { UpsertMortgageRequest, Mortgage } from '@/types/mortgage'
import type { DashboardSummary } from '@/types/dashboard'
import type {
  Residence,
  CreateResidenceRequest,
  UpdateResidenceRequest,
  ResidenceFilter,
} from '@/types/residence'
import type {
  TimeEntry,
  CreateTimeEntryRequest,
  UpdateTimeEntryRequest,
  TimeEntryFilter,
} from '@/types/time_entry'
import type { PropertyQBOMappingInput, QBOClass, QBOStatus, QBOSyncResult } from '@/types/qbo'
import type { DocumentMeta } from '@/types/document'

const client = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

export const propertyApi = {
  list: (): Promise<Property[]> => {
    // console.log('propertyApi list')
    return client.get('/properties').then((res) => res.data)
  },

  get: (id: number): Promise<Property> => {
    return client.get(`/properties/${id}`).then((res) => res.data)
  },

  create: (data: CreatePropertyRequest): Promise<Property> => {
    return client.post('/properties', data).then((res) => res.data)
  },

  update: (id: number, data: CreatePropertyRequest): Promise<Property> => {
    return client.put(`/properties/${id}`, data).then((res) => res.data)
  },

  delete: (id: number): Promise<void> => {
    return client.delete(`/properties/${id}`).then(() => undefined)
  },

  updateValue: (id: number, data: UpdatePropertyValueRequest): Promise<Property> => {
    return client.put(`/properties/${id}/value`, data).then((res) => res.data)
  },
}

export const mortgageApi = {
  get: (propertyId: number): Promise<Mortgage> => {
    return client.get(`/properties/${propertyId}/mortgage`).then((res) => res.data)
  },

  upsert: (propertyId: number, data: UpsertMortgageRequest): Promise<Mortgage> => {
    return client.put(`/properties/${propertyId}/mortgage`, data).then((res) => res.data)
  },

  delete: (propertyId: number): Promise<void> => {
    return client.delete(`/properties/${propertyId}/mortgage`).then(() => undefined)
  },
}

export const dashboardApi = {
  summary: (): Promise<DashboardSummary> => {
    return client.get('/dashboard/summary').then((res) => res.data)
  },
}

export const qboApi = {
  status: (): Promise<QBOStatus> => {
    return client.get('/integrations/qbo/status').then((res) => res.data)
  },

  connect: (): Promise<{ authorize_url: string }> => {
    return client.get('/integrations/qbo/connect').then((res) => res.data)
  },

  disconnect: (): Promise<void> => {
    return client.delete('/integrations/qbo/disconnect').then(() => undefined)
  },

  classes: (): Promise<QBOClass[]> => {
    return client.get('/integrations/qbo/classes').then((res) => res.data)
  },

  updateMappings: (mappings: PropertyQBOMappingInput[]): Promise<QBOStatus> => {
    return client.put('/integrations/qbo/mappings', { mappings }).then((res) => res.data)
  },

  sync: (): Promise<QBOSyncResult> => {
    return client.post('/integrations/qbo/sync').then((res) => res.data)
  },
}

export const timeEntryApi = {
  list: (filter?: TimeEntryFilter): Promise<TimeEntry[]> => {
    const params = new URLSearchParams()
    if (filter?.property_id !== undefined) {
      if (filter.property_id === null || filter.property_id === 0) {
        params.append('property_id', 'null')
      } else {
        params.append('property_id', filter.property_id.toString())
      }
    }
    if (filter?.start_date) {
      params.append('start_date', filter.start_date)
    }
    if (filter?.end_date) {
      params.append('end_date', filter.end_date)
    }
    if (filter?.category) {
      params.append('category', filter.category)
    }
    const queryString = params.toString()
    const url = queryString ? `/time-entries?${queryString}` : '/time-entries'
    return client.get(url).then((res) => res.data)
  },

  get: (id: number): Promise<TimeEntry> => {
    return client.get(`/time-entries/${id}`).then((res) => res.data)
  },

  create: (data: CreateTimeEntryRequest): Promise<TimeEntry> => {
    return client.post('/time-entries', data).then((res) => res.data)
  },

  update: (id: number, data: UpdateTimeEntryRequest): Promise<TimeEntry> => {
    return client.put(`/time-entries/${id}`, data).then((res) => res.data)
  },

  delete: (id: number): Promise<void> => {
    return client.delete(`/time-entries/${id}`).then(() => undefined)
  },

  export: (filter?: TimeEntryFilter): Promise<Blob> => {
    const params = new URLSearchParams()
    if (filter?.property_id !== undefined) {
      if (filter.property_id === null || filter.property_id === 0) {
        params.append('property_id', 'null')
      } else {
        params.append('property_id', filter.property_id.toString())
      }
    }
    if (filter?.start_date) {
      params.append('start_date', filter.start_date)
    }
    if (filter?.end_date) {
      params.append('end_date', filter.end_date)
    }
    if (filter?.category) {
      params.append('category', filter.category)
    }
    const queryString = params.toString()
    const url = queryString ? `/time-entries/export?${queryString}` : '/time-entries/export'
    return client.get(url, { responseType: 'blob' }).then((res) => res.data)
  },
}

export const residenceApi = {
  list: (filter?: ResidenceFilter): Promise<Residence[]> => {
    const params = new URLSearchParams()
    if (filter?.property_id !== undefined) {
      params.append('property_id', filter.property_id.toString())
    }
    if (filter?.is_active !== undefined) {
      params.append('is_active', filter.is_active.toString())
    }
    const queryString = params.toString()
    const url = queryString ? `/residences?${queryString}` : '/residences'
    return client.get(url).then((res) => res.data)
  },

  get: (id: number): Promise<Residence> => {
    return client.get(`/residences/${id}`).then((res) => res.data)
  },

  create: (data: CreateResidenceRequest): Promise<Residence> => {
    return client.post('/residences', data).then((res) => res.data)
  },

  update: (id: number, data: UpdateResidenceRequest): Promise<Residence> => {
    return client.put(`/residences/${id}`, data).then((res) => res.data)
  },
}

export const documentApi = {
  get: (residenceId: number): Promise<DocumentMeta> => {
    return client.get(`/residences/${residenceId}/document`).then((res) => res.data)
  },

  upload: (residenceId: number, file: File): Promise<DocumentMeta> => {
    const form = new FormData()
    form.append('file', file)
    return client
      .post(`/residences/${residenceId}/document`, form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then((res) => res.data)
  },

  download: (residenceId: number): Promise<Blob> => {
    return client
      .get(`/residences/${residenceId}/document/download`, { responseType: 'blob' })
      .then((res) => res.data)
  },
}
