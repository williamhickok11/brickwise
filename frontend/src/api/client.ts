import axios from 'axios'
import type { Property, CreatePropertyRequest } from '@/types/property'
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
