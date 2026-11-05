// Shared enumerations and helpers mirroring API_CONTRACT.md.

export const TINGKAT_PENERIMAAN = [
  'Sangat Terbuka',
  'Terbuka',
  'Netral',
  'Kurang Terbuka',
  'Tertutup'
]

export const METODE = ['Halaqoh', 'Program', 'Antara', 'Dakwah Fardiyah']

export const STATUS = ['Berhasil', 'Proses', 'Pending', 'Gagal']

export const UPL_OPTIONS = [
  { value: '', label: 'Tidak ada' },
  { value: 'UPLS', label: 'UPLS' },
  { value: 'UPLM', label: 'UPLM' }
]

export const PROVIDERS = [
  { value: 'local', label: 'Local' },
  { value: 'supabase', label: 'Supabase' },
  { value: 'google_sheets', label: 'Google Sheets' }
]

export function providerLabel(value) {
  const found = PROVIDERS.find((item) => item.value === value)
  return found ? found.label : value
}

// Maps a tingkat_penerimaan value to a badge tone used by the stylesheet.
export function tingkatTone(value) {
  switch (value) {
    case 'Sangat Terbuka':
      return 'tone-success'
    case 'Terbuka':
      return 'tone-info'
    case 'Netral':
      return 'tone-neutral'
    case 'Kurang Terbuka':
      return 'tone-warning'
    case 'Tertutup':
      return 'tone-danger'
    default:
      return 'tone-neutral'
  }
}

export function statusTone(value) {
  switch (value) {
    case 'Berhasil':
      return 'tone-success'
    case 'Proses':
      return 'tone-info'
    case 'Pending':
      return 'tone-warning'
    case 'Gagal':
      return 'tone-danger'
    default:
      return 'tone-neutral'
  }
}
