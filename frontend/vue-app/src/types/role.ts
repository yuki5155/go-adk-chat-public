// Role types matching backend
export type UserRole = 'user' | 'subscriber' | 'admin' | 'root'
export type RequestStatus = 'pending' | 'approved' | 'rejected'
export type RoleStatus = 'active' | 'suspended'

// Role Request DTO
export interface RoleRequest {
  request_id: string
  user_id: string
  user_email: string
  requested_role: UserRole
  status: RequestStatus
  requested_at: string
  processed_at?: string
  processed_by?: string
  notes?: string
}

// User Role DTO
export interface UserRoleData {
  user_id: string
  email: string
  role: UserRole
  status: RoleStatus
  created_at: string
  updated_at: string
}

// API Response types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
}

export interface RoleRequestListResponse {
  requests: RoleRequest[]
  count: number
}

export interface UserRoleListResponse {
  users: UserRoleData[]
  count: number
}

// Request payloads
export interface RequestRolePayload {
  requested_role: 'user' | 'subscriber'
}

export interface ApproveRequestPayload {
  request_id: string
  notes?: string
}

export interface RejectRequestPayload {
  request_id: string
  notes: string
}
