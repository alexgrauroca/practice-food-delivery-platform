#!/bin/bash
set -e

# Create temporary directory for comparison
TEMP_DIR=$(mktemp -d)
cp -r clients/ "$TEMP_DIR/"

# Generate clients using existing make targets
for service in services/*/; do
    if [ -f "${service}Makefile" ] && grep -q "generate-clients:" "${service}Makefile"; then
        echo "Generating clients for ${service}..."
        SERVICE_NAME=$(basename "$service")

        # Create the client directory if it doesn't exist
        mkdir -p "clients/${SERVICE_NAME}"

        # Generate the client using the service's make target
        (cd "$service" && make generate-clients)
    fi
done

# Compare with original
if ! diff -r clients/ "$TEMP_DIR/clients/" > /dev/null; then
    echo "::error::Generated clients are out of date. Please run 'make generate-clients' in the relevant services and commit the changes"
    echo "Differences found:"
    diff -r clients/ "$TEMP_DIR/clients/" || true
    exit 1
fi

echo "Clients are up to date!"