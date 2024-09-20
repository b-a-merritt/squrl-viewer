package handlers

import (
	"strconv"
	"strings"

	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/b-a-merritt/squrlviewer/app/utils"
	"github.com/labstack/echo/v4"
)

var numItemsPerPage = 1000

func getTableItems(table string, orderBy *squrl.OrderBy, offset int) (string, []any, error) {
	s := squrl.
		New(table).
		Select("*")

	if orderBy != nil {
		s = s.OrderBy([]squrl.OrderBy{*orderBy})
	}

	return s.
		Limit(numItemsPerPage).
		Offset(offset).
		Query()
}

type TableHandler struct {
	cfg   *config.Config
	query func(table string, orderBy *squrl.OrderBy, offset int) (string, []any, error)
}

func NewTableHandler(cfg *config.Config) *TableHandler {
	return &TableHandler{
		cfg:   cfg,
		query: getTableItems,
	}
}

func (h *TableHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")
	orderByParam := c.QueryParam("orderBy")
	sortOrderParam := c.QueryParam("sortOrder")
	offsetParam := c.QueryParam("page")

	sortOrder := squrl.ASC
	offset := 0
	var orderByTerm *squrl.OrderBy

	if orderByParam != "" {
		if strings.ToUpper(sortOrderParam) == "DESC" {
			sortOrder = squrl.DESC
		}

		order := squrl.OrderBy{
			Table: table,
			Field: orderByParam,
			Order: sortOrder,
		}
		orderByTerm = &order
	}

	if offsetParam == "" {
		offset = 0
	} else {
		var err error
		offset, err = strconv.Atoi(offsetParam)

		if err != nil {
			return err
		} else {
			offset *= numItemsPerPage
		}
	}

	query, params, err := h.query(table, orderByTerm, offset)
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

	comp := templates.Table(resultList)
	return comp.Render(c.Request().Context(), c.Response())
}
