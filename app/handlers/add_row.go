package handlers

import (
	"slices"

	"github.com/b-a-merritt/squrl"
	"github.com/b-a-merritt/squrlviewer/app/config"
	"github.com/b-a-merritt/squrlviewer/app/templates"
	"github.com/labstack/echo/v4"
)

type AddRowHandler struct {
	cfg            *config.Config
	query          func(name string) (string, []any, error)
	fkQuery        func(name string) (string, []any, error)
	fkOptionsQuery func(name string, field string) (string, []any, error)
}

func NewAddRowHandler(cfg *config.Config) *AddRowHandler {
	return &AddRowHandler{
		cfg: cfg,
		query: func(name string) (string, []any, error) {
			return squrl.
				New(name).
				Select("*").
				Limit(1).
				Query()
		},
		fkQuery: func(name string) (string, []any, error) {
			return squrl.
				New("table_constraints").
				SetSchema("information_schema").
				Select().
				Join(squrl.JoinTerm{
					Fields:   []string{"column_name"},
					JoinType: squrl.InnerJoin,
					On: squrl.JoinTables{
						Left:  "constraint_name",
						Right: "constraint_name",
					},
					Pk: "constraint_name",
					Tables: squrl.JoinTables{
						Left:  "table_constraints",
						Right: "key_column_usage",
					},
				}).
				Join(squrl.JoinTerm{
					Fields:   []string{"table_name", "column_name"},
					JoinType: squrl.InnerJoin,
					On: squrl.JoinTables{
						Left:  "constraint_name",
						Right: "constraint_name",
					},
					Pk: "constraint_name",
					Tables: squrl.JoinTables{
						Left:  "table_constraints",
						Right: "constraint_column_usage",
					},
				}).
				Where([]squrl.WhereTerm{
					{
						Equals: "FOREIGN KEY",
						Field:  "constraint_type",
						Table:  "table_constraints",
					},
					{
						Equals: "public",
						Field:  "table_schema",
						Table:  "table_constraints",
					},
					{
						Equals: name,
						Field:  "table_name",
						Table:  "table_constraints",
					},
				}).
				Query()
		},
		fkOptionsQuery: func(name, field string) (string, []any, error) {
			return squrl.
				New(name).
				Select(field).
				Query()
		},
	}
}

func (h *AddRowHandler) ServeHTTP(c echo.Context) error {
	table := c.Param("table")

	columns := make([]templates.Column, 0)

	query, params, err := h.query(table)
	if err != nil {
		return err
	}
	rows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer rows.Close()
	types, _ := rows.ColumnTypes()
	for _, ptr := range types {
		colType := ptr.DatabaseTypeName()
		nullable, _ := ptr.Nullable()
		length, _ := ptr.Length()

		if colType == "" {
			colType = "UNKNOWN"
		}

		column := templates.Column{
			Name:     ptr.Name(),
			ColType:  colType,
			Size:     length,
			Nullable: nullable,
		}
		columns = append(columns, column)
	}

	query, params, err = h.fkQuery(table)
	if err != nil {
		return err
	}
	fkRows, err := h.cfg.DB.Query(query, params...)
	if err != nil {
		return err
	}
	defer fkRows.Close()
	for fkRows.Next() {
		var fk templates.ForiegnKey
		err := fkRows.Scan(&fk.ColName, &fk.FTableName, &fk.FColName)
		if err != nil {
			return nil
		}
		idx := slices.IndexFunc(columns, func(col templates.Column) bool {
			return fk.ColName == col.Name
		})
		if idx != -1 {
			columns[idx].Fk = &fk
		}
	}

	for i, col := range columns {
		if col.Fk == nil {
			continue
		}
		options := make([]interface{}, 0)
		query, params, err = h.fkOptionsQuery(col.Fk.FTableName, col.Fk.FColName)
		if err != nil {
			return nil
		}
		optionRows, err := h.cfg.DB.Query(query, params...)
		if err != nil {
			return err
		}
		for optionRows.Next() {
			var option interface{}
			err := optionRows.Scan(&option)
			if err != nil {
				return nil
			}
			options = append(options, option)
		}
		columns[i].Fk.Options = options
		optionRows.Close()
	}

	comp := templates.AddRow(table, columns)
	return comp.Render(c.Request().Context(), c.Response())
}
