#!/bin/bash
set -e

# Generate clients for all services that have the generate-clients make target
for service in services/*/; do
    if [ -f "${service}Makefile" ]; then
        if grep -q "generate-clients:" "${service}Makefile"; then
            (cd "$service" && make generate-clients)
        fi
    fi
done

# Check if there are any changes
if [[ -n $(git status -s clients/) ]]; then
    echo "changes_detected=true" >> "$GITHUB_OUTPUT"
    echo "Changes detected in clients/"
else
    echo "changes_detected=false" >> "$GITHUB_OUTPUT"
    echo "No changes detected in clients/"
fi