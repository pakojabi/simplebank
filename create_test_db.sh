#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE simple_bank_test;
	GRANT ALL PRIVILEGES ON DATABASE simple_bank_test TO root;
EOSQL