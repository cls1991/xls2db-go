/*
	core logic for export xls sheet to mysql table.
*/

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/cls1991/xls2db-go/resource"
)

func main() {
	user, passwd, database := "root", "flyfishdb", "xls2db"
    db, err := gorm.Open(
        "mysql", 
        fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, database),
    )
    if err != nil {
        panic(err)
    }
    defer db.Close()
    // debug mode
	db.LogMode(true)

	// test resource package
	r := resource.New("sample", "sid", 0, 2)
	r.ImportData(db, "data/sample.xlsx")
}
