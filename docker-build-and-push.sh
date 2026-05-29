#!/usr/bin/env bash
set -e
docker build . -t srselivan/course-service:latest
docker push srselivan/course-service:latest