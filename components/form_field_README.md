## Usage

```go
    // Form Fields
    b.DivClass("section").R(
        b.H2().T("Form Field Components"),
        b.Form("action", "#", "method", "post").R(
            FormField{
                Label:       "Email Address",
                Name:        "email",
                Type:        "email",
                Placeholder: "you@example.com",
                Required:    true,
                HelpText:    "We'll never share your email.",
            }.Render(b),
            FormField{
                Label:       "Password",
                Name:        "password",
                Type:        "password",
                Required:    true,
            }.Render(b),
            FormField{
                Label: "Username",
                Name:  "username",
                Type:  "text",
                Value: "invalid user!",
                Error: "Username can only contain letters and numbers.",
            }.Render(b),
            b.DivClass("form-actions").R(
                b.ButtonClass("btn btn-primary", "type", "submit").T("Submit"),
            ),
        ),
    ),
```

### Example CSS

```css

/* Form */
.form-field {
    margin-bottom: 1.25rem;
}

.form-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
}

.required { color: #e74c3c; }

.form-input {
    width: 100%;
    padding: 0.625rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
}

.form-input:focus {
    outline: none;
    border-color: #3498db;
    box-shadow: 0 0 0 3px rgba(52,152,219,0.1);
}

.has-error .form-input {
    border-color: #e74c3c;
}

.form-error {
    display: block;
    color: #e74c3c;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.form-help {
    display: block;
    color: #666;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.form-actions {
    margin-top: 1.5rem;
}

.btn {
    padding: 0.625rem 1.25rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
}

.btn-primary {
    background: #3498db;
    color: white;
}

.btn-primary:hover {
    background: #2980b9;
}
```