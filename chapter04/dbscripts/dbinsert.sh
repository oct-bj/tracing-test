#!/usr/bin/env bash
docker exec -i mysql56 mysql -uroot -pmysqlpwd < ./database.sql
