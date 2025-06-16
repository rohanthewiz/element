# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Element is a Go library for generating HTML programmatically without templates. It uses a builder pattern and leverages Go's function execution order to create HTML structure naturally.

## Common Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestElementBasic ./...
```

### Building Examples
```bash
# Build the simple example
cd example/simple_element_example && go build

# Build interfaces example
cd example/interfaces && go build

# Build builder functions example
cd example/using_builder_funcs && go build
```

### Module Management
```bash
# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Code Architecture

### Core Components

1. **Builder Pattern** (`builder.go`, `builder_elements.go`)
   - Central to the library's API
   - Maintains an underlying `strings.Builder` for efficient HTML generation
   - Single-pass rendering with minimal memory allocation

2. **Element Structure** (`element.go`)
   - Core Element struct that represents HTML elements
   - Elements are opened immediately when created
   - Closing is deferred until `.R()` or `.T()` is called

3. **Component System** (`component.go`)
   - Interface for reusable HTML components
   - Components implement `Render(b *Builder) any`
   - Enables composition of complex UI elements

4. **Debug Mode** (`element_debug.go`)
   - Helps identify issues with unpaired attributes or malformed elements
   - Adds `data-ele-id` attributes to elements in debug mode
   - Access via `DebugSet()`, `DebugShow()`, and `DebugClear()`

### Key Design Principles

1. **Function Execution Order**: The library leverages Go's natural function and argument execution order (AST) to match HTML's tree structure.

2. **Termination Methods**:
   - `.R(children...)` - Renders children then closes the element
   - `.T(text...)` - For elements with only text content
   - Single-tag elements (like `<br>`) don't require termination

3. **No Reflection**: Pure compiled Go without reflection, annotations, or "weird rules" for maximum performance.

### Testing Approach

- Unit tests are in `*_test.go` files alongside source files
- Tests focus on HTML output correctness
- Debug mode testing verifies error detection capabilities
- Example applications serve as integration tests

## Important Notes

- Always use the Builder pattern (`element.NewBuilder()`) rather than creating Elements directly
- The library has been production-tested for 7+ years at ccswm.org
- Performance is a key consideration - avoid adding features that require reflection or multiple passes
- Keep the API light and unobtrusive