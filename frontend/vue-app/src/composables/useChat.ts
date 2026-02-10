import { ref, computed } from 'vue'
import type {
  Thread,
  ThreadDetail,
  Message,
  ThreadListResponse,
  SendMessageResponse,
  CreateThreadPayload,
  SendMessagePayload,
  ApiResponse,
  Model,
  ModelListResponse,
} from '../types/chat'

// Backend URL from environment
const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
// Streaming URL - uses Lambda Function URLs via CloudFront for true SSE streaming
const streamingUrl = import.meta.env.VITE_STREAMING_URL || backendUrl

// Shared state
const threads = ref<Thread[]>([])
const currentThread = ref<ThreadDetail | null>(null)
const isLoading = ref(false)
const isStreaming = ref(false)
const streamingContent = ref('')
const error = ref<string | null>(null)
const nextKey = ref<string | null>(null)
const availableModels = ref<Model[]>([])
const defaultModel = ref<string>('')

export function useChat() {
  // Computed
  const hasMoreThreads = computed(() => !!nextKey.value)

  // Fetch available models
  async function fetchModels(): Promise<Model[]> {
    try {
      const response = await fetch(`${backendUrl}/api/chat/models`, {
        method: 'GET',
        credentials: 'include',
      })

      const data: ApiResponse<ModelListResponse> = await response.json()

      if (response.ok && data.success && data.data) {
        availableModels.value = data.data.models
        defaultModel.value = data.data.default
        return data.data.models
      } else {
        console.error('Failed to fetch models:', data.error?.message)
        return []
      }
    } catch (err) {
      console.error('Fetch models error:', err)
      return []
    }
  }

  // Create a new thread
  async function createThread(title?: string, model?: string): Promise<Thread | null> {
    isLoading.value = true
    error.value = null

    try {
      const payload: CreateThreadPayload = {}
      if (title) {
        payload.title = title
      }
      if (model) {
        payload.model = model
      }

      const response = await fetch(`${backendUrl}/api/chat/threads`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(payload),
      })

      const data: ApiResponse<Thread> = await response.json()

      if (response.ok && data.success && data.data) {
        // Add to beginning of threads list
        threads.value.unshift(data.data)
        return data.data
      } else {
        error.value = data.error?.message || 'Failed to create thread'
        return null
      }
    } catch (err) {
      console.error('Create thread error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return null
    } finally {
      isLoading.value = false
    }
  }

  // List threads
  async function listThreads(loadMore = false): Promise<Thread[]> {
    isLoading.value = true
    error.value = null

    try {
      let url = `${backendUrl}/api/chat/threads?limit=20`
      if (loadMore && nextKey.value) {
        url += `&last_key=${encodeURIComponent(nextKey.value)}`
      }

      const response = await fetch(url, {
        method: 'GET',
        credentials: 'include',
      })

      const data: ApiResponse<ThreadListResponse> = await response.json()

      if (response.ok && data.success && data.data) {
        if (loadMore) {
          threads.value.push(...data.data.threads)
        } else {
          threads.value = data.data.threads
        }
        nextKey.value = data.data.next_key || null
        return data.data.threads
      } else {
        error.value = data.error?.message || 'Failed to fetch threads'
        return []
      }
    } catch (err) {
      console.error('List threads error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return []
    } finally {
      isLoading.value = false
    }
  }

  // Select a thread and load its messages
  async function selectThread(threadId: string): Promise<ThreadDetail | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/chat/threads/${threadId}`, {
        method: 'GET',
        credentials: 'include',
      })

      const data: ApiResponse<ThreadDetail> = await response.json()

      if (response.ok && data.success && data.data) {
        currentThread.value = data.data
        return data.data
      } else {
        error.value = data.error?.message || 'Failed to fetch thread'
        return null
      }
    } catch (err) {
      console.error('Select thread error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return null
    } finally {
      isLoading.value = false
    }
  }

  // Send a message (non-streaming)
  async function sendMessage(threadId: string, content: string): Promise<SendMessageResponse | null> {
    isLoading.value = true
    error.value = null

    try {
      const payload: SendMessagePayload = { content }

      const response = await fetch(`${backendUrl}/api/chat/threads/${threadId}/message`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(payload),
      })

      const data: ApiResponse<SendMessageResponse> = await response.json()

      if (response.ok && data.success && data.data) {
        // Add messages to current thread
        if (currentThread.value && currentThread.value.thread_id === threadId) {
          currentThread.value.messages.push(data.data.message)
          currentThread.value.messages.push(data.data.response)
          currentThread.value.message_count += 2
          currentThread.value.last_message = data.data.response.content
        }

        // Update thread in list
        const threadIndex = threads.value.findIndex(t => t.thread_id === threadId)
        if (threadIndex !== -1) {
          threads.value[threadIndex].message_count += 2
          threads.value[threadIndex].last_message = data.data.response.content
          threads.value[threadIndex].updated_at = new Date().toISOString()
        }

        return data.data
      } else {
        error.value = data.error?.message || 'Failed to send message'
        return null
      }
    } catch (err) {
      console.error('Send message error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return null
    } finally {
      isLoading.value = false
    }
  }

  // Send a message with streaming
  async function sendMessageStream(
    threadId: string,
    content: string,
    onChunk?: (chunk: string) => void
  ): Promise<boolean> {
    isStreaming.value = true
    streamingContent.value = ''
    error.value = null

    // Add user message to current thread immediately
    if (currentThread.value && currentThread.value.thread_id === threadId) {
      const userMessage: Message = {
        message_id: `temp-${Date.now()}`,
        role: 'user',
        content,
        timestamp: new Date().toISOString(),
      }
      // Initialize messages array if null
      if (!currentThread.value.messages) {
        currentThread.value.messages = []
      }
      currentThread.value.messages.push(userMessage)
    }

    try {
      const payload: SendMessagePayload = { content }

      const response = await fetch(`${streamingUrl}/api/chat/threads/${threadId}/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(payload),
      })

      if (!response.ok) {
        const data = await response.json()
        error.value = data.error?.message || 'Failed to send message'
        return false
      }

      const reader = response.body?.getReader()
      if (!reader) {
        error.value = 'Streaming not supported'
        return false
      }

      const decoder = new TextDecoder()
      let fullResponse = ''
      let currentEvent = 'message' // Track current SSE event type

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const chunk = decoder.decode(value, { stream: true })
        const lines = chunk.split('\n')

        for (const line of lines) {
          // Skip SSE comments (lines starting with :)
          if (line.startsWith(':')) {
            continue
          }
          if (line.startsWith('event: ')) {
            // Update current event type
            currentEvent = line.slice(7).trim()
          } else if (line.startsWith('data: ')) {
            const data = line.slice(6)

            if (currentEvent === 'error') {
              error.value = data
            } else if (currentEvent === 'done') {
              // Ignore done event data (contains message IDs)
            } else {
              // Regular message data - show immediately with typewriter effect
              fullResponse += data
              // Apply typewriter effect for each chunk
              for (const char of data) {
                streamingContent.value += char
                await new Promise(resolve => setTimeout(resolve, 8))
              }
              if (onChunk) {
                onChunk(data)
              }
            }
          } else if (line === '') {
            // Empty line resets event type to default
            currentEvent = 'message'
          }
        }
      }

      // Add assistant message to current thread
      if (currentThread.value && currentThread.value.thread_id === threadId) {
        const assistantMessage: Message = {
          message_id: `temp-${Date.now()}-assistant`,
          role: 'assistant',
          content: fullResponse,
          timestamp: new Date().toISOString(),
        }
        // Initialize messages array if null
        if (!currentThread.value.messages) {
          currentThread.value.messages = []
        }
        currentThread.value.messages.push(assistantMessage)
        currentThread.value.message_count = (currentThread.value.message_count || 0) + 2
        currentThread.value.last_message = fullResponse
      }

      // Update thread in list
      const threadIndex = threads.value.findIndex(t => t.thread_id === threadId)
      if (threadIndex !== -1) {
        threads.value[threadIndex].message_count += 2
        threads.value[threadIndex].last_message = fullResponse.slice(0, 100)
        threads.value[threadIndex].updated_at = new Date().toISOString()
      }

      return true
    } catch (err) {
      console.error('Stream message error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return false
    } finally {
      isStreaming.value = false
      streamingContent.value = ''
    }
  }

  // Delete a thread
  async function deleteThread(threadId: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${backendUrl}/api/chat/threads/${threadId}`, {
        method: 'DELETE',
        credentials: 'include',
      })

      const data: ApiResponse<{ message: string }> = await response.json()

      if (response.ok && data.success) {
        // Remove from threads list
        threads.value = threads.value.filter(t => t.thread_id !== threadId)

        // Clear current thread if it was deleted
        if (currentThread.value?.thread_id === threadId) {
          currentThread.value = null
        }

        return true
      } else {
        error.value = data.error?.message || 'Failed to delete thread'
        return false
      }
    } catch (err) {
      console.error('Delete thread error:', err)
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return false
    } finally {
      isLoading.value = false
    }
  }

  // Clear current thread selection
  function clearCurrentThread(): void {
    currentThread.value = null
  }

  // Clear error
  function clearError(): void {
    error.value = null
  }

  return {
    // State
    threads,
    currentThread,
    isLoading,
    isStreaming,
    streamingContent,
    error,
    hasMoreThreads,
    availableModels,
    defaultModel,

    // Actions
    fetchModels,
    createThread,
    listThreads,
    selectThread,
    sendMessage,
    sendMessageStream,
    deleteThread,
    clearCurrentThread,
    clearError,
  }
}
