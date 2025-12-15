// Package main implements a DAG editor web application using Element and rweb.
// It provides a visual editor for creating Airflow-style DAGs with node connections.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
)

var store = NewStore()

func main() {
	// Initialize with a sample DAG
	store.SaveDAG(CreateSampleDAG())

	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8080",
		Verbose: true,
	})

	s.Use(rweb.RequestInfo)

	// Page routes
	s.Get("/", homeHandler)
	s.Get("/editor/{id}", editorHandler)

	// API routes
	s.Get("/api/dags", apiListDAGs)
	s.Get("/api/dags/{id}", apiGetDAG)
	s.Post("/api/dags", apiCreateDAG)
	s.Put("/api/dags/{id}", apiUpdateDAG)
	s.Delete("/api/dags/{id}", apiDeleteDAG)
	s.Get("/api/dags/{id}/export", apiExportDAG)

	fmt.Println("DAG Editor Server")
	fmt.Println("=================")
	fmt.Println("Visit: http://localhost:8080")
	log.Fatal(s.Run())
}

// homeHandler renders the DAG list page
func homeHandler(c rweb.Context) error {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	dags := store.ListDAGs()

	b.Html().R(
		renderHead(b, "DAG Editor - Home"),
		b.Body().R(
			b.DivClass("app").R(
				// Header
				b.Header().R(
					b.DivClass("container header-content").R(
						b.H1().T("DAG Editor"),
						b.ButtonClass("btn btn-primary", "onclick", "createNewDAG()").T("+ New DAG"),
					),
				),
				// Main content
				b.Main().R(
					b.DivClass("container").R(
						b.H2().T("Your DAGs"),
						b.DivClass("dag-grid").R(
							func() (x any) {
								if len(dags) == 0 {
									b.DivClass("empty-state").R(
										b.P().T("No DAGs yet. Create your first DAG to get started."),
									)
								} else {
									for _, dag := range dags {
										renderDAGCard(b, dag)
									}
								}
								return
							}(),
						),
					),
				),
			),
		),
		b.Script().T(homePageScript()),
	)

	return c.WriteHTML(b.String())
}

// renderDAGCard renders a card for a DAG in the list
func renderDAGCard(b *element.Builder, dag *DAG) {
	b.DivClass("dag-card", "data-id", dag.ID).R(
		b.DivClass("dag-card-header").R(
			b.H3().T(dag.Name),
			b.SpanClass("dag-schedule").T(func() string {
				if dag.Schedule == "" || dag.Schedule == "None" {
					return "Manual"
				}
				return dag.Schedule
			}()),
		),
		b.DivClass("dag-card-body").R(
			b.P().T(dag.Description),
			b.DivClass("dag-stats").R(
				b.Span().F("%d nodes", len(dag.Nodes)),
				b.Span().T(" | "),
				b.Span().F("%d connections", len(dag.Edges)),
			),
			b.DivClass("dag-tags").R(
				func() (x any) {
					for _, tag := range dag.Tags {
						b.SpanClass("tag").T(tag)
					}
					return
				}(),
			),
		),
		b.DivClass("dag-card-actions").R(
			b.AClass("btn btn-sm", "href", "/editor/"+dag.ID).T("Edit"),
			b.ButtonClass("btn btn-sm btn-secondary", "onclick", fmt.Sprintf("exportDAG('%s')", dag.ID)).T("Export"),
			b.ButtonClass("btn btn-sm btn-danger", "onclick", fmt.Sprintf("deleteDAG('%s')", dag.ID)).T("Delete"),
		),
	)
}

// editorHandler renders the DAG editor page
func editorHandler(c rweb.Context) error {
	dagID := c.Request().Param("id")

	dag, ok := store.GetDAG(dagID)
	if !ok {
		return c.WriteHTML("<h1>DAG not found</h1>")
	}

	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	dagJSON, _ := json.Marshal(dag)

	b.Html().R(
		renderHead(b, "Edit: "+dag.Name),
		b.Body().R(
			b.DivClass("editor-app").R(
				// Toolbar
				b.DivClass("toolbar").R(
					b.DivClass("toolbar-left").R(
						b.AClass("btn btn-sm", "href", "/").T("< Back"),
						b.SpanClass("dag-name").T(dag.Name),
					),
					b.DivClass("toolbar-center").R(
						b.ButtonClass("tool-btn active", "data-tool", "select", "title", "Select").T("S"),
						b.ButtonClass("tool-btn", "data-tool", "connect", "title", "Connect Nodes").T("C"),
						b.ButtonClass("tool-btn", "data-tool", "delete", "title", "Delete").T("D"),
					),
					b.DivClass("toolbar-right").R(
						b.ButtonClass("btn btn-sm", "onclick", "addNode('k8s_pod_operator')").T("+ K8s Task"),
						b.ButtonClass("btn btn-sm", "onclick", "addNode('external_sensor')").T("+ External Sensor"),
						b.ButtonClass("btn btn-sm btn-primary", "onclick", "saveDAG()").T("Save"),
						b.ButtonClass("btn btn-sm btn-secondary", "onclick", "exportDAG()").T("Export"),
					),
				),
				// Main editor area
				b.DivClass("editor-main").R(
					// Canvas
					b.DivClass("canvas-container").R(
						b.Div("id", "canvas", "class", "canvas").R(),
						renderCanvasSVG(b),
					),
					// Properties panel
					b.DivClass("properties-panel").R(
						b.DivClass("panel-header").R(
							b.H3().T("Properties"),
							b.DivClass("panel-tabs").R(
								b.ButtonClass("panel-tab active", "data-panel", "node").T("Node"),
								b.ButtonClass("panel-tab", "data-panel", "dag").T("DAG"),
								b.ButtonClass("panel-tab", "data-panel", "code").T("Code"),
							),
						),
						b.DivClass("panel-content").R(
							// Node properties
							b.Div("id", "node-panel", "class", "panel-section active").R(
								b.DivClass("empty-selection").R(
									b.P().T("Select a node to edit its properties"),
								),
								b.Div("id", "node-form", "class", "hidden").R(),
							),
							// DAG properties
							b.Div("id", "dag-panel", "class", "panel-section").R(
								renderDAGPropertiesForm(b, dag),
							),
							// Code panel with Monaco
							b.Div("id", "code-panel", "class", "panel-section").R(
								b.Div("id", "monaco-editor", "class", "monaco-container").R(),
							),
						),
					),
				),
				// Dependencies output
				b.DivClass("dependencies-panel").R(
					b.H4().T("Dependencies"),
					b.Pre("id", "dependencies-output").R(),
				),
			),
			// Hidden data
			b.Script().T(fmt.Sprintf("window.dagData = %s;", string(dagJSON))),
			// Monaco Editor loader
			b.Script("src", "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/loader.min.js").R(),
			// Editor script
			b.Script().T(editorScript()),
		),
	)

	return c.WriteHTML(b.String())
}

