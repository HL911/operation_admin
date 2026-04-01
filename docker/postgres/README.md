# Local PostgreSQL Import

This setup starts PostgreSQL 16 and imports the dump from `/Users/leon/Downloads/schema_pincermarket_20260401.sql` during the first initialization of the database volume.

## Start

```bash
docker compose up -d
```

## Connection Defaults

- Host: `127.0.0.1`
- Port: `15432`
- Database: `pincermarket`
- User: `postgres`
- Password: `postgres`

## Override the Dump Path

```bash
PINCERMARKET_SQL_DUMP=/absolute/path/to/your.sql docker compose up -d
```

## Re-import from Scratch

The initialization scripts only run when the `postgres_data` volume is empty.

```bash
docker compose down -v
docker compose up -d
```

## Notes

- The provided dump is UTF-16 LE, so the init script converts it to UTF-8 before running `psql`.
- The inspected dump looks like a schema export. No `COPY` or `INSERT` statements were found, so this import may create database structure without application rows.
