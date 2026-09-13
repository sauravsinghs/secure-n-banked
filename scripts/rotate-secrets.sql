-- Run after a credential leak to rotate the DB password and clear sessions.
-- Usage (existing volume, old password still works):
--   psql -v new_password='NEW_PASSWORD' "postgresql://root:OLD_PASSWORD@localhost:5432/simple_bank?sslmode=disable" -f scripts/rotate-secrets.sql

ALTER USER root WITH PASSWORD :'new_password';

TRUNCATE TABLE sessions;
