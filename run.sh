#!/bin/zsh

trap '' INT

export HYVE_HTTP_PORT=8080
export HYVE_HTTPS_PORT=4443
export DATA_DIR=./.data

npm --prefix ./ui run dev &
go run ./examples/base/main.go &

if [ "$(docker ps -f "name=hyve-pebble" | wc -l | awk '{print $1}')" -eq 2 ]; then
  # Pebble already running
elif [ "$(docker ps -a -f "name=hyve-pebble" | wc -l | awk '{print $1}')" -eq 2 ]; then
  docker start hyve-pebble > /dev/null
else
  docker run \
    -p 14000:14000 \
    -p 15000:15000 \
    -e "PEBBLE_VA_NOSLEEP=1" \
    --mount src=$(pwd)/.pebble,target=/pebble,type=bind \
    --name hyve-pebble \
    -d ghcr.io/letsencrypt/pebble -config /pebble/config.json -dnsserver 8.8.8.8:53 > /dev/null
fi

wait
