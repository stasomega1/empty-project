package projectname

import (
	"fmt"
	"net/http"
	"project/inetrnal/app/services"
	"project/inetrnal/app/store"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	logDateFormat = "02.01.2006 15:04:05"
)

func Start(config *Config) error {
	//db
	db, err := configureDB(config.DbConfig)
	if err != nil {
		return fmt.Errorf("configureDB: %v", err)
	}
	defer db.Close()
	//redis
	redisClient, err := configureRedis(config.RedisConfig)
	if err != nil {
		return fmt.Errorf("configureRedis: %v", err)
	}
	defer redisClient.Close()
	//rabbit
	rabbitPublisher, err := configureRabbit(config.MqConfig.GetConnectionUrl())
	if err != nil {
		return fmt.Errorf("configureRabbit: %v", err)
	}
	defer rabbitPublisher.StopPublishing()
	//store
	myStore := store.NewStore(db, redisClient, rabbitPublisher)
	//logger
	logger, err := configureLogger(config.LogLevel)
	if err != nil {
		return fmt.Errorf("configureLogger: %v", err)
	}
	//service
	projectnameService := services.NewProjectnameServiceService(myStore, logger, config.Domain)
	schedulerService := services.NewSchedulerService(myStore, logger, config.Domain)
	//servers
	srv := newServer(projectnameService, schedulerService, logger, config.ErrLevel)
	healthServer := configureHealthServer(db)
	//scheduler
	go srv.SchedulerService.Start()

	errgr := errgroup.Group{}
	errgr.Go(func() error {
		srv.Logger.Infof("Start health server on port %d", config.HealthCheckPort)
		return fmt.Errorf("healthServer.Start: %v",
			http.ListenAndServe(fmt.Sprintf(":%d", config.HealthCheckPort), healthServer))
	})
	errgr.Go(func() error {
		srv.Logger.Infof("Start api server on port %d", config.AppPort)
		return fmt.Errorf("healthServer.Start: %v",
			http.ListenAndServe(fmt.Sprintf(":%d", config.AppPort), srv))
	})

	return errgr.Wait()
}

func configureDB(cnf *DbConfig) (*sqlx.DB, error) {
	databaseUrl := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		cnf.Host,
		cnf.Port,
		cnf.DbName,
		cnf.SSLMode,
		cnf.Username,
		cnf.Password)
	db, err := sqlx.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cnf.MaxIdleCon)
	db.SetMaxOpenConns(cnf.MaxOpenCon)
	return db, nil
}

func configureRedis(cnf *RedisConfig) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%d", cnf.Host, cnf.Port)
	redisConfig := &redis.Options{
		Addr:     address,
		Password: cnf.Password,
		DB:       cnf.Db,
	}
	return redis.NewClient(redisConfig), nil
}

func configureRabbit(rabbitUrl string) (*store.RabbitPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(rabbitUrl, amqp091.Config{})
	if err != nil {
		return nil, fmt.Errorf("configureRabbit: %v", err)
	}

	return &store.RabbitPublisher{Publisher: publisher}, nil
}

func configureLogger(level string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	customFormatter := &logrus.TextFormatter{}
	customFormatter.TimestampFormat = logDateFormat
	customFormatter.FullTimestamp = true
	logger.SetFormatter(customFormatter)
	return logger, nil
}

func configureHealthServer(DB *sqlx.DB) *HealthServer {
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(1000))
	health.AddReadinessCheck("database-ready", healthcheck.DatabasePingCheck(DB.DB, 3*time.Second))
	return NewHealthServer(health)
}
