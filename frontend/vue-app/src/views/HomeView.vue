<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface CookieTestResult {
  sent: boolean
  cookies: string
  timestamp: string
}

// 環境変数からバックエンドURLを取得（デフォルトはlocalhost）
const defaultBackendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
const backendUrl = ref(defaultBackendUrl)
const setCookieResponse = ref('')
const cookieTestResult = ref<CookieTestResult | null>(null)
const currentCookies = ref('')
const testStatus = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const errorMessage = ref('')

// ページロード時に現在のCookieを取得
onMounted(() => {
  updateCurrentCookies()
  console.log('Backend URL:', backendUrl.value)
})

// 現在のCookieを更新
const updateCurrentCookies = () => {
  currentCookies.value = document.cookie || '(No cookies set)'
}

// Set-Cookieをテスト
const testSetCookie = async () => {
  testStatus.value = 'loading'
  errorMessage.value = ''
  setCookieResponse.value = ''

  try {
    const response = await fetch(`${backendUrl.value}/api/set-cookie`, {
      method: 'GET',
      credentials: 'include', // Cookieを含める
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    setCookieResponse.value = JSON.stringify(data, null, 2)

    // Set-Cookieヘッダーを表示
    const setCookieHeader = response.headers.get('Set-Cookie')
    if (setCookieHeader) {
      setCookieResponse.value += `\n\nSet-Cookie Header:\n${setCookieHeader}`
    }

    testStatus.value = 'success'
    updateCurrentCookies()
  } catch (error) {
    testStatus.value = 'error'
    errorMessage.value = error instanceof Error ? error.message : 'Unknown error occurred'
  }
}

// Cookieがバックエンドに送信されているか確認
const testCookieSending = async () => {
  testStatus.value = 'loading'
  errorMessage.value = ''
  cookieTestResult.value = null

  try {
    const response = await fetch(`${backendUrl.value}/api/check-cookie`, {
      method: 'GET',
      credentials: 'include', // Cookieを含める
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    cookieTestResult.value = {
      sent: data.cookieReceived || false,
      cookies: data.cookies || '(No cookies received)',
      timestamp: new Date().toISOString(),
    }

    testStatus.value = 'success'
  } catch (error) {
    testStatus.value = 'error'
    errorMessage.value = error instanceof Error ? error.message : 'Unknown error occurred'
  }
}

// Cookieをクリア
const clearCookies = () => {
  // すべてのCookieを削除
  document.cookie.split(';').forEach((cookie) => {
    const eqPos = cookie.indexOf('=')
    const name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie
    document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/'
  })
  updateCurrentCookies()
  cookieTestResult.value = null
  setCookieResponse.value = ''
  testStatus.value = 'idle'
}
</script>

<template>
  <div class="home-page">
    <h1>Session & Cookie Test</h1>
    <p class="text-muted">
      このページでは、バックエンドからのSet-Cookieとブラウザから送信されるCookieをテストできます。
    </p>

    <!-- Backend URL設定 -->
    <section class="section">
      <h2 class="section-title">Backend URL</h2>
      <div class="form-group">
        <label class="form-label">Backend API URL:</label>
        <input v-model="backendUrl" type="text" class="form-input" placeholder="http://localhost:8080" />
      </div>
    </section>

    <!-- 現在のCookie表示 -->
    <section class="section">
      <h2 class="section-title">Current Cookies</h2>
      <div class="info-card">
        <p><strong>ブラウザに保存されているCookie:</strong></p>
        <div class="code-block">{{ currentCookies }}</div>
        <div class="button-group mt-2">
          <button class="btn btn-info" @click="updateCurrentCookies">Refresh</button>
          <button class="btn btn-danger" @click="clearCookies">Clear All Cookies</button>
        </div>
      </div>
    </section>

    <!-- Set-Cookieテスト -->
    <section class="section">
      <h2 class="section-title">1. Test Set-Cookie from Backend</h2>
      <p>
        バックエンドの <code>/api/set-cookie</code>
        エンドポイントを呼び出して、Set-CookieヘッダーでCookieを設定します。
      </p>
      <div class="button-group mt-2">
        <button class="btn btn-primary" :disabled="testStatus === 'loading'" @click="testSetCookie">
          {{ testStatus === 'loading' ? 'Testing...' : 'Test Set-Cookie' }}
        </button>
      </div>

      <div v-if="setCookieResponse" class="info-card mt-2">
        <p>
          <span class="status-indicator status-success"></span>
          <strong>Response:</strong>
        </p>
        <div class="code-block">{{ setCookieResponse }}</div>
      </div>

      <div v-if="testStatus === 'error' && errorMessage" class="info-card mt-2">
        <p>
          <span class="status-indicator status-error"></span>
          <strong class="text-danger">Error:</strong>
        </p>
        <div class="code-block text-danger">{{ errorMessage }}</div>
      </div>
    </section>

    <!-- Cookie送信テスト -->
    <section class="section">
      <h2 class="section-title">2. Test Cookie Sending to Backend</h2>
      <p>
        バックエンドの <code>/api/check-cookie</code>
        エンドポイントを呼び出して、Cookieが正しく送信されているか確認します。
      </p>
      <div class="button-group mt-2">
        <button class="btn btn-primary" :disabled="testStatus === 'loading'" @click="testCookieSending">
          {{ testStatus === 'loading' ? 'Testing...' : 'Check Cookie Sending' }}
        </button>
      </div>

      <div v-if="cookieTestResult" class="info-card mt-2">
        <p>
          <span
            class="status-indicator"
            :class="cookieTestResult.sent ? 'status-success' : 'status-error'"
          ></span>
          <strong>Result:</strong>
        </p>
        <div class="code-block">
          Cookie Sent: {{ cookieTestResult.sent ? 'Yes ✓' : 'No ✗' }}
          Received Cookies: {{ cookieTestResult.cookies }}
          Timestamp: {{ cookieTestResult.timestamp }}
        </div>
      </div>

      <div v-if="testStatus === 'error' && errorMessage" class="info-card mt-2">
        <p>
          <span class="status-indicator status-error"></span>
          <strong class="text-danger">Error:</strong>
        </p>
        <div class="code-block text-danger">{{ errorMessage }}</div>
      </div>
    </section>

    <!-- 手順説明 -->
    <section class="section">
      <h2 class="section-title">Test Instructions</h2>
      <ol>
        <li>まず「Test Set-Cookie」ボタンをクリックして、バックエンドからCookieを設定します。</li>
        <li>「Current Cookies」セクションで、Cookieが正しく設定されたことを確認します。</li>
        <li>
          「Check Cookie Sending」ボタンをクリックして、設定されたCookieがバックエンドに送信されているか確認します。
        </li>
        <li>テストが終わったら「Clear All Cookies」でCookieをクリアできます。</li>
      </ol>
    </section>
  </div>
</template>
