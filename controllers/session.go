package controllers

import (
	"context"
	"net/http"

	"github.com/capella/cdive/models"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type contextKey string

const UserContextKey contextKey = "user"
const SessionsUserIDKey string = "id"
const SessionsName string = "user-session"

func getContextUser(ctx context.Context) (user *models.User) {
	if user, ok := ctx.Value(UserContextKey).(*models.User); ok && user != nil {
		return user
	}
	return nil
}

func (c *controllers) store() sessions.Store {
	return sessions.NewCookieStore([]byte(c.Config.Server.Secret))
}

func (c *controllers) session(r *http.Request) (*sessions.Session, error) {
	return c.store().Get(r, SessionsName)
}

func (c *controllers) sessionUserID(r *http.Request) *uint {
	session, err := c.session(r)
	if session == nil || err != nil {
		return nil
	}
	rawID, ok := session.Values[SessionsUserIDKey]
	if !ok {
		return nil
	}
	if id, ok := rawID.(uint); ok {
		return &id
	}
	return nil
}

func (c *controllers) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := c.sessionUserID(r); id != nil {
			user := &models.User{
				Model: gorm.Model{ID: *id},
			}
			c.DB.First(user)
			ctx := context.WithValue(r.Context(), UserContextKey, user)
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
	err = decoder.Decode(&login, r.PostForm)
	if err != nil {
		logrus.Error(err)
	}

	user := &models.User{
		Email: login.Email,
	}
	err = user.SetPassword(login.Password, c.Config.Server.Secret)
	if err != nil {
		logrus.Error(err)
	}

	result := c.DB.Where(user).First(&user)

	if result.Error != nil {
		formErrors = append(formErrors, "Username and password not found.")
	} else if login.Email == user.Email {
		session, err := c.session(r)
		if err != nil {
			logrus.Error(err)
		}
		session.Values[SessionsUserIDKey] = user.ID
		err = session.Save(r, w)
		if err != nil {
			logrus.Error(err)
			formErrors = append(formErrors, "Fail to save session Cookie.")
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	c.renderTemplate("login", formErrors, nil)(w, r)
}

func (c *controllers) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := c.session(r)
	if session != nil && err == nil {
		// MaxAge<0 means delete cookie immediately.
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			logrus.Error(err)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c contextKey) String() string {
	return string(c)
}
