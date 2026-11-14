// Reading and writing spreadsheet files for the batch upload.
//
// A CSV file is read without any library, so the upload also works without
// internet. An Excel file needs SheetJS, which is loaded from the CDN the
// same way the previous application did it (old-app/index.html:7).

// The columns of the template, the same as the previous application
// (old-app/index.html:992).
export const TEMPLATE_HEADERS = [
  'NAMA TARGET',
  'NO WHATSAPP',
  'ALAMAT',
  'RT',
  'RW',
  'KECAMATAN',
  'KELURAHAN',
  'UPLS',
  'UPLM',
  'AKTIFITAS',
  'TINGKAT PENERIMAAN',
  'METODE',
  'STATUS',
  'KETERANGAN'
]

// The example row of the template.
export const TEMPLATE_SAMPLE = [
  'Contoh Nama Target',
  '081234567890',
  'Contoh alamat',
  '01',
  '02',
  'Pamulang',
  'Pamulang Barat',
  '',
  '',
  'Karyawan Swasta',
  'Terbuka',
  'Dakwah Fardiyah',
  'Proses',
  'Contoh keterangan'
]

// The column of the template and the field of the record it fills.
//
// The columns that carry the owner of the record are not mapped: the server
// fills them from the session.
const COLUMN_FIELD = {
  'NAMA TARGET': 'nama',
  'NO WHATSAPP': 'no_wa',
  ALAMAT: 'alamat',
  RT: 'rt',
  RW: 'rw',
  KECAMATAN: 'kecamatan',
  KELURAHAN: 'kelurahan',
  UPLS: 'upls',
  UPLM: 'uplm',
  AKTIFITAS: 'aktivitas',
  AKTIVITAS: 'aktivitas',
  'TINGKAT PENERIMAAN': 'tingkat_penerimaan',
  METODE: 'metode',
  STATUS: 'status',
  KETERANGAN: 'keterangan'
}

// The SheetJS build that reads an Excel file.
const SHEETJS_URL = 'https://cdn.sheetjs.com/xlsx-0.20.3/package/dist/xlsx.full.min.js'

// sheetJSPromise holds the load of the library, so it is loaded once.
let sheetJSPromise = null

// loadSheetJS loads the library that reads an Excel file.
//
// It is loaded from the CDN, like the previous application did, so the bundle
// of this application stays small.
function loadSheetJS() {
  if (typeof window !== 'undefined' && window.XLSX) {
    return Promise.resolve(window.XLSX)
  }
  if (sheetJSPromise) return sheetJSPromise
  sheetJSPromise = new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = SHEETJS_URL
    script.async = true
    script.onload = () => {
      if (window.XLSX) {
        resolve(window.XLSX)
        return
      }
      reject(new Error('Pustaka pembaca Excel tidak termuat.'))
    }
    script.onerror = () => {
      reject(new Error('Pustaka pembaca Excel gagal dimuat. Simpan file sebagai CSV lalu coba lagi.'))
    }
    document.head.appendChild(script)
  })
  return sheetJSPromise
}

// normalizeHeader folds a header into the shape the mapping uses.
function normalizeHeader(value) {
  return String(value == null ? '' : value)
    .replace(/^﻿/, '')
    .trim()
    .toUpperCase()
    .replace(/\s+/g, ' ')
}

// parseCSV reads a CSV file into its cells.
//
// A cell may be quoted, and a quoted cell may hold a comma.
export function parseCSV(text) {
  const rows = []
  let row = []
  let cell = ''
  let quoted = false
  const source = text.replace(/^﻿/, '')

  for (let i = 0; i < source.length; i++) {
    const char = source[i]
    if (quoted) {
      if (char === '"') {
        if (source[i + 1] === '"') {
          cell += '"'
          i++
        } else {
          quoted = false
        }
      } else {
        cell += char
      }
      continue
    }
    if (char === '"') {
      quoted = true
      continue
    }
    if (char === ',') {
      row.push(cell)
      cell = ''
      continue
    }
    if (char === '\r') continue
    if (char === '\n') {
      row.push(cell)
      rows.push(row)
      row = []
      cell = ''
      continue
    }
    cell += char
  }
  if (cell !== '' || row.length > 0) {
    row.push(cell)
    rows.push(row)
  }
  return rows.filter((r) => r.some((c) => String(c).trim() !== ''))
}

// readExcelFile reads an Excel file into its cells.
async function readExcelFile(file) {
  const XLSX = await loadSheetJS()
  const data = await file.arrayBuffer()
  const wb = XLSX.read(data, { type: 'array' })
  const sheetName = wb.SheetNames.includes('TEMPLATE') ? 'TEMPLATE' : wb.SheetNames[0]
  const ws = wb.Sheets[sheetName]
  const rows = XLSX.utils.sheet_to_json(ws, { header: 1, defval: '', raw: false })
  return rows.filter((r) => r.some((c) => String(c == null ? '' : c).trim() !== ''))
}

// parseSpreadsheet reads a file into rows of records.
//
// The file may be a CSV or an Excel file. The first row holds the headers.
export async function parseSpreadsheet(file) {
  const name = String(file.name || '').toLowerCase()
  let cells
  if (name.endsWith('.csv') || name.endsWith('.txt')) {
    cells = parseCSV(await file.text())
  } else {
    cells = await readExcelFile(file)
  }
  if (!cells.length) {
    throw new Error('File tidak memiliki data.')
  }
  const headers = cells[0].map((h) => normalizeHeader(h))
  return cells.slice(1).map((row) => {
    const record = {}
    headers.forEach((header, index) => {
      const field = COLUMN_FIELD[header]
      if (!field) return
      const value = row[index] == null ? '' : String(row[index]).trim()
      if (value !== '') record[field] = value
    })
    return record
  })
}

// toCSV writes cells as the text of a CSV file.
export function toCSV(cells) {
  return cells
    .map((row) =>
      row
        .map((cell) => {
          const value = cell == null ? '' : String(cell)
          return /[",\n\r]/.test(value) ? `"${value.replace(/"/g, '""')}"` : value
        })
        .join(',')
    )
    .join('\r\n')
}

// downloadTemplate offers the template of the upload as a CSV file.
//
// The file carries a byte order mark, so Excel reads it in UTF-8.
export function downloadTemplate() {
  const text = '﻿' + toCSV([TEMPLATE_HEADERS, TEMPLATE_SAMPLE])
  const blob = new Blob([text], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'Template_Upload_Pendaftaran_TARBIYAH.csv'
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}
