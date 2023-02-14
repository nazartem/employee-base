#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=4112
SERVERADDR=https://localhost:${SERVERPORT}

# Add some employees
curl -iL -w "\n" --cacert cert.pem -H "Authorization: Basic am9lOjEyMzQ=" -X POST -H "Content-Type: application/json" --data '{"firstName":"Joe","lastName":"Black", "email":"joeblack@ya.ru"}' ${SERVERADDR}/employee/
curl -iL -w "\n" --cacert cert.pem -H "Authorization: Basic am9lOjEyMzQ=" -X POST -H "Content-Type: application/json" --data '{"firstName":"Nick","lastName":"Stevens", "email":"nicksteve@ya.ru"}' ${SERVERADDR}/employee/
curl -iL -w "\n" --cacert cert.pem -H "Authorization: Basic am9lOjEyMzQ=" -X POST -H "Content-Type: application/json" --data '{"firstName":"Liam","lastName":"Black", "email":"liamblack@ya.ru"}' ${SERVERADDR}/employee/

# Get employees by lastName
curl -iL -w "\n" --cacert cert.pem ${SERVERADDR}/employee/Black/

# Get employees by id
curl -iL -w "\n" --cacert cert.pem -X GET ${SERVERADDR}/employee/1/

# В заголовке строка joe:1234 (логин и пароль для методов, требующих аутентификацию) представлена в кодировке base64
# Можно получить, выполнив команду: $ echo -n "joe:1234" | base64