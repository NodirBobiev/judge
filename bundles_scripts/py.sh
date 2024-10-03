#!/bin/bash

image="python:3.12-alpine"
container="python-temp-container"
destdir="./fs_bundles/python3.12"

docker rm $container > /dev/null 2>&1

date \
&&
echo "⏳ pulling $image..." \
&&
docker pull $image \
&&
echo "⏳ running $container..." \
&&
docker run -d --name $container $image \
&&
echo "⏳ exporting container filesystem..." \
&&
mkdir -p $destdir && docker export $container | tar -x -C $destdir \
&&
echo "⏳ removing $container..." \
&&
docker rm $container \
&&
echo "✅done" \
&&
date
