# Element Components Example

This example demonstrates reusable UI components built with the Element library.

## Running the Example

```bash
cd examples/components
go run .
```

Then visit: http://localhost:8080

## Components Included

### Table Component

Renders an HTML table with customizable headers and data rows.

```go
table := Table{
    Headers:  []string{"ID", "Name", "Status"},
    Rows:     [][]any{
        {1, "Alice", "Active"},
        {2, "Bob", "Pending"},
    },
    Striped:  true,
    Bordered: true,
}

b.Div().R(
    table.Render(b),
)
```

### Card Component

A styled card container with optional header and footer.

```go
card := Card{
    Title:  "My Card",
    Body:   "Card content goes here.",
    Footer: "Last updated: Today",
}
```

### Nav Component

A navigation bar with links.

```go
nav := Nav{
    Brand: "My Site",
    Items: []NavItem{
        {Label: "Home", Href: "/", Active: true},
        {Label: "About", Href: "/about"},
    },
}
```

### Alert Component

Styled notification messages with different severity levels.

```go
alert := Alert{
    Type:        AlertSuccess,
    Title:       "Success!",
    Message:     "Your changes have been saved.",
    Dismissible: true,
}
```

### Breadcrumb Component

Navigation breadcrumb trail.

```go
breadcrumb := Breadcrumb{
    Items: []BreadcrumbItem{
        {Label: "Home", Href: "/"},
        {Label: "Products", Href: "/products"},
        {Label: "Current Page", Href: ""},
    },
}
```

### DefinitionList Component

A definition list for term/definition pairs.

```go
dl := DefinitionList{
    Items: []Definition{
        {Term: "API", Definition: "Application Programming Interface"},
        {Term: "SDK", Definition: "Software Development Kit"},
    },
}
```

### Pagination Component

Page navigation controls.

```go
pagination := Pagination{
    CurrentPage: 3,
    TotalPages:  10,
    BaseURL:     "/page/%d",
    ShowFirst:   true,
    ShowLast:    true,
}
```

### FormField Component

Form field with label, input, and optional help/error text.

```go
field := FormField{
    Label:       "Email",
    Name:        "email",
    Type:        "email",
    Required:    true,
    HelpText:    "We'll never share your email.",
}
```

## Using Components Together

Render multiple components with `element.RenderComponents`:

```go
b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)

b.Div().R(
    element.RenderComponents(b,
        nav,
        breadcrumb,
        Alert{Type: AlertInfo, Message: "Welcome!"},
        table,
    ),
)
```

## Creating Your Own Components

Implement the `element.Component` interface:

```go
type MyComponent struct {
    Title string
    Items []string
}

func (m MyComponent) Render(b *element.Builder) (x any) {
    b.DivClass("my-component").R(
        b.H3().T(m.Title),
        b.Ul().R(
            element.ForEach(m.Items, func(item string) {
                b.Li().T(item)
            }),
        ),
    )
    return
}
```

All components in this example have no external dependencies beyond the Element library itself, making them easy to copy and adapt for your own projects.
