#!/bin/bash
set -e

if ! command -v java &> /dev/null; then
    echo "Java is required but not found. Please install Java first:"
    echo "sudo apt-get update && sudo apt-get install -y default-jre"
    exit 1
fi

if ! command -v npx &> /dev/null; then
    echo "npm/npx is required but not installed. Please install Node.js and npm first."
    exit 1
fi

# Configure npm to use a different directory for global installations
if [ ! -d "$HOME/.npm-global" ]; then
    echo "Creating global npm directory in user space..."
    mkdir "$HOME/.npm-global"
    npm config set prefix "$HOME/.npm-global"

    # Determine which shell config file to use
    SHELL_CONFIG=""
    if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "/bin/zsh" ] || [ "$SHELL" = "/usr/bin/zsh" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "/bin/bash" ] || [ "$SHELL" = "/usr/bin/bash" ]; then
        SHELL_CONFIG="$HOME/.bashrc"
    fi

    # Add to PATH if shell config file is found and path is not already there
    if [ -n "$SHELL_CONFIG" ] && ! grep -q ".npm-global/bin" "$SHELL_CONFIG"; then
        echo "Adding npm-global to PATH in $SHELL_CONFIG"
        echo 'export PATH="$HOME/.npm-global/bin:$PATH"' >> "$SHELL_CONFIG"
        export PATH="$HOME/.npm-global/bin:$PATH"
    fi
fi

# Install openapi-generator-cli if not already installed
if ! npm list -g @openapitools/openapi-generator-cli &> /dev/null; then
    echo "Installing @openapitools/openapi-generator-cli..."
    npm install -g @openapitools/openapi-generator-cli
fi