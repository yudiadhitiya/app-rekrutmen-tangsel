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
    { key: 'halaqoh', label: 'Halaqoh', value: find('Halaqoh').total, tone: 'green' },
    { key: 'antara', label: 'Prog. Antara', value: find('Program Antara').total, tone: 'orange' },
    { key: 'fardiyah', label: 'Dakwah Fardiyah', value: find('Dakwah Fardiyah').total, tone: 'purple' }
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

// A part to whole bar of the registration outcome.
const outcomeSegments = computed(() => {
  const totals = methodTotals.value
  const target = totals.target || 0
  const parts = [
    { key: 'berhasil', label: 'Berhasil', value: totals.berhasil, tone: 'success' },
    { key: 'proses', label: 'Proses', value: totals.proses, tone: 'warning' },
    { key: 'gagalPending', label: 'Gagal/Pending', value: totals.gagalPending, tone: 'danger' }
  ]
  return parts.map((part) => ({
    ...part,
    percent: percentOf(part.value, target)
  }))
})

// The scope line under the cards.
const scopeText = computed(() => {
  const scope = data.value.scope || 'Rekapitulasi'
  return `${scope} • Total ${formatNumber(data.value.total)} target`
})

// The share of one method card against the grand total.
function shareOf(value) {
  return percentOf(value, data.value.total)
}

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

// A semantic tone for a status or acceptance label, so the colors stay
// consistent across the dashboard.
function toneOf(name) {
  const text = String(name || '').toLowerCase()
  if (text.includes('berhasil') || text.includes('diterima') || text.includes('lolos')) return 'success'
  if (text.includes('gagal') || text.includes('tolak') || text.includes('batal')) return 'danger'
  if (text.includes('proses') || text.includes('menunggu') || text.includes('pending')) return 'warning'
  if (text.includes('hold') || text.includes('tunda')) return 'neutral'
  return 'info'
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
  <section class="dash" :aria-busy="loading">
    <div class="page-head">
      <div>
        <h1 class="page-head__title">Dashboard Rekapitulasi</h1>
        <p class="page-head__subtitle">Ringkasan pendaftaran anggota pelapor di wilayah Kota Tangerang Selatan.</p>
      </div>
      <div class="page-head__actions">
        <div v-if="mayChooseFilter" class="dash-segment" role="group" aria-label="Filter dashboard">
          <button
            v-for="item in filters"
            :key="item.value"
            class="dash-segment__item"
            :class="{ 'dash-segment__item--active': filter === item.value }"
            type="button"
            :aria-pressed="filter === item.value"
            @click="filter = item.value; onFilterChange()"
          >
            {{ item.label }}
          </button>
        </div>
        <button class="btn btn--secondary" type="button" :disabled="loading" @click="load">
          <svg
            class="dash-spin"
            :class="{ 'dash-spin--active': loading }"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            aria-hidden="true"
          >
            <path d="M21 12a9 9 0 1 1-2.64-6.36" />
            <path d="M21 3v6h-6" />
          </svg>
          {{ loading ? 'Memuat...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div class="stack">
      <div v-if="errorMessage" class="alert alert--error" role="alert">{{ errorMessage }}</div>

      <div class="dash-cards">
        <article
          v-for="card in cards"
          :key="card.label"
          class="dash-card"
          :class="`dash-card--${card.tone}`"
        >
          <div class="dash-card__head">
            <span class="dash-card__icon" aria-hidden="true">
              <svg v-if="card.key === 'total'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9.5" cy="7" r="4" />
                <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
              </svg>
              <svg v-else-if="card.key === 'halaqoh'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
              </svg>
              <svg v-else-if="card.key === 'antara'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 5h16" />
                <path d="M4 12h10" />
                <path d="M4 19h6" />
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="9" />
                <path d="M12 3v9l6 4" />
              </svg>
            </span>
            <span>{{ card.label }}</span>
          </div>
          <div class="dash-card__number">{{ formatNumber(card.value) }}</div>
          <div class="dash-card__foot">
            {{ card.key === 'total' ? 'seluruh wilayah' : formatPercent(shareOf(card.value)) + ' dari total' }}
          </div>
        </article>
      </div>

      <p class="dash-scope">{{ scopeText }}</p>

      <div class="card dash-outcome">
        <div class="card__header">
          <h2 class="card__title">Outcome Pendaftaran</h2>
          <span class="dash-outcome__total">{{ formatNumber(methodTotals.target) }} target</span>
        </div>
        <div class="card__body">
          <div
            class="dash-stack"
            role="img"
            :aria-label="`Distribusi outcome: ${outcomeSegments.map((item) => `${item.label} ${item.percent}%`).join(', ')}`"
          >
            <span
              v-for="item in outcomeSegments"
              :key="item.key"
              class="dash-stack__seg"
              :class="`dash-stack__seg--${item.tone}`"
              :style="{ width: Math.min(100, item.percent) + '%' }"
            ></span>
          </div>
          <ul class="dash-legend">
            <li v-for="item in outcomeSegments" :key="item.key" class="dash-legend__item">
              <span class="dash-legend__dot" :class="`dash-legend__dot--${item.tone}`" aria-hidden="true"></span>
              <span class="dash-legend__label">{{ item.label }}</span>
              <span class="dash-legend__value">{{ formatNumber(item.value) }}</span>
              <span class="dash-legend__percent">{{ formatPercent(item.percent) }}</span>
            </li>
          </ul>
        </div>
      </div>

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
                <span
                  class="dash-row__bar"
                  role="progressbar"
                  :aria-label="`${item.name}: ${formatPercent(item.percent)} dari total target`"
                  aria-valuemin="0"
                  aria-valuemax="100"
                  :aria-valuenow="Math.round(item.percent)"
                >
                  <span
                    class="dash-row__bar-fill"
                    :class="`dash-row__bar-fill--${toneOf(item.name)}`"
                    :style="{ width: Math.min(100, item.percent) + '%' }"
                  ></span>
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
                <span
                  class="dash-row__bar"
                  role="progressbar"
                  :aria-label="`${item.name}: ${formatPercent(item.percent)} dari total target`"
                  aria-valuemin="0"
                  aria-valuemax="100"
                  :aria-valuenow="Math.round(item.percent)"
                >
                  <span
                    class="dash-row__bar-fill"
                    :class="`dash-row__bar-fill--${toneOf(item.name)}`"
                    :style="{ width: Math.min(100, item.percent) + '%' }"
                  ></span>
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
          <table class="data-table dash-table">
            <caption class="visually-hidden">
              Ringkasan metode pendaftaran beserta jumlah target, berhasil, proses, dan gagal atau pending.
            </caption>
            <thead>
              <tr>
                <th scope="col">Metode Pendaftaran</th>
                <th scope="col" class="dash-table__num">Jumlah Target</th>
                <th scope="col" class="dash-table__num">Berhasil</th>
                <th scope="col" class="dash-table__num">Proses</th>
                <th scope="col" class="dash-table__num">Gagal/Pending</th>
                <th scope="col" class="dash-table__num">% dari Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in data.methodSummary" :key="item.name">
                <td class="data-table__strong">{{ item.name }}</td>
                <td class="dash-table__num">{{ formatNumber(item.total) }}</td>
                <td class="dash-table__num">{{ formatNumber(item.berhasil) }}</td>
                <td class="dash-table__num">{{ formatNumber(item.proses) }}</td>
                <td class="dash-table__num">{{ formatNumber(item.gagalPending) }}</td>
                <td class="dash-table__num">{{ formatPercent(item.percent) }}</td>
              </tr>
              <tr v-if="data.methodSummary.length === 0">
                <td colspan="6" class="table-empty">Belum ada data.</td>
              </tr>
            </tbody>
            <tfoot v-if="data.methodSummary.length > 0">
              <tr>
                <td>TOTAL</td>
                <td class="dash-table__num">{{ formatNumber(methodTotals.target) }}</td>
                <td class="dash-table__num">{{ formatNumber(methodTotals.berhasil) }}</td>
                <td class="dash-table__num">{{ formatNumber(methodTotals.proses) }}</td>
                <td class="dash-table__num">{{ formatNumber(methodTotals.gagalPending) }}</td>
                <td class="dash-table__num">100%</td>
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
          <table class="data-table dash-table">
            <caption class="visually-hidden">
              Rekapitulasi target, berhasil, dan proses per kecamatan beserta persentase.
            </caption>
            <thead>
              <tr>
                <th scope="col">Kecamatan</th>
                <th scope="col" class="dash-table__num">Target</th>
                <th scope="col" class="dash-table__num">Berhasil</th>
                <th scope="col" class="dash-table__num">Proses</th>
                <th scope="col" class="dash-table__num">% Berhasil</th>
                <th scope="col" class="dash-table__num">% Proses</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in data.kecamatanSummary" :key="item.kecamatan">
                <td class="data-table__strong">{{ item.kecamatan }}</td>
                <td class="dash-table__num">{{ formatNumber(item.target) }}</td>
                <td class="dash-table__num">{{ formatNumber(item.berhasil) }}</td>
                <td class="dash-table__num">{{ formatNumber(item.proses) }}</td>
                <td class="dash-table__num">{{ formatPercent(item.persenBerhasil) }}</td>
                <td class="dash-table__num">{{ formatPercent(item.persenProses) }}</td>
              </tr>
              <tr v-if="data.kecamatanSummary.length === 0">
                <td colspan="6" class="table-empty">Belum ada data.</td>
              </tr>
            </tbody>
            <tfoot v-if="data.kecamatanSummary.length > 0">
              <tr>
                <td>TOTAL</td>
                <td class="dash-table__num">{{ formatNumber(kecamatanTotals.target) }}</td>
                <td class="dash-table__num">{{ formatNumber(kecamatanTotals.berhasil) }}</td>
                <td class="dash-table__num">{{ formatNumber(kecamatanTotals.proses) }}</td>
                <td class="dash-table__num">{{ formatPercent(kecamatanTotals.persenBerhasil) }}</td>
                <td class="dash-table__num">{{ formatPercent(kecamatanTotals.persenProses) }}</td>
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
          <table class="data-table dash-table">
            <caption class="visually-hidden">Sepuluh nama dengan jumlah pendaftaran terbanyak.</caption>
            <thead>
              <tr>
                <th scope="col" style="width: 60px">#</th>
                <th scope="col">Nama</th>
                <th scope="col" class="dash-table__num">Jumlah</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in data.top10" :key="item.nama">
                <td>
                  <span class="dash-rank" :class="{ 'dash-rank--top': index < 3 }">{{ index + 1 }}</span>
                </td>
                <td class="data-table__strong">{{ item.nama }}</td>
                <td class="dash-table__num">{{ formatNumber(item.total) }}</td>
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
