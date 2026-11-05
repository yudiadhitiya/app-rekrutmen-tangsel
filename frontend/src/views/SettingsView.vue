<script setup>
import { computed, onMounted, ref } from 'vue'
import { api } from '../api'
import { PROVIDERS, providerLabel } from '../constants'
import { toast } from '../toast'

const settings = ref({})
const provider = ref('')
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')

const PROVIDER_DESCRIPTIONS = {
  local: 'Penyimpanan lokal pada server aplikasi.',
  supabase: 'PostgreSQL terkelola di Supabase.',
  google_sheets: 'Google Sheets sebagai sumber data.'
}

const currentProvider = computed(() => provider.value || settings.value.provider || '')

async function loadSettings() {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api.settings.get()
    settings.value = data || {}
    provider.value = data && data.provider ? data.provider : ''
  } catch (err) {
    errorMessage.value = err.message || 'Gagal memuat pengaturan.'
  } finally {
    loading.value = false
  }
}

async function saveProvider() {
  if (!provider.value) {
    toast.error('Pilih provider terlebih dahulu.')
    return
  }
  if (provider.value === settings.value.provider) {
    toast.info('Provider yang dipilih sudah aktif.')
    return
  }
  saving.value = true
  try {
    const data = await api.settings.update(provider.value)
    settings.value = data && typeof data === 'object' ? { ...settings.value, ...data } : settings.value
    toast.success(`Provider berhasil diubah ke ${providerLabel(provider.value)}.`)
    await loadSettings()
  } catch (err) {
    toast.error(err.message || 'Gagal mengubah provider.')
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <section>
    <div class="page-head">
      <div>
        <h1 class="page-head__title">Pengaturan</h1>
        <p class="page-head__subtitle">Kelola sumber penyimpanan data aplikasi. Hanya admin yang dapat mengubah.</p>
      </div>
    </div>

    <div class="stack">
      <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Provider Penyimpanan</h2>
        </div>
        <div class="card__body">
          <div v-if="loading" class="loading-block">
            <span class="spinner spinner--dark" aria-hidden="true"></span>
            Memuat pengaturan...
          </div>

          <template v-else>
            <div class="setting-row">
              <span class="setting-row__label">Provider aktif saat ini</span>
              <span class="badge tone-info">{{ providerLabel(currentProvider) || 'Tidak diketahui' }}</span>
            </div>

            <div class="provider-grid" style="margin-top: 16px">
              <button
                v-for="item in PROVIDERS"
                :key="item.value"
                type="button"
                class="provider-option"
                :class="{ 'provider-option--active': provider === item.value }"
                :disabled="saving"
                @click="provider = item.value"
              >
                <div class="provider-option__name">{{ item.label }}</div>
                <div class="provider-option__desc">
                  {{ PROVIDER_DESCRIPTIONS[item.value] || '' }}
                </div>
              </button>
            </div>

            <div class="alert alert--info" style="margin-top: 16px">
              Mengganti provider tidak memindahkan data secara otomatis. Pastikan data sudah diekspor sebelum
              perpindahan digunakan.
            </div>

            <div class="form-actions">
              <button class="btn btn--secondary" type="button" :disabled="saving" @click="loadSettings">
                Muat ulang
              </button>
              <button class="btn btn--primary" type="button" :disabled="saving" @click="saveProvider">
                <span v-if="saving" class="spinner" aria-hidden="true"></span>
                {{ saving ? 'Menyimpan...' : 'Simpan Provider' }}
              </button>
            </div>
          </template>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Informasi Provider</h2>
        </div>
        <div class="card__body">
          <div v-for="(value, key) in settings" :key="key" class="setting-row">
            <span class="setting-row__label">{{ key }}</span>
            <span class="setting-row__value">{{ typeof value === 'object' ? JSON.stringify(value) : value }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
