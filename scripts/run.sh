#!/bin/bash
docker network create gantry-dev-networks
docker-compose up --build -d