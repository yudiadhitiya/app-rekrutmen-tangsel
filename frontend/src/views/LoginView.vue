<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login } from '../auth'

const route = useRoute()
const router = useRouter()

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const submitting = ref(false)

async function handleSubmit() {
  errorMessage.value = ''
  if (!username.value.trim() || !password.value) {
    errorMessage.value = 'Username dan kata sandi wajib diisi.'
    return
  }
  submitting.value = true
  try {
    await login(username.value.trim(), password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (err) {
    errorMessage.value = err.message || 'Gagal masuk. Silakan coba lagi.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-card__brand">
        <span class="app-bar__logo">RT</span>
        <h1 class="login-card__title">Rekrutmen Tangsel</h1>
      </div>
      <p class="login-card__subtitle">Silakan masuk untuk mengelola data rekrutmen.</p>

      <form class="login-card__form" @submit.prevent="handleSubmit">
        <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

        <div class="field">
          <label class="field__label" for="login-username">Username</label>
          <input
            id="login-username"
            v-model="username"
            class="input"
            type="text"
            name="username"
            autocomplete="username"
            placeholder="Masukkan username"
            :disabled="submitting"
          />
        </div>

        <div class="field">
          <label class="field__label" for="login-password">Kata Sandi</label>
          <input
            id="login-password"
            v-model="password"
            class="input"
            type="password"
            name="password"
            autocomplete="current-password"
            placeholder="Masukkan kata sandi"
            :disabled="submitting"
          />
        </div>

        <button class="btn btn--primary btn--block" type="submit" :disabled="submitting">
          <span v-if="submitting" class="spinner" aria-hidden="true"></span>
          {{ submitting ? 'Memproses...' : 'Masuk' }}
        </button>
      </form>
    </div>
  </div>
</template>
