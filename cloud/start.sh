#!/bin/bash

if [ -z "$ENV" ]; then
  echo "Missing ENV to load file env.\$ENV"
  exit
fi

rm config
ln -s config.$ENV config
env $(cat env.$ENV) docker-compose -p analytics up -d $@
