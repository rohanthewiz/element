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

## DebugShow Enhancement Plan

### Current State (as of 2025-06-21)

The `DebugShow()` function in `element_debug.go` has been enhanced with:
1. **Tabbed interface** - HTML view (default) and Markdown view
2. **HTML view** - Shows styled table with issues, copy functionality for element IDs, and "Clear Issues" button
3. **Markdown view** - Shows markdown table format with a blue "Copy" button in top-right corner
4. **Terminal output** - Always outputs markdown table to console
5. **Uses Class builder methods** - e.g., `DivClass()`, `ButtonClass()`, `TableClass()`

### Issue Deduplication Plan

**Problem**: Multiple identical issues can appear for the same file location when elements are created in loops or repeated code patterns. Element IDs are randomly generated on each page refresh, so they cannot be used for deduplication.

**Deduplication Key Components**:
1. **File location** - `el.location` (e.g., "element_test.go:123")
2. **Element name** - `el.name` (e.g., "div", "span")
3. **Issue text** - The actual issue message from `el.issues[]`
4. **Function context** - `el.function` (the function where the element was created)

**Implementation Strategy**:
1. Create a deduplication map with composite key: `location + "|" + name + "|" + issueText`
2. Track count of duplicates for each unique issue
3. Display deduplicated issues with occurrence count
4. In the concerns map iteration, group by the deduplication key
5. Show format like: "**div** tag not closed (3 occurrences)"

**Data Structures Needed**:
```go
type DeduplicatedIssue struct {
    Key      string   // Composite key for deduplication
    Element  Element  // First occurrence of the element
    Count    int      // Number of occurrences
    IssueText string  // The issue message
}
```

**UI Changes**:
- Add a new column "Count" or incorporate count into existing columns
- In markdown view, show count in parentheses
- In HTML view, could add a badge or number indicator

**Code Locations to Modify**:
1. `DebugShow()` function around lines 150-460
2. The loop that builds HTML table (around line 425)
3. The loop that builds markdown content (around line 164)
4. Consider adding a toggle for deduplication on/off

**Testing Considerations**:
- Test with elements created in loops
- Test with identical issues in different functions
- Ensure deduplication doesn't hide important context
- Verify both HTML and Markdown views show counts correctly