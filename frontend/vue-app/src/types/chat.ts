// Chat types matching backend DTOs

export interface Thread {
  thread_id: string
  title: string
  model: string
  status: 'active' | 'archived'
  message_count: number
  last_message: string
  created_at: string
  updated_at: string
}

export interface Message {
  message_id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: string
}

export interface ThreadDetail extends Thread {
  messages: Message[]
}

export interface ThreadListResponse {
  threads: Thread[]
  next_key?: string
}

export interface SendMessageResponse {
  message: Message
  response: Message
}

export interface CreateThreadPayload {
  title?: string
  model?: string
}

// Model types from API
export interface Model {
  id: string
  name: string
  description: string
}

export interface ModelListResponse {
  models: Model[]
  default: string
}

export interface SendMessagePayload {
  content: string
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
}
