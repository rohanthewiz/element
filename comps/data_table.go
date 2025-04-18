package comps

import (
	_ "embed"
	"strings"

	"github.com/rohanthewiz/element"
)

//go:embed data_table.css
var dataTableCSS string

//go:embed data_table.js
var dataTableJS string

// DataTable is a component for representing rows of data.
// Use NewDataTable function for properly creating a new DataTable instance.
type DataTable struct {
	Id string
	DataTableOptions
}

type DataTableOptions struct {
	CustomStyles       string
	ZeroInternalStyles bool
	CustomJS           string
	ZeroInternalJS     bool
}

// NewDataTable returns a new data table component
func NewDataTable(dtOpts DataTableOptions) *DataTable {
	return &DataTable{
		Id:               "ele-data-table", // + genRandId(), // Create a unique ID for the data table
		DataTableOptions: dtOpts,
	}
}

func (tbl DataTable) Render(b *element.Builder) (x any) {
	dataTableJS = strings.ReplaceAll(dataTableJS, "$DATA_TABLE_ID", tbl.Id)
	e, t := b.Vars()

	b.Body().R(
		b.Wrap(func() {
			if len(tbl.DataTableOptions.CustomStyles) > 0 {
				b.Style().R(t(tbl.DataTableOptions.CustomStyles))
			}
			if !tbl.ZeroInternalStyles {
				b.Style().R(t(dataTableCSS))
			}

			if len(tbl.DataTableOptions.CustomJS) > 0 {
				b.Script().R(t(tbl.DataTableOptions.CustomJS))
			}
			if !tbl.ZeroInternalJS {
				b.Script().R(t(dataTableJS))
			}
		}),

		e("div", "class", "data-table-container").R(
			// Loading overlay
			e("div", "class", "loading-overlay").R(
				e("div", "class", "spinner"),
			),

			// Table header
			e("div", "class", "table-header").R(
				e("h2", "class", "table-title").R(t("Customer Data")),
				e("div", "class", "table-controls").R(
					e("div", "class", "search-container").R(
						e("input", "type", "text", "class", "search-input", "placeholder", "Search...", "id", "searchInput"),
					),
					e("select", "class", "per-page-selector", "id", "perPageSelector").R(
						e("option", "value", "5").R(t("5 per page")),
						e("option", "value", "10", "selected").R(t("10 per page")),
						e("option", "value", "25").R(t("25 per page")),
						e("option", "value", "50").R(t("50 per page")),
					),
					e("button", "class", "export-btn", "id", "exportBtn").R(t("Export CSV")),
				),
			),

			// Bulk actions
			e("div", "class", "bulk-actions", "id", "bulkActions").R(
				e("div", "class", "bulk-info").R(
					e("span", "id", "selectedCount").R(t("0")),
					t(" items selected"),
				),
				e("div", "class", "bulk-btns").R(
					e("button", "class", "bulk-btn", "style", "background-color: #f44336; color: white;", "id", "bulkDeleteBtn").R(t("Delete Selected")),
					e("button", "class", "bulk-btn", "style", "background-color: #4CAF50; color: white;", "id", "bulkStatusBtn").R(t("Change Status")),
					e("button", "class", "bulk-btn", "style", "background-color: #607d8b; color: white;", "id", "bulkCancelBtn").R(t("Cancel")),
				),
			),

			// Data table
			e("table", "id", tbl.Id).R(
				e("thead").R(
					e("tr").R(
						e("th", "width", "40px").R(
							e("input", "type", "checkbox", "id", "selectAll"),
						),
						e("th", "data-column", "id", "class", "sortable").R(t("ID")),
						e("th", "data-column", "name", "class", "sortable").R(t("Name")),
						e("th", "data-column", "email", "class", "sortable").R(t("Email")),
						e("th", "data-column", "status", "class", "sortable").R(t("Status")),
						e("th", "data-column", "country", "class", "sortable").R(t("Country")),
						e("th", "data-column", "created", "class", "sortable").R(t("Created Date")),
						e("th", "data-column", "actions", "width", "120px").R(t("Actions")),
					),
				),
				e("tbody"),
			),

			// Pagination
			e("div", "class", "pagination").R(
				e("div", "class", "pagination-info").R(
					t("Showing "),
					e("span", "id", "showingStart").R(t("1")),
					t(" to "),
					e("span", "id", "showingEnd").R(t("10")),
					t(" of "),
					e("span", "id", "totalItems").R(t("0")),
					t(" entries"),
				),
				e("div", "class", "pagination-controls", "id", "paginationControls"),
			),

			// Table footer
			e("div", "class", "table-footer").R(
				e("div", "class", "footer-info").R(
					t("Updated: "),
					e("span", "id", "lastUpdated"),
				),
				e("div", "class", "footer-actions").R(
					e("button", "class", "column-toggle-btn", "id", "columnToggleBtn").R(t("Toggle Columns")),
				),
			),

			// Column toggle dropdown
			e("div", "class", "column-toggle-dropdown", "id", "columnToggleDropdown"),
		),

		// Add/Edit Modal
		e("div", "class", "modal", "id", "userModal").R(
			e("div", "class", "modal-content").R(
				e("div", "class", "modal-header").R(
					e("h3", "class", "modal-title", "id", "modalTitle").R(t("Add New User")),
					e("button", "class", "modal-close", "id", "modalClose").R(t("×")),
				),
				e("div", "class", "modal-body").R(
					e("form", "id", "userForm").R(
						e("input", "type", "hidden", "id", "userId"),
						createFormGroup(b, "Name", "text", "userName", true),
						createFormGroup(b, "Email", "email", "userEmail", true),
						createStatusSelect(b),
						createCountrySelect(b),
					),
				),
				e("div", "class", "modal-footer").R(
					e("button", "class", "action-btn", "style", "background-color: #9e9e9e;", "id", "cancelBtn").R(t("Cancel")),
					e("button", "class", "action-btn", "style", "background-color: #4CAF50;", "id", "saveBtn").R(t("Save")),
				),
			),
		),

		// View Details Modal
		createViewModal(b),

		// Confirmation Modal
		createConfirmModal(b),

		// Toast Notification
		e("div", "class", "toast", "id", "toast"),
	)

	return b.String()
}

