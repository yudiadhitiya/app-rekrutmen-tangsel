<script setup>
import { computed, onMounted, ref } from 'vue'
import { api } from '../api'
import { isMemberOnly } from '../auth'

// The dashboard is the high level view of the recruitment data. It was part
// of the previous application and is restored here, so the managers keep the
// same numbers they used to see.

const loading = ref(false)
const errorMessage = ref('')
const data = ref({
  scope: '',
  total: 0,
  methods: [],
  methodSummary: [],
  acceptance: [],
  status: [],
  statusDetail: [],
  kecamatanSummary: [],
  top10: []
})

// Only a member may choose the scope of the dashboard. A district
// administrator always sees its own district and a full administrator
// always sees every record.
const mayChooseFilter = computed(() => isMemberOnly())
const filter = ref('ALL')

const filters = [
  { value: 'ALL', label: 'Semua Kecamatan' },
  { value: 'KECAMATAN', label: 'Data Kecamatan Saya' },
  { value: 'ANGGOTA', label: 'Data Anggota Saya' }
]

// The four cards of the head of the dashboard.
const cards = computed(() => {
  const methods = data.value.methods || []
  const find = (name) => methods.find((item) => item.name === name) || { total: 0 }
  return [
    { key: 'total', label: 'Total Target', value: data.value.total, tone: 'blue' },
    { key: 'metode', label: 'Halaqoh', value: find('Halaqoh').total, tone: 'green' },
    { key: 'metode', label: 'Prog. Antara', value: find('Program Antara').total, tone: 'orange' },
    { key: 'metode', label: 'Dakwah Fardiyah', value: find('Dakwah Fardiyah').total, tone: 'purple' }
  ]
})

// The total of a breakdown list, used by the foot of the tables.
function sumOf(items, field) {
  return (items || []).reduce((sum, item) => sum + Number(item[field] || 0), 0)
}

const methodTotals = computed(() => {
  const items = data.value.methodSummary || []
  return {
    target: sumOf(items, 'total'),
    berhasil: sumOf(items, 'berhasil'),
    proses: sumOf(items, 'proses'),
    gagalPending: sumOf(items, 'gagalPending')
  }
})

const kecamatanTotals = computed(() => {
  const items = data.value.kecamatanSummary || []
  const target = sumOf(items, 'target')
  const berhasil = sumOf(items, 'berhasil')
  const proses = sumOf(items, 'proses')
  return {
    target,
    berhasil,
    proses,
    persenBerhasil: percentOf(berhasil, target),
    persenProses: percentOf(proses, target)
  }
})

// The scope line under the cards.
const scopeText = computed(() => {
  const scope = data.value.scope || 'Rekapitulasi'
  return `${scope} • Total ${formatNumber(data.value.total)} target`
})

function formatNumber(value) {
  return Number(value || 0).toLocaleString('id-ID')
}

// A percentage with one decimal, unless it is a whole number.
function formatPercent(value) {
  const number = Number(value || 0)
  return Number.isInteger(number) ? `${number}%` : `${number.toFixed(1)}%`
}

// The share of the value in the total.
function percentOf(value, total) {
  if (!total) return 0
  return Math.round((Number(value || 0) / total) * 1000) / 10
}

async function load() {
  loading.value = true
  errorMessage.value = ''
  try {
    const query = mayChooseFilter.value ? { filter: filter.value } : {}
    const result = await api.dashboard(query)
    data.value = {
      scope: (result && result.scope) || '',
      total: Number((result && result.total) || 0),
      methods: (result && result.methods) || [],
      methodSummary: (result && result.methodSummary) || [],
      acceptance: (result && result.acceptance) || [],
      status: (result && result.status) || [],
      statusDetail: (result && result.statusDetail) || [],
      kecamatanSummary: (result && result.kecamatanSummary) || [],
      top10: (result && result.top10) || []
    }
  } catch (err) {
    errorMessage.value = err.message || 'Gagal memuat dashboard.'
  } finally {
    loading.value = false
  }
}

function onFilterChange() {
  load()
}

onMounted(load)
</script>

