const CONFIG = {
  DATA_SHEET: 'PENDAFTARAN',
  USERS_SHEET: 'USERS',
  REGISTRATIONS_SHEET: 'REGISTRATIONS',
  SESSION_TTL: 21600, // 6 hours
  ROLES: {
    ANGGOTA: 'ANGGOTA',
    ADMIN_KECAMATAN: 'ADMIN KECAMATAN',
    ADMIN_DPD: 'ADMIN PTD',
    SYSTEM_ADMIN: 'SYSTEM ADMIN'
  },
  HEADERS: [
    'ID','NAMA ANGGOTA','NAMA TARGET','NO WHATSAPP','NAMA KECAMATAN',
    'ALAMAT','RT','RW','KECAMATAN','KELURAHAN','UPLS','UPLM',
    'AKTIFITAS','TINGKAT PENERIMAAN','METODE','STATUS','KETERANGAN',
    'CREATED_AT','UPDATED_AT','CREATED_BY'
  ],
  USER_HEADERS: [
    'USERNAME','PASSWORD_HASH','ROLE','NAMA ANGGOTA','NAMA KECAMATAN','AKTIF',
    'SECURITY_Q1','SECURITY_A1_HASH','SECURITY_Q2','SECURITY_A2_HASH'
  ],
  REGISTRATION_HEADERS: [
    'ID','USERNAME','PASSWORD_HASH','NAMA ANGGOTA','NAMA KECAMATAN',
    'STATUS','CREATED_AT','UPDATED_AT','KECAMATAN_APPROVED_BY',
    'DPD_APPROVED_BY','REJECTED_BY','REJECTION_REASON'
  ]
};

const TANGSEL_WILAYAH = {
  'Ciputat': ['Cipayung','Ciputat','Jombang','Sawah Baru','Sawah Lama','Serua','Serua Indah'],
  'Ciputat Timur': ['Cempaka Putih','Cireundeu','Pisangan','Pondok Ranji','Rengas','Rempoa'],
  'Pamulang': ['Bambu Apus','Benda Baru','Kedaung','Pamulang Barat','Pamulang Timur','Pondok Benda','Pondok Cabe Ilir','Pondok Cabe Udik'],
  'Pondok Aren': ['Jurang Mangu Barat','Jurang Mangu Timur','Pondok Aren','Pondok Betung','Pondok Jaya','Pondok Kacang Barat','Pondok Kacang Timur','Pondok Karya','Pondok Pucung','Perigi','Perigi Baru'],
  'Serpong': ['Buaran','Ciater','Cilenggang','Lengkong Gudang','Lengkong Gudang Timur','Lengkong Wetan','Rawa Buntu','Rawa Mekar Jaya','Serpong'],
  'Serpong Utara': ['Jelupang','Lengkong Karya','Pakualam','Paku Jaya','Pakulonan','Pondok Jagung','Pondok Jagung Timur'],
  'Setu': ['Babakan','Bakti Jaya','Kademangan','Kranggan','Muncul','Setu']
};

const AKTIFITAS_OPTIONS = [
  'Pegawai Negeri Sipil (PNS) / ASN','Pegawai BUMN / BUMD','Karyawan Swasta','TNI / Polri',
  'Wirawasta / Pengusaha / Pedagang','Profesional','Petani / Peternak / Nelayan',
  'Buruh / Pekerja Harian Lepas','Pekerja Sektor Informal','Pelajar / Mahasiswa',
  'Ibu Rumah Tangga','Pensiunan','Tidak Bekerja / Pengangguran',
  'Pejabat Negara / Anggota Dewan','Lainnya'
];

function getWilayahTangsel() { return TANGSEL_WILAYAH; }


function doGet() {
  return HtmlService.createHtmlOutputFromFile('Index')
    .setTitle('Manajemen Pendaftaran TARBIYAH')
    .setXFrameOptionsMode(HtmlService.XFrameOptionsMode.ALLOWALL);
}

function setupSheets() {
  const ss = SpreadsheetApp.getActive();
  let data = ss.getSheetByName(CONFIG.DATA_SHEET);
  if (!data) data = ss.insertSheet(CONFIG.DATA_SHEET);

  const oldHeaders = data.getLastRow() ? data.getRange(1,1,1,Math.max(data.getLastColumn(),1)).getValues()[0].map(String) : [];
  const isLegacy = oldHeaders.indexOf('CREATED_BY') >= 0 && oldHeaders.indexOf('KECAMATAN') < 0;

  if (isLegacy && data.getLastRow() > 1) {
    const oldRows = data.getRange(2,1,data.getLastRow()-1,Math.min(data.getLastColumn(),14)).getValues();
    const migrated = oldRows.map(r => [
      r[0]||'',r[1]||'',r[2]||'',r[3]||'',r[4]||'',r[5]||'',
      '','','','','','',
      r[6]||'',r[7]||'',r[8]||'',r[9]||'',r[10]||'',r[11]||'',r[12]||'',r[13]||''
    ]);
    data.clearContents();
    data.getRange(1,1,1,CONFIG.HEADERS.length).setValues([CONFIG.HEADERS]);
    if(migrated.length) data.getRange(2,1,migrated.length,CONFIG.HEADERS.length).setValues(migrated);
  } else {
    CONFIG.HEADERS.forEach((header,i)=>data.getRange(1,i+1).setValue(header));
  }

  let users=ss.getSheetByName(CONFIG.USERS_SHEET);
  if(!users) users=ss.insertSheet(CONFIG.USERS_SHEET);
  CONFIG.USER_HEADERS.forEach((header,i)=>users.getRange(1,i+1).setValue(header));

  let registrations=ss.getSheetByName(CONFIG.REGISTRATIONS_SHEET);
  if(!registrations) registrations=ss.insertSheet(CONFIG.REGISTRATIONS_SHEET);
  CONFIG.REGISTRATION_HEADERS.forEach((header,i)=>registrations.getRange(1,i+1).setValue(header));

  return 'Setup selesai. Struktur PENDAFTARAN sudah disiapkan untuk versi baru.';
}

