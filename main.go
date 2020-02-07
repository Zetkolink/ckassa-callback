package main

import (
	"ckassa-callback/dataStore/sqlDb"
	"ckassa-callback/endpoints/rest"
	"ckassa-callback/pkg/graceful"
	"ckassa-callback/pkg/logger"
	"ckassa-callback/pkg/middlewares"
	"ckassa-callback/usecases/cardCallback"
	"ckassa-callback/usecases/paymentCallback"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

func main() {

	cfg := loadConfig()

	lg := logger.New(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", cfg.DbConn.Host, cfg.DbConn.User, cfg.DbConn.Name, cfg.DbConn.Pass)

	db, err := sqlDb.Connect(dbUri, cfg.DbDriver)
	if err != nil {
		lg.Fatalf("DB connection failed %s", err.Error())
	}

	ccs := cardCallback.NewCardCallbacksStore(db)
	pcs := paymentCallback.NewPaymentCallbacksStore(db)

	restHandler := rest.New(lg, pcs, ccs)

	srv := setupServer(cfg, lg, restHandler)
	lg.Infof("listening for requests on %s...", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil {
		lg.Fatalf("http server exited: %s", err)
	}

}

func setupServer(cfg *Config, lg logger.Logger, rest http.Handler) *graceful.Server {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(rest)

	handler := middlewares.WithRecovery(lg, router)

	srv := graceful.NewServer(handler, cfg.GracefulTimeout*time.Second, os.Interrupt)
	srv.Log = lg.Errorf
	srv.Addr = cfg.Addr
	return srv
}

type Config struct {
	LogLevel        string        `yaml:"logLevel"`
	LogFormat       string        `yaml:"logFormat"`
	DbConn          Db            `yaml:"dbConn"`
	DbDriver        string        `yaml:"dbDriver"`
	GracefulTimeout time.Duration `yaml:"timeout"`
	Addr            string        `yaml:"addr"`
}

type Db struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

func loadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Read config error")
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
