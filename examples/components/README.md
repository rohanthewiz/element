# Element Components Example

This example demonstrates the ready-made UI components in
`github.com/rohanthewiz/element/components` — including the interactive ones
that need **no JavaScript** (Tabs, Accordion, Modal, Dropdown).

## Running the Example

```bash
cd examples/components
go run .
```

Then visit: http://localhost:8080

## Components Included

| Component | Highlights |
|---|---|
| `Table` | Column alignment/classes, caption, footer rows, empty message, per-row class hook, `element.Component` values in cells, responsive wrapper |
| `Tabs` | CSS-only radio-button pattern — zero JS; emits its own scoped show/hide rules |
| `Accordion` | Native `<details>`/`<summary>`; `Exclusive` keeps one section open via the `name` attribute |
| `Modal` | Native Popover API — browser handles open/close, Esc, click-outside; style the overlay with `.modal::backdrop` |
| `Dropdown` | Native `<details>` disclosure menu with links and dividers |
| `Nav`, `Breadcrumb`, `Pagination` | Navigation with proper `aria-current`/`rel` attributes |
| `Alert`, `Badge`, `Button`, `ProgressBar` | Status and action elements with variant classes |
| `Card`, `DefinitionList` | Content containers |
| `FormField`, `SelectField`, `TextAreaField`, `CheckboxField`, `RadioGroup` | Form controls with labels, required markers, help text, and error display wired up with `aria-describedby`/`aria-invalid` |

## Usage Pattern

Import the package and render components inside any element tree:

```go
import "github.com/rohanthewiz/element/components"

components.Table{
    Caption: "Q3 user activity",
    Columns: []components.Column{
        {Header: "ID", Align: "right"},
        {Header: "Name"},
        {Header: "Status", Align: "center"},
    },
    Rows: [][]any{
        // Cells can be scalars or any element.Component
        {1, "Alice", components.Badge{Text: "Active", Type: components.AlertSuccess}},
    },
    Striped: true,
    Hover:   true,
}.Render(b)
```

Composition works everywhere a body is accepted — tab panels, accordion
sections, modal bodies, alert bodies, and table cells all take an
`element.Component`.

For one-off inline components, adapt a function with `element.CompFunc`:

```go
element.CompFunc(func(b *element.Builder) (x any) {
    b.H1().T("Hello!")
    return
})
```

## Styling

Components emit semantic class names (`.table`, `.tabs`, `.accordion-item`,
`.modal`, `.form-field`, ...) and leave the cosmetics to your stylesheet.
`styles.go` in this example contains a complete, copy-paste-friendly
stylesheet covering every component.

The only CSS a component generates itself is the structural show/hide rules
for `Tabs`, which are scoped to that instance's `ID`.

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

Note: component text is written as-is (not HTML-escaped), consistent with the
core library. Escape untrusted user input (e.g. with `html.EscapeString`)
before placing it in component fields.
