# CLAUDE.md

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
cd examples/simple_element_example && go build

# Build interfaces example
cd examples/interfaces && go build
```

### Module Management
```bash
# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Code Architecture
See ai_docs/ARCHITECTURE.md

### Core Components

1. **Builder Pattern** (`builder.go`, `builder_elements.go`)
   - Central to the library's API
   - Maintains an underlying `strings.Builder` for efficient HTML generation
   - Single-pass rendering with minimal memory allocation

2. **Element Structure** (`element.go`)
   - Core Element struct that represents HTML elements
   - Elements are opened immediately when created
   - Closing is deferred until `.R()` or `.T()` or `.F()` is called

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
   - `.F(format, params...)` - For elements only text content -- formatted
   - Single-tag elements (like `<br>`) don't require termination, however so the developer doesn't have to track single vs paired elements, it is recommended to close single tag elements with `.R()`

3. **No Reflection**: Pure compiled Go without reflection, annotations, or "weird rules" for maximum performance.

### Testing Approach
- 
- Example applications serve as integration tests

## Important Notes

- Keep the API light and unobtrusive
- Always use the Builder pattern (`element.NewBuilder()`) rather than creating Elements directly
- Performance is a key consideration
  - avoid adding features that require reflection or multiple passes
  - Prefer to use the Bytes versions of functions to write to the response buffer without string conversion 
