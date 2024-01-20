package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

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
}

func (c *controllers) renderTemplate(
	templateName string,
	formErrors []string,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getContextUser(r.Context())
		if user != nil {
			user = nil
		}

		viewData := &templateData{
			DB:         c.DB,
			User:       user,
			CSRFField:  csrf.TemplateField(r),
			Config:     *c.Config,
			FormErrors: formErrors,
		}

		tmpl := template.Must(
			template.New("layout.html").
				Funcs(template.FuncMap{"form": r.PostFormValue}).
				ParseFiles(
					"views/layout.html",
					fmt.Sprintf("views/%s", templateName),
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
	c.renderTemplate("index.html", nil)(w, r)
}

func NewController(db *gorm.DB, config *Config) *controllers {
	return &controllers{
		DB:     db,
		Config: config,
	}
}
