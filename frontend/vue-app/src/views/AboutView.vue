<script setup lang="ts">
</script>

<template>
  <div class="about-page">
    <h1>About This Application</h1>
    <section class="section">
      <h2>概要</h2>
      <p>
        このアプリケーションは、Go言語で構築されたバックエンドとVue.js 3で構築されたフロントエンドで構成されています。
      </p>
      <p>主な目的は、Cookie/Sessionの動作をテストし、Google認証の実装を検証することです。</p>
    </section>

    <section class="section">
      <h2>技術スタック</h2>
      <ul>
        <li><strong>Frontend:</strong> Vue 3 + TypeScript + Vite</li>
        <li><strong>Backend:</strong> Go + Echo Framework</li>
        <li><strong>Database:</strong> DynamoDB Local</li>
        <li><strong>Container:</strong> Docker + Docker Compose</li>
      </ul>
    </section>

    <section class="section">
      <h2>機能</h2>
      <ul>
        <li>Set-Cookieヘッダーのテスト</li>
        <li>Cookie送信の確認</li>
        <li>セッション管理のテスト</li>
        <li>Google OAuth 2.0認証（実装予定）</li>
      </ul>
    </section>
  </div>
</template>

<style scoped>
.about-page {
  max-width: 1200px;
  margin: 0 auto;
}

.about-page h1 {
  margin-bottom: 2rem;
}

.about-page ul {
  list-style-position: inside;
  margin-left: 1rem;
}

.about-page li {
  margin-bottom: 0.5rem;
}
</style>
