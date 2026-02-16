<template>
  <div class="thread-sidebar">
    <div class="sidebar-header">
      <h3>Conversations</h3>
      <button @click="showCreateModal = true" class="new-thread-btn" :disabled="isLoading">
        + New
      </button>
    </div>

    <div class="threads-list" v-if="threads.length > 0">
      <div
        v-for="thread in threads"
        :key="thread.thread_id"
        class="thread-item"
        :class="{ active: currentThread?.thread_id === thread.thread_id }"
        @click="handleSelectThread(thread.thread_id)"
      >
        <div class="thread-info">
          <div class="thread-title">{{ thread.title || 'New Conversation' }}</div>
          <div class="thread-preview">{{ thread.last_message || 'No messages yet' }}</div>
          <div class="thread-meta">
            <span class="model-badge">{{ getModelName(thread.model) }}</span>
            <span class="message-count">{{ thread.message_count }} msgs</span>
            <span class="thread-date">{{ formatDate(thread.updated_at) }}</span>
          </div>
        </div>
        <button
          @click.stop="handleDeleteThread(thread.thread_id)"
          class="delete-btn"
          title="Delete thread"
        >
          x
        </button>
      </div>

      <button
        v-if="hasMoreThreads"
        @click="handleLoadMore"
        class="load-more-btn"
        :disabled="isLoading"
      >
        {{ isLoading ? 'Loading...' : 'Load More' }}
      </button>
    </div>

    <div v-else class="empty-state">
      <p>No conversations yet</p>
      <button @click="showCreateModal = true" class="start-btn">Start a conversation</button>
    </div>

    <!-- Create Thread Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-content">
        <h4>New Conversation</h4>

        <div class="form-group">
          <label for="model-select">Select Model</label>
          <select id="model-select" v-model="selectedModel" :disabled="availableModels.length === 0">
            <option
              v-for="model in availableModels"
              :key="model.id"
              :value="model.id"
            >
              {{ model.name }} - {{ model.description }}
            </option>
          </select>
        </div>

        <div class="modal-actions">
          <button @click="showCreateModal = false" class="cancel-btn">Cancel</button>
          <button @click="handleCreateThread" class="create-btn" :disabled="isLoading || availableModels.length === 0">
            {{ isLoading ? 'Creating...' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useChat } from '../composables/useChat'

const {
  threads,
  currentThread,
  isLoading,
  hasMoreThreads,
  availableModels,
  defaultModel,
  fetchModels,
  createThread,
  listThreads,
  selectThread,
  deleteThread,
} = useChat()

const emit = defineEmits<{
  (e: 'thread-selected', threadId: string): void
  (e: 'thread-created', threadId: string): void
}>()

// Modal state
const showCreateModal = ref(false)
const selectedModel = ref('')

onMounted(async () => {
  // Fetch models and threads in parallel
  await Promise.all([fetchModels(), listThreads()])
  // Set default model after fetching
  if (defaultModel.value) {
    selectedModel.value = defaultModel.value
  }
})

async function handleCreateThread() {
  // Title is auto-generated from the first message
  const thread = await createThread(undefined, selectedModel.value)
  if (thread) {
    await selectThread(thread.thread_id)
    emit('thread-created', thread.thread_id)
    // Reset modal state
    showCreateModal.value = false
    selectedModel.value = defaultModel.value
  }
}

function getModelName(modelId: string): string {
  const model = availableModels.value.find(m => m.id === modelId)
  return model?.name.replace('Gemini ', '') || modelId || 'Unknown'
}

async function handleSelectThread(threadId: string) {
  await selectThread(threadId)
  emit('thread-selected', threadId)
}

async function handleDeleteThread(threadId: string) {
  if (confirm('Are you sure you want to delete this conversation?')) {
    await deleteThread(threadId)
  }
}

async function handleLoadMore() {
  await listThreads(true)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
  } else if (diffDays === 1) {
    return 'Yesterday'
  } else if (diffDays < 7) {
    return date.toLocaleDateString('en-US', { weekday: 'short' })
  } else {
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }
}
</script>

<style scoped>
.thread-sidebar {
  width: 280px;
  background: #f8fafc;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1rem;
  color: #1e293b;
}

.new-thread-btn {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background 0.2s;
}

.new-thread-btn:hover:not(:disabled) {
  background: #2563eb;
}

.new-thread-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

.threads-list {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

.thread-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  margin-bottom: 0.25rem;
}

.thread-item:hover {
  background: #e2e8f0;
}

.thread-item.active {
  background: #dbeafe;
  border-left: 3px solid #3b82f6;
}

.thread-info {
  flex: 1;
  min-width: 0;
}

.thread-title {
  font-weight: 500;
  color: #1e293b;
  font-size: 0.875rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.thread-preview {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.thread-meta {
  display: flex;
  gap: 0.5rem;
  font-size: 0.625rem;
  color: #94a3b8;
  margin-top: 0.25rem;
}

.delete-btn {
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 1rem;
  line-height: 1;
  opacity: 0;
  transition: all 0.2s;
}

.thread-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: #fee2e2;
  color: #ef4444;
}

.load-more-btn {
  width: 100%;
  padding: 0.75rem;
  margin-top: 0.5rem;
  background: transparent;
  border: 1px dashed #cbd5e1;
  border-radius: 6px;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s;
}

.load-more-btn:hover:not(:disabled) {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.load-more-btn:disabled {
  cursor: not-allowed;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  text-align: center;
}

.empty-state p {
  color: #64748b;
  margin-bottom: 1rem;
}

.start-btn {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.start-btn:hover {
  background: #2563eb;
}

.model-badge {
  background: #e0f2fe;
  color: #0369a1;
  padding: 0.125rem 0.375rem;
  border-radius: 4px;
  font-size: 0.625rem;
  font-weight: 500;
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-content h4 {
  margin: 0 0 1.25rem 0;
  font-size: 1.125rem;
  color: #1e293b;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #475569;
  margin-bottom: 0.5rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.625rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 0.875rem;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.cancel-btn {
  background: #f1f5f9;
  color: #475569;
  border: none;
  padding: 0.625rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: background 0.2s;
}

.cancel-btn:hover {
  background: #e2e8f0;
}

.create-btn {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.625rem 1.25rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.create-btn:hover:not(:disabled) {
  background: #2563eb;
}

.create-btn:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .thread-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 280px;
    z-index: 20;
    transform: translateX(-100%);
    transition: transform 0.25s ease;
    box-shadow: none;
  }

  .thread-sidebar.sidebar-open {
    transform: translateX(0);
    box-shadow: 4px 0 16px rgba(0, 0, 0, 0.15);
  }

  .delete-btn {
    opacity: 1;
  }
}
</style>
