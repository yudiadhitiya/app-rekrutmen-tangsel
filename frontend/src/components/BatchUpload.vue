<script setup>
// Uploading a batch of records from a spreadsheet.
//
// The previous application let a member (role ANGGOTA) upload a file of
// records (old-app/code.gs:732, old-app/index.html:322). The file is read in
// the browser and the server validates it again before saving anything.

import { computed, ref } from 'vue'
import { api } from '../api'
import { downloadTemplate, parseSpreadsheet } from '../excel'
import { toast } from '../toast'

// The same limits as the server (api/batch.go).
const MAX_ROWS = 1000
const MAX_SIZE = 5 * 1024 * 1024

const emit = defineEmits(['uploaded'])

const fileInput = ref(null)
const uploading = ref(false)
const fileName = ref('')
const rowCount = ref(0)

// The outcome of the last upload.
const summary = ref(null)
// The rows the server refused, with the reason of every problem.
const rowErrors = ref([])
const errorTotal = ref(0)
const errorMessage = ref('')

const skippedRows = computed(() => (summary.value && summary.value.skipped) || [])

function openPicker() {
  if (!uploading.value) fileInput.value && fileInput.value.click()
}

// The template of the upload, so the columns are written the right way.
function getTemplate() {
  try {
    downloadTemplate()
  } catch (err) {
    toast.error('Gagal membuat berkas template.')
  }
}

async function onFileSelected(event) {
  const input = event.target
  const file = input.files && input.files[0]
  // The same file may be chosen again after a failure.
  input.value = ''
  if (!file) return

  rowErrors.value = []
  errorTotal.value = 0
  errorMessage.value = ''
  summary.value = null
  fileName.value = file.name

  if (file.size > MAX_SIZE) {
    toast.error('Berkas terlalu besar. Maksimal 5 MB.')
    return
  }

  let rows = []
  try {
    rows = await parseSpreadsheet(file)
  } catch (err) {
    toast.error(err.message || 'Gagal membaca berkas.')
    return
  }

  if (!rows.length) {
    toast.error('Berkas tidak berisi data.')
    return
  }
  if (rows.length > MAX_ROWS) {
    toast.error(`Maksimal ${MAX_ROWS} baris per upload. Berkas berisi ${rows.length} baris.`)
    return
  }

  rowCount.value = rows.length
  uploading.value = true
  try {
    const data = await api.uploadBatch(rows)
    const created = Number((data && data.created) || 0)
    summary.value = { created, skipped: (data && data.skipped) || [] }
    toast.success(`${created} data berhasil disimpan.`)
    emit('uploaded')
  } catch (err) {
    const details = err.details || {}
    rowErrors.value = details.errors || []
    errorTotal.value = Number(details.error_count || rowErrors.value.length)
    errorMessage.value = err.message || 'Upload gagal.'
    toast.error(errorMessage.value)
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div class="card">
    <div class="card__header">
      <h2 class="card__title">Upload Data Massal</h2>
      <button class="btn btn--ghost btn--sm" type="button" @click="getTemplate">Unduh Template</button>
    </div>

    <div class="card__body">
      <p class="field__hint">
        Unduh template, isi datanya, lalu unggah berkasnya. Maksimal 1.000 baris
        per upload dan 5 MB. Bila satu baris salah, tidak ada baris yang disimpan.
      </p>

      <div class="upload-actions">
        <button class="btn btn--primary btn--sm" type="button" :disabled="uploading" @click="openPicker">
          Pilih Berkas
        </button>
        <span v-if="fileName" class="muted">{{ fileName }}</span>
      </div>

      <input
        ref="fileInput"
        class="visually-hidden"
        type="file"
        accept=".xlsx,.xls,.csv"
        @change="onFileSelected"
      />

      <div v-if="uploading" class="loading-block loading-block--tight">
        <span class="spinner spinner--dark" aria-hidden="true"></span>
        Mengunggah {{ rowCount }} baris...
      </div>

      <div v-if="errorMessage" class="alert alert--error" role="alert">
        {{ errorMessage }}
      </div>

      <div v-if="rowErrors.length" class="upload-errors">
        <p class="upload-errors__title">
          {{ errorTotal }} baris bermasalah. Baris yang ditampilkan dibatasi 30.
        </p>
        <ul class="upload-errors__list">
          <li v-for="(item, index) in rowErrors" :key="index">
            <span class="upload-errors__row">Baris {{ item.row }}</span>
            <span v-if="item.field" class="upload-errors__field">{{ item.field }}</span>
            <span>{{ item.reason }}</span>
          </li>
        </ul>
      </div>

      <div v-if="summary" class="alert alert--success" role="status">
        {{ summary.created }} data tersimpan.
      </div>

      <div v-if="skippedRows.length" class="upload-errors">
        <p class="upload-errors__title">
          {{ skippedRows.length }} baris dilewati karena nomor WhatsApp sudah terdaftar.
        </p>
        <ul class="upload-errors__list">
          <li v-for="item in skippedRows" :key="item.row">
            <span class="upload-errors__row">Baris {{ item.row }}</span>
            <span class="upload-errors__field">{{ item.nama }}</span>
            <span>{{ item.no_wa }} — {{ item.reason }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>
