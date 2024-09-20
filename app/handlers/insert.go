package handlers

import (
	"fmt"
	"net/http"

	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/b-a-merritt/squrlviewer/app/utils"
	"github.com/labstack/echo/v4"
)

type InsertHandler struct {
	cfg       *config.Config
	colsQuery func(name string) (string, []any, error)
	countQuery func(name string) (string, []any, error)
	query     func(string, map[string]interface{}, []string) (string, []any, error)
}

func NewInsertHandler(cfg *config.Config) *InsertHandler {
	return &InsertHandler{
		cfg: cfg,
		colsQuery: func(name string) (string, []any, error) {
			return squrl.
				New(name).
				Select("*").
				Query()
		},
		countQuery: func(name string) (string, []any, error) {
			return squrl.
				New(name).
				Select(squrl.Count("*", name)).
				Query()
		},
		query: func(name string, values map[string]interface{}, returning []string) (string, []any, error) {
			return squrl.
				New(name).
				Insert(values).
				Returning(returning...).
				Query()
		},
	}
}

func (h *InsertHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")

	query, params, err := h.colsQuery(table)
	if err != nil {
		return err
	}
	rows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, _ := rows.Columns()

	var count int
	query, params, err = h.countQuery(table)
	if err != nil {
		return err
	}
	countRows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer countRows.Close()

	for countRows.Next() {
		countRows.Scan(&count)
	}

	var formValue interface{}
	values := make(map[string]interface{})
	returning := make([]string, 0)

	for _, col := range cols {
		formValue = c.FormValue(col)
		if formValue == "" && col == "id" {
			continue
		} else if formValue == "" {
			c.HTML(http.StatusBadRequest, fmt.Sprintf("<p>%v is missing from request body</p>", col))
		}
		values[col] = formValue
		returning = append(returning, col)
	}

	if len(values) == 0 {
		c.HTML(http.StatusBadRequest, "<p>the request body was empty</p>")
	}

	query, params, err = h.query(table, values, returning)
	if err != nil {
		return err
	}

	insertRows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer insertRows.Close()

	i := 0
	row := make([]interface{}, len(cols))
	rowPtr := make([]interface{}, len(cols))

	for i = range row {
		rowPtr[i] = &row[i]
	}

	for insertRows.Next() {
		err = insertRows.Scan(rowPtr...)
		if err != nil {
			return err
		}
		r := make([]interface{}, len(cols))

		for j, val := range row {
			r[j] = utils.FormatFromDB(val)
		}
	}

	comp := templates.Row(count, row)
	return comp.Render(c.Request().Context(), c.Response())
}
