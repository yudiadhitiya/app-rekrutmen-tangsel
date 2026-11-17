<script setup>
import { nextTick, onMounted, ref } from 'vue'
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
const errors = ref({})

// Password visibility.
const loginPasswordVisible = ref(false)
const registerPasswordVisible = ref(false)

// Element refs for focus management.
const usernameInput = ref(null)
const passwordInput = ref(null)

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

function clearFieldError(field) {
  if (!errors.value[field]) return
  const next = { ...errors.value }
  delete next[field]
  errors.value = next
}

function focusField(field) {
  nextTick(() => {
    const target = field === 'password' ? passwordInput.value : usernameInput.value
    if (target) target.focus()
  })
}

async function handleSubmit() {
  errorMessage.value = ''
  const next = {}
  if (!username.value.trim()) next.username = 'Username wajib diisi.'
  if (!password.value) next.password = 'Kata sandi wajib diisi.'
  errors.value = next

  const firstInvalid = Object.keys(next)[0]
  if (firstInvalid) {
    focusField(firstInvalid)
    return
  }

  submitting.value = true
  try {
    const result = await login(username.value.trim(), password.value)
    if (result && result.migrated) {
      toast.success('Akun berhasil dimigrasikan. Username dan kata sandi yang sama tetap berlaku.')
    }
    // The dashboard is the landing page of every account. The `redirect`
    // query is only kept when it points at a page the account may open.
    const target = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
    await router.replace(target)
  } catch (err) {
    errorMessage.value = err.message || 'Gagal masuk. Silakan coba lagi.'
    passwordInput.value?.focus()
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
  errors.value = {}
  loadRegistrationOptions()
}

function showLogin() {
  mode.value = 'login'
  errorMessage.value = ''
  errors.value = {}
}

async function handleRegister() {
  errorMessage.value = ''
  const next = {}
  if (!username.value.trim()) next.username = 'Username wajib diisi.'
  if (!password.value) next.password = 'Kata sandi wajib diisi.'
  if (!nama.value.trim()) next.nama = 'Nama anggota wajib diisi.'
  if (!kecamatan.value) next.kecamatan = 'Kecamatan wajib dipilih.'
  if (!question1.value || !question2.value) next.question = 'Dua pertanyaan keamanan wajib dipilih.'
  else if (question1.value === question2.value) next.question = 'Pilih dua pertanyaan keamanan yang berbeda.'
  if (!answer1.value.trim() || !answer2.value.trim()) next.answer = 'Dua jawaban keamanan wajib diisi.'
  errors.value = next

  const firstInvalid = Object.keys(next)[0]
  if (firstInvalid) {
    focusField(firstInvalid === 'password' ? 'password' : 'username')
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
  <div class="auth">
    <aside class="auth__aside" aria-hidden="true">
      <div class="auth__brand">
        <span class="auth__logo">RT</span>
        <span class="auth__brand-name">Rekrutmen Tangsel</span>
      </div>

      <div class="auth__aside-body">
        <h2 class="auth__headline">Kelola data rekrutmen dengan satu pintu.</h2>
        <p class="auth__lead">
          Aplikasi terpadu untuk officer, verifikator, dan pelaporan recruitment worker
          se-Kota Tangerang Selatan.
        </p>

        <ul class="auth__benefits">
          <li class="auth__benefit">
            <span class="auth__benefit-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
                <path d="M16 3.13a4 4 0 0 1 0 7.75" />
              </svg>
            </span>
            <span>
              <strong>Wilayah kerja per kecamatan</strong>
              <span>Pembagian territory dan borough yang jelas antar anggota.</span>
            </span>
          </li>
          <li class="auth__benefit">
            <span class="auth__benefit-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
                <path d="m9 12 2 2 4-4" />
              </svg>
            </span>
            <span>
              <strong>Data terlindungi</strong>
              <span>Akses berbasis peran dengan riwayat aktivitas tercatat.</span>
            </span>
          </li>
          <li class="auth__benefit">
            <span class="auth__benefit-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M3 3v18h18" />
                <path d="m7 14 4-4 3 3 5-6" />
              </svg>
            </span>
            <span>
              <strong>Rekapitulasi cepat</strong>
              <span>Dashboard, ekspor data, dan impor massal dalam satu klik.</span>
            </span>
          </li>
        </ul>
      </div>

      <p class="auth__foot">© {{ new Date().getFullYear() }} Pemerintah Kota Tangerang Selatan</p>
    </aside>

    <main class="auth__main">
      <div class="auth__card">
        <div class="auth__tabs" role="tablist" aria-label="Pilih formulir">
          <button
            class="auth__tab"
            :class="{ 'auth__tab--active': mode === 'login' }"
            type="button"
            role="tab"
            :aria-selected="mode === 'login'"
            @click="showLogin"
          >
            Masuk
          </button>
          <button
            class="auth__tab"
            :class="{ 'auth__tab--active': mode === 'register' }"
            type="button"
            role="tab"
            :aria-selected="mode === 'register'"
            @click="showRegister"
          >
            Daftar
          </button>
        </div>

        <template v-if="mode === 'login'">
          <h1 class="auth__title">Selamat datang kembali</h1>
          <p class="auth__subtitle">Silakan masuk untuk mengelola data rekrutmen.</p>

          <form class="auth__form" novalidate @submit.prevent="handleSubmit">
            <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

            <div class="field">
              <label class="field__label" for="login-username">Username</label>
              <div class="input-group">
                <span class="input-group__icon">
                  <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
                    <circle cx="12" cy="7" r="4" />
                  </svg>
                </span>
                <input
                  id="login-username"
                  ref="usernameInput"
                  v-model="username"
                  class="input"
                  :class="{ 'input--error': errors.username }"
                  type="text"
                  name="username"
                  autocomplete="username"
                  autocapitalize="none"
                  spellcheck="false"
                  placeholder="Masukkan username"
                  :aria-invalid="errors.username ? 'true' : 'false'"
                  :aria-describedby="errors.username ? 'login-username-error' : undefined"
                  :disabled="submitting"
                  @input="clearFieldError('username')"
                />
              </div>
              <p v-if="errors.username" id="login-username-error" class="field__error">
                {{ errors.username }}
              </p>
            </div>

            <div class="field">
              <label class="field__label" for="login-password">Kata Sandi</label>
              <div class="input-group">
                <span class="input-group__icon">
                  <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="4" y="10" width="16" height="11" rx="2" />
                    <path d="M8 10V7a4 4 0 0 1 8 0v3" />
                  </svg>
                </span>
                <input
                  id="login-password"
                  ref="passwordInput"
                  v-model="password"
                  class="input input--has-action"
                  :class="{ 'input--error': errors.password }"
                  :type="loginPasswordVisible ? 'text' : 'password'"
                  name="password"
                  autocomplete="current-password"
                  placeholder="Masukkan kata sandi"
                  :aria-invalid="errors.password ? 'true' : 'false'"
                  :aria-describedby="errors.password ? 'login-password-error' : undefined"
                  :disabled="submitting"
                  @input="clearFieldError('password')"
                />
                <button
                  class="input-group__action"
                  type="button"
                  :aria-pressed="loginPasswordVisible"
                  :aria-label="loginPasswordVisible ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'"
                  @click="loginPasswordVisible = !loginPasswordVisible"
                >
                  <svg v-if="!loginPasswordVisible" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                  <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M10.7 5.1A9.9 9.9 0 0 1 12 5c6.5 0 10 7 10 7a17.6 17.6 0 0 1-3.4 4.4" />
                    <path d="M6.6 6.6A17.6 17.6 0 0 0 2 12s3.5 7 10 7a9.8 9.8 0 0 0 5.4-1.6" />
                    <path d="M3 3l18 18" />
                  </svg>
                </button>
              </div>
              <p v-if="errors.password" id="login-password-error" class="field__error">
                {{ errors.password }}
              </p>
            </div>

            <button class="btn btn--primary btn--block auth__submit" type="submit" :disabled="submitting">
              <span v-if="submitting" class="spinner" aria-hidden="true"></span>
              {{ submitting ? 'Memproses...' : 'Masuk' }}
            </button>

            <p class="auth__switch">
              Belum punya akun?
              <button type="button" @click="showRegister">Daftar anggota</button>
            </p>

            <p class="auth__note">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <rect x="4" y="10" width="16" height="11" rx="2" />
                <path d="M8 10V7a4 4 0 0 1 8 0v3" />
              </svg>
              Gunakan akun resmi dari administrator. Aktivitas Anda tercatat pada log sistem.
            </p>
          </form>
        </template>

        <template v-else>
          <h1 class="auth__title">Daftar akun anggota</h1>
          <p class="auth__subtitle">
            Setelah registrasi berhasil, akun langsung aktif dan dapat digunakan untuk masuk.
          </p>

          <form class="auth__form" novalidate @submit.prevent="handleRegister">
            <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

            <div v-if="!registrationAvailable" class="alert alert--warning" role="alert">
              Pendaftaran anggota belum tersedia. Silakan hubungi administrator.
            </div>

            <div class="field">
              <label class="field__label" for="reg-username">Username</label>
              <div class="input-group">
                <span class="input-group__icon">
                  <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
                    <circle cx="12" cy="7" r="4" />
                  </svg>
                </span>
                <input
                  id="reg-username"
                  v-model="username"
                  class="input"
                  :class="{ 'input--error': errors.username }"
                  type="text"
                  name="username"
                  autocomplete="username"
                  autocapitalize="none"
                  spellcheck="false"
                  placeholder="Minimal 4 karakter"
                  :aria-invalid="errors.username ? 'true' : 'false'"
                  :disabled="submitting || !registrationAvailable"
                  @input="clearFieldError('username')"
                />
              </div>
              <p v-if="errors.username" class="field__error">{{ errors.username }}</p>
            </div>

            <div class="field">
              <label class="field__label" for="reg-password">Kata Sandi</label>
              <div class="input-group">
                <span class="input-group__icon">
                  <svg width="17" height="17" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="4" y="10" width="16" height="11" rx="2" />
                    <path d="M8 10V7a4 4 0 0 1 8 0v3" />
                  </svg>
                </span>
                <input
                  id="reg-password"
                  v-model="password"
                  class="input input--has-action"
                  :class="{ 'input--error': errors.password }"
                  :type="registerPasswordVisible ? 'text' : 'password'"
                  name="password"
                  autocomplete="new-password"
                  placeholder="Minimal 6 karakter"
                  :aria-invalid="errors.password ? 'true' : 'false'"
                  :disabled="submitting || !registrationAvailable"
                  @input="clearFieldError('password')"
                />
                <button
                  class="input-group__action"
                  type="button"
                  :aria-pressed="registerPasswordVisible"
                  :aria-label="registerPasswordVisible ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'"
                  :disabled="submitting || !registrationAvailable"
                  @click="registerPasswordVisible = !registerPasswordVisible"
                >
                  <svg v-if="!registerPasswordVisible" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                  <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M10.7 5.1A9.9 9.9 0 0 1 12 5c6.5 0 10 7 10 7a17.6 17.6 0 0 1-3.4 4.4" />
                    <path d="M6.6 6.6A17.6 17.6 0 0 0 2 12s3.5 7 10 7a9.8 9.8 0 0 0 5.4-1.6" />
                    <path d="M3 3l18 18" />
                  </svg>
                </button>
              </div>
              <p v-if="errors.password" class="field__error">{{ errors.password }}</p>
            </div>

            <div class="field">
              <label class="field__label" for="reg-nama">Nama Anggota</label>
              <input
                id="reg-nama"
                v-model="nama"
                class="input"
                :class="{ 'input--error': errors.nama }"
                type="text"
                name="nama"
                placeholder="Nama lengkap"
                :aria-invalid="errors.nama ? 'true' : 'false'"
                :disabled="submitting || !registrationAvailable"
                @input="clearFieldError('nama')"
              />
              <p v-if="errors.nama" class="field__error">{{ errors.nama }}</p>
            </div>

            <div class="field">
              <label class="field__label" for="reg-kecamatan">Kecamatan</label>
              <select
                id="reg-kecamatan"
                v-model="kecamatan"
                class="select"
                :class="{ 'select--error': errors.kecamatan }"
                :aria-invalid="errors.kecamatan ? 'true' : 'false'"
                :disabled="submitting || !registrationAvailable"
                @change="clearFieldError('kecamatan')"
              >
                <option value="">Pilih kecamatan</option>
                <option v-for="item in kecamatanOptions" :key="item" :value="item">{{ item }}</option>
              </select>
              <p v-if="errors.kecamatan" class="field__error">{{ errors.kecamatan }}</p>
            </div>

            <div class="field">
              <label class="field__label" for="reg-question1">Pertanyaan Keamanan 1</label>
              <select
                id="reg-question1"
                v-model="question1"
                class="select"
                :class="{ 'select--error': errors.question }"
                :aria-invalid="errors.question ? 'true' : 'false'"
                :disabled="submitting || !registrationAvailable"
                @change="clearFieldError('question')"
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
                @input="clearFieldError('answer')"
              />
            </div>

            <div class="field">
              <label class="field__label" for="reg-question2">Pertanyaan Keamanan 2</label>
              <select
                id="reg-question2"
                v-model="question2"
                class="select"
                :class="{ 'select--error': errors.question }"
                :aria-invalid="errors.question ? 'true' : 'false'"
                :disabled="submitting || !registrationAvailable"
                @change="clearFieldError('question')"
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
                @input="clearFieldError('answer')"
              />
              <p v-if="errors.question" class="field__error">{{ errors.question }}</p>
              <p v-if="errors.answer" class="field__error">{{ errors.answer }}</p>
            </div>

            <button
              class="btn btn--primary btn--block auth__submit"
              type="submit"
              :disabled="submitting || !registrationAvailable"
            >
              <span v-if="submitting" class="spinner" aria-hidden="true"></span>
              {{ submitting ? 'Memproses...' : 'Daftar & Aktifkan Akun' }}
            </button>

            <p class="auth__switch">
              Sudah punya akun?
              <button type="button" @click="showLogin">Kembali ke login</button>
            </p>
          </form>
        </template>
      </div>
    </main>
  </div>
</template>