/**
 * Backward compatibility:
 * role lama "TARBIYAH KECAMATAN" dianggap sebagai "ADMIN KECAMATAN".
 */
function normalizeRole_(role) {
  const r = String(role || '').trim().toUpperCase();
  if (r === 'TARBIYAH KECAMATAN' || r === 'ADMIN KECAMATAN' || r === 'ADMIN_KECAMATAN') return CONFIG.ROLES.ADMIN_KECAMATAN;
  if (r === 'ADMIN PTD' || r === 'ADMIN_DPD') return CONFIG.ROLES.ADMIN_DPD;
  if (r === 'SYSTEM ADMIN' || r === 'SYSTEM_ADMIN') return CONFIG.ROLES.SYSTEM_ADMIN;
  if (r === 'ANGGOTA') return CONFIG.ROLES.ANGGOTA;
  return r;
}

function createUser(username, password, role, namaAnggota, namaKecamatan) {
  if (!username || !password || !role) {
    throw new Error('Username, password, dan role wajib diisi.');
  }

  const normalizedRole = normalizeRole_(role);
  const validRoles = [
    CONFIG.ROLES.ANGGOTA,
    CONFIG.ROLES.ADMIN_KECAMATAN,
    CONFIG.ROLES.ADMIN_DPD,
    CONFIG.ROLES.SYSTEM_ADMIN
  ];
  if (validRoles.indexOf(normalizedRole) === -1) {
    throw new Error('Role tidak valid: ' + role);
  }

  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!sh) throw new Error('Jalankan setupSheets() terlebih dahulu.');

  const values = sh.getDataRange().getValues();
  const exists = values.slice(1).some(r =>
    String(r[0]).trim().toLowerCase() === String(username).trim().toLowerCase()
  );
  if (exists) throw new Error('Username sudah digunakan.');

  sh.appendRow([
    String(username).trim(),
    hashPassword_(password),
    normalizedRole,
    namaAnggota || '',
    namaKecamatan || '',
    'Y'
  ]);

  return 'User berhasil dibuat.';
}


/**
 * Public - daftar KECAMATAN yang dapat dipilih pada form Self-Registration.
 * Hanya nama KECAMATAN dari Admin KECAMATAN aktif yang ditampilkan.
 */
function getRegistrationKecamatans() {
  return Object.keys(TANGSEL_WILAYAH).sort();
}

/**
 * Public - self-registration anggota.
 * Akun langsung aktif setelah self-registration dan dapat digunakan untuk login.
 */
function selfRegister(username, password, namaAnggota, namaKecamatan, q1, a1, q2, a2) {
  username = String(username || '').trim();
  password = String(password || '');
  namaAnggota = String(namaAnggota || '').trim();
  namaKecamatan = String(namaKecamatan || '').trim();
  q1=String(q1||'').trim(); a1=String(a1||'').trim(); q2=String(q2||'').trim(); a2=String(a2||'').trim();

  if (!username || !password || !namaAnggota || !namaKecamatan || !q1 || !a1 || !q2 || !a2) throw new Error('Semua data dan dua pertanyaan-jawaban keamanan wajib diisi.');
  if (q1 === q2) throw new Error('Pilih dua pertanyaan keamanan yang berbeda.');
  if (username.length < 4) throw new Error('Username minimal 4 karakter.');
  if (password.length < 6) throw new Error('Password minimal 6 karakter.');

  const allowedKecamatans = getRegistrationKecamatans();
  if (!allowedKecamatans.some(d => normalize_(d) === normalize_(namaKecamatan))) 
    throw new Error('KECAMATAN tidak valid. Silakan pilih Kecamatan dari daftar yang tersedia.');

  const users = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!users) throw new Error('Sheet USERS belum tersedia. Jalankan setupSheets() terlebih dahulu.');
  const userRows = users.getDataRange().getValues();
  if (userRows.slice(1).some(r => normalize_(r[0]) === normalize_(username))) throw new Error('Username sudah digunakan.');

  ensureSecurityColumns_();
  users.appendRow([username, hashPassword_(password), CONFIG.ROLES.ANGGOTA, namaAnggota, namaKecamatan, 'Y', q1, hashPassword_(normalizeAnswer_(a1)), q2, hashPassword_(normalizeAnswer_(a2))]);
  return {status:'ACTIVE', message:'Registrasi berhasil. Akun langsung aktif dan dapat digunakan untuk login.'};
}

/**
 * Fungsi approval lama tidak lagi digunakan; self-registration langsung aktif.
 */
function formatDateSafe_(value) {
  if (!value) return '';
  if (value instanceof Date) {
    return Utilities.formatDate(
      value,
      Session.getScriptTimeZone(),
      'dd/MM/yyyy HH:mm'
    );
  }
  return String(value);
}

function login(username, password) {
  ensureSecurityColumns_();
  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!sh) throw new Error('Sheet USERS belum tersedia.');

  const rows = sh.getDataRange().getValues();
  const hash = hashPassword_(password);

  for (let i = 1; i < rows.length; i++) {
    const [u, ph, role, anggota, kecamatan, active] = rows[i];

    if (
      String(u).trim().toLowerCase() === String(username).trim().toLowerCase() &&
      String(ph) === hash &&
      String(active).trim().toUpperCase() === 'Y'
    ) {
      const token = Utilities.getUuid();
      const session = {
        username: String(u),
        role: normalizeRole_(role),
        anggota: String(anggota || ''),
        kecamatan: String(kecamatan || '')
      };

      CacheService.getScriptCache().put(
        'sess_' + token,
        JSON.stringify(session),
        CONFIG.SESSION_TTL
      );

      return { token, user: session, requiresSecuritySetup: !(String(rows[i][6]||'').trim() && String(rows[i][7]||'').trim() && String(rows[i][8]||'').trim() && String(rows[i][9]||'').trim()) };
    }
  }

  throw new Error('Username atau password salah, atau akun tidak aktif.');
}

