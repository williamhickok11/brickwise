export interface QBOClass {
  id: string
  name: string
}

export interface PropertyQBOLink {
  property_id: number
  qbo_class_id: string
  qbo_class_name: string
  updated_at: string
}

export interface QBOStatus {
  configured: boolean
  connected: boolean
  realm_id?: string
  connected_at?: string
  last_synced_at?: string
  period_start?: string
  period_end?: string
  mappings: PropertyQBOLink[]
  class_required: boolean
  message?: string
}

export interface QBOSyncResult {
  synced_at: string
  period_start: string
  period_end: string
  updated: number
}

export interface PropertyQBOMappingInput {
  property_id: number
  qbo_class_id: string
  qbo_class_name: string
}
