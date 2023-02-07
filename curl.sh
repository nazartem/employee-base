#!/usr/bin/bash

set -eux
set -o pipefail

SERVERPORT=4112
SERVERADDR=localhost:${SERVERPORT}

# Start by deleting all existing employees on the server
curl -iL -w "\n" -X DELETE ${SERVERADDR}/employee/

# Add some employees
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"firstName":"Joe","lastName":"Black", "email":"joeblack@ya.ru"}' ${SERVERADDR}/employee/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"firstName":"Nick","lastName":"Stevens", "email":"nicksteve@ya.ru"}' ${SERVERADDR}/employee/
curl -iL -w "\n" -X POST -H "Content-Type: application/json" --data '{"firstName":"Liam","lastName":"Black", "email":"liamblack@ya.ru"}' ${SERVERADDR}/employee/

# Get employees by lastName
curl -iL -w "\n" ${SERVERADDR}/employee/Black

# Get employees by id
curl -iL -w "\n" -X GET ${SERVERADDR}/employee/1/