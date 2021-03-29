package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

// App config struct
type Config struct {
	IsDebug         bool
	RabbitMQ        RabbitMQ
	Exchange        Exchange
	Queue           Queue
	Telegram        Telegram
	Postgres        Postgres
	GoogleMapApiKey string
}

//RabbitMQ
type RabbitMQ struct {
	Host           string
	Port           string
	User           string
	Password       string
	WorkerPoolSize int
}

// Postgres
type Postgres struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
}

//Telegram
type Telegram struct {
	Token string
}

//Exchanges
type Exchange struct {
	CommandUsersGet    string
	CommandUsersCreate string
}

//Queues
type Queue struct {
	UsersRepoCommandUsersGet            string
	UsersRepoCommandUsersCreate         string
	UsersRepoCommandUsersUpdateTimezone string
	ParsingCommandTextParse             string
	ReminderCommandReminderAdd          string
	ReminderCommandRemindersGet          string
	BotEventNotifySend                  string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	
	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}
	
	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	
	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	
	if err := parseEnv(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func GetConfigPath(isDebug string) string {
	if isDebug == "false" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

func parseEnv(cfg *Config) error {
	if os.Getenv("RabbitMQPassword") == "" {
		return errors.New("rabbitmq password not found")
	} else {
		cfg.RabbitMQ.Password = os.Getenv("RabbitMQPassword")
	}
	
	if os.Getenv("IsDebug") == "" {
		return errors.New("IsDebug not found in .env file")
	} else {
		if os.Getenv("IsDebug") != "true" && os.Getenv("IsDebug") != "false" {
			return errors.New("incorrect format IsDebug in .env file")
		}
		if os.Getenv("IsDebug") == "true" {
			cfg.IsDebug = true
		}
		if os.Getenv("IsDebug") == "false" {
			cfg.IsDebug = false
		}
	}
	
	if os.Getenv("TGToken") == "" {
		return errors.New("tg token not found")
	} else {
		cfg.Telegram.Token = os.Getenv("TGToken")
	}
	if os.Getenv("GoogleMapApiKey") == "" {
		return errors.New("tg token not found")
	} else {
		cfg.GoogleMapApiKey = os.Getenv("GoogleMapApiKey")
	}
	
	if os.Getenv("PostgresqlPassword") == "" {
		return errors.New("postgresql password not found")
	} else {
		cfg.Postgres.PostgresqlPassword = os.Getenv("PostgresqlPassword")
	}
	
	return nil
}