<template>
  <section>
    <div class="page-head">
      <div>
        <h1 class="page-head__title">Dashboard Rekapitulasi</h1>
        <p class="page-head__subtitle">Ringkasan pendaftaran anggota pelapor di wilayah Kota Tangerang Selatan.</p>
      </div>
      <div class="page-head__actions">
        <select
          v-if="mayChooseFilter"
          v-model="filter"
          class="select select--compact"
          aria-label="Filter dashboard"
          @change="onFilterChange"
        >
          <option v-for="item in filters" :key="item.value" :value="item.value">{{ item.label }}</option>
        </select>
        <button class="btn btn--secondary" type="button" :disabled="loading" @click="load">
          {{ loading ? 'Memuat...' : '↻ Refresh' }}
        </button>
      </div>
    </div>

    <div class="stack">
      <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

      <div class="dash-cards">
        <article v-for="card in cards" :key="card.label" class="dash-card" :class="`dash-card--${card.tone}`">
          <div class="dash-card__head">{{ card.label }}</div>
          <div class="dash-card__number">{{ formatNumber(card.value) }}</div>
          <div class="dash-card__foot">orang</div>
        </article>
      </div>

      <p class="dash-scope">{{ scopeText }}</p>

      <div class="dash-grid">
        <div class="card">
          <div class="card__header">
            <h2 class="card__title">Berdasarkan Tingkat Penerimaan</h2>
          </div>
          <div class="card__body card__body--flush">
            <template v-if="data.acceptance.length === 0">
              <div class="empty-state">
                <div class="empty-state__title">Belum ada data</div>
                <p>Belum ada data rekapitulasi untuk ditampilkan.</p>
              </div>
            </template>
            <template v-else>
              <div v-for="item in data.acceptance" :key="item.name" class="dash-row">
                <span class="dash-row__label">{{ item.name }}</span>
                <span class="dash-row__value">{{ formatNumber(item.total) }}</span>
                <span class="dash-row__percent">{{ formatPercent(item.percent) }}</span>
                <span class="dash-row__bar">
                  <span class="dash-row__bar-fill" :style="{ width: Math.min(100, item.percent) + '%' }"></span>
                </span>
              </div>
            </template>
          </div>
        </div>

        <div class="card">
          <div class="card__header">
            <h2 class="card__title">Berdasarkan Status</h2>
          </div>
          <div class="card__body card__body--flush">
            <template v-if="data.status.length === 0">
              <div class="empty-state">
                <div class="empty-state__title">Belum ada data</div>
                <p>Belum ada data rekapitulasi untuk ditampilkan.</p>
              </div>
            </template>
            <template v-else>
              <div v-for="item in data.status" :key="item.name" class="dash-row">
                <span class="dash-row__label">{{ item.name }}</span>
                <span class="dash-row__value">{{ formatNumber(item.total) }}</span>
                <span class="dash-row__percent">{{ formatPercent(item.percent) }}</span>
                <span class="dash-row__bar">
                  <span class="dash-row__bar-fill" :style="{ width: Math.min(100, item.percent) + '%' }"></span>
                </span>
              </div>
            </template>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Ringkasan Metode Pendaftaran</h2>
        </div>
        <div class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th>Metode Pendaftaran</th>
                <th>Jumlah Target</th>
                <th>Berhasil</th>
                <th>Proses</th>
                <th>Gagal/Pending</th>
                <th>% dari Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in data.methodSummary" :key="item.name">
                <td class="data-table__strong">{{ item.name }}</td>
                <td>{{ formatNumber(item.total) }}</td>
                <td>{{ formatNumber(item.berhasil) }}</td>
                <td>{{ formatNumber(item.proses) }}</td>
                <td>{{ formatNumber(item.gagalPending) }}</td>
                <td>{{ formatPercent(item.percent) }}</td>
              </tr>
              <tr v-if="data.methodSummary.length === 0">
                <td colspan="6" class="table-empty">Belum ada data.</td>
              </tr>
            </tbody>
            <tfoot v-if="data.methodSummary.length > 0">
              <tr>
                <td>TOTAL</td>
                <td>{{ formatNumber(methodTotals.target) }}</td>
                <td>{{ formatNumber(methodTotals.berhasil) }}</td>
                <td>{{ formatNumber(methodTotals.proses) }}</td>
                <td>{{ formatNumber(methodTotals.gagalPending) }}</td>
                <td>100%</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Summary per Kecamatan</h2>
        </div>
        <div class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th>Kecamatan</th>
                <th>Target</th>
                <th>Berhasil</th>
                <th>Proses</th>
                <th>% Berhasil</th>
                <th>% Proses</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in data.kecamatanSummary" :key="item.kecamatan">
                <td class="data-table__strong">{{ item.kecamatan }}</td>
                <td>{{ formatNumber(item.target) }}</td>
                <td>{{ formatNumber(item.berhasil) }}</td>
                <td>{{ formatNumber(item.proses) }}</td>
                <td>{{ formatPercent(item.persenBerhasil) }}</td>
                <td>{{ formatPercent(item.persenProses) }}</td>
              </tr>
            </tbody>
            <tfoot>
              <tr>
                <td>TOTAL</td>
                <td>{{ formatNumber(kecamatanTotals.target) }}</td>
                <td>{{ formatNumber(kecamatanTotals.berhasil) }}</td>
                <td>{{ formatNumber(kecamatanTotals.proses) }}</td>
                <td>{{ formatPercent(kecamatanTotals.persenBerhasil) }}</td>
                <td>{{ formatPercent(kecamatanTotals.persenProses) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <div class="card">
        <div class="card__header">
          <h2 class="card__title">Top 10 Nama</h2>
        </div>
        <div class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 60px">#</th>
                <th>Nama</th>
                <th style="text-align: right">Jumlah</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in data.top10" :key="item.nama">
                <td>{{ index + 1 }}</td>
                <td class="data-table__strong">{{ item.nama }}</td>
                <td class="text-right">{{ formatNumber(item.total) }}</td>
              </tr>
              <tr v-if="data.top10.length === 0">
                <td colspan="3" class="table-empty">Belum ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </section>
</template>
