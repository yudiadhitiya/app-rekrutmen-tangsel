#!/usr/bin/env bash
# Menjalankan backend untuk pengembangan dengan konfigurasi dari .env.local.
#
#   cd backend && ./run-dev.sh
#
# Backend membaca konfigurasi dari environment, jadi file .env.local di-source
# di sini (file itu sendiri diabaikan git).
set -euo pipefail

cd "$(dirname "$0")"

if [ ! -f .env.local ]; then
  echo "backend/.env.local tidak ada" >&2
  exit 1
fi

set -a
# shellcheck disable=SC1091
. ./.env.local
set +a

export PATH="$PATH:/usr/local/go/bin"

if [ -n "${GOOGLE_KEY_FILE:-}" ] && [ -f "$GOOGLE_KEY_FILE" ]; then
  perms=$(stat -c '%a' "$GOOGLE_KEY_FILE")
  if [ "$perms" != "600" ]; then
    echo "izin file kredensial $GOOGLE_KEY_FILE adalah $perms, backend menolak selain 600" >&2
    echo "perbaiki dengan: chmod 600 $GOOGLE_KEY_FILE" >&2
    exit 1
  fi
elif [ -n "${GOOGLE_KEY_FILE:-}" ]; then
  echo "file kredensial tidak ditemukan: $GOOGLE_KEY_FILE" >&2
  echo "Google Sheets tidak akan bisa diakses; login warisan dimatikan" >&2
fi

exec go run ./cmd/server
