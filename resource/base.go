/*
	Base Data Source.
*/

package resource

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"github.com/jinzhu/gorm"
	"github.com/xuri/excelize"
	"github.com/cls1991/xls2db-go/model"
)

type resource struct {
	sheetName string
	uniqueKey string
	headerIndex int
	contentIndex int
}

func New(sheetName string, uniqueKey string, headerIndex int, contentIndex int) resource {
	r := resource {sheetName, uniqueKey, headerIndex, contentIndex}
	return r
}

func (r resource) ImportData(db *gorm.DB, xlsxName string) {
	xlsx, err := excelize.OpenFile(xlsxName)
	if err != nil {
		panic(err)
	}
	rows := xlsx.GetRows(r.sheetName)
	if len(rows) == 0 {
		panic(fmt.Sprintf("Sheet %s not found", r.sheetName))
	}
	var uHeaders []string
	var uIndex = -1
	for line, row := range rows {
		if line == r.headerIndex {
			uHeaders = row
			for k, v := range uHeaders {
				if v == r.uniqueKey {
					uIndex = k
				}
			}
			if uIndex == -1 {
				panic(fmt.Sprintf("unique key %s not defined", r.uniqueKey))
			}
		}else if line >= r.contentIndex {
			var m model.SampleModel
			query := db.Where(fmt.Sprintf("%s = ?", r.uniqueKey), row[uIndex])
			if query.RecordNotFound() == false {
				query.First(&m)
			}
			elems := reflect.ValueOf(&m).Elem()
			for k := range row {
				header := strings.Title(uHeaders[k])
				field := elems.FieldByName(header)
				if field.IsValid() == false {
					panic(fmt.Sprintf("Field %s not valid", header))
				}
				// parse int
				if field.Kind() == reflect.Int {
					n, err := strconv.ParseInt(row[k], 10, 64)
					if err != nil {
						panic(err)
					}
					field.SetInt(n)
				}else {           // string
					field.SetString(row[k])
				}
			}
			db.Save(&m)
		}
	}
}
