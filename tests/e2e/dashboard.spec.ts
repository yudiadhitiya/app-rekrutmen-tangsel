import { test, expect } from '@playwright/test'

// Checks of the dashboard of the recruitment data, against a running server
// that serves the built frontend.
//
// The accounts used here are test accounts, so the built-in accounts stay
// untouched for the other checks.
const BASE = process.env.E2E_BASE_URL || 'http://localhost:8100'

// A number that changes on every run, so a repeated run does not meet the
// records of the previous one.
const RUN = String(Date.now()).slice(-7)
const wa = (n: number) => `628${RUN}${n}`

async function login(page, username: string, password: string) {
  await page.goto(`${BASE}/login`)
  await page.locator('#login-username').fill(username)
  await page.locator('#login-password').fill(password)
  await page.getByRole('button', { name: 'Masuk' }).click()
  await page.waitForURL((url) => !url.pathname.includes('/login'))
  await page.waitForLoadState('networkidle')
}

// createRecord stores a record with the session of the signed in account.
async function createRecord(
  page,
  noWA: string,
  nama: string,
  kecamatan: string,
  metode: string,
  status: string
) {
  const response = await page.request.post(`${BASE}/api/rekrutmen`, {
    data: {
      nama,
      no_wa: noWA,
      kecamatan,
      kelurahan: '',
      aktivitas: 'Wiraswasta',
      tingkat_penerimaan: 'Terbuka',
      metode,
      status,
      keterangan: 'e2e dashboard'
    }
  })
  expect(response.status()).toBe(201)
}

test.describe.configure({ mode: 'serial' })

test('Dashboard terbuka dari navigasi dan menampilkan semua panel', async ({ page }) => {
  await login(page, 'pwadmin', 'pwadmin123')

  await page.getByRole('link', { name: 'Dashboard' }).click()
  await expect(page).toHaveURL(/\/dashboard/)
  await expect(page.getByRole('heading', { name: 'Dashboard Rekapitulasi' })).toBeVisible()

  // The four cards of the head.
  await expect(page.getByText('Total Target').first()).toBeVisible()
  await expect(page.getByText('Halaqoh').first()).toBeVisible()
  await expect(page.getByText('Prog. Antara').first()).toBeVisible()
  await expect(page.getByText('Dakwah Fardiyah').first()).toBeVisible()

  // The panels of the body.
  await expect(page.getByRole('heading', { name: 'Berdasarkan Tingkat Penerimaan' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Berdasarkan Status' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Ringkasan Metode Pendaftaran' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Summary per Kecamatan' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Top 10 Nama' })).toBeVisible()

  // Every district of Tangsel is listed, so the panel keeps its layout.
  await expect(page.locator('.dash-card')).toHaveCount(4)
})

test('Admin penuh melihat seluruh data pada dashboard', async ({ page }) => {
  await login(page, 'pwadmin', 'pwadmin123')
  await createRecord(page, wa(1), 'Dashboard Admin Satu', 'Serpong', 'Halaqoh', 'Berhasil')

  await page.goto(`${BASE}/dashboard`)
  await page.waitForLoadState('networkidle')

  await expect(page.locator('.dash-scope')).toContainText('Seluruh data pendaftaran')
  await expect(page.locator('.dash-card__number').first()).not.toHaveText('0')
  await expect(page.locator('body')).toContainText('Dashboard Admin Satu')
})

test('Anggota memilih cakupan dashboard sendiri', async ({ page }) => {
  await login(page, 'pwpetugas', 'pwpetugas123')
  await createRecord(page, wa(2), 'Dashboard Anggota Satu', 'Ciputat', 'Program Antara', 'Proses')

  await page.goto(`${BASE}/dashboard`)
  await page.waitForLoadState('networkidle')

  // Only the member may choose the scope of the dashboard.
  const filter = page.getByLabel('Filter dashboard')
  await expect(filter).toBeVisible()
  await expect(filter.locator('option')).toHaveCount(3)

  await expect(page.locator('.dash-scope')).toContainText('Dashboard seluruh kecamatan')

  await filter.selectOption('KECAMATAN')
  await page.waitForLoadState('networkidle')
  await expect(page.locator('.dash-scope')).toContainText('Dashboard KECAMATAN')

  await filter.selectOption('ANGGOTA')
  await page.waitForLoadState('networkidle')
  await expect(page.locator('.dash-scope')).toContainText('Dashboard data anggota Anda')
})

test('Admin kecamatan hanya melihat kecamatannya', async ({ page }) => {
  await login(page, 'pwkec', 'pwkec12345')
  await page.goto(`${BASE}/dashboard`)
  await page.waitForLoadState('networkidle')

  await expect(page.locator('.dash-scope')).toContainText('KECAMATAN Serpong')
  // The scope of the account is fixed, so no filter is offered.
  await expect(page.getByLabel('Filter dashboard')).toHaveCount(0)
})