/** Tambahkan kolom keamanan ke sheet USERS tanpa mengubah kolom lama. */
function ensureSecurityColumns_() {
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if(!sh) throw new Error('Sheet USERS belum tersedia.');
  const headers=CONFIG.USER_HEADERS;
  headers.forEach((h,i)=>{ if(String(sh.getRange(1,i+1).getValue()||'').trim()!==h) sh.getRange(1,i+1).setValue(h); });
}
function normalizeAnswer_(value){ return String(value||'').trim().toLowerCase().replace(/\s+/g,' '); }
function getSecurityQuestions(username){
  ensureSecurityColumns_(); username=String(username||'').trim();
  if(!username) throw new Error('Masukkan username terlebih dahulu.');
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(username)){
    if(String(values[i][5]||'').toUpperCase()!=='Y') throw new Error('Antum tidak punya akses ke Aplikasi. Hubungi System Admin untuk konfirmasi');
    const q1=String(values[i][6]||''), q2=String(values[i][8]||'');
    if(!q1||!q2||!values[i][7]||!values[i][9]) return {requiresSetup:true};
    return {requiresSetup:false, questions:[q1,q2]};
  }
  throw new Error('Username tidak ditemukan.');
}
function hasSecurityQuestions(token){
  const me=getSession_(token); ensureSecurityColumns_();
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(me.username)) return !!(values[i][6]&&values[i][7]&&values[i][8]&&values[i][9]);
  throw new Error('User tidak ditemukan.');
}
function saveMySecurityQuestions(token,q1,a1,q2,a2){
  const me=getSession_(token); q1=String(q1||'').trim(); q2=String(q2||'').trim(); a1=String(a1||'').trim(); a2=String(a2||'').trim();
  if(!q1||!q2||!a1||!a2) throw new Error('Dua pertanyaan dan jawaban wajib diisi.');
  if(q1===q2) throw new Error('Pilih dua pertanyaan yang berbeda.');
  ensureSecurityColumns_(); const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(me.username)){
    sh.getRange(i+1,7,1,4).setValues([[q1,hashPassword_(normalizeAnswer_(a1)),q2,hashPassword_(normalizeAnswer_(a2))]]);
    return true;
  }
  throw new Error('User tidak ditemukan.');
}
function verifySecurityAnswers(username,a1,a2){
  ensureSecurityColumns_(); username=String(username||'').trim();
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(username)){
    const r=values[i];
    if(String(r[5]||'').toUpperCase()!=='Y') throw new Error('Antum tidak punya akses ke Aplikasi. Hubungi System Admin untuk konfirmasi');
    if(!r[6]||!r[7]||!r[8]||!r[9]) throw new Error('Pertanyaan keamanan belum lengkap. Silakan login dan lengkapi pertanyaan keamanan.');
    if(hashPassword_(normalizeAnswer_(a1))!==String(r[7]) || hashPassword_(normalizeAnswer_(a2))!==String(r[9])) throw new Error('Jawaban pertanyaan keamanan salah.');
    const recoveryToken=Utilities.getUuid();
    CacheService.getScriptCache().put('recover_'+recoveryToken, String(r[0]), 600);
    return {recoveryToken:recoveryToken};
  }
  throw new Error('Antum tidak punya akses ke Aplikasi. Hubungi System Admin untuk konfirmasi');
}
function changePasswordWithRecovery(recoveryToken,newPassword){
  newPassword=String(newPassword||''); if(newPassword.length<6) throw new Error('Password baru minimal 6 karakter.');
  const cache=CacheService.getScriptCache(), key='recover_'+String(recoveryToken||''), username=cache.get(key);
  if(!username) throw new Error('Sesi pemulihan sudah berakhir. Silakan ulangi proses Lupa Password.');
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(username)){
    if(String(values[i][5]||'').toUpperCase()!=='Y') throw new Error('Antum tidak punya akses ke Aplikasi. Hubungi System Admin untuk konfirmasi');
    sh.getRange(i+1,2).setValue(hashPassword_(newPassword)); cache.remove(key); return true;
  }
  throw new Error('User tidak ditemukan.');
}

function logout(token) {
  if (token) CacheService.getScriptCache().remove('sess_' + token);
  return true;
}

function getMe(token) {
  return getSession_(token);
}

function getSession_(token) {
  if (!token) throw new Error('Sesi tidak ditemukan. Silakan login kembali.');

  const raw = CacheService.getScriptCache().get('sess_' + token);
  if (!raw) throw new Error('Sesi sudah berakhir. Silakan login kembali.');

  return JSON.parse(raw);
}

function requireRole_(me, allowedRoles) {
  if (!allowedRoles.includes(me.role)) {
    throw new Error('Anda tidak memiliki hak akses untuk tindakan ini.');
  }
}

/**
 * RBAC data:
 * - ANGGOTA   : hanya data yang dibuat/diinput oleh username sendiri.
 * - ADMIN KECAMATAN : hanya data dengan NAMA KECAMATAN yang sama dengan KECAMATAN user.
 * - ADMIN DPD : semua data.
 * - SYSTEM ADMIN: semua data (read-only untuk data pendaftaran).
 */
function getRecruitments(token) {
  // Backward-compatible: ambil seluruh data sesuai hak akses.
  const me = getSession_(token);
  return getRecruitmentsPaged(token, 1, 1000000, ' ').rows;
}

