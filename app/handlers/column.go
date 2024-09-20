package handlers

import (
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/labstack/echo/v4"
)

type ColumnSortHandler struct {
}

func NewColumnSortHandler(cfg *config.Config) *ColumnSortHandler {
	return &ColumnSortHandler{}
}

func (h *ColumnSortHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")
	column := c.Param("column")
	sortOrderParam := c.QueryParam("sortOrder")

	if sortOrderParam == "ASC" {
		comp := templates.ColumnHeaderCellDesc(table, column)
		return comp.Render(c.Request().Context(), c.Response())
	} else if sortOrderParam == "DESC" {
		comp := templates.ColumnHeaderCellNone(table, column)
		return comp.Render(c.Request().Context(), c.Response())
	} else {
		comp := templates.ColumnHeaderCellAsc(table, column)
		return comp.Render(c.Request().Context(), c.Response())
	}
}
