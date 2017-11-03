/*
	core logic for export xls sheet to mysql table.
*/

package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"net/http"
	"html/template"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/cls1991/xls2db-go/resource"
)

type Option struct {
	Name string
	Value string
}

type PageVariables struct {
	Status string
	Options []Option
	Message string
}

var options []Option
var db *gorm.DB

func main() {
	user, passwd, database := "root", "flyfishdb", "xls2db"
	var err error
    db, err = gorm.Open(
        "mysql", 
        fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, database),
    )
    if err != nil {
		log.Print("connect to mysql err: ", err)
    }
    defer db.Close()
    // debug mode
	db.LogMode(true)

	// init data source mappings
	options = []Option {
		Option {"sample(sample.xlsx)", "sample"},
	}

	log.Print("web server is listening at 0.0.0.0:5001...")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":5001", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = fmt.Sprintf("template/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	if err := t.Execute(w, data); err != nil {
		log.Print("template executing error: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	MyPageVariables := PageVariables {
		Options: options,
	}
	renderTemplate(w, "index.tmpl", MyPageVariables)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "POST" {
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Print("get form file err: ", err)
		}
		defer file.Close()
		dirName := "tmp/"
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			err = os.Mkdir(dirName, 0755)
			if err != nil {
				log.Print("create directory err: ", err)
			}
		}
		filePath := dirName + handler.Filename
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Print("write file err: ", err)
		}
		defer f.Close()
		io.Copy(f, file)
		log.Print(filePath)
		model := r.FormValue("model")
		// @TODO: auto mapping with resource
		// is ugly and dirty
		if model == "sample" {
			r := resource.New("sample", "sid", 0, 2)
			r.ImportData(db, filePath)
		}

		MyPageVariables := PageVariables {
			Status: "成功",
			Options: options,
		}
		renderTemplate(w, "index.tmpl", MyPageVariables)
	} else {
		log.Print("unknown http method: ", r.Method)
	}
}
