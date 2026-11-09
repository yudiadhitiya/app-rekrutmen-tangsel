<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
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

function setError(field, message) {
  if (message) errors[field] = message
  else delete errors[field]
}

function clearErrors() {
  Object.keys(errors).forEach((key) => delete errors[key])
  formError.value = ''
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
    return
  }

  saving.value = true
  try {
    const proceed = await checkDuplicateWa()
    if (!proceed) return

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
    } else if (err instanceof ApiError && err.status === 409) {
      setError('no_wa', 'No. WhatsApp sudah digunakan oleh data lain.')
      toast.error(err.message || 'No. WhatsApp sudah digunakan oleh data lain.')
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

      <!-- Identitas -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__num">1</span>
          <h3>Identitas</h3>
        </legend>
        <div class="form-section__body">
          <div class="field">
            <label class="field__label" for="f-nama">
              Nama <span class="field__required">*</span>
            </label>
            <input
              id="f-nama"
              v-model="form.nama"
              class="input"
              :class="{ 'input--error': errors.nama }"
              type="text"
              autocomplete="name"
              placeholder="Nama lengkap"
              :disabled="saving"
            />
            <span v-if="errors.nama" class="field__error">{{ errors.nama }}</span>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-no-wa">
                No. WhatsApp <span class="field__required">*</span>
              </label>
              <input
                id="f-no-wa"
                v-model="form.no_wa"
                class="input"
                :class="{ 'input--error': errors.no_wa }"
                type="text"
                inputmode="numeric"
                placeholder="628123456789"
                :disabled="saving"
              />
              <span class="field__hint">Hanya angka 0-9, contoh: 628123456789</span>
              <span v-if="errors.no_wa" class="field__error">{{ errors.no_wa }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-alamat">Alamat</label>
              <textarea
                id="f-alamat"
                v-model="form.alamat"
                class="textarea"
                :class="{ 'textarea--error': errors.alamat }"
                placeholder="Alamat lengkap"
                :disabled="saving"
              ></textarea>
              <span v-if="errors.alamat" class="field__error">{{ errors.alamat }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Wilayah -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__num">2</span>
          <h3>Wilayah</h3>
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
                :disabled="saving"
              />
              <span v-if="errors.rt" class="field__error">{{ errors.rt }}</span>
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
                :disabled="saving"
              />
              <span v-if="errors.rw" class="field__error">{{ errors.rw }}</span>
            </div>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-kecamatan">
                Kecamatan <span class="field__required">*</span>
              </label>
              <select
                id="f-kecamatan"
                v-model="form.kecamatan"
                class="select"
                :class="{ 'select--error': errors.kecamatan }"
                :disabled="saving"
              >
                <option value="">-- Pilih Kecamatan --</option>
                <option v-for="item in wilayah.kecamatan" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
              <span v-if="errors.kecamatan" class="field__error">{{ errors.kecamatan }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-kelurahan">Kelurahan</label>
              <select
                id="f-kelurahan"
                v-model="form.kelurahan"
                class="select"
                :class="{ 'select--error': errors.kelurahan }"
                :disabled="saving || !form.kecamatan"
              >
                <option value="">{{ form.kecamatan ? '-- Pilih Kelurahan --' : 'Pilih kecamatan terlebih dahulu' }}</option>
                <option v-for="item in kelurahanOptions" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
              <span v-if="errors.kelurahan" class="field__error">{{ errors.kelurahan }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Informasi Rekrutmen -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__num">3</span>
          <h3>Informasi Rekrutmen</h3>
        </legend>
        <div class="form-section__body">
          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-aktivitas">
                Aktivitas <span class="field__required">*</span>
              </label>
              <select
                id="f-aktivitas"
                v-model="form.aktivitas"
                class="select"
                :class="{ 'select--error': errors.aktivitas }"
                :disabled="saving"
              >
                <option value="">-- Pilih Aktivitas --</option>
                <option v-for="item in aktivitasOptions" :key="item.id" :value="item.nama">
                  {{ item.nama }}
                </option>
              </select>
              <span v-if="errors.aktivitas" class="field__error">{{ errors.aktivitas }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-tingkat">
                Tingkat Penerimaan <span class="field__required">*</span>
              </label>
              <select
                id="f-tingkat"
                v-model="form.tingkat_penerimaan"
                class="select"
                :class="{ 'select--error': errors.tingkat_penerimaan }"
                :disabled="saving"
              >
                <option value="">-- Pilih Tingkat Penerimaan --</option>
                <option v-for="item in TINGKAT_PENERIMAAN" :key="item" :value="item">{{ item }}</option>
              </select>
              <span v-if="errors.tingkat_penerimaan" class="field__error">{{ errors.tingkat_penerimaan }}</span>
            </div>
          </div>

          <div class="form-row">
            <div class="field">
              <label class="field__label" for="f-metode">
                Metode <span class="field__required">*</span>
              </label>
              <select
                id="f-metode"
                v-model="form.metode"
                class="select"
                :class="{ 'select--error': errors.metode }"
                :disabled="saving"
              >
                <option value="">-- Pilih Metode --</option>
                <option v-for="item in METODE" :key="item" :value="item">{{ item }}</option>
              </select>
              <span v-if="errors.metode" class="field__error">{{ errors.metode }}</span>
            </div>

            <div class="field">
              <label class="field__label" for="f-status">
                Status <span class="field__required">*</span>
              </label>
              <select
                id="f-status"
                v-model="form.status"
                class="select"
                :class="{ 'select--error': errors.status }"
                :disabled="saving"
              >
                <option value="">-- Pilih Status --</option>
                <option v-for="item in STATUS" :key="item" :value="item">{{ item }}</option>
              </select>
              <span v-if="errors.status" class="field__error">{{ errors.status }}</span>
            </div>
          </div>
        </div>
      </fieldset>

      <!-- Keterangan -->
      <fieldset class="form-section">
        <legend class="form-section__head">
          <span class="form-section__num">4</span>
          <h3>Keterangan</h3>
        </legend>
        <div class="form-section__body">
          <div class="field">
            <label class="field__label" for="f-keterangan">Keterangan</label>
            <textarea
              id="f-keterangan"
              v-model="form.keterangan"
              class="textarea"
              :class="{ 'textarea--error': errors.keterangan }"
              placeholder="Catatan tambahan (opsional)"
              :disabled="saving"
            ></textarea>
            <span v-if="errors.keterangan" class="field__error">{{ errors.keterangan }}</span>
          </div>

          <div class="field">
            <span class="field__label">UPL</span>
            <div class="radio-group">
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
                  :disabled="saving"
                />
                <span>{{ option.label }}</span>
              </label>
            </div>
            <span v-if="errors.upl" class="field__error">{{ errors.upl }}</span>
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
