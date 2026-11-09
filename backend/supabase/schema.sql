-- Skema Supabase untuk Aplikasi Rekrutmen Tangsel.
-- Jalankan di SQL Editor Supabase (Dashboard > SQL > New query).
--
-- Kredensial diambil dari environment backend sebagai service-role key, sehingga
-- Row Level Security tidak membatasi akses backend. RLS di bawah tetap aktif
-- untuk akses langsung dari browser, dan menolak semua request anonim.

create table if not exists public.rekrutmen (
  id                  uuid primary key default gen_random_uuid(),
  nama                text not null check (char_length(trim(nama)) > 0),
  no_wa               text not null unique check (no_wa ~ '^[0-9]+$'),
  alamat              text not null default '',
  rt                  text not null default '' check (rt = '' or rt ~ '^[0-9]+$'),
  rw                  text not null default '' check (rw = '' or rw ~ '^[0-9]+$'),
  kecamatan           text not null check (
                        kecamatan in ('Serpong','Serpong Utara','Pondok Aren','Ciputat',
                                      'Ciputat Timur','Pamulang','Setu')),
  kelurahan           text not null default '',
  aktivitas           text not null,
  tingkat_penerimaan  text not null check (
                        tingkat_penerimaan in ('Sangat Terbuka','Terbuka','Netral',
                                               'Kurang Terbuka','Tertutup')),
  metode              text not null check (
                        metode in ('Halaqoh','Program Antara','Dakwah Fardiyah')),
  status              text not null check (
                        status in ('Berhasil','Proses','Pending','Gagal')),
  keterangan          text not null default '',
  upl                 text not null default '' check (upl in ('','UPLS','UPLM')),
  created_at          timestamptz not null default now(),
  updated_at          timestamptz not null default now(),
  created_by          text not null default '',
  updated_by          text not null default ''
);

create table if not exists public.aktivitas (
  id      uuid primary key default gen_random_uuid(),
  nama    text not null unique check (char_length(trim(nama)) > 0),
  aktif   boolean not null default true
);

create index if not exists rekrutmen_kecamatan_idx on public.rekrutmen (kecamatan);
create index if not exists rekrutmen_kelurahan_idx on public.rekrutmen (kelurahan);
create index if not exists rekrutmen_metode_idx    on public.rekrutmen (metode);
create index if not exists rekrutmen_status_idx    on public.rekrutmen (status);
create index if not exists rekrutmen_created_idx   on public.rekrutmen (created_at desc);
create index if not exists rekrutmen_aktivitas_idx on public.rekrutmen (aktivitas);

-- Migrasi enum metode: "Program" dan "Antara" digabung menjadi satu pilihan
-- "Program Antara". Constraint check lama dilepas DULUAN, karena nilai baru
-- "Program Antara" belum diizinkan oleh constraint lama sehingga UPDATE
-- akan gagal jika constraint masih aktif.
alter table public.rekrutmen drop constraint if exists rekrutmen_metode_check;

update public.rekrutmen set metode = 'Program Antara'
  where metode in ('Program', 'Antara');

alter table public.rekrutmen
  add constraint rekrutmen_metode_check
  check (metode in ('Halaqoh','Program Antara','Dakwah Fardiyah'));

-- updated_at diperbarui otomatis setiap kali baris berubah.
create or replace function public.set_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

drop trigger if exists rekrutmen_set_updated_at on public.rekrutmen;
create trigger rekrutmen_set_updated_at
  before update on public.rekrutmen
  for each row execute function public.set_updated_at();

-- Row Level Security: aplikasi hanya mengakses lewat backend (service role),
-- jadi akses langsung dari browser ditolak.
alter table public.rekrutmen enable row level security;
alter table public.aktivitas enable row level security;

drop policy if exists "no anon access rekrutmen" on public.rekrutmen;
create policy "no anon access rekrutmen" on public.rekrutmen for all to anon using (false);

drop policy if exists "no anon access aktivitas" on public.aktivitas;
create policy "no anon access aktivitas" on public.aktivitas for all to anon using (false);

-- Master aktivitas awal (dapat dikelola Admin lewat aplikasi afterwards).
insert into public.aktivitas (nama) values
  ('Agen Asuransi'),
  ('Buruh Pabrik'),
  ('Driver Ojek Online'),
  ('Driver Taksi Online'),
  ('Guru'),
  ('Ibu Rumah Tangga'),
  ('Karyawan Swasta'),
  ('Pedagang'),
  ('Pedagang Online'),
  ('Pelajar / Mahasiswa'),
  ('Pensiunan'),
  ('Petani'),
  ('Perawat'),
  ('Pedagang Keeping'),
  ('Wiraswasta'),
  ('Ustadz / Daiyah'),
  ('Tukang Bangun'),
  ('Pengusaha')
on conflict (nama) do nothing;
