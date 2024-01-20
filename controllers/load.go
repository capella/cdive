package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"reflect"

	"github.com/capella/cdive/models"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var decoder = schema.NewDecoder()

type controllers struct {
	DB     *gorm.DB
	Config *Config
}

type templateData struct {
	DB   *gorm.DB
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
			DB:         c.DB,
			User:       user,
			CSRFField:  csrf.TemplateField(r),
			Config:     *c.Config,
			FormErrors: formErrors,
			Controller: data,
		}

		tmpl := template.Must(
			template.New("layout.gohtmltmpl").
				Funcs(template.FuncMap{
					"form":      r.PostFormValue,
					"hasValues": hasValues,
				}).
				ParseFiles(
					"views/layout.gohtmltmpl",
					fmt.Sprintf("views/%s.gohtmltmpl", templateName),
				),
		)

		err := tmpl.Execute(w, viewData)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
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