/**
 * Mengambil data pendaftaran secara server-side dengan pagination dan search.
 * RBAC tetap diterapkan SEBELUM data dikirim ke browser.
 *
 * page     : nomor halaman mulai 1
 * pageSize : jumlah data per halaman (10/25/50/100)
 * keyword  : pencarian pada field data pendaftaran
 */
function getRecruitmentsPaged(token, page, pageSize, keyword) {
  const me = getSession_(token);

  const sh = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(CONFIG.DATA_SHEET);
  if (!sh) throw new Error('Sheet PENDAFTARAN tidak ditemukan.');

  page = Math.max(1, parseInt(page, 10) || 1);
  pageSize = parseInt(pageSize, 10) || 10;
  if ([10, 25, 50, 100].indexOf(pageSize) === -1) pageSize = 10;

  keyword = String(keyword || '').trim().toLowerCase();

  const lastRow = sh.getLastRow();
  const lastCol = Math.max(sh.getLastColumn(), CONFIG.HEADERS.length);

  if (lastRow < 2) {
    return { rows: [], total: 0, page: 1, pageSize: pageSize, totalPages: 1 };
  }

  const values = sh.getRange(1, 1, lastRow, lastCol).getValues();
  const headers = values[0].map(h => String(h).trim());
  const matched = [];

  const myUsername = normalize_(me.username);
  const myKecamatan = normalize_(me.kecamatan);

  for (let i = 1; i < values.length; i++) {
    const row = values[i];

    const isEmptyRow = row.every(value =>
      String(value == null ? '' : value).trim() === ''
    );
    if (isEmptyRow) continue;

    const namaTarget = String(row[2] == null ? '' : row[2]).trim();
    if (!namaTarget) continue;

    const createdBy = normalize_(row[19]);
    const rowKecamatan = normalize_(row[4]);

    let allowed = false;

    switch (me.role) {
      case CONFIG.ROLES.ANGGOTA:
        allowed = createdBy === myUsername;
        break;
      case CONFIG.ROLES.ADMIN_KECAMATAN:
        allowed = rowKecamatan === myKecamatan;
        break;
      case CONFIG.ROLES.ADMIN_DPD:
      case CONFIG.ROLES.SYSTEM_ADMIN:
        allowed = true;
        break;
      default:
        allowed = false;
    }

    if (!allowed) continue;

    // Search server-side agar browser tidak perlu menerima seluruh data.
    if (keyword) {
      const searchable = [
        row[1], row[2], row[3], row[4], row[5],
        row[6], row[7], row[8], row[9], row[10], row[11], row[12], row[13], row[14], row[15], row[16]
      ].map(v => String(v == null ? '' : v)).join(' ').toLowerCase();

      if (searchable.indexOf(keyword) === -1) continue;
    }

    // Pastikan ID tersedia.
    if (!row[0]) {
      const newId = Utilities.getUuid();
      sh.getRange(i + 1, 1).setValue(newId);
      row[0] = newId;
    }

    const obj = {};
    headers.forEach((header, j) => {
      let value = row[j];
      if (value instanceof Date) {
        value = Utilities.formatDate(
          value,
          Session.getScriptTimeZone(),
          'dd/MM/yyyy HH:mm:ss'
        );
      }
      obj[header] = value;
    });

    matched.push(obj);
  }

  const total = matched.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  if (page > totalPages) page = totalPages;

  const startIndex = (page - 1) * pageSize;
  const pagedRows = matched.slice(startIndex, startIndex + pageSize);

  return {
    rows: pagedRows,
    total: total,
    page: page,
    pageSize: pageSize,
    totalPages: totalPages
  };
}


/**
 * Dashboard rekapitulasi pendaftaran.
 * Data yang dihitung sudah dibatasi sesuai RBAC user:
 * - ANGGOTA      : data yang dibuat username sendiri
 * - ADMIN KECAMATAN    : data pada KECAMATAN user
 * - ADMIN DPD    : seluruh data
 * - SYSTEM ADMIN : seluruh data
 */
function getDashboard(token, dashboardFilter) {
  const me = getSession_(token);
  const sh = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(CONFIG.DATA_SHEET);
  if (!sh) throw new Error('Sheet PENDAFTARAN tidak ditemukan.');

  // Filter tambahan hanya berlaku untuk dashboard ANGGOTA.
  // Dashboard hanya mengembalikan agregat, bukan baris data pendaftaran.
  let selectedFilter = '';
  if (me.role === CONFIG.ROLES.ANGGOTA) {
    selectedFilter = String(dashboardFilter || 'ALL').trim().toUpperCase();
    if (!['ANGGOTA', 'KECAMATAN', 'ALL'].includes(selectedFilter)) selectedFilter = 'ALL';
  }

  const lastRow = sh.getLastRow();
  if (lastRow < 2) {
    return buildDashboardResult_(me, [], selectedFilter);
  }

  const lastCol = Math.max(sh.getLastColumn(), CONFIG.HEADERS.length);
  const values = sh.getRange(1, 1, lastRow, lastCol).getValues();
  const rows = [];

  const myUsername = normalize_(me.username);
  const myKecamatan = normalize_(me.kecamatan);

  for (let i = 1; i < values.length; i++) {
    const row = values[i];
    const isEmptyRow = row.every(value => String(value == null ? '' : value).trim() === '');
    if (isEmptyRow) continue;

    const namaTarget = String(row[2] == null ? '' : row[2]).trim();
    if (!namaTarget) continue;

    const createdBy = normalize_(row[19]);
    const rowKecamatan = normalize_(row[4]);
    let allowed = false;

    switch (me.role) {
      case CONFIG.ROLES.ANGGOTA:
        if (selectedFilter === 'ANGGOTA') {
          allowed = createdBy === myUsername;
        } else if (selectedFilter === 'KECAMATAN') {
          allowed = rowKecamatan === myKecamatan;
        } else {
          // ALL: seluruh kecamatan, khusus agregasi dashboard.
          allowed = true;
        }
        break;
      case CONFIG.ROLES.ADMIN_KECAMATAN:
        allowed = rowKecamatan === myKecamatan;
        break;
      case CONFIG.ROLES.ADMIN_DPD:
      case CONFIG.ROLES.SYSTEM_ADMIN:
        allowed = true;
        break;
      default:
        allowed = false;
    }

    if (allowed) rows.push(row);
  }

  return buildDashboardResult_(me, rows, selectedFilter);
}