// Helper functions for creating form elements
func createFormGroup(b *element.Builder, label, inputType, id string, required bool) (x any) {
	e, t := b.Vars()

	requiredAttr := ""
	if required {
		requiredAttr = "required"
	}

	e("div", "class", "form-group").R(
		e("label", "class", "form-label", "for", id).R(t(label)),
		e("input", "type", inputType, "class", "form-input", "id", id, requiredAttr),
		e("div", "class", "error-message", "id", id+"Error"),
	)
	return
}

func createStatusSelect(b *element.Builder) (x any) {
	e, t := b.Vars()

	e("div", "class", "form-group").R(
		e("label", "class", "form-label", "for", "userStatus").R(t("Status")),
		e("select", "class", "form-select", "id", "userStatus").R(
			e("option", "value", "active").R(t("Active")),
			e("option", "value", "inactive").R(t("Inactive")),
			e("option", "value", "pending").R(t("Pending")),
		),
	)

	return
}

func createCountrySelect(b *element.Builder) (x any) {
	e, t := b.Vars()

	countries := []string{"USA", "Canada", "UK", "Germany", "France", "Japan", "Australia", "Brazil", "India", "China"}
	options := make([]*element.Element, len(countries))

	e("div", "class", "form-group").R(
		e("label", "class", "form-label", "for", "userCountry").R(t("Country")),
		e("select", "class", "form-select", "id", "userCountry").R(
			func() (x any) {
				for _, option := range options {
					e("option").R(option)
				}
				return
			}(),
		),
	)
	return
}

func createViewModal(b *element.Builder) (x any) {
	e, t := b.Vars()

	e("div", "class", "modal", "id", "viewModal").R(
		e("div", "class", "modal-content").R(
			e("div", "class", "modal-header").R(
				e("h3", "class", "modal-title").R(t("User Details")),
				e("button", "class", "modal-close", "id", "viewModalClose").R(t("×")),
			),
			e("div", "class", "modal-body", "id", "viewModalBody"),
			e("div", "class", "modal-footer").R(
				e("button", "class", "action-btn", "style", "background-color: #9e9e9e;", "id", "closeViewBtn").R(t("Close")),
			),
		),
	)
	return
}

func createConfirmModal(b *element.Builder) (x any) {
	e, t := b.Vars()

	e("div", "class", "modal", "id", "confirmModal").R(
		e("div", "class", "modal-content").R(
			e("div", "class", "modal-header").R(
				e("h3", "class", "modal-title").R(t("Confirm Action")),
				e("button", "class", "modal-close", "id", "confirmModalClose").R(t("×")),
			),
			e("div", "class", "modal-body").R(
				e("p", "id", "confirmMessage").R(t("Are you sure you want to perform this action?")),
			),
			e("div", "class", "modal-footer").R(
				e("button", "class", "action-btn", "style", "background-color: #9e9e9e;", "id", "cancelConfirmBtn").R(t("Cancel")),
				e("button", "class", "action-btn", "style", "background-color: #f44336;", "id", "confirmBtn").R(t("Confirm")),
			),
		),
	)
	return
}
