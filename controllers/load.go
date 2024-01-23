package controllers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"reflect"

	spring "github.com/Masterminds/sprig/v3"
	"github.com/capella/cdive/models"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(false)
}

type controllers struct {
	DB     *gorm.DB
	Config *Config
}

type templateData struct {
	User *models.User

	CSRFField template.HTML
	Config    Config

	Form       url.Values
	FormErrors []string
	Controller any
}

func (c *controllers) renderTemplate(
	templateName string,
	formErrors []string,
	data any,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getContextUser(r.Context())

		viewData := &templateData{
			User:       user,
			CSRFField:  csrf.TemplateField(r),
			Config:     *c.Config,
			FormErrors: formErrors,
			Controller: data,
		}

		renderErr := renderTemplate(templateName, viewData, w, r)
		if renderErr == nil {
			return
		}

		if tmpl, err := template.ParseFiles("views/error.html"); err == nil {
			_ = tmpl.Execute(w, renderErr)
		} else {
			logrus.Error("err")
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func renderTemplate(
	templateName string,
	data *templateData,
	w io.Writer,
	r *http.Request,
) (err error) {
	files := []string{
		"views/layout.html",
		"views/modal.html",
		"views/inputs/text-input.html",
		fmt.Sprintf("views/%s.html", templateName),
	}
	functions := template.FuncMap{
		"form":      r.PostFormValue,
		"hasValues": hasValues,
	}

	tmpl, err := template.
		New("layout.html").
		Funcs(functions).
		Funcs(spring.FuncMap()).
		ParseFiles(files...)
	if err != nil {
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	return
}

func (c *controllers) Home(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate("index", nil, nil)(w, r)
}

func NewController(db *gorm.DB, config *Config) *controllers {
	return &controllers{
		DB:     db,
		Config: config,
	}
}

func hasValues(name string, data any) bool {
	items := reflect.ValueOf(data)
	if items.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		if item.Kind() == reflect.Struct {
			v := reflect.Indirect(item)
			for j := 0; j < v.NumField(); j++ {
				if (v.Type().Field(j).Name) == name {
					return true
				}
			}
		}
	}
	return false
}