function buildDashboardResult_(me, rows, dashboardFilter) {
  const total = rows.length;

  // Metode: urutan utama dibuat konsisten dengan pilihan pada form.
  const methodOrder = ['Halaqah', 'Program Antara', 'Dakwah Fardliyah'];
  const methodMap = {};
  methodOrder.forEach(name => methodMap[name] = 0);
  let methodOther = 0;

  rows.forEach(row => {
    const method = String(row[14] == null ? '' : row[14]).trim();
    if (methodOrder.indexOf(method) >= 0) methodMap[method]++;
    else methodOther++;
  });

  const methods = methodOrder.map(name => ({
    name: name,
    total: methodMap[name],
    percent: total ? roundPercent_(methodMap[name], total) : 0
  }));

  // Ringkasan metode x status untuk dashboard referensi.
  // Field lama (methods) tetap dipertahankan agar tidak memengaruhi fungsi existing.
  const methodSummaryMap = {};
  methodOrder.forEach(name => methodSummaryMap[name] = {
    total: 0,
    berhasil: 0,
    proses: 0,
    gagalPending: 0
  });
  let methodSummaryOther = { total: 0, berhasil: 0, proses: 0, gagalPending: 0 };

  rows.forEach(row => {
    const method = String(row[14] == null ? '' : row[14]).trim();
    const statusValue = String(row[15] == null ? '' : row[15]).trim();
    const bucket = methodSummaryMap[method] || methodSummaryOther;
    bucket.total++;
    if (statusValue === 'Berhasil') bucket.berhasil++;
    else if (statusValue === 'Proses') bucket.proses++;
    else bucket.gagalPending++;
  });

  const methodSummary = methodOrder.map(name => ({
    name: name,
    total: methodSummaryMap[name].total,
    berhasil: methodSummaryMap[name].berhasil,
    proses: methodSummaryMap[name].proses,
    gagalPending: methodSummaryMap[name].gagalPending,
    percent: total ? roundPercent_(methodSummaryMap[name].total, total) : 0
  }));
  if (methodSummaryOther.total) {
    methodSummary.push({
      name: 'Lainnya / Belum Diisi',
      total: methodSummaryOther.total,
      berhasil: methodSummaryOther.berhasil,
      proses: methodSummaryOther.proses,
      gagalPending: methodSummaryOther.gagalPending,
      percent: total ? roundPercent_(methodSummaryOther.total, total) : 0
    });
  }
  if (methodOther) {
    methods.push({
      name: 'Lainnya / Belum Diisi',
      total: methodOther,
      percent: total ? roundPercent_(methodOther, total) : 0
    });
  }

  // Tingkat penerimaan mengikuti nilai pada form.
  const acceptanceOrder = ['Sangat Terbuka', 'Terbuka', 'Netral', 'Kurang Terbuka', 'Tertutup'];
  const acceptanceMap = {};
  acceptanceOrder.forEach(name => acceptanceMap[name] = 0);
  let acceptanceOther = 0;

  rows.forEach(row => {
    const value = String(row[13] == null ? '' : row[13]).trim();
    if (acceptanceOrder.indexOf(value) >= 0) acceptanceMap[value]++;
    else acceptanceOther++;
  });

  const acceptance = acceptanceOrder.map(name => ({
    name: name,
    total: acceptanceMap[name],
    percent: total ? roundPercent_(acceptanceMap[name], total) : 0
  }));
  if (acceptanceOther) {
    acceptance.push({
      name: 'Belum Diisi / Lainnya',
      total: acceptanceOther,
      percent: total ? roundPercent_(acceptanceOther, total) : 0
    });
  }

  // Status aplikasi disederhanakan menjadi 4 kelompok agar mirip dashboard referensi.
  // Aktif = berhasil; tiga status pendekatan = proses; Tidak Lanjut = gagal.
  const status = { berhasil: 0, proses: 0, pending: 0, gagal: 0 };
  const statusDetail = {};
  rows.forEach(row => {
    const value = String(row[15] == null ? '' : row[15]).trim();
    if (value) statusDetail[value] = (statusDetail[value] || 0) + 1;

    if (value === 'Berhasil') status.berhasil++;
    else if (value === 'Gagal') status.gagal++;
    else if (value === 'Proses') status.proses++;
    else status.pending++;
  });

  const topMap={};
  rows.forEach(row=>{ const n=String(row[1]||'').trim()||'(Tanpa Nama)'; topMap[n]=(topMap[n]||0)+1; });
  const top10=Object.keys(topMap).map(n=>({nama:n,total:topMap[n]})).sort((a,b)=>b.total-a.total||a.nama.localeCompare(b.nama)).slice(0,10);

  return {
    role: me.role,
    kecamatanSummary: buildKecamatanSummary_(rows),
    scope: me.role === CONFIG.ROLES.ANGGOTA
      ? (dashboardFilter === 'ANGGOTA'
          ? 'Dashboard data anggota Anda'
          : dashboardFilter === 'KECAMATAN'
            ? 'Dashboard KECAMATAN ' + (me.kecamatan || '-')
            : 'Dashboard seluruh kecamatan')
      : me.role === CONFIG.ROLES.ADMIN_KECAMATAN
        ? 'Data KECAMATAN ' + (me.kecamatan || '-')
        : 'Seluruh data pendaftaran',
    total: total,
    methods: methods,
    methodSummary: methodSummary,
    acceptance: acceptance,
    status: [
      { name: 'Berhasil', total: status.berhasil, percent: total ? roundPercent_(status.berhasil, total) : 0 },
      { name: 'Proses', total: status.proses, percent: total ? roundPercent_(status.proses, total) : 0 },
      { name: 'Pending', total: status.pending, percent: total ? roundPercent_(status.pending, total) : 0 },
      { name: 'Gagal', total: status.gagal, percent: total ? roundPercent_(status.gagal, total) : 0 }
    ],
    statusDetail: Object.keys(statusDetail).map(name => ({
      name: name,
      total: statusDetail[name],
      percent: total ? roundPercent_(statusDetail[name], total) : 0
    })),
    top10: top10
  };
}


