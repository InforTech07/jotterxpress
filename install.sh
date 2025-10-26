#!/bin/bash

# JotterXpress Installation Script
# This script installs JotterXpress to your system

set -e

echo "📝 Installing JotterXpress..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

# Build the application
echo "🔨 Building JotterXpress..."
go build -o bin/jx cmd/jotterxpress/main.go

if [ ! -f "bin/jx" ]; then
    echo "❌ Build failed. Please check for errors."
    exit 1
fi

# Install to /usr/local/bin
echo "📦 Installing to /usr/local/bin..."
sudo cp bin/jx /usr/local/bin/

# Make sure it's executable
sudo chmod +x /usr/local/bin/jx

# Test installation
if command -v jx &> /dev/null; then
    echo "✅ Installation successful!"
    echo ""
    echo "🎉 JotterXpress is now installed!"
    echo ""
    echo "Usage examples:"
    echo "  jx \"This is my note\"     # Create a new note"
    echo "  jx list                  # List today's notes"
    echo "  jx list-date 2024-01-15 # List notes for specific date"
    echo "  jx --help               # Show help"
    echo ""
    echo "📁 Notes are stored in: ~/.jotterxpress/notes/"
    echo ""
    echo "Try it out: jx \"Hello, JotterXpress!\""
else
    echo "❌ Installation failed. jx command not found."
    exit 1
fi
