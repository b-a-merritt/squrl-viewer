package handlers

import (
	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/labstack/echo/v4"
)

type ListTablesPageHandler struct {
	cfg   *config.Config
	query func() (string, []any, error)
}

func NewListTablesPageHandler(cfg *config.Config) *ListTablesPageHandler {
	return &ListTablesPageHandler{
		cfg: cfg,
		query: func() (string, []any, error) {
			return squrl.
				New("tables").
				SetSchema("information_schema").
				Select("table_name").
				Where([]squrl.WhereTerm{{
					Table:  "tables",
					Field:  "table_schema",
					Equals: "public",
				}}).
				OrderBy([]squrl.OrderBy{
					{Field: "table_name", Table: "tables", Order: squrl.ASC},
				}).
				Query()
		},
	}
}

func (h *ListTablesPageHandler) ServeHTTP(c echo.Context) error {
	var res string
	var tables []string

	query, params, err := h.query()
	if err != nil {
		c.JSON(422, "An error occurred")
		return err
	}

	rows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		c.JSON(422, "An error occurred")
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res)
		tables = append(tables, res)
	}

	comp := templates.ListTables(tables)
	layout := templates.Layout(comp, "Home", []string{})
	return layout.Render(c.Request().Context(), c.Response())
}
