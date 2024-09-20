package handlers

import (
	"fmt"

	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/labstack/echo/v4"
)

func countTableRows(table string) (string, []any, error) {
	return squrl.
		New(table).
		SetSchema("public").
		Select(squrl.Count("*", table)).
		Query()
}

type CountTableHandler struct {
	cfg   *config.Config
	query func(name string) (string, []any, error)
}

func NewCountTableHandler(cfg *config.Config) *CountTableHandler {
	return &CountTableHandler{
		cfg:   cfg,
		query: countTableRows,
	}
}

func (h *CountTableHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")

	var res int

	query, params, err := h.query(table)
	if err != nil {
		return err
	}

	rows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&res)
	}

	comp := templates.CountTables(fmt.Sprintf("%v", res))
	return comp.Render(c.Request().Context(), c.Response())
}
