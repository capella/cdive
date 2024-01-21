package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/capella/cdive/controllers"
	"github.com/capella/cdive/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// Set the correct path for templates
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestUnlogUserRedirect(t *testing.T) {
	dbFile, _ := os.CreateTemp("", "cdive-test-")
	dbConn := sqlite.Open(dbFile.Name())
	db, _ := gorm.Open(dbConn)
	models.AutoMigrate(db, "test")
	defer os.RemoveAll(dbFile.Name())

	c := controllers.NewController(db, &controllers.Config{})
	router := c.Router()
	server := httptest.NewServer(router)
	defer server.Close()

	res, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	println()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected %d status code and got %d", http.StatusFound, res.StatusCode)
	}

	url := res.Request.URL.Path
	if url != "/login" {
		t.Errorf("expected %s to ends with /login", url)
	}
}
