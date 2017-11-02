/*
    mysql data table model, mapping with xls sheet.
*/

package main

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type SampleModel struct {
    gorm.Model
    sid int `sql:"auto_increment"`
    shortcut string
    content string
}

func (SampleModel) TableName() string {
    return "tb_sample_template"
}

func main() {
    user, passwd := "root", "flyfishdb"
    db, err := gorm.Open(
        "mysql", 
        fmt.Sprintf("%s:%s@/xls2db?charset=utf8&parseTime=True&loc=Local", user, passwd),
    )
    if err != nil {
        panic(err)
    }
    defer db.Close()
    db.LogMode(true)

    db.AutoMigrate(&SampleModel{})
    var sample SampleModel
    db.First(&sample)
    fmt.Println(sample)
}
