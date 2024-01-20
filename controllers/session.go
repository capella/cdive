package controllers

import (
	"context"
	"net/http"

	"github.com/capella/cdive/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ContextUserKey string = "user"
const SessionsIDKey string = "id"

func getContextUser(ctx context.Context) (user *models.User) {
	if user, ok := ctx.Value(ContextUserKey).(*models.User); ok && user != nil {
		return user
	}
	return nil
}

func (c *controllers) SessionMiddleware(next http.Handler) http.Handler {
	store := sessions.NewCookieStore([]byte(c.Config.Server.Secret))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, ContextUserKey)
		if id, ok := session.Values[SessionsIDKey].(uint); ok {
			user := &models.User{
				Model: gorm.Model{ID: id},
			}
			c.DB.First(user)
			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (c *controllers) Internal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := getContextUser(r.Context()); user != nil {
			c.DB.First(user)
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}

func (c *controllers) Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := getContextUser(r.Context()); user != nil && user.Admin {
			c.DB.First(user)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Admin permission requires", http.StatusForbidden)
		}
	})
}

type LoginFields struct {
	Password string
	Email    string
}

func (c *controllers) Login(w http.ResponseWriter, r *http.Request) {
	if user := getContextUser(r.Context()); user != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	c.renderTemplate("login", nil, nil)(w, r)
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
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	c.renderTemplate("login", formErrors, nil)(w, r)
}

func (c *controllers) Logout(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte(c.Config.Server.Secret))
	session, _ := store.Get(r, "user")
	session.Values = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
