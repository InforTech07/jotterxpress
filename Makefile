.PHONY: build install clean test run help

# Default target
help:
	@echo "📝 JotterXpress - CLI Note Taking Tool"
	@echo ""
	@echo "Available commands:"
	@echo "  build    - Build the application"
	@echo "  install  - Install the application to /usr/local/bin"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help message"

# Build the application
build:
	@echo "🔨 Building JotterXpress..."
	go build -o bin/jx cmd/jotterxpress/main.go
	@echo "✅ Build complete! Binary available at bin/jx"

# Install the application
install: build
	@echo "📦 Installing JotterXpress to /usr/local/bin..."
	sudo cp bin/jx /usr/local/bin/
	@echo "✅ Installation complete! You can now use 'jx' from anywhere"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	@echo "✅ Clean complete!"

# Run the application
run:
	@echo "🚀 Running JotterXpress..."
	go run cmd/jotterxpress/main.go $(ARGS)

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Development setup
dev-setup:
	@echo "⚙️  Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✅ Development setup complete!"

# Quick build and run
quick: build
	@echo "🚀 Running JotterXpress..."
	./bin/jx $(ARGS)
