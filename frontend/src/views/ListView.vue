<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'
import { isAdmin } from '../auth'
import { METODE, STATUS, tingkatTone, statusTone, normalizeMetode } from '../constants'
import { toast } from '../toast'
import { confirmDialog } from '../confirm'

const router = useRouter()

// Only an admin may delete a record; an officer sees and edits their own.
const admin = computed(() => isAdmin())

const rows = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(25)
const loading = ref(false)
const deletingId = ref('')
const errorMessage = ref('')

const wilayah = ref({ kecamatan: [] })
const wilayahLoaded = ref(false)

const filters = reactive({
  q: '',
  kecamatan: '',
  kelurahan: '',
  metode: '',
  status: ''
})

// Debounced search term so typing does not fire a request per keystroke.
const searchTerm = ref('')
let searchTimer = null
watch(
  () => filters.q,
  (value) => {
    if (searchTimer) clearTimeout(searchTimer)
    searchTimer = setTimeout(() => {
      searchTerm.value = value.trim()
      page.value = 1
    }, 300)
  }
)

const activeFilters = computed(() => ({
  q: searchTerm.value,
  kecamatan: filters.kecamatan,
  kelurahan: filters.kelurahan,
  metode: filters.metode,
  status: filters.status
}))

// The empty list means two different things: the officer has not entered
// anything yet, or the filter does not match any record.
const hasActiveFilter = computed(() =>
  Object.values(activeFilters.value).some((value) => value !== '')
)

const emptyTitle = computed(() =>
  !admin.value && !hasActiveFilter.value ? 'Belum ada data Anda' : 'Tidak ada data'
)

const emptyMessage = computed(() =>
  !admin.value && !hasActiveFilter.value
    ? 'Data yang Anda input akan tampil di sini.'
    : 'Belum ada data rekrutmen yang sesuai dengan filter yang dipilih.'
)

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))

const pageStart = computed(() => (total.value === 0 ? 0 : (page.value - 1) * limit.value + 1))
const pageEnd = computed(() => Math.min(page.value * limit.value, total.value))

const selectedKecamatan = computed(() =>
  wilayah.value.kecamatan.find((item) => item.id === filters.kecamatan)
)

const kelurahanOptions = computed(() => {
  if (!selectedKecamatan.value) return []
  return selectedKecamatan.value.kelurahan || []
})

async function loadWilayah() {
  try {
    wilayah.value = await api.wilayah()
    if (!wilayah.value) wilayah.value = { kecamatan: [] }
  } catch (err) {
    wilayah.value = { kecamatan: [] }
  } finally {
    wilayahLoaded.value = true
  }
}

async function loadRows() {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api.rekrutmen.list({
      ...activeFilters.value,
      page: page.value,
      limit: limit.value
    })
    // Fold the retired Metode values into the current option so old records
    // still match the "Program Antara" filter and show the right label.
    rows.value = ((data && data.data) || []).map((row) => ({
      ...row,
      metode: normalizeMetode(row.metode)
    }))
    total.value = Number((data && data.total) || 0)
    if (page.value > totalPages.value) page.value = totalPages.value
  } catch (err) {
    rows.value = []
    total.value = 0
    errorMessage.value = err.message || 'Gagal memuat data rekrutmen.'
  } finally {
    loading.value = false
  }
}

function resetFilters() {
  if (searchTimer) clearTimeout(searchTimer)
  filters.q = ''
  searchTerm.value = ''
  filters.kecamatan = ''
  filters.kelurahan = ''
  filters.metode = ''
  filters.status = ''
  page.value = 1
}

// Changing kecamatan always resets the dependent kelurahan filter.
watch(
  () => filters.kecamatan,
  () => {
    filters.kelurahan = ''
  }
)

// Single watcher keeps every list interaction in sync with the backend.
watch(
  () => JSON.stringify({ ...activeFilters.value, page: page.value, limit: limit.value }),
  () => {
    loadRows()
  }
)

function goToPage(next) {
  if (next < 1 || next > totalPages.value) return
  page.value = next
}

function editRecord(record) {
  router.push({ path: '/ubah/' + encodeURIComponent(record.id) })
}

