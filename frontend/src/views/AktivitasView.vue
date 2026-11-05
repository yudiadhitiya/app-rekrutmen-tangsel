<script setup>
import { onMounted, reactive, ref } from 'vue'
import { api } from '../api'
import { isAdmin } from '../auth'
import { toast } from '../toast'
import { confirmDialog } from '../confirm'

const items = ref([])
const loading = ref(false)
const errorMessage = ref('')

const newName = ref('')
const adding = ref(false)

const editingId = ref(null)
const editForm = reactive({ nama: '', aktif: true })
const savingEdit = ref(false)

const admin = isAdmin()

async function loadItems() {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api.aktivitas.list()
    items.value = Array.isArray(data) ? data : []
  } catch (err) {
    items.value = []
    errorMessage.value = err.message || 'Gagal memuat master aktivitas.'
  } finally {
    loading.value = false
  }
}

async function addItem() {
  const nama = newName.value.trim()
  if (!nama) {
    toast.error('Nama aktivitas wajib diisi.')
    return
  }
  adding.value = true
  try {
    await api.aktivitas.create(nama)
    newName.value = ''
    toast.success(`Aktivitas "${nama}" berhasil ditambahkan.`)
    await loadItems()
  } catch (err) {
    toast.error(err.message || 'Gagal menambahkan aktivitas.')
  } finally {
    adding.value = false
  }
}

function startEdit(item) {
  editingId.value = item.id
  editForm.nama = item.nama || ''
  editForm.aktif = item.aktif !== false
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit(item) {
  const nama = editForm.nama.trim()
  if (!nama) {
    toast.error('Nama aktivitas wajib diisi.')
    return
  }
  savingEdit.value = true
  try {
    await api.aktivitas.update(item.id, { nama, aktif: !!editForm.aktif })
    toast.success(`Aktivitas "${nama}" berhasil diperbarui.`)
    editingId.value = null
    await loadItems()
  } catch (err) {
    toast.error(err.message || 'Gagal memperbarui aktivitas.')
  } finally {
    savingEdit.value = false
  }
}

async function removeItem(item) {
  const confirmed = await confirmDialog({
    title: 'Hapus aktivitas?',
    message: `Aktivitas "${item.nama}" akan dihapus dari master. Data rekrutmen yang memakai aktivitas ini tidak akan bisa dihapus.`,
    confirmText: 'Ya, hapus'
  })
  if (!confirmed) return

  try {
    await api.aktivitas.remove(item.id)
    toast.success(`Aktivitas "${item.nama}" berhasil dihapus.`)
    await loadItems()
  } catch (err) {
    if (err.status === 409) {
      toast.error('Aktivitas ini masih digunakan oleh data rekrutmen sehingga tidak bisa dihapus.')
    } else {
      toast.error(err.message || 'Gagal menghapus aktivitas.')
    }
  }
}

onMounted(loadItems)
</script>

<template>
  <section>
    <div class="page-head">
      <div>
        <h1 class="page-head__title">Master Aktivitas</h1>
        <p class="page-head__subtitle">
          Daftar kategori pekerjaan/aktivitas yang dipakai pada form rekrutmen.
        </p>
      </div>
    </div>

    <div class="stack">
      <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

      <div v-if="admin" class="card">
        <div class="card__header">
          <h2 class="card__title">Tambah Aktivitas</h2>
        </div>
        <div class="card__body">
          <form class="inline-form" @submit.prevent="addItem">
            <input
              v-model="newName"
              class="input"
              type="text"
              placeholder="Contoh: Wiraswasta"
              maxlength="120"
              :disabled="adding"
            />
            <button class="btn btn--primary" type="submit" :disabled="adding">
              <span v-if="adding" class="spinner" aria-hidden="true"></span>
              {{ adding ? 'Menyimpan...' : 'Tambah' }}
            </button>
          </form>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Daftar Aktivitas</h2>
          <span class="result-count">{{ items.length }} aktivitas</span>
        </div>

        <div v-if="loading" class="loading-block">
          <span class="spinner spinner--dark" aria-hidden="true"></span>
          Memuat data...
        </div>

        <div v-else-if="items.length === 0" class="empty-state">
          <div class="empty-state__title">Belum ada aktivitas</div>
          <p>{{ admin ? 'Tambahkan aktivitas pertama Anda di atas.' : 'Admin belum menambahkan aktivitas.' }}</p>
        </div>

        <div v-else class="master-list">
          <div v-for="item in items" :key="item.id" class="master-item">
            <template v-if="editingId === item.id">
              <div class="inline-form" style="flex: 1 1 auto">
                <input
                  v-model="editForm.nama"
                  class="input"
                  type="text"
                  placeholder="Nama aktivitas"
                  :disabled="savingEdit"
                />
                <label class="checkbox-inline">
                  <input v-model="editForm.aktif" type="checkbox" :disabled="savingEdit" />
                  Aktif
                </label>
              </div>
              <div class="master-item__actions">
                <button class="btn btn--secondary btn--sm" type="button" :disabled="savingEdit" @click="cancelEdit">
                  Batal
                </button>
                <button class="btn btn--primary btn--sm" type="button" :disabled="savingEdit" @click="saveEdit(item)">
                  {{ savingEdit ? 'Menyimpan...' : 'Simpan' }}
                </button>
              </div>
            </template>

            <template v-else>
              <div>
                <div class="master-item__name">{{ item.nama }}</div>
                <div class="master-item__meta">{{ item.aktif === false ? 'Nonaktif' : 'Aktif' }}</div>
              </div>
              <div v-if="admin" class="master-item__actions">
                <button class="btn btn--secondary btn--sm" type="button" @click="startEdit(item)">Ubah</button>
                <button class="btn btn--danger-soft btn--sm" type="button" @click="removeItem(item)">Hapus</button>
              </div>
              <span v-else class="badge" :class="item.aktif === false ? 'tone-neutral' : 'tone-success'">
                {{ item.aktif === false ? 'Nonaktif' : 'Aktif' }}
              </span>
            </template>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
