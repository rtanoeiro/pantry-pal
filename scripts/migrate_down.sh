#!/bin/bash
# Trying to set up an env file for the variables above make the scipt fails
ENV=$1

goose -dir sql/schema sqlite3 ./data/pantry_pal_${ENV}.db down-to 0