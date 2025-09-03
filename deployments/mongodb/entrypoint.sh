#!/bin/bash
set -e

# Create temporary directory for combined scripts
mkdir -p /docker-entrypoint-initdb.d

# Find all service directories in mongodb-init
counter=0
for service_dir in /mongodb-init/*; do
    if [ -d "$service_dir" ]; then
        # Extract service name from path
        service_name=$(basename "$service_dir")

        # Copy all .js files with a prefix based on the order of discovery
        for script in "$service_dir"/*.js; do
            if [ -f "$script" ]; then
                filename=$(basename "$script")
                cp "$script" "/docker-entrypoint-initdb.d/${counter}${filename}"
            fi
        done

        echo "Processed init scripts for service: $service_name"
        counter=$((counter + 1))
    fi
done

# Execute the original MongoDB entrypoint
exec docker-entrypoint.sh mongod