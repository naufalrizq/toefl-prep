#!/usr/bin/env bash
# Bootstraps a disposable test database for integration / contract tests and
# for running the server under k6.
#
# Usage:
#   ./scripts/db-test-setup.sh [create|reset]
set -euo pipefail

DB=toefl_test
OWNER=toefl
PASSWORD=toefl

psql -U postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$OWNER'" | grep -q 1 ||
  psql -U postgres -qc "CREATE ROLE $OWNER LOGIN PASSWORD '$PASSWORD'"

case "${1:-reset}" in
  create)
    psql -U postgres -qc "CREATE DATABASE $DB OWNER $OWNER TEMPLATE template0" 2>/dev/null ||
      echo "database $DB already exists (skipping)"
    ;;
  reset)
    psql -U postgres -qc "DROP DATABASE IF EXISTS $DB"
    psql -U postgres -qc "CREATE DATABASE $DB OWNER $OWNER TEMPLATE template0"
    ;;
esac

echo "DATABASE_URL=postgres://$OWNER:$PASSWORD@localhost:5432/$DB"