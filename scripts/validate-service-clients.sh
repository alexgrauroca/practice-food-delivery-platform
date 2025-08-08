#!/bin/bash
set -e

# Create temporary directory for comparison
TEMP_DIR=$(mktemp -d)
cp -r clients/ "$TEMP_DIR/"

# Generate clients using existing make targets
for service in services/*/; do
    if [ -f "${service}Makefile" ] && grep -q "generate-clients:" "${service}Makefile"; then
        echo "Generating clients for ${service}..."
        (cd "$service" && make generate-clients)
    fi
done

# Compare with original
if ! diff -r clients/ "$TEMP_DIR/clients/" > /dev/null; then
    echo "::error::Generated clients are out of date. Please run 'make generate-clients' in the relevant services and commit the changes"
    exit 1
fi

echo "Clients are up to date!"