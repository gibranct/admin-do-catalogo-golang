package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	usecase "github.com.br/gibranct/admin_do_catalogo/internal/usecases"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config
	logger   *slog.Logger
	useCases usecase.UseCases
}

func main() {
	cfg, err := GetEnvs()

	if err != nil {
		panic(err)
	}

	flag.IntVar(&cfg.port, "port", cfg.port, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", cfg.db.dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", cfg.db.maxOpenConns, "Postgres max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", cfg.db.maxIdleConns, "Postgres max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", cfg.db.maxIdleTime, "Postgres max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(*cfg)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool stablished")

	app := &application{
		logger:   logger,
		config:   *cfg,
		useCases: usecase.NewUseCases(db),
	}

	app.server()
}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxIdleConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func GetEnvs() (*config, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))

	fmt.Println(os.Getenv("DB_DSN"))

	if err != nil {
		return nil, err
	}

	maxOpenConn, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONN"))

	if err != nil {
		return nil, err
	}

	maxIdleConn, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONN"))

	if err != nil {
		return nil, err
	}

	maxIdleTime, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_TIME_IN_MINUTES"))

	if err != nil {
		return nil, err
	}

	return &config{
		port: port,
		db: struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  time.Duration
		}{
			dsn:          os.Getenv("DB_DSN"),
			maxOpenConns: maxOpenConn,
			maxIdleConns: maxIdleConn,
			maxIdleTime:  time.Duration(maxIdleTime) * time.Minute,
		},
	}, nil
}
