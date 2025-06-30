# user-manager
# clean architecture

# Cobra
*** make command ***
# cobra-cli add command_name

# goose migration
# need => go get github.com/jackc/pgx/v5
$ goose create add_some_column sql
$ Created new file: 20170506082420_add_some_column.sql

$ goose -s create add_some_column sql
$ Created new file: 00001_add_some_column.sql


$ goose create fetch_user_data go
$ Created new file: 20170506082421_fetch_user_data.go