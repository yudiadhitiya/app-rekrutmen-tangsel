<script setup>
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api, ApiError } from '../api'
import { METODE, STATUS, TINGKAT_PENERIMAAN, UPL_OPTIONS, normalizeMetode } from '../constants'
import { toast } from '../toast'
import { confirmDialog } from '../confirm'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  nama: '',
  no_wa: '',
  alamat: '',
  rt: '',
  rw: '',
  kecamatan: '',
  kelurahan: '',
  aktivitas: '',
  tingkat_penerimaan: '',
  metode: '',
  status: '',
  keterangan: '',
  upl: ''
})

const errors = reactive({})
const formError = ref('')
const loading = ref(false)
const saving = ref(false)

const wilayah = ref({ kecamatan: [] })
const aktivitasMaster = ref([])

const selectedKecamatan = computed(() =>
  wilayah.value.kecamatan.find((item) => item.id === form.kecamatan)
)

const kelurahanOptions = computed(() => {
  if (!selectedKecamatan.value) return []
  return selectedKecamatan.value.kelurahan || []
})

const aktivitasOptions = computed(() => aktivitasMaster.value.filter((item) => item.aktif !== false))

// How many fields still need attention.
const errorCount = computed(() => Object.keys(errors).length)

// The element id of every field, so the first failing field can take the
// focus after a failed submit.
const FIELD_IDS = {
  nama: 'f-nama',
  no_wa: 'f-no-wa',
  alamat: 'f-alamat',
  rt: 'f-rt',
  rw: 'f-rw',
  kecamatan: 'f-kecamatan',
  kelurahan: 'f-kelurahan',
  aktivitas: 'f-aktivitas',
  tingkat_penerimaan: 'f-tingkat',
  metode: 'f-metode',
  status: 'f-status',
  keterangan: 'f-keterangan',
  upl: 'f-upl'
}

function setError(field, message) {
  if (message) errors[field] = message
  else delete errors[field]
}

function clearErrors() {
  Object.keys(errors).forEach((key) => delete errors[key])
  formError.value = ''
}

// Moves the focus to the first field that failed the validation.
function focusFirstError() {
  const field = Object.keys(errors)[0]
  if (!field) return
  nextTick(() => {
    const element = document.getElementById(FIELD_IDS[field] || '')
    if (!element) return
    element.focus()
    element.scrollIntoView({ block: 'center', behavior: 'smooth' })
  })
}

// Changing kecamatan resets the dependent kelurahan selection.
watch(
  () => form.kecamatan,
  (value, previous) => {
    if (!value || value === previous) return
    const allowed = kelurahanOptions.value || []
    if (form.kelurahan && !allowed.some((item) => item.id === form.kelurahan)) {
      form.kelurahan = ''
    }
    setError('kelurahan', '')
  }
)

function validate() {
  clearErrors()
  let valid = true

  const fail = (field, message) => {
    setError(field, message)
    valid = false
  }

  if (!form.nama.trim()) fail('nama', 'Nama wajib diisi.')

  if (!form.no_wa.trim()) {
    fail('no_wa', 'No. WhatsApp wajib diisi.')
  } else if (!/^\d+$/.test(form.no_wa.trim())) {
    fail('no_wa', 'No. WhatsApp hanya boleh berisi angka 0-9.')
  }

  if (form.rt.trim() && !/^\d+$/.test(form.rt.trim())) fail('rt', 'RT hanya boleh berisi angka 0-9.')
  if (form.rw.trim() && !/^\d+$/.test(form.rw.trim())) fail('rw', 'RW hanya boleh berisi angka 0-9.')

  if (!form.kecamatan) fail('kecamatan', 'Kecamatan wajib dipilih.')

  if (form.kelurahan) {
    const allowed = kelurahanOptions.value.some((item) => item.id === form.kelurahan)
    if (!allowed) fail('kelurahan', 'Kelurahan harus berada di dalam kecamatan yang dipilih.')
  }

  if (!form.aktivitas) fail('aktivitas', 'Aktivitas wajib dipilih.')
  if (!form.tingkat_penerimaan) fail('tingkat_penerimaan', 'Tingkat penerimaan wajib dipilih.')
  if (!form.metode) fail('metode', 'Metode wajib dipilih.')
  if (!form.status) fail('status', 'Status wajib dipilih.')

  return valid
}

// The WhatsApp number is stored in the 62 format, so a local number that
// starts with 0 or 8 is rewritten when the field loses the focus.
function normalizeWa() {
  const digits = form.no_wa.replace(/\D+/g, '')
  if (!digits) {
    form.no_wa = ''
    return
  }
  if (digits.startsWith('0')) form.no_wa = `62${digits.slice(1)}`
  else if (digits.startsWith('8')) form.no_wa = `62${digits}`
  else form.no_wa = digits
}

function payload() {
  return {
    nama: form.nama.trim(),
    no_wa: form.no_wa.trim(),
    alamat: form.alamat.trim(),
    rt: form.rt.trim(),
    rw: form.rw.trim(),
    kecamatan: form.kecamatan,
    kelurahan: form.kelurahan,
    aktivitas: form.aktivitas,
    tingkat_penerimaan: form.tingkat_penerimaan,
    metode: form.metode,
    status: form.status,
    keterangan: form.keterangan.trim(),
    upl: form.upl || ''
  }
}

async function checkDuplicateWa() {
  const number = form.no_wa.trim()
  if (!number) return true
  try {
    const result = await api.rekrutmen.checkWa(number)
    if (!result || !result.exists) return true
    if (isEdit.value && result.id && result.id === route.params.id) return true
    const confirmed = await confirmDialog({
      title: 'Nomor WhatsApp sudah terdaftar',
      message:
        `Nomor WhatsApp ${number} sudah digunakan` +
        (result.nama ? ` oleh "${result.nama}"` : ' oleh data lain') +
        '. Tetap simpan data ini?',
      confirmText: 'Tetap simpan',
      cancelText: 'Batal',
      tone: 'primary'
    })
    if (!confirmed) setError('no_wa', 'Nomor WhatsApp sudah digunakan oleh data lain.')
    return confirmed
  } catch (err) {
    // A failed pre-check must not block saving; the server still validates uniqueness.
    return true
  }
}

async function handleSubmit() {
  formError.value = ''
  if (!validate()) {
    toast.error('Periksa kembali isian yang ditandai merah.')
    focusFirstError()
    return
  }

  saving.value = true
  try {
    const proceed = await checkDuplicateWa()
    if (!proceed) {
      focusFirstError()
      return
    }

    if (isEdit.value) {
      await api.rekrutmen.update(route.params.id, payload())
      toast.success('Data rekrutmen berhasil diperbarui.')
    } else {
      await api.rekrutmen.create(payload())
      toast.success('Data rekrutmen berhasil disimpan.')
    }
    router.push('/')
  } catch (err) {
    if (err instanceof ApiError && err.fields && Object.keys(err.fields).length) {
      Object.entries(err.fields).forEach(([field, message]) => setError(field, message))
      toast.error('Penyimpanan gagal. Periksa isian yang ditandai.')
      focusFirstError()
    } else if (err instanceof ApiError && err.status === 409) {
      setError('no_wa', 'No. WhatsApp sudah digunakan oleh data lain.')
      toast.error(err.message || 'No. WhatsApp sudah digunakan oleh data lain.')
      focusFirstError()
    } else {
      formError.value = err.message || 'Gagal menyimpan data rekrutmen.'
      toast.error(formError.value)
    }
  } finally {
    saving.value = false
  }
}

function cancel() {
  router.push('/')
}

async function loadMasters() {
  try {
    const data = await api.wilayah()
    wilayah.value = data && data.kecamatan ? data : { kecamatan: [] }
  } catch (err) {
    wilayah.value = { kecamatan: [] }
    formError.value = err.message || 'Gagal memuat data wilayah.'
  }

  try {
    const data = await api.aktivitas.list()
    aktivitasMaster.value = Array.isArray(data) ? data : []
  } catch (err) {
    aktivitasMaster.value = []
  }
}

async function loadRecord() {
  loading.value = true
  formError.value = ''
  try {
    const data = await api.rekrutmen.get(route.params.id)
    form.nama = data.nama || ''
    form.no_wa = data.no_wa || ''
    form.alamat = data.alamat || ''
    form.rt = data.rt || ''
    form.rw = data.rw || ''
    form.kecamatan = data.kecamatan || ''
    form.kelurahan = data.kelurahan || ''
    form.aktivitas = data.aktivitas || ''
    form.tingkat_penerimaan = data.tingkat_penerimaan || ''
    form.metode = normalizeMetode(data.metode || '')
    form.status = data.status || ''
    form.keterangan = data.keterangan || ''
    form.upl = data.upl || ''
  } catch (err) {
    formError.value = err.message || 'Gagal memuat data rekrutmen.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadMasters()
  if (isEdit.value) await loadRecord()
})
</script>

<template>
  <section>
    <div class="page-head">
      <div>
        <h1 class="page-head__title">{{ isEdit ? 'Ubah Data Rekrutmen' : 'Tambah Data Rekrutmen' }}</h1>
        <p class="page-head__subtitle">
          Lengkapi bagian di bawah ini. Field bertanda <span class="field__required">*</span> wajib diisi.
        </p>
      </div>
      <RouterLink to="/" class="btn btn--secondary">Kembali ke Daftar</RouterLink>
    </div>

    <div v-if="loading" class="card">
      <div class="loading-block">
        <span class="spinner spinner--dark" aria-hidden="true"></span>
        Memuat data...
      </div>
    </div>

    <form v-else class="stack" novalidate @submit.prevent="handleSubmit">
      <div v-if="formError" class="alert alert--error" role="alert">{{ formError }}</div>

      <div v-if="errorCount" class="alert alert--error form-summary" role="alert">
        <strong>{{ errorCount }} isian perlu diperbaiki.</strong>
        <span>Periksa kembali isian yang ditandai merah di bawah, lalu simpan ulang.</span>
      </div>

      <!-- Identitas -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__icon" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2" />
              <circle cx="12" cy="7" r="4" />
            </svg>
          </span>
          <span class="form-section__titles">
            <span class="form-section__title">Identitas</span>
            <span class="form-section__hint">Nama dan kontak calon pelapor</span>
          </span>
        </legend>
        <div class="form-section__body">
          <div class="field">
            <label class="field__label" for="f-nama">
              Nama <span class="field__required" aria-hidden="true">*</span>
            </label>
            <input
              id="f-nama"
              v-model="form.nama"
              class="input"
              :class="{ 'input--error': errors.nama }"
              type="text"
              autocomplete="name"
              placeholder="Nama lengkap"
              :aria-invalid="errors.nama ? 'true' : 'false'"
              :aria-describedby="errors.nama ? 'f-nama-error' : undefined"
              :disabled="saving"
              @input="setError('nama', '')"
            />
            <span v-if="errors.nama" id="f-nama-error" class="field__error">{{ errors.nama }}</span>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-no-wa">
                No. WhatsApp <span class="field__required" aria-hidden="true">*</span>
              </label>
              <input
                id="f-no-wa"
                v-model="form.no_wa"
                class="input"
                :class="{ 'input--error': errors.no_wa }"
                type="tel"
                inputmode="numeric"
                autocomplete="tel"
                placeholder="628123456789"
                :aria-invalid="errors.no_wa ? 'true' : 'false'"
                :aria-describedby="errors.no_wa ? 'f-no-wa-hint f-no-wa-error' : 'f-no-wa-hint'"
                :disabled="saving"
                @input="setError('no_wa', '')"
                @blur="normalizeWa"
              />
              <span id="f-no-wa-hint" class="field__hint">
                Hanya angka 0-9. Awalan 0 atau 8 otomatis diubah menjadi format 62.
              </span>
              <span v-if="errors.no_wa" id="f-no-wa-error" class="field__error">{{ errors.no_wa }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-alamat">Alamat</label>
              <textarea
                id="f-alamat"
                v-model="form.alamat"
                class="textarea"
                :class="{ 'textarea--error': errors.alamat }"
                rows="2"
                placeholder="Alamat lengkap"
                :aria-invalid="errors.alamat ? 'true' : 'false'"
                :aria-describedby="errors.alamat ? 'f-alamat-error' : undefined"
                :disabled="saving"
                @input="setError('alamat', '')"
              ></textarea>
              <span v-if="errors.alamat" id="f-alamat-error" class="field__error">{{ errors.alamat }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Wilayah -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__icon" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 21s7-5.5 7-11a7 7 0 1 0-14 0c0 5.5 7 11 7 11z" />
              <circle cx="12" cy="10" r="2.5" />
            </svg>
          </span>
          <span class="form-section__titles">
            <span class="form-section__title">Wilayah</span>
            <span class="form-section__hint">Wilayah penempatan perekrutmen</span>
          </span>
        </legend>
        <div class="form-section__body">
          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-rt">RT</label>
              <input
                id="f-rt"
                v-model="form.rt"
                class="input"
                :class="{ 'input--error': errors.rt }"
                type="text"
                inputmode="numeric"
                placeholder="003"
                :aria-invalid="errors.rt ? 'true' : 'false'"
                :aria-describedby="errors.rt ? 'f-rt-error' : undefined"
                :disabled="saving"
                @input="setError('rt', '')"
              />
              <span v-if="errors.rt" id="f-rt-error" class="field__error">{{ errors.rt }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-rw">RW</label>
              <input
                id="f-rw"
                v-model="form.rw"
                class="input"
                :class="{ 'input--error': errors.rw }"
                type="text"
                inputmode="numeric"
                placeholder="005"
                :aria-invalid="errors.rw ? 'true' : 'false'"
                :aria-describedby="errors.rw ? 'f-rw-error' : undefined"
                :disabled="saving"
                @input="setError('rw', '')"
              />
              <span v-if="errors.rw" id="f-rw-error" class="field__error">{{ errors.rw }}</span>
            </div>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-kecamatan">
                Kecamatan <span class="field__required" aria-hidden="true">*</span>
              </label>
              <select
                id="f-kecamatan"
                v-model="form.kecamatan"
                class="select"
                :class="{ 'select--error': errors.kecamatan }"
                :aria-invalid="errors.kecamatan ? 'true' : 'false'"
                :aria-describedby="errors.kecamatan ? 'f-kecamatan-error' : undefined"
                :disabled="saving"
                @change="setError('kecamatan', '')"
              >
                <option value="">-- Pilih Kecamatan --</option>
                <option v-for="item in wilayah.kecamatan" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
              <span v-if="errors.kecamatan" id="f-kecamatan-error" class="field__error">{{ errors.kecamatan }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-kelurahan">Kelurahan</label>
              <select
                id="f-kelurahan"
                v-model="form.kelurahan"
                class="select"
                :class="{ 'select--error': errors.kelurahan }"
                :aria-invalid="errors.kelurahan ? 'true' : 'false'"
                :aria-describedby="errors.kelurahan ? 'f-kelurahan-error' : undefined"
                :disabled="saving || !form.kecamatan"
              >
                <option value="">{{ form.kecamatan ? '-- Pilih Kelurahan --' : 'Pilih kecamatan terlebih dahulu' }}</option>
                <option v-for="item in kelurahanOptions" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
              <span v-if="errors.kelurahan" id="f-kelurahan-error" class="field__error">{{ errors.kelurahan }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Informasi Rekrutmen -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__icon" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M9 3h6v3H9z" />
              <path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2" />
              <path d="M9 12h6M9 16h4" />
            </svg>
          </span>
          <span class="form-section__titles">
            <span class="form-section__title">Informasi Rekrutmen</span>
            <span class="form-section__hint">Klasifikasi dan status perekrutmen</span>
          </span>
        </legend>
        <div class="form-section__body">
          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-aktivitas">
                Aktivitas <span class="field__required" aria-hidden="true">*</span>
              </label>
              <select
                id="f-aktivitas"
                v-model="form.aktivitas"
                class="select"
                :class="{ 'select--error': errors.aktivitas }"
                :aria-invalid="errors.aktivitas ? 'true' : 'false'"
                :aria-describedby="errors.aktivitas ? 'f-aktivitas-error' : undefined"
                :disabled="saving"
                @change="setError('aktivitas', '')"
              >
                <option value="">-- Pilih Aktivitas --</option>
                <option v-for="item in aktivitasOptions" :key="item.id" :value="item.nama">
                  {{ item.nama }}
                </option>
              </select>
              <span v-if="errors.aktivitas" id="f-aktivitas-error" class="field__error">{{ errors.aktivitas }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-tingkat">
                Tingkat Penerimaan <span class="field__required" aria-hidden="true">*</span>
              </label>
              <select
                id="f-tingkat"
                v-model="form.tingkat_penerimaan"
                class="select"
                :class="{ 'select--error': errors.tingkat_penerimaan }"
                :aria-invalid="errors.tingkat_penerimaan ? 'true' : 'false'"
                :aria-describedby="errors.tingkat_penerimaan ? 'f-tingkat-error' : 'f-tingkat-hint'"
                :disabled="saving"
                @change="setError('tingkat_penerimaan', '')"
              >
                <option value="">-- Pilih Tingkat Penerimaan --</option>
                <option v-for="item in TINGKAT_PENERIMAAN" :key="item" :value="item">{{ item }}</option>
              </select>
              <span id="f-tingkat-hint" class="field__hint">Tingkat keterbukaan wilayah terhadap kandidat baru.</span>
              <span v-if="errors.tingkat_penerimaan" id="f-tingkat-error" class="field__error">
                {{ errors.tingkat_penerimaan }}
              </span>
            </div>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-metode">
                Metode <span class="field__required" aria-hidden="true">*</span>
              </label>
              <select
                id="f-metode"
                v-model="form.metode"
                class="select"
                :class="{ 'select--error': errors.metode }"
                :aria-invalid="errors.metode ? 'true' : 'false'"
                :aria-describedby="errors.metode ? 'f-metode-error' : undefined"
                :disabled="saving"
                @change="setError('metode', '')"
              >
                <option value="">-- Pilih Metode --</option>
                <option v-for="item in METODE" :key="item" :value="item">{{ item }}</option>
              </select>
              <span v-if="errors.metode" id="f-metode-error" class="field__error">{{ errors.metode }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-status">
                Status <span class="field__required" aria-hidden="true">*</span>
              </label>
              <select
                id="f-status"
                v-model="form.status"
                class="select"
                :class="{ 'select--error': errors.status }"
                :aria-invalid="errors.status ? 'true' : 'false'"
                :aria-describedby="errors.status ? 'f-status-error' : undefined"
                :disabled="saving"
                @change="setError('status', '')"
              >
                <option value="">-- Pilih Status --</option>
                <option v-for="item in STATUS" :key="item" :value="item">{{ item }}</option>
              </select>
              <span v-if="errors.status" id="f-status-error" class="field__error">{{ errors.status }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Catatan -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__icon" aria-hidden="true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 20h9" />
              <path d="M16.5 3.5a2.1 2.1 0 0 1 3 3L7 19l-4 1 1-4z" />
            </svg>
          </span>
          <span class="form-section__titles">
            <span class="form-section__title">Catatan</span>
            <span class="form-section__hint">Keterangan tambahan dan klasifikasi UPL</span>
          </span>
        </legend>
        <div class="form-section__body">
          <div class="field">
            <label class="field__label" for="f-keterangan">Keterangan</label>
            <textarea
              id="f-keterangan"
              v-model="form.keterangan"
              class="textarea"
              :class="{ 'textarea--error': errors.keterangan }"
              rows="3"
              maxlength="280"
              placeholder="Catatan tambahan (opsional)"
              :aria-invalid="errors.keterangan ? 'true' : 'false'"
              :aria-describedby="errors.keterangan ? 'f-keterangan-error f-keterangan-count' : 'f-keterangan-count'"
              :disabled="saving"
              @input="setError('keterangan', '')"
            ></textarea>
            <span id="f-keterangan-count" class="field__count">{{ form.keterangan.length }}/280</span>
            <span v-if="errors.keterangan" id="f-keterangan-error" class="field__error">{{ errors.keterangan }}</span>
          </div>

          <div class="field">
            <span class="field__label" id="f-upl-label">Klasifikasi UPL</span>
            <div class="radio-group" role="radiogroup" aria-labelledby="f-upl-label">
              <label
                v-for="option in UPL_OPTIONS"
                :key="option.value"
                class="radio-option"
                :class="{ 'radio-option--active': form.upl === option.value }"
              >
                <input
                  type="radio"
                  name="upl"
                  :value="option.value"
                  v-model="form.upl"
                  :aria-invalid="errors.upl ? 'true' : 'false'"
                  :aria-describedby="errors.upl ? 'f-upl-error' : 'f-upl-hint'"
                  :disabled="saving"
                />
                <span>{{ option.label }}</span>
              </label>
            </div>
            <span id="f-upl-hint" class="field__hint">Opsional — pilih salah satu klasifikasi, atau biarkan kosong.</span>
            <span v-if="errors.upl" id="f-upl-error" class="field__error">{{ errors.upl }}</span>
          </div>
        </div>
      </fieldset>

      <div class="form-actions">
        <button class="btn btn--secondary" type="button" :disabled="saving" @click="cancel">Batal</button>
        <button class="btn btn--primary" type="submit" :disabled="saving">
          <span v-if="saving" class="spinner" aria-hidden="true"></span>
          {{ saving ? 'Menyimpan...' : 'Simpan' }}
        </button>
      </div>
    </form>
  </section>
</template>
