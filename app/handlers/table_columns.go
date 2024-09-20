package handlers

import (
	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/labstack/echo/v4"
)

type TableColumnsHandler struct {
	cfg   *config.Config
	query func(table string) (string, []any, error)
}

func NewTableColumnsHandler(cfg *config.Config) *TableColumnsHandler {
	return &TableColumnsHandler{
		cfg: cfg,
		query: func(table string) (string, []any, error) {
			return squrl.
				New("columns").
				SetSchema("information_schema").
				Select("column_name").
				Where([]squrl.WhereTerm{
					{
						Table:  "columns",
						Field:  "table_schema",
						Equals: "public",
					},
					{
						Table:  "columns",
						Field:  "table_name",
						Equals: table,
					},
				}).
				Query()
		},
	}
}

func (h *TableColumnsHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")

	var res string
	var tables []string

	query, params, err := h.query(table)
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
