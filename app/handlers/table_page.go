package handlers

import (
	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/b-a-merritt/squrlviewer/app/utils"
	"github.com/labstack/echo/v4"
)

type TablePageHandler struct {
	cfg   *config.Config
	query func(table string, orderBy *squrl.OrderBy, offset int) (string, []any, error)
}

func NewTablePageHandler(cfg *config.Config) *TablePageHandler {
	return &TablePageHandler{
		cfg:   cfg,
		query: getTableItems,
	}
}

func (h *TablePageHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")

	query, params, err := h.query(table, nil, 0)
	if err != nil {
		return err
	}

	rows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	resultList := [][]interface{}{}

	i := 0
	cols, _ := rows.Columns()
	row := make([]interface{}, len(cols))
	rowPtr := make([]interface{}, len(cols))

	for i = range row {
		rowPtr[i] = &row[i]
	}

	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			return err
		}
		r := make([]interface{}, len(cols))

		for j, val := range row {
			r[j] = utils.FormatFromDB(val)
		}

		resultList = append(resultList, r)
	}

	comp := templates.TablePage(table, cols, resultList)
	layout := templates.Layout(comp, table, []string{table})
	return layout.Render(c.Request().Context(), c.Response())
}
