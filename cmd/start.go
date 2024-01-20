package cmd

import (
	"net/http"
	"time"

	"github.com/capella/cdive/controllers"
	"github.com/capella/cdive/models"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the http server.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := controllers.LoadConfig()
		if err != nil {
			logrus.Panic(err)
		}

		var dbConn gorm.Dialector
		if config.Server.DSN != nil {
			dbConn = mysql.Open(*config.Server.DSN)
		} else {
			dbConn = sqlite.Open("test.db")
		}

		db, err := gorm.Open(dbConn)
		if err != nil {
			logrus.Panic(err)
		}

		go models.AutoMigrate(db, config.Server.Secret)

		c := controllers.NewController(db, &config)
		router := c.Router()

		logrus.WithField("address", config.Server.Address).Info("starting server")
		csrfRouter := csrf.Protect([]byte(config.Server.Secret))(router)
		srv := &http.Server{
			Handler: csrfRouter,
			Addr:    config.Server.Address,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		logrus.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
