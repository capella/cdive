package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/capella/cdive/models"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
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

func (c *controllers) renderTemplate(templateName string, formErrors []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		viewData := &templateData{
			DB:         c.DB,
			User:       nil,
			CSRFField:  csrf.TemplateField(r),
			Config:     *c.Config,
			FormErrors: formErrors,
		}
		store := sessions.NewCookieStore([]byte(c.Config.Server.Secret))
		session, err := store.Get(r, "user")
		if session != nil && err != nil {
			c.DB.First(viewData.User, session.Values["id"])
		}

		tmpl := template.Must(
			template.New("layout.html").
				Funcs(template.FuncMap{"form": r.PostFormValue}).
				ParseFiles(
					"views/layout.html",
					fmt.Sprintf("views/%s", templateName),
				),
		)

		err = tmpl.Execute(w, viewData)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

type LoginFields struct {
	Password string
	Email    string
}

func (c *controllers) Login(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate("login.html", nil)(w, r)
}

func (c *controllers) LoginPOST(w http.ResponseWriter, r *http.Request) {
	formErrors := []string{}

	err := r.ParseForm()
	if err != nil {
		logrus.Error(err)
		return
	}

	var login LoginFields
	decoder.Decode(&login, r.PostForm)

	user := &models.User{
		Email: models.S(login.Email),
	}
	err = user.SetPassword(login.Password, c.Config.Server.Secret)
	if err != nil {
		logrus.Error(err)
	}

	result := c.DB.Where(user).First(&user)

	if result.Error != nil {
		formErrors = append(formErrors, "Username and password not found.")
	} else if login.Email == user.Email.String {
		store := sessions.NewCookieStore([]byte(c.Config.Server.Secret))
		session, _ := store.Get(r, "user")
		session.Values["id"] = user.ID
		http.Redirect(w, r, "/", http.StatusFound)
	}

	c.renderTemplate("login.html", formErrors)(w, r)
}

func NewController(db *gorm.DB, config *Config) *controllers {
	return &controllers{
		DB:     db,
		Config: config,
	}
}
