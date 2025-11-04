#!/bin/bash
FILE="../db_data_dev/bd_nahue_v2.sql"
CONTAINER="2025-utn-tpgrupal-grupo1-db-1"
DBUSER="postgres"
DBNAME="db_nahue"

docker exec -i $CONTAINER pg_restore \
  -U $DBUSER \
  -d $DBNAME \
  --clean --if-exists --no-owner --no-privileges < "$FILE"

echo "✅ Restore completado"
