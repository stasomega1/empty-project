package projectname

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"log"
)

var (
	ErrParseConfig = errors.New("error parse config for profile")
)

type Config struct {
	AppPort         int    `toml:"app_port" env:"APP_PORT,required"`
	HealthCheckPort int    `toml:"health_port" env:"HEALTH_CHECK_PORT,required"`
	LogLevel        string `toml:"log_level" env:"LOG_LEVEL,required"`
	ErrLevel        string `toml:"err_level" ENV:"ERR_LEVEL,required"`
	Domain          string `toml:"domain" env:"DOMAIN,required"`
	DbConfig        *DbConfig
	RedisConfig     *RedisConfig
	MqConfig        *MqConfig

	//TimeZone      string `env:"TZ,required"`
	//JwtConfig     *JwtConfig
	//FilenetConfig *FilenetConfig
}

type DbConfig struct {
	Host       string `toml:"host" env:"DB_HOST,required"`
	Port       int    `toml:"port" env:"DB_PORT,required"`
	DbName     string `toml:"db_name" env:"DB_NAME,required"`
	SSLMode    string `toml:"db_ssl_mode" env:"DB_SSLMODE,required"`
	Username   string `toml:"db_username" env:"DB_USERNAME_DSR,required"`
	Password   string `toml:"db_password" env:"DB_PASSWORD_DSR,required"`
	MaxIdleCon int    `toml:"max_idle_con" env:"DB_MAX_IDLE_CON,required"`
	MaxOpenCon int    `toml:"max_open_con" env:"DB_MAX_OPEN_CON,required"`
}

type JwtConfig struct {
	JwtKeyStr           string `toml:"jwt_key" env:"JWT_KEY,required"`
	JwtTokenLiveMinutes int    `toml:"jwt_token_live" env:"JWT_TOKEN_LIVE_MINUTES,required"`
	JwtKey              []byte
}

func (j *JwtConfig) SetJwtKey() {
	j.JwtKey = []byte(j.JwtKeyStr)
}

type FilenetConfig struct {
	Host     string `toml:"host" env:"FILENET_HOSTNAME,required"`
	Username string `toml:"username" env:"FILENET_USERNAME,required"`
	Password string `toml:"password" env:"FILENET_PASSWORD,required"`
	Sitename string `toml:"sitename" env:"FILENET_SITENAME,required"`
}

type RedisConfig struct {
	Host     string `toml:"host" env:"REDIS_HOST,required"`
	Port     int    `toml:"port" env:"REDIS_PORT,required"`
	Password string `toml:"pass" env:"REDIS_PASSWORD,required"`
	Db       int    `toml:"db" env:"REDIS_DB,required"`
}

type MqConfig struct {
	Host     string `toml:"host" env:"RABBIT_HOST,required"`
	Port     int    `toml:"port" env:"RABBIT_PORT,required"`
	Username string `toml:"user" env:"RABBIT_USR,required"`
	Password string `toml:"pass" env:"RABBIT_PASS,required"`
	Vhost    string `toml:"vhost" env:"RABBIT_VHOST,required"`
}

func (mq *MqConfig) GetConnectionUrl() string {
	if mq.Vhost != "" {
		return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", mq.Username, mq.Password, mq.Host, mq.Port, mq.Vhost)
	} else {
		return fmt.Sprintf("amqp://%s:%s@%s:%d", mq.Username, mq.Password, mq.Host, mq.Port)
	}
}

func NewConfig() *Config {
	return &Config{
		AppPort:         8080,
		HealthCheckPort: 8086,
		LogLevel:        "debug",
		DbConfig:        &DbConfig{},
		MqConfig:        &MqConfig{},
		RedisConfig:     &RedisConfig{},
		//JwtConfig:       &JwtConfig{},
		//FilenetConfig:   &FilenetConfig{},
	}
}

func ParseConfig(configPath, profile string) (*Config, error) {
	log.Printf("Load from profile %s", profile)

	cnf := NewConfig()

	switch profile {
	case "docker":
		if err := env.Parse(cnf); err != nil {
			log.Printf("error parsing config, error is %v", err)
			return nil, fmt.Errorf("parse config err: %v", err)
		}
	case "local":
		if _, err := toml.DecodeFile(configPath, cnf); err != nil {
			return nil, ErrParseConfig
		}
	}

	//cnf.JwtConfig.SetJwtKey()

	return cnf, nil

}
