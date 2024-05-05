#!/bin/bash
set -eu
DB_NAME=mail
MIGRATION_NAME=${1:-}

# Reset the shadow database
encore db reset --shadow $DB_NAME

# ent executes Go code without initializing Encore when generating migrations,
# so configure the Encore runtime to be aware that this is expected.
export ENCORERUNTIME_NOPANIC=1

# Generate the migration
atlas migrate diff $MIGRATION_NAME --env local --dev-url "$(encore db conn-uri --shadow $DB_NAME)&search_path=public"
