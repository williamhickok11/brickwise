export interface PropertyBreakdown {
  property_id: number
  name: string
  current_value: number
  mortgage_balance: number
  equity: number
  monthly_rent: number
  monthly_pi: number
  net_cash_flow: number
  actual_cash_flow: number | null
  cash_flow_source: 'quickbooks' | 'estimate' | string
}

export interface DashboardSummary {
  as_of: string
  portfolio_value: number
  total_mortgage_balance: number
  total_equity: number
  monthly_revenue: number
  monthly_mortgage_pi: number
  estimated_cash_flow: number
  actual_cash_flow: number | null
  qbo_connected: boolean
  qbo_synced_at?: string
  cash_flow_period_start?: string
  cash_flow_period_end?: string
  properties: PropertyBreakdown[]
}
