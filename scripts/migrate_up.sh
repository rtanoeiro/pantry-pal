#!/bin/bash
# Trying to set up an env file for the variables above make the scipt fails
goose -dir sql/schema sqlite3 ./data/pantry_pal.db up