/**
 * Summary per KECAMATAN untuk dashboard.
 * Target = jumlah target/data pada KECAMATAN tersebut.
 * Berhasil = STATUS "Berhasil".
 * Proses = STATUS "Proses" + "Pending".
 * Persentase dihitung terhadap Target.
 * Semua KECAMATAN Tangsel tetap ditampilkan agar layout dashboard konsisten.
 */
function buildKecamatanSummary_(rows) {
  const order = Object.keys(TANGSEL_WILAYAH);
  const map = {};
  order.forEach(kecamatan => {
    map[normalize_(kecamatan)] = {
      kecamatan: kecamatan,
      target: 0,
      berhasil: 0,
      proses: 0
    };
  });

  (rows || []).forEach(row => {
    const kecamatan = String(row[4] == null ? '' : row[4]).trim();
    const key = normalize_(kecamatan);

    // Hanya hitung KECAMATAN yang terdaftar di master wilayah.
    if (!map[key]) return;

    map[key].target++;

    const status = String(row[15] == null ? '' : row[15]).trim();
    if (status === 'Berhasil') {
      map[key].berhasil++;
    } else if (status === 'Proses' || status === 'Pending') {
      map[key].proses++;
    }
  });

  return order.map(kecamatan => {
    const item = map[normalize_(kecamatan)];
    return {
      kecamatan: item.kecamatan,
      target: item.target,
      berhasil: item.berhasil,
      proses: item.proses,
      persenBerhasil: item.target ? roundPercent_(item.berhasil, item.target) : 0,
      persenProses: item.target ? roundPercent_(item.proses, item.target) : 0
    };
  });
}

function roundPercent_(value, total) {
  if (!total) return 0;
  return Math.round((value / total) * 1000) / 10;
}

function normalize_(value) {
  return String(value || '')
    .replace(/\u00A0/g, ' ')
    .trim()
    .toLowerCase();
}


/**
 * Upload bulk data dari template Excel.
 * Hanya ANGGOTA yang boleh melakukan input bulk.
 * NAMA ANGGOTA, NAMA KECAMATAN (pemilik), ID, timestamp dan CREATED_BY
 * selalu diisi oleh server dari session agar tidak dapat dipalsukan dari Excel.
 */
function uploadRecruitments(token, items) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.ANGGOTA]);

  if (!Array.isArray(items) || !items.length) {
    throw new Error('Tidak ada data yang dapat di-upload.');
  }
  if (items.length > 1000) {
    throw new Error('Maksimal 1.000 baris per upload.');
  }

  const prepared = [];
  const errors = [];

  items.forEach((raw, index) => {
    const rowNumber = index + 2;
    const p = {
      'NAMA TARGET': String(raw['NAMA TARGET'] || '').trim(),
      'NO WHATSAPP': String(raw['NO WHATSAPP'] || '').trim(),
      'ALAMAT': String(raw['ALAMAT'] || '').trim(),
      'RT': String(raw['RT'] || '').trim(),
      'RW': String(raw['RW'] || '').trim(),
      'KECAMATAN': String(raw['KECAMATAN'] || '').trim(),
      'KELURAHAN': String(raw['KELURAHAN'] || '').trim(),
      'UPLS': String(raw['UPLS'] || '').trim(),
      'UPLM': String(raw['UPLM'] || '').trim(),
      'AKTIFITAS': String(raw['AKTIFITAS'] || '').trim(),
      'TINGKAT PENERIMAAN': String(raw['TINGKAT PENERIMAAN'] || '').trim(),
      'METODE': String(raw['METODE'] || '').trim(),
      'STATUS': String(raw['STATUS'] || '').trim(),
      'KETERANGAN': String(raw['KETERANGAN'] || '').trim()
    };

    try {
      validatePayload_(p);
      prepared.push(p);
    } catch (e) {
      errors.push('Baris ' + rowNumber + ': ' + e.message);
    }
  });

  if (errors.length) {
    throw new Error(
      'Validasi upload gagal. Tidak ada data yang disimpan.\n\n' +
      errors.slice(0, 30).join('\n') +
      (errors.length > 30 ? '\n... dan ' + (errors.length - 30) + ' error lainnya.' : '')
    );
  }

  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.DATA_SHEET);
  if (!sh) throw new Error('Sheet PENDAFTARAN tidak ditemukan.');

  const lock = LockService.getScriptLock();
  lock.waitLock(30000);

  try {
    const now = new Date();
    const values = prepared.map(p => [
      Utilities.getUuid(),
      me.anggota,
      p['NAMA TARGET'],
      p['NO WHATSAPP'],
      me.kecamatan,
      p['ALAMAT'],
      p['RT'],
      p['RW'],
      p['KECAMATAN'],
      p['KELURAHAN'],
      p['UPLS'],
      p['UPLM'],
      p['AKTIFITAS'],
      p['TINGKAT PENERIMAAN'],
      p['METODE'],
      p['STATUS'],
      p['KETERANGAN'],
      now,
      now,
      me.username
    ]);

    const startRow = sh.getLastRow() + 1;
    sh.getRange(startRow, 1, values.length, CONFIG.HEADERS.length).setValues(values);
    return {
      count: values.length,
      message: values.length + ' data berhasil ditambahkan.'
    };
  } finally {
    lock.releaseLock();
  }
}

/**
 * Hanya ANGGOTA yang dapat input data baru.
 * Data pemilik otomatis berasal dari session, bukan dari client.
 */
function saveRecruitment(token, payload) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.ANGGOTA]);
  validatePayload_(payload);
  payload['NAMA ANGGOTA'] = me.anggota;
  payload['NAMA KECAMATAN'] = me.kecamatan;
  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.DATA_SHEET);
  if (!sh) throw new Error('Sheet PENDAFTARAN tidak ditemukan.');
  const now = new Date(), id = Utilities.getUuid();
  sh.appendRow([
    id,payload['NAMA ANGGOTA'],payload['NAMA TARGET'],payload['NO WHATSAPP'],payload['NAMA KECAMATAN'],
    payload['ALAMAT']||'',payload['RT']||'',payload['RW']||'',payload['KECAMATAN'],payload['KELURAHAN'],
    payload['UPLS']||'',payload['UPLM']||'',payload['AKTIFITAS'],payload['TINGKAT PENERIMAAN'],
    payload['METODE'],payload['STATUS'],payload['KETERANGAN']||'',now,now,me.username
  ]);
  return id;
}

function updateRecruitment(token, payload) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.ANGGOTA, CONFIG.ROLES.ADMIN_KECAMATAN]);
  validatePayload_(payload);
  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.DATA_SHEET);
  const data = sh.getDataRange().getValues();
  const id = String(payload.ID || '');

  for (let i=1;i<data.length;i++) {
    if (String(data[i][0]) !== id) continue;
    const owner=normalize_(data[i][19]), rowKecamatan=normalize_(data[i][4]);
    if (me.role===CONFIG.ROLES.ANGGOTA && owner!==normalize_(me.username)) throw new Error('Anda hanya dapat mengubah data yang Anda input sendiri.');
    if (me.role===CONFIG.ROLES.ADMIN_KECAMATAN && rowKecamatan!==normalize_(me.kecamatan)) throw new Error('Anda hanya dapat mengubah data pada KECAMATAN Anda.');

    const row=i+1;
    sh.getRange(row,2,1,15).setValues([[
      String(data[i][1]||''),payload['NAMA TARGET'],payload['NO WHATSAPP'],String(data[i][4]||''),
      payload['ALAMAT']||'',payload['RT']||'',payload['RW']||'',payload['KECAMATAN'],payload['KELURAHAN'],
      payload['UPLS']||'',payload['UPLM']||'',payload['AKTIFITAS'],payload['TINGKAT PENERIMAAN'],
      payload['METODE'],payload['STATUS']
    ]]);
    sh.getRange(row,17).setValue(payload['KETERANGAN']||'');
    sh.getRange(row,19).setValue(new Date());
    return true;
  }
  throw new Error('Data tidak ditemukan.');
}

function deleteRecruitment(token, id) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.ANGGOTA]);
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.DATA_SHEET), data=sh.getDataRange().getValues();
  for(let i=1;i<data.length;i++) if(String(data[i][0])===String(id)){
    if(normalize_(data[i][19])!==normalize_(me.username)) throw new Error('Anda hanya dapat menghapus data yang Anda input sendiri.');
    sh.deleteRow(i+1); return true;
  }
  throw new Error('Data tidak ditemukan.');
}

function getRecruitmentsForExport(token) {
  const me=getSession_(token), sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.DATA_SHEET);
  if(!sh) throw new Error('Sheet PENDAFTARAN tidak ditemukan.');
  const values=sh.getDataRange().getValues(), headers=values.length?values[0].map(String):CONFIG.HEADERS.slice(), rows=[];
  for(let i=1;i<values.length;i++){
    const row=values[i]; if(!row[0]) continue;
    const owner=normalize_(row[19]), rowKecamatan=normalize_(row[4]);
    const allowed=me.role===CONFIG.ROLES.ANGGOTA?owner===normalize_(me.username):me.role===CONFIG.ROLES.ADMIN_KECAMATAN?rowKecamatan===normalize_(me.kecamatan):[CONFIG.ROLES.ADMIN_DPD,CONFIG.ROLES.SYSTEM_ADMIN].includes(me.role);
    if(allowed) rows.push(row.map(v=>v instanceof Date?formatDateSafe_(v):(v==null?'':v)));
  }
  return {headers,rows,filename:'Pendaftaran_TARBIYAH_'+Utilities.formatDate(new Date(),Session.getScriptTimeZone(),'yyyyMMdd_HHmmss')+'.xls'};
}

/**
 * SYSTEM ADMIN - daftar user.
 * Password hash tidak pernah dikirim ke browser.
 */
function getUsers(token) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.SYSTEM_ADMIN]);

  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!sh) throw new Error('Sheet USERS belum tersedia.');

  const values = sh.getDataRange().getValues();
  const result = [];

  for (let i = 1; i < values.length; i++) {
    const r = values[i];
    if (!String(r[0] || '').trim()) continue;

    result.push({
      username: String(r[0]),
      role: normalizeRole_(r[2]),
      anggota: String(r[3] || ''),
      kecamatan: String(r[4] || ''),
      active: String(r[5] || '').toUpperCase() === 'Y'
    });
  }

  return result;
}

/**
 * SYSTEM ADMIN - reset password user.
 */
function resetPassword(token, username, newPassword) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.SYSTEM_ADMIN]);

  if (!username || !newPassword || String(newPassword).length < 6) {
    throw new Error('Username dan password baru minimal 6 karakter wajib diisi.');
  }

  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!sh) throw new Error('Sheet USERS belum tersedia.');

  const values = sh.getDataRange().getValues();

  for (let i = 1; i < values.length; i++) {
    if (normalize_(values[i][0]) === normalize_(username)) {
      sh.getRange(i + 1, 2).setValue(hashPassword_(newPassword));

      // Invalidasi sesi user yang di-reset agar login ulang menggunakan password baru.
      invalidateUserSessions_(String(values[i][0]));

      return 'Password berhasil di-reset.';
    }
  }

  throw new Error('User tidak ditemukan.');
}

/**
 * SYSTEM ADMIN - aktif/nonaktif user.
 * Akun SYSTEM ADMIN sendiri tidak boleh dinonaktifkan.
 */
function setUserActive(token, username, active) {
  const me = getSession_(token);
  requireRole_(me, [CONFIG.ROLES.SYSTEM_ADMIN]);

  const targetUsername = String(username || '').trim();
  if (!targetUsername) throw new Error('Username wajib diisi.');

  if (normalize_(targetUsername) === normalize_(me.username) && !active) {
    throw new Error('System Admin yang sedang login tidak dapat dinonaktifkan.');
  }

  const sh = SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET);
  if (!sh) throw new Error('Sheet USERS belum tersedia.');

  const values = sh.getDataRange().getValues();

  for (let i = 1; i < values.length; i++) {
    if (normalize_(values[i][0]) === normalize_(targetUsername)) {
      sh.getRange(i + 1, 6).setValue(active ? 'Y' : 'N');

      if (!active) invalidateUserSessions_(String(values[i][0]));

      return active ? 'User diaktifkan.' : 'User dinonaktifkan.';
    }
  }

  throw new Error('User tidak ditemukan.');
}

function invalidateUserSessions_(username) {
  // CacheService tidak menyediakan pencarian key wildcard.
  // Session akan kedaluwarsa maksimal sesuai SESSION_TTL.
  // Login baru dengan password lama tetap tidak dapat dilakukan.
  console.log('Password/aktivitas user berubah untuk: ' + username);
}

function validatePayload_(p) {
  ['NAMA TARGET','NO WHATSAPP','AKTIFITAS','TINGKAT PENERIMAAN','METODE','STATUS','KECAMATAN','KELURAHAN'].forEach(k=>{
    if(!p[k] || String(p[k]).trim()==='') throw new Error(k+' wajib diisi.');
  });
  const kec=String(p.KECAMATAN).trim(), kel=String(p.KELURAHAN).trim();
  const upls=String(p.UPLS||'').trim(), uplm=String(p.UPLM||'').trim();
  if (upls && uplm) throw new Error('UPLS dan UPLM hanya boleh memilih salah satu.');
  if (upls && upls !== 'UPLS') throw new Error('Nilai UPLS tidak valid.');
  if (uplm && uplm !== 'UPLM') throw new Error('Nilai UPLM tidak valid.');
  if(!Object.prototype.hasOwnProperty.call(TANGSEL_WILAYAH,kec)) throw new Error('Kecamatan tidak valid.');
  if(TANGSEL_WILAYAH[kec].indexOf(kel)===-1) throw new Error('Kelurahan tidak sesuai dengan Kecamatan.');
  if(AKTIFITAS_OPTIONS.indexOf(String(p.AKTIFITAS).trim())===-1) throw new Error('Aktifitas tidak valid.');
}

function changePassword(token,currentPassword,newPassword){
  const me=getSession_(token);
  currentPassword=String(currentPassword||''); newPassword=String(newPassword||'');
  if(!currentPassword||!newPassword) throw new Error('Password lama dan password baru wajib diisi.');
  if(newPassword.length<6) throw new Error('Password baru minimal 6 karakter.');
  if(currentPassword===newPassword) throw new Error('Password baru harus berbeda dari password lama.');
  const sh=SpreadsheetApp.getActive().getSheetByName(CONFIG.USERS_SHEET), values=sh.getDataRange().getValues();
  for(let i=1;i<values.length;i++) if(normalize_(values[i][0])===normalize_(me.username)){
    if(String(values[i][1])!==hashPassword_(currentPassword)) throw new Error('Password lama salah.');
    sh.getRange(i+1,2).setValue(hashPassword_(newPassword)); return 'Password berhasil diubah.';
  }
  throw new Error('User tidak ditemukan.');
}

function hashPassword_(password) {
  const bytes = Utilities.computeDigest(
    Utilities.DigestAlgorithm.SHA_256,
    String(password),
    Utilities.Charset.UTF_8
  );

  return bytes
    .map(b => ('0' + (b & 0xff).toString(16)).slice(-2))
    .join('');
}

/**
 * Test user - opsional.
 * Jalankan sekali jika membutuhkan akun pengujian.
 */
function buatUserTest() {
  try {
    createUser(
      'admin_kecamatan_pamulang',
      'Pamulang123!',
      'ADMIN KECAMATAN',
      'Admin KECAMATAN Pamulang',
      'KECAMATAN Pamulang'
    );
  } catch (e) {
    console.log(e.message);
  }

  try {
    createUser(
      'admin_dpd',
      'AdminDPD123!',
      'ADMIN PTD',
      'Admin PTD',
      ''
    );
  } catch (e) {
    console.log(e.message);
  }

  try {
    createUser(
      'system_admin',
      'SystemAdmin123!',
      'SYSTEM ADMIN',
      'System Administrator',
      ''
    );
  } catch (e) {
    console.log(e.message);
  }

  try {
    createUser(
      'anggota_pamulang',
      'Anggota123!',
      'ANGGOTA',
      'Nama Anggota Test',
      'KECAMATAN Pamulang'
    );
  } catch (e) {
    console.log(e.message);
  }

  return 'User test selesai.';
}
