<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../api'
import { login } from '../auth'
import { toast } from '../toast'

const route = useRoute()
const router = useRouter()

const mode = ref('login')

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const submitting = ref(false)

// Registration form.
const nama = ref('')
const kecamatan = ref('')
const question1 = ref('')
const answer1 = ref('')
const question2 = ref('')
const answer2 = ref('')
const questions = ref([])
const kecamatanOptions = ref([])
const registrationAvailable = ref(true)

async function handleSubmit() {
  errorMessage.value = ''
  if (!username.value.trim() || !password.value) {
    errorMessage.value = 'Username dan kata sandi wajib diisi.'
    return
  }
  submitting.value = true
  try {
    const result = await login(username.value.trim(), password.value)
    if (result && result.migrated) {
      toast.success('Akun berhasil dimigrasikan. Username dan kata sandi yang sama tetap berlaku.')
    }
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (err) {
    errorMessage.value = err.message || 'Gagal masuk. Silakan coba lagi.'
  } finally {
    submitting.value = false
  }
}

async function loadRegistrationOptions() {
  if (questions.value.length) return
  try {
    const data = await api.register.options()
    questions.value = data.questions || []
    kecamatanOptions.value = data.kecamatan || []
    registrationAvailable.value = data.available !== false
  } catch (err) {
    registrationAvailable.value = false
    errorMessage.value = err.message || 'Formulir pendaftaran tidak dapat dimuat.'
  }
}

function showRegister() {
  mode.value = 'register'
  errorMessage.value = ''
  loadRegistrationOptions()
}

function showLogin() {
  mode.value = 'login'
  errorMessage.value = ''
}

async function handleRegister() {
  errorMessage.value = ''
  if (!username.value.trim()) {
    errorMessage.value = 'Username wajib diisi.'
    return
  }
  if (!password.value) {
    errorMessage.value = 'Kata sandi wajib diisi.'
    return
  }
  if (!nama.value.trim()) {
    errorMessage.value = 'Nama anggota wajib diisi.'
    return
  }
  if (!kecamatan.value) {
    errorMessage.value = 'Kecamatan wajib dipilih.'
    return
  }
  if (!question1.value || !question2.value) {
    errorMessage.value = 'Dua pertanyaan keamanan wajib dipilih.'
    return
  }
  if (question1.value === question2.value) {
    errorMessage.value = 'Pilih dua pertanyaan keamanan yang berbeda.'
    return
  }
  if (!answer1.value.trim() || !answer2.value.trim()) {
    errorMessage.value = 'Dua jawaban keamanan wajib diisi.'
    return
  }

  submitting.value = true
  try {
    await api.register.submit({
      username: username.value.trim(),
      password: password.value,
      nama: nama.value.trim(),
      kecamatan: kecamatan.value,
      question1: question1.value,
      answer1: answer1.value.trim(),
      question2: question2.value,
      answer2: answer2.value.trim()
    })
    toast.success('Registrasi berhasil. Akun langsung aktif, silakan masuk.')
    password.value = ''
    answer1.value = ''
    answer2.value = ''
    showLogin()
  } catch (err) {
    errorMessage.value = err.message || 'Gagal mendaftar. Silakan coba lagi.'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (route.query.register) showRegister()
})
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-card__brand">
        <span class="app-bar__logo">RT</span>
        <h1 class="login-card__title">Rekrutmen Tangsel</h1>
      </div>

      <template v-if="mode === 'login'">
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

          <button class="btn btn--ghost btn--block" type="button" @click="showRegister">
            Belum punya akun? Daftar anggota
          </button>
        </form>
      </template>

      <template v-else>
        <p class="login-card__subtitle">
          Daftarkan diri sebagai anggota. Setelah registrasi berhasil, akun langsung aktif dan dapat
          digunakan untuk masuk.
        </p>

        <form class="login-card__form" @submit.prevent="handleRegister">
          <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

          <div v-if="!registrationAvailable" class="alert alert--warning" role="alert">
            Pendaftaran anggota belum tersedia. Silakan hubungi administrator.
          </div>

          <div class="field">
            <label class="field__label" for="reg-username">Username</label>
            <input
              id="reg-username"
              v-model="username"
              class="input"
              type="text"
              name="username"
              autocomplete="username"
              placeholder="Minimal 4 karakter"
              :disabled="submitting || !registrationAvailable"
            />
          </div>

          <div class="field">
            <label class="field__label" for="reg-password">Kata Sandi</label>
            <input
              id="reg-password"
              v-model="password"
              class="input"
              type="password"
              name="password"
              autocomplete="new-password"
              placeholder="Minimal 6 karakter"
              :disabled="submitting || !registrationAvailable"
            />
          </div>

          <div class="field">
            <label class="field__label" for="reg-nama">Nama Anggota</label>
            <input
              id="reg-nama"
              v-model="nama"
              class="input"
              type="text"
              name="nama"
              placeholder="Nama lengkap"
              :disabled="submitting || !registrationAvailable"
            />
          </div>

          <div class="field">
            <label class="field__label" for="reg-kecamatan">Kecamatan</label>
            <select
              id="reg-kecamatan"
              v-model="kecamatan"
              class="select"
              :disabled="submitting || !registrationAvailable"
            >
              <option value="">Pilih kecamatan</option>
              <option v-for="item in kecamatanOptions" :key="item" :value="item">{{ item }}</option>
            </select>
          </div>

          <div class="field">
            <label class="field__label" for="reg-question1">Pertanyaan Keamanan 1</label>
            <select
              id="reg-question1"
              v-model="question1"
              class="select"
              :disabled="submitting || !registrationAvailable"
            >
              <option value="">Pilih pertanyaan</option>
              <option v-for="item in questions" :key="item" :value="item">{{ item }}</option>
            </select>
          </div>

          <div class="field">
            <label class="field__label" for="reg-answer1">Jawaban 1</label>
            <input
              id="reg-answer1"
              v-model="answer1"
              class="input"
              type="text"
              autocomplete="off"
              placeholder="Jawaban pertanyaan keamanan"
              :disabled="submitting || !registrationAvailable"
            />
          </div>

          <div class="field">
            <label class="field__label" for="reg-question2">Pertanyaan Keamanan 2</label>
            <select
              id="reg-question2"
              v-model="question2"
              class="select"
              :disabled="submitting || !registrationAvailable"
            >
              <option value="">Pilih pertanyaan</option>
              <option v-for="item in questions" :key="item" :value="item">{{ item }}</option>
            </select>
          </div>

          <div class="field">
            <label class="field__label" for="reg-answer2">Jawaban 2</label>
            <input
              id="reg-answer2"
              v-model="answer2"
              class="input"
              type="text"
              autocomplete="off"
              placeholder="Jawaban pertanyaan keamanan"
              :disabled="submitting || !registrationAvailable"
            />
          </div>

          <button
            class="btn btn--primary btn--block"
            type="submit"
            :disabled="submitting || !registrationAvailable"
          >
            <span v-if="submitting" class="spinner" aria-hidden="true"></span>
            {{ submitting ? 'Memproses...' : 'Daftar & Aktifkan Akun' }}
          </button>

          <button class="btn btn--ghost btn--block" type="button" @click="showLogin">
            Kembali ke login
          </button>
        </form>
      </template>
    </div>
  </div>
</template>
