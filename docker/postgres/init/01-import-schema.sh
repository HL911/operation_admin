#!/bin/bash
set -Eeuo pipefail

DUMP_FILE="/seed/schema_pincermarket_20260401.sql"
IMPORT_FILE="/tmp/pincermarket-import.sql"

if [[ ! -f "${DUMP_FILE}" ]]; then
  echo "Expected SQL dump at ${DUMP_FILE}, but it was not found."
  exit 1
fi

echo "Preparing PostgreSQL dump from ${DUMP_FILE}"

BOM_HEX="$(LC_ALL=C head -c 2 "${DUMP_FILE}" | od -An -tx1 | tr -d ' \n')"

if [[ "${BOM_HEX}" == "fffe" ]]; then
  echo "Detected UTF-16 LE dump, converting to UTF-8 before import"
  iconv -f UTF-16LE -t UTF-8 "${DUMP_FILE}" > "${IMPORT_FILE}"
else
  cp "${DUMP_FILE}" "${IMPORT_FILE}"
fi

echo "Importing dump into ${POSTGRES_DB}"
psql \
  --username "${POSTGRES_USER}" \
  --dbname "${POSTGRES_DB}" \
  -v ON_ERROR_STOP=1 \
  --file "${IMPORT_FILE}"

echo "PostgreSQL import completed"