async function deleteRecord(record) {
  const confirmed = await confirmDialog({
    title: 'Hapus data rekrutmen?',
    message: `Data "${record.nama}" (${record.no_wa}) akan dihapus permanen. Tindakan ini tidak dapat dibatalkan.`,
    confirmText: 'Ya, hapus'
  })
  if (!confirmed) return

  deletingId.value = record.id
  try {
    await api.rekrutmen.remove(record.id)
    toast.success(`Data "${record.nama}" berhasil dihapus.`)
    if (rows.value.length === 1 && page.value > 1) page.value -= 1
    else await loadRows()
  } catch (err) {
    toast.error(err.message || 'Gagal menghapus data rekrutmen.')
  } finally {
    deletingId.value = ''
  }
}

onMounted(() => {
  loadWilayah()
  loadRows()
})
</script>

<template>
  <section>
    <div class="page-head">
      <div>
        <h1 class="page-head__title">Data Rekrutmen</h1>
        <p class="page-head__subtitle">Kelola seluruh data rekrutmen wilayah Kota Tangerang Selatan.</p>
      </div>
      <RouterLink to="/tambah" class="btn btn--primary">+ Tambah Rekrutmen</RouterLink>
    </div>

    <div class="stack">
      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Pencarian &amp; Filter</h2>
          <button class="btn btn--ghost btn--sm" type="button" @click="resetFilters">Reset</button>
        </div>
        <div class="card__body">
          <div class="filters">
            <div class="field filters__search">
              <label class="field__label" for="filter-q">Cari</label>
              <input
                id="filter-q"
                v-model="filters.q"
                class="input"
                type="search"
                placeholder="Cari nama atau No. WhatsApp..."
              />
            </div>

            <div class="field">
              <label class="field__label" for="filter-kecamatan">Kecamatan</label>
              <select id="filter-kecamatan" v-model="filters.kecamatan" class="select">
                <option value="">Semua Kecamatan</option>
                <option v-for="item in wilayah.kecamatan" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
            </div>

            <div class="field">
              <label class="field__label" for="filter-kelurahan">Kelurahan</label>
              <select
                id="filter-kelurahan"
                v-model="filters.kelurahan"
                class="select"
                :disabled="!filters.kecamatan"
              >
                <option value="">{{ filters.kecamatan ? 'Semua Kelurahan' : 'Pilih kecamatan dulu' }}</option>
                <option v-for="item in kelurahanOptions" :key="item.id" :value="item.id">{{ item.nama }}</option>
              </select>
            </div>

            <div class="field">
              <label class="field__label" for="filter-metode">Metode</label>
              <select id="filter-metode" v-model="filters.metode" class="select">
                <option value="">Semua Metode</option>
                <option v-for="item in METODE" :key="item" :value="item">{{ item }}</option>
              </select>
            </div>

            <div class="field">
              <label class="field__label" for="filter-status">Status</label>
              <select id="filter-status" v-model="filters.status" class="select">
                <option value="">Semua Status</option>
                <option v-for="item in STATUS" :key="item" :value="item">{{ item }}</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Daftar Rekrutmen</h2>
          <span class="result-count">{{ total }} data ditemukan</span>
        </div>

        <div v-if="errorMessage" class="card__body">
          <div class="alert alert--error" role="alert">{{ errorMessage }}</div>
        </div>

        <div v-if="loading" class="loading-block">
          <span class="spinner spinner--dark" aria-hidden="true"></span>
          Memuat data...
        </div>

        <template v-else-if="rows.length === 0">
          <div class="empty-state">
            <div class="empty-state__title">{{ emptyTitle }}</div>
            <p>{{ emptyMessage }}</p>
          </div>
        </template>

        <template v-else>
          <!-- Wide screens: table -->
          <div class="table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Nama</th>
                  <th>No. WhatsApp</th>
                  <th>Kecamatan</th>
                  <th>Kelurahan</th>
                  <th>Aktivitas</th>
                  <th>Tingkat Penerimaan</th>
                  <th>Metode</th>
                  <th>Status</th>
                  <th>UPL</th>
                  <th style="text-align: right">Aksi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in rows" :key="row.id">
                  <td class="data-table__strong">{{ row.nama }}</td>
                  <td>{{ row.no_wa }}</td>
                  <td>{{ row.kecamatan }}</td>
                  <td>{{ row.kelurahan || '-' }}</td>
                  <td>{{ row.aktivitas }}</td>
                  <td>
                    <span class="badge" :class="tingkatTone(row.tingkat_penerimaan)">{{ row.tingkat_penerimaan }}</span>
                  </td>
                  <td>{{ row.metode }}</td>
                  <td>
                    <span class="badge badge--dot" :class="statusTone(row.status)">{{ row.status }}</span>
                  </td>
                  <td>
                    <span v-if="row.upl" class="badge tone-neutral">{{ row.upl }}</span>
                    <span v-else class="muted">-</span>
                  </td>
                  <td>
                    <div class="data-table__actions">
                      <button class="btn btn--secondary btn--sm" type="button" @click="editRecord(row)">Ubah</button>
                      <button
                        v-if="admin"
                        class="btn btn--danger-soft btn--sm"
                        type="button"
                        :disabled="deletingId === row.id"
                        @click="deleteRecord(row)"
                      >
                        {{ deletingId === row.id ? '...' : 'Hapus' }}
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Narrow screens: cards -->
          <div class="recruiter-cards">
            <article v-for="row in rows" :key="row.id" class="recruiter-card">
              <div class="recruiter-card__head">
                <div>
                  <div class="recruiter-card__name">{{ row.nama }}</div>
                  <div class="recruiter-card__wa">{{ row.no_wa }}</div>
                </div>
                <span class="badge badge--dot" :class="statusTone(row.status)">{{ row.status }}</span>
              </div>

              <div class="recruiter-card__badges">
                <span class="badge" :class="tingkatTone(row.tingkat_penerimaan)">{{ row.tingkat_penerimaan }}</span>
                <span class="badge tone-neutral">{{ row.metode }}</span>
                <span v-if="row.upl" class="badge tone-info">{{ row.upl }}</span>
              </div>

              <div class="recruiter-card__meta">
                <div class="recruiter-card__meta-item">
                  <span class="recruiter-card__meta-label">Kecamatan</span>
                  <span class="recruiter-card__meta-value">{{ row.kecamatan }}</span>
                </div>
                <div class="recruiter-card__meta-item">
                  <span class="recruiter-card__meta-label">Kelurahan</span>
                  <span class="recruiter-card__meta-value">{{ row.kelurahan || '-' }}</span>
                </div>
                <div class="recruiter-card__meta-item">
                  <span class="recruiter-card__meta-label">Aktivitas</span>
                  <span class="recruiter-card__meta-value">{{ row.aktivitas }}</span>
                </div>
                <div class="recruiter-card__meta-item">
                  <span class="recruiter-card__meta-label">Alamat</span>
                  <span class="recruiter-card__meta-value">{{ row.alamat || '-' }}</span>
                </div>
              </div>

              <div class="recruiter-card__actions">
                <button class="btn btn--secondary btn--sm" type="button" @click="editRecord(row)">Ubah</button>
                <button
                  v-if="admin"
                  class="btn btn--danger-soft btn--sm"
                  type="button"
                  :disabled="deletingId === row.id"
                  @click="deleteRecord(row)"
                >
                  Hapus
                </button>
              </div>
            </article>
          </div>

          <div class="pagination">
            <span class="pagination__info">
              Menampilkan {{ pageStart }}-{{ pageEnd }} dari {{ total }} data
            </span>

            <div class="pagination__controls">
              <select v-model.number="limit" class="select select--compact" aria-label="Jumlah baris per halaman">
                <option :value="10">10</option>
                <option :value="25">25</option>
                <option :value="50">50</option>
                <option :value="100">100</option>
              </select>
              <button class="btn btn--secondary btn--sm" type="button" :disabled="page <= 1" @click="goToPage(page - 1)">
                Sebelumnya
              </button>
              <span class="pagination__page">{{ page }} / {{ totalPages }}</span>
              <button
                class="btn btn--secondary btn--sm"
                type="button"
                :disabled="page >= totalPages"
                @click="goToPage(page + 1)"
              >
                Berikutnya
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>
  </section>
</template>
