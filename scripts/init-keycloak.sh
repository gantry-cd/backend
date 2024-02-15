#! /bin/bash

# Wait for Keycloak to start
until $(curl --output /dev/null --silent --head --fail http://localhost:8080); do
    printf '.'
    sleep 5
done

# create initial admin user
.bin/init-keycloak/init