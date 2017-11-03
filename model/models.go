/*
    mysql data table model, mapping with xls sheet.
*/

package model

type SampleModel struct {
    Sid int `gorm:"primary_key"`
    Shortcut string
    Content string
}

func (SampleModel) TableName() string {
    return "tb_sample_template"
}
