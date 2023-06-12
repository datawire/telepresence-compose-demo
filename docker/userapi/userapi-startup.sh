#!/bin/sh

# Maximum number of connection attempts

WAIT_SECONDS=5
MAX_ATTEMPTS=10

attempt=1
while [ $attempt -le $MAX_ATTEMPTS ]
do
    echo "Attempt $attempt: Checking PostgreSQL connection..."

    # Try connecting to the PostgreSQL server
    export PGPASSWORD=${DB_PASSWORD}
    psql -h ${DB_HOST} -p 5432 -U ${DB_USERNAME} -c "SELECT 1" >/dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "PostgreSQL is ready!"
        ./userapi
    fi

    echo "PostgreSQL is not yet ready. Retrying in $WAIT_SECONDS seconds..."
    sleep $WAIT_SECONDS
    attempt=$((attempt + 1))
done


