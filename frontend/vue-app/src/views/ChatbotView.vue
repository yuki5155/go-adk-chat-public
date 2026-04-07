<template>
  <div class="chatbot-container">
    <!-- Mobile sidebar backdrop -->
    <div v-if="sidebarOpen" class="sidebar-backdrop" @click="sidebarOpen = false"></div>

    <!-- Thread Sidebar -->
    <ThreadSidebar
      :class="{ 'sidebar-open': sidebarOpen }"
      @thread-selected="handleThreadSelected"
      @thread-created="handleThreadCreated"
    />

    <!-- Main Chat Area -->
    <div class="chat-main">
      <div class="chatbot-header">
        <button class="sidebar-toggle" @click="sidebarOpen = !sidebarOpen" aria-label="Toggle sidebar">
          &#9776;
        </button>
        <div>
          <h1>AI Chatbot</h1>
          <p class="subtitle">Subscriber Exclusive Feature</p>
        </div>
      </div>

      <div class="chat-window">
        <div class="messages-container" ref="messagesContainer">
          <!-- Welcome message when no thread selected -->
          <div v-if="!currentThread" class="welcome-message">
            <div class="bot-avatar">AI</div>
            <div class="message-content">
              <p><strong>Welcome to the AI Chatbot!</strong></p>
              <p>Select a conversation from the sidebar or start a new one to begin chatting.</p>
            </div>
          </div>

          <!-- Empty thread message -->
          <div v-else-if="!currentThread.messages?.length" class="welcome-message">
            <div class="bot-avatar">AI</div>
            <div class="message-content">
              <p><strong>New Conversation</strong></p>
              <p>Ask me anything and I'll help you out!</p>
            </div>
          </div>

          <!-- Chat messages -->
          <div
            v-for="message in currentThread?.messages || []"
            :key="message.message_id"
            class="message"
            :class="message.role"
          >
            <div class="message-avatar">
              {{ message.role === 'user' ? 'You' : 'AI' }}
            </div>
            <div class="message-bubble">
              <div
                v-if="message.role === 'assistant'"
                class="message-text markdown-body"
                v-html="renderMarkdown(message.content)"
              />
              <p v-else class="message-text">{{ message.content }}</p>
              <span class="message-time">{{ formatTime(message.timestamp) }}</span>
            </div>
          </div>

          <!-- Streaming response -->
          <div v-if="currentThread && isStreaming" class="message assistant">
            <div class="message-avatar">AI</div>
            <div class="message-bubble">
              <!-- Show thinking animation when waiting for first token -->
              <div v-if="!streamingContent && activeTools.length === 0" class="thinking-indicator">
                <span class="thinking-text">Thinking</span>
                <span class="thinking-dots">
                  <span>.</span><span>.</span><span>.</span>
                </span>
              </div>
              <!-- Tool activity indicators -->
              <div v-if="activeTools.length > 0" class="tool-indicators">
                <div
                  v-for="tool in activeTools"
                  :key="tool.name"
                  class="tool-pill"
                  :class="tool.status"
                >
                  <span class="tool-icon">{{ tool.status === 'running' ? '&#9881;' : '&#10003;' }}</span>
                  <span class="tool-name">{{ tool.name }}</span>
                </div>
              </div>
              <!-- Show streaming content once it arrives -->
              <p v-if="streamingContent" class="message-text">{{ streamingContent }}<span class="cursor">|</span></p>
            </div>
          </div>

          <!-- Loading indicator (only when thread is selected and waiting for response) -->
          <div v-if="currentThread && isLoading && !isStreaming" class="message assistant">
            <div class="message-avatar">AI</div>
            <div class="message-bubble">
              <div class="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        </div>

        <!-- Input area -->
        <div class="input-container">
          <textarea
            v-model="userInput"
            @keydown="handleKeyDown"
            @compositionstart="handleCompositionStart"
            @compositionend="handleCompositionEnd"
            placeholder="Type your message here... (Press Enter to send, Shift+Enter for new line)"
            rows="3"
            :disabled="!currentThread || isLoading || isStreaming"
          ></textarea>
          <button
            @click="handleSubmit"
            :disabled="!userInput.trim() || !currentThread || isLoading || isStreaming || isComposing"
            class="send-button"
          >
            <span v-if="!isLoading && !isStreaming">Send</span>
            <span v-else>Sending...</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useChat } from '../composables/useChat'
import ThreadSidebar from '../components/ThreadSidebar.vue'
import { renderMarkdown } from '../utils/markdown'

const router = useRouter()
const { user } = useAuth()
const {
  currentThread,
  isLoading,
  isStreaming,
  streamingContent,
  activeTools,
  sendMessageStream,
} = useChat()

const userInput = ref('')
const messagesContainer = ref<HTMLElement | null>(null)
const sidebarOpen = ref(false)
const isComposing = ref(false)
let compositionJustEnded = false
let conversionKeydownHandled = false

// Check access on mount
onMounted(() => {
  if (user.value && !['subscriber', 'admin', 'root'].includes(user.value.role)) {
    router.push('/dashboard')
  }
})

// Scroll to bottom when messages change
watch(
  () => currentThread.value?.messages?.length,
  () => {
    nextTick(() => scrollToBottom())
  }
)

// Scroll when streaming
watch(streamingContent, () => {
  nextTick(() => scrollToBottom())
})

function handleCompositionStart() {
  isComposing.value = true
  compositionJustEnded = false
  conversionKeydownHandled = false
}

function handleCompositionEnd() {
  // Immediately clear the button-disabled state; composition is done.
  isComposing.value = false
  // If the conversion keydown already fired before this event
  // (keydown-before-compositionend order), the next Enter should send freely.
  // Otherwise set the one-shot flag to block the upcoming conversion keydown
  // (compositionend-before-keydown order, common in Chrome).
  if (!conversionKeydownHandled) {
    compositionJustEnded = true
  }
  conversionKeydownHandled = false
}

function handleKeyDown(event: KeyboardEvent) {
  if (event.key !== 'Enter') return

  // Shift+Enter → insert line break (let the browser's default run)
  if (event.shiftKey) return

  // Other modifier combos (Ctrl/Alt/Meta+Enter) → ignore
  if (event.ctrlKey || event.altKey || event.metaKey) return

  // Prevent the textarea's default newline for all plain Enter presses.
  // Called before IME checks so a blocked IME-Enter never inserts a newline.
  event.preventDefault()

  // event.isComposing / keyCode 229: keydown fired during active composition
  // compositionJustEnded: compositionend already fired before this keydown
  if (event.isComposing || event.keyCode === 229 || compositionJustEnded) {
    compositionJustEnded = false
    conversionKeydownHandled = true
    return
  }

  handleSubmit()
}

async function handleSubmit() {
  if (!userInput.value.trim() || !currentThread.value || isLoading.value || isStreaming.value) {
    return
  }

  const content = userInput.value.trim()
  userInput.value = ''

  await sendMessageStream(currentThread.value.thread_id, content)
  scrollToBottom()
}

function handleThreadSelected(_threadId: string) {
  sidebarOpen.value = false
  nextTick(() => scrollToBottom())
}

function handleThreadCreated(_threadId: string) {
  sidebarOpen.value = false
  userInput.value = ''
}

function formatTime(timestamp: string): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
</script>

<style scoped>
.chatbot-container {
  display: flex;
  height: calc(100vh - 64px);
  background: #f1f5f9;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  min-width: 0;
}

.chatbot-header {
  text-align: center;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
}

.sidebar-toggle {
  display: none;
  background: none;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 1.25rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  color: #475569;
}

.sidebar-toggle:hover {
  background: #f1f5f9;
}

.sidebar-backdrop {
  display: none;
}

.chatbot-header h1 {
  font-size: 1.5rem;
  color: #1e293b;
  margin: 0 0 0.25rem 0;
}

.subtitle {
  color: #3b82f6;
  font-weight: 600;
  font-size: 0.875rem;
  margin: 0;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.welcome-message {
  display: flex;
  gap: 1rem;
  padding: 1.5rem;
  background: #f0f9ff;
  border-radius: 12px;
  border-left: 4px solid #3b82f6;
}

.welcome-message .bot-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.welcome-message .message-content {
  flex: 1;
}

.welcome-message p {
  margin: 0 0 0.5rem 0;
}

.welcome-message p:last-child {
  margin: 0;
  color: #64748b;
}

.message {
  display: flex;
  gap: 0.75rem;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  flex-shrink: 0;
}

.message.user .message-avatar {
  background: #3b82f6;
  color: white;
}

.message.assistant .message-avatar {
  background: #e2e8f0;
  color: #475569;
}

.message-bubble {
  max-width: 70%;
  padding: 0.75rem 1rem;
  border-radius: 12px;
}

.message.user .message-bubble {
  background: #3b82f6;
  color: white;
  border-radius: 12px 12px 0 12px;
}

.message.assistant .message-bubble {
  background: #f1f5f9;
  color: #1e293b;
  border-radius: 12px 12px 12px 0;
}

.message-text {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.5;
}

.message-time {
  font-size: 0.625rem;
  opacity: 0.7;
  margin-top: 0.5rem;
  display: block;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 0.25rem 0;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #64748b;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    opacity: 0.3;
    transform: translateY(0);
  }
  30% {
    opacity: 1;
    transform: translateY(-4px);
  }
}

.thinking-indicator {
  display: flex;
  align-items: center;
  gap: 2px;
  color: #64748b;
  font-size: 0.875rem;
}

.thinking-text {
  animation: pulse 2s ease-in-out infinite;
}

.thinking-dots span {
  animation: dotBounce 1.4s ease-in-out infinite;
  display: inline-block;
}

.thinking-dots span:nth-child(1) {
  animation-delay: 0s;
}

.thinking-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.thinking-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

@keyframes dotBounce {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  30% {
    transform: translateY(-3px);
    opacity: 1;
  }
}

.tool-indicators {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.tool-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.25rem 0.75rem;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 500;
  background: #e0e7ff;
  color: #4338ca;
}

.tool-pill.running {
  animation: toolPulse 1.5s ease-in-out infinite;
}

.tool-pill.completed {
  background: #d1fae5;
  color: #065f46;
}

.tool-icon {
  font-size: 0.8rem;
}

.tool-pill.running .tool-icon {
  animation: spin 1.5s linear infinite;
}

@keyframes toolPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.cursor {
  animation: blink 1s step-end infinite;
  color: #3b82f6;
  font-weight: bold;
}

@keyframes blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}

.input-container {
  border-top: 1px solid #e2e8f0;
  padding: 1rem;
  display: flex;
  gap: 0.75rem;
  background: #f9fafb;
}

.input-container textarea {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-family: inherit;
  font-size: 0.875rem;
  resize: none;
  transition: border-color 0.2s;
}

.input-container textarea:focus {
  outline: none;
  border-color: #3b82f6;
}

.input-container textarea:disabled {
  background: #f1f5f9;
  cursor: not-allowed;
}

.send-button {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
  align-self: flex-end;
}

.send-button:hover:not(:disabled) {
  background: #2563eb;
}

.send-button:disabled {
  background: #94a3b8;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .sidebar-toggle {
    display: block;
  }

  .sidebar-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 10;
  }

  .chat-main {
    padding: 0.75rem;
  }

  .chatbot-header h1 {
    font-size: 1.25rem;
  }

  .message-bubble {
    max-width: 85%;
  }

  .input-container {
    flex-direction: column;
  }

  .send-button {
    width: 100%;
  }
}

</style>
