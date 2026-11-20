import { test, expect } from '@playwright/test'
import fs from 'fs'
import os from 'os'
import path from 'path'

// Checks of the batch upload of the member (level 1), against a running
// server that serves the built frontend.
//
// The files are written to a temporary folder, so nothing enters the repo.
const BASE = process.env.E2E_BASE_URL || 'http://localhost:8100'

// A number that changes on every run, so a repeated run does not meet the
// records of the previous one.
const RUN = String(Date.now()).slice(-7)
const wa = (n: number) => `628${RUN}${n}`

const HEADERS =
  'NAMA TARGET,NO WHATSAPP,ALAMAT,RT,RW,KECAMATAN,KELURAHAN,UPLS,UPLM,AKTIFITAS,TINGKAT PENERIMAAN,METODE,STATUS,KETERANGAN'

function csvRow(nama: string, noWA: string, metode = 'Halaqoh') {
  return `${nama},${noWA},Jl. Melati No. 5,003,005,Serpong,Buaran,,,Wiraswasta,Terbuka,${metode},Berhasil,upload`
}

// writeCsv puts a spreadsheet in a temporary folder and returns its path.
function writeCsv(name: string, content: string) {
  const file = path.join(os.tmpdir(), `e2e-${name}-${Date.now()}.csv`)
  fs.writeFileSync(file, content, 'utf8')
  return file
}

async function login(page, username: string, password: string) {
  await page.goto(`${BASE}/login`)
  await page.locator('#login-username').fill(username)
  await page.locator('#login-password').fill(password)
  await page.getByRole('button', { name: 'Masuk' }).click()
  await page.waitForURL((url) => !url.pathname.includes('/login'))
  await page.waitForLoadState('networkidle')
}

test.describe.configure({ mode: 'serial' })

test('Anggota melihat kartu upload dan bisa mengunduh template', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  const card = page.getByText('Upload Data Massal')
  await expect(card).toBeVisible()

  const templateButton = page.getByRole('button', { name: 'Unduh Template' })
  await expect(templateButton).toBeVisible()

  const download = await Promise.all([page.waitForEvent('download'), templateButton.click()])
  const file = await download[0].path()
  const content = fs.readFileSync(file, 'utf8')
  expect(content).toContain('NAMA TARGET')
  expect(content).toContain('NO WHATSAPP')
  expect(download[0].suggestedFilename()).toContain('.csv')
})

test('Anggota penuh (level 3) tidak melihat kartu upload', async ({ page }) => {
  await login(page, 'admin', 'e2e-admin-Pass1')
  await expect(page.getByText('Upload Data Massal')).toHaveCount(0)
})

test('Upload CSV yang valid menyimpan data dan menampilkan ringkasan', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  // The member already has records from the earlier runs, so the count is
  // compared before and after the upload.
  const before = await page.locator('.data-table tbody tr').count()

  const file = writeCsv(
    'valid',
    [HEADERS, csvRow('E2E Valid Satu', `${wa(1)}`), csvRow('E2E Valid Dua', `${wa(2)}`)].join('\n')
  )

  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.alert--success')).toContainText('2 data tersimpan.')
  await expect(page.locator('.data-table tbody tr')).toHaveCount(before + 2)
  await expect(page.locator('.data-table tbody tr td').first()).toHaveText(/E2E Valid Satu|E2E Valid Dua/)
  await expect(page.locator('.data-table tbody')).toContainText('E2E Valid Satu')
  await expect(page.locator('.data-table tbody')).toContainText('E2E Valid Dua')
})

test('Baris yang salah tidak disimpan dan nomor barisnya dilaporkan', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  const before = await page.locator('.data-table tbody tr').count()

  const file = writeCsv(
    'salah',
    [HEADERS, csvRow('E2E Baik', `${wa(3)}`), csvRow('E2E Salah', 'bukan-angka')].join('\n')
  )
  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.alert--error')).toContainText('validasi upload gagal')
  await expect(page.locator('.upload-errors__list li').first()).toContainText('Baris 3')
  await expect(page.locator('.upload-errors__list li').first()).toContainText('no_wa')

  // Nothing is saved, so the table still holds the previous rows.
  await expect(page.locator('.data-table tbody tr')).toHaveCount(before)
})

test('Nomor yang sama di dalam satu file menggagalkan upload', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  const file = writeCsv(
    'ganda',
    [HEADERS, csvRow('E2E Sama A', `${wa(4)}`), csvRow('E2E Sama B', `${wa(4)}`)].join('\n')
  )
  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.alert--error')).toContainText('validasi upload gagal')
  await expect(page.locator('.upload-errors__list li').first()).toContainText('Baris 3')
  await expect(page.locator('.data-table tbody')).not.toContainText(wa(4))
})

test('Nomor yang sudah terdaftar dilewati dan dilaporkan', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  // The first upload stores the number.
  const first = writeCsv(
    'pertama',
    [HEADERS, csvRow('E2E Sudah Ada', `${wa(5)}`)].join('\n')
  )
  await page.locator('input[type=file]').setInputFiles(first)
  await expect(page.locator('.alert--success')).toContainText('1 data tersimpan.')

  // The second upload meets the number again.
  const file = writeCsv(
    'duplikat',
    [HEADERS, csvRow('E2E Sudah Ada', `${wa(5)}`), csvRow('E2E Baru', `${wa(7)}`)].join('\n')
  )
  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.alert--success')).toContainText('1 data tersimpan.')
  await expect(page.locator('.upload-errors__title')).toContainText('1 baris dilewati')
  await expect(page.locator('.upload-errors__list')).toContainText('nomor whatsapp sudah terdaftar')
})

test('Berkas dengan lebih dari 1.000 baris ditolak di sisi klien', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  const rows = [HEADERS]
  for (let i = 0; i < 1001; i += 1) rows.push(csvRow(`E2E Volume ${i}`, `6288${String(i).padStart(7, '0')}`))
  const file = writeCsv('volume', rows.join('\n'))

  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.toast--error')).toContainText('Maksimal 1000 baris per upload')
})

test('Ejaan metode dari aplikasi lama tetap diterima', async ({ page }) => {
  await login(page, 'petugas', 'e2e-petugas-Pass1')

  const file = writeCsv(
    'warisan',
    [HEADERS, csvRow('E2E Ejaan Lama', `${wa(6)}`, 'Dakwah Fardliyah')].join('\n')
  )
  await page.locator('input[type=file]').setInputFiles(file)

  await expect(page.locator('.alert--success')).toContainText('1 data tersimpan.')
  await expect(page.locator('.data-table tbody')).toContainText('Dakwah Fardiyah')
})
