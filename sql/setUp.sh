#!/bin/sh

# Following script sets up PostgresSQL database with provided name as provided user.

# $1 -> database_name
# $2 -> username, defaults to postgres
# $3 -> populate_with_test_data, defaults to 0. To run it you have to provide 1 as 3rd argument.

if [ $# -lt 1 ]
then
    echo Usage: ./setUp.sh database_name "[username]" "[populate_with_test]"
    exit
fi

populateDB=0

#TODO: Is this condition correct?
if [ -z "$3" ]
then
    populateDB=$3
fi

psql -d $1 -f ./extensions.sql
psql -d $1 -f ./create_tables.sql
psql -d $1 -f ./procedures.sql

if [ $3 -eq 1 ]
then
    psql -d $1 -f ./test_commands.sql
fi
