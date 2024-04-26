#!/bin/bash

docker pull python:3.12-alpine
&&
docker run -d --name python-temp-container python:3.12-alpine
&&
mkdir -p ./fs_bundles/python3.12 && docker export python-temp-container | tar -x -C ./fs_bundles/python3.12
&&
docker rm python-temp-container
