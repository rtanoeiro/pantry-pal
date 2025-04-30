# /bin/bash

if [ -f .env ]; then
    source .env
fi

echo $DATABASE_URL

cd sql/schema

goose turso $DATABASE_URL up