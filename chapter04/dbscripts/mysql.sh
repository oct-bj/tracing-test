#!/usr/bin/env bash
docker run -d --name mysql56 -p3306:3306 -eMYSQL_ROOT_PASSWORD=mysqlpwd mysql:5.6
