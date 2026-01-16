<template>
  <div class="chatbot-container">
    <div class="chatbot-header">
      <h1>AI Chatbot</h1>
      <p class="subtitle">Subscriber Exclusive Feature</p>
    </div>

    <div class="chat-window">
      <div class="messages-container" ref="messagesContainer">
        <!-- Welcome message -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="bot-avatar">🤖</div>
          <div class="message-content">
            <p><strong>Welcome to the AI Chatbot!</strong></p>
            <p>This is a subscriber-exclusive feature. Ask me anything and I'll help you out!</p>
          </div>
        </div>

        <!-- Chat messages -->
        <div
          v-for="(message, index) in messages"
          :key="index"
          class="message"
          :class="message.role"
        >
          <div class="message-avatar">
            {{ message.role === 'user' ? '👤' : '🤖' }}
          </div>
          <div class="message-content">
            <p>{{ message.content }}</p>
            <span class="message-time">{{ formatTime(message.timestamp) }}</span>
          </div>
        </div>

        <!-- Loading indicator -->
        <div v-if="isLoading" class="message bot">
          <div class="message-avatar">🤖</div>
          <div class="message-content">
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
          @keydown.enter.prevent="handleSubmit"
          placeholder="Type your message here... (Press Enter to send)"
          rows="3"
          :disabled="isLoading"
        ></textarea>
        <button
          @click="handleSubmit"
          :disabled="!userInput.trim() || isLoading"
          class="send-button"
        >
          <span v-if="!isLoading">Send</span>
          <span v-else>Sending...</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

interface Message {
  role: 'user' | 'bot'
  content: string
  timestamp: Date
}

const router = useRouter()
const { user } = useAuth()
const messages = ref<Message[]>([])
const userInput = ref('')
const isLoading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

// Check access on mount
onMounted(() => {
  if (user.value && !['subscriber', 'admin', 'root'].includes(user.value.role)) {
    router.push('/dashboard')
  }
})

async function handleSubmit() {
  if (!userInput.value.trim() || isLoading.value) return

  const userMessage: Message = {
    role: 'user',
    content: userInput.value.trim(),
    timestamp: new Date(),
  }

  messages.value.push(userMessage)
  userInput.value = ''
  isLoading.value = true

  // Scroll to bottom
  await nextTick()
  scrollToBottom()

  // Simulate AI response (replace with actual API call)
  setTimeout(() => {
    const botMessage: Message = {
      role: 'bot',
      content: generateResponse(userMessage.content),
      timestamp: new Date(),
    }
    messages.value.push(botMessage)
    isLoading.value = false
    nextTick(() => scrollToBottom())
  }, 1500)
}

function generateResponse(input: string): string {
  // This is a mock response - replace with actual AI API call
  const responses = [
    "That's an interesting question! As a demo chatbot, I'm here to show you the subscriber-exclusive feature.",
    "I understand what you're asking. In a production environment, this would connect to a real AI service.",
    "Great question! This chatbot demonstrates how subscriber-only features work in this application.",
    "Thank you for your message. This is a placeholder response - integrate with your preferred AI service!",
  ]
  return responses[Math.floor(Math.random() * responses.length)]
}

function formatTime(date: Date): string {
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
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
  height: calc(100vh - 4rem);
  display: flex;
  flex-direction: column;
}

.chatbot-header {
  text-align: center;
  margin-bottom: 2rem;
}

.chatbot-header h1 {
  font-size: 2rem;
  color: #333;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  color: #3b82f6;
  font-weight: 600;
  margin: 0;
}

.chat-window {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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
  font-size: 2.5rem;
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
  color: #666;
}

.message {
  display: flex;
  gap: 1rem;
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
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  flex-shrink: 0;
}

.message.user .message-avatar {
  background: #3b82f6;
}

.message.bot .message-avatar {
  background: #f3f4f6;
}

.message-content {
  flex: 1;
  max-width: 70%;
}

.message.user .message-content {
  background: #3b82f6;
  color: white;
  border-radius: 12px 12px 0 12px;
  padding: 1rem;
}

.message.bot .message-content {
  background: #f3f4f6;
  color: #333;
  border-radius: 12px 12px 12px 0;
  padding: 1rem;
}

.message-content p {
  margin: 0 0 0.5rem 0;
}

.message-content p:last-child {
  margin: 0;
}

.message-time {
  font-size: 0.75rem;
  opacity: 0.7;
  margin-top: 0.5rem;
  display: block;
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 0.5rem;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #666;
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
    transform: translateY(-10px);
  }
}

.input-container {
  border-top: 1px solid #e0e0e0;
  padding: 1.5rem;
  display: flex;
  gap: 1rem;
  background: #f9fafb;
}

.input-container textarea {
  flex: 1;
  padding: 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-family: inherit;
  font-size: 1rem;
  resize: none;
  transition: border-color 0.3s;
}

.input-container textarea:focus {
  outline: none;
  border-color: #3b82f6;
}

.input-container textarea:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.send-button {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 1rem 2rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  align-self: flex-end;
}

.send-button:hover:not(:disabled) {
  background: #2563eb;
  transform: translateY(-2px);
}

.send-button:disabled {
  background: #cbd5e1;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 768px) {
  .chatbot-container {
    padding: 1rem;
  }

  .message-content {
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
