/*
 * Copyright (c) 2021 BlueStorm
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package datatable

import (
	"fmt"

	"github.com/oarkflow/pkg/str"
)

func (dt *DataTable) GroupBy(query string) *DataTable {
	if dt.Count == 0 {
		return dt
	}
	exp := ParseExpr([]byte(query))
	count := len(exp.GroupExpr)
	if count == 0 {
		return dt
	}
	dataTable := &DataTable{Name: dt.Name}
	for _, item := range exp.GroupExpr {
		for _, column := range dt.Columns {
			if column.Name == item.Name {
				dataTable.Columns = append(dataTable.Columns, column)
				break
			}
		}
	}
	if len(dataTable.Columns) == 0 {
		return dt
	}
	var field string
	var groupRows = make(map[string][]map[string]any)

	for _, dr := range dt.Rows {
		row := make(map[string]any)
		var fieldValue string
		var fields []string
		for _, column := range dataTable.Columns {
			field = column.Name
			if !str.Contains(fields, field) {
				fields = append(fields, field)
				value := dr[field]
				fieldValue = ToString(value)
				row[field] = fieldValue
			}
			if _, ok := groupRows[fieldValue]; !ok {
				mp := make([]map[string]any, 0)
				groupRows[fieldValue] = mp
				row["$GroupKey$"] = fieldValue
				dataTable.Rows = append(dataTable.Rows, row)
			}
			groupRows[fieldValue] = append(groupRows[fieldValue], dr)
		}
	}

	for _, dr := range dataTable.Rows {
		key := fmt.Sprintf("%v", dr[field])
		dr[key] = groupRows[ToString(dr["$GroupKey$"])]
		delete(dr, "$GroupKey$")
		delete(dr, field)
	}
	return dataTable
}
