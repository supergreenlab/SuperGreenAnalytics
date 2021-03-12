#!/bin/bash

docker build -t analytics-dev . -f Dockerfile.dev
docker run  --name=analytics --network=supergreencloud_back-tier -p 8080:8080 --rm -it -v $(pwd)/config:/etc/analytics -v $(pwd):/app analytics-dev
docker rmi analytics-dev