// renderHead renders the HTML head with styles
func renderHead(b *element.Builder, title string) any {
	return b.Head().R(
		b.Meta("charset", "utf-8"),
		b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1"),
		b.Title().T(title),
		b.Style().T(getStyles()),
	)
}

// renderCanvasSVG renders the SVG layer for connections
func renderCanvasSVG(b *element.Builder) any {
	return b.Ele("svg", "id", "connections-svg", "class", "connections-svg").R(
		b.Ele("defs").R(
			b.Ele("marker", "id", "arrowhead", "markerWidth", "10", "markerHeight", "7",
				"refX", "9", "refY", "3.5", "orient", "auto").R(
				b.Ele("polygon", "points", "0 0, 10 3.5, 0 7", "fill", "#666").R(),
			),
		),
	)
}

// renderDAGPropertiesForm renders the DAG properties form
func renderDAGPropertiesForm(b *element.Builder, dag *DAG) any {
	return b.Form("id", "dag-form").R(
		b.DivClass("form-group").R(
			b.Label("for", "dag-name").T("Name"),
			b.Input("type", "text", "id", "dag-name", "name", "name", "value", dag.Name, "class", "form-input"),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-description").T("Description"),
			b.TextArea("id", "dag-description", "name", "description", "class", "form-input", "rows", "2").T(dag.Description),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-schedule").T("Schedule (cron or None)"),
			b.Input("type", "text", "id", "dag-schedule", "name", "schedule", "value", dag.Schedule, "class", "form-input", "placeholder", "None or */30 * * * *"),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-concurrency").T("Concurrency"),
			b.Input("type", "number", "id", "dag-concurrency", "name", "concurrency", "value", fmt.Sprintf("%d", dag.Concurrency), "class", "form-input"),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-tags").T("Tags (comma-separated)"),
			b.Input("type", "text", "id", "dag-tags", "name", "tags", "value", joinTags(dag.Tags), "class", "form-input"),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-owner").T("Owner"),
			b.Input("type", "text", "id", "dag-owner", "name", "owner", "value", dag.DefaultArgs.Owner, "class", "form-input"),
		),
		b.DivClass("form-group").R(
			b.Label("for", "dag-email").T("Email"),
			b.Input("type", "email", "id", "dag-email", "name", "email", "value", dag.DefaultArgs.Email, "class", "form-input"),
		),
	)
}

func joinTags(tags []string) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += ", "
		}
		result += t
	}
	return result
}

// API Handlers

func apiListDAGs(c rweb.Context) error {
	dags := store.ListDAGs()
	return c.WriteJSON(dags)
}

func apiGetDAG(c rweb.Context) error {
	id := c.Request().Param("id")
	dag, ok := store.GetDAG(id)
	if !ok {
		return c.Status(404).WriteJSON(map[string]string{"error": "DAG not found"})
	}
	return c.WriteJSON(dag)
}

func apiCreateDAG(c rweb.Context) error {
	var dag DAG
	if err := json.Unmarshal(c.Request().Body(), &dag); err != nil {
		return c.Status(400).WriteJSON(map[string]string{"error": err.Error()})
	}
	if dag.ID == "" {
		dag.ID = fmt.Sprintf("dag-%d", time.Now().UnixNano())
	}
	store.SaveDAG(&dag)
	return c.Status(201).WriteJSON(dag)
}

func apiUpdateDAG(c rweb.Context) error {
	id := c.Request().Param("id")
	var dag DAG
	if err := json.Unmarshal(c.Request().Body(), &dag); err != nil {
		return c.Status(400).WriteJSON(map[string]string{"error": err.Error()})
	}
	dag.ID = id
	store.SaveDAG(&dag)
	return c.WriteJSON(dag)
}

func apiDeleteDAG(c rweb.Context) error {
	id := c.Request().Param("id")
	if store.DeleteDAG(id) {
		return c.WriteJSON(map[string]string{"status": "deleted"})
	}
	return c.Status(404).WriteJSON(map[string]string{"error": "DAG not found"})
}

func apiExportDAG(c rweb.Context) error {
	id := c.Request().Param("id")
	dag, ok := store.GetDAG(id)
	if !ok {
		return c.Status(404).WriteJSON(map[string]string{"error": "DAG not found"})
	}

	code := dag.GeneratePythonCode("your-registry.com", "airflow-image")

	c.Response().SetHeader("Content-Type", "text/plain")
	c.Response().SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s.py", dag.ID))
	_, err := c.Response().Write([]byte(code))
	return err
}
