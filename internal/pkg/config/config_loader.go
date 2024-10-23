package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"

	"github.com/josephlbailey/alert-service/internal/pkg/path"
)

func LoadConfig[T interface{}](name string, env ...string) (config T) {

	p, err := path.Determine("config")
	if err != nil {
		log.Printf("Unable to determine config path: %v\n", err)
	}

	c := viper.New()
	c.AddConfigPath(p + "/config")

	c.SetConfigType("yaml")
	c.SetConfigName(name)

	log.Println("Loading config...")
	if err := c.ReadInConfig(); err == nil {
		if err := c.Unmarshal(&config); err != nil {
			log.Printf("Unable to unmarshal config %v\n", err)
		}
	} else {
		log.Printf("Unable to load config: %v\n", err)
	}

	if env != nil && len(env) > 0 {
		e := viper.New()
		e.AddConfigPath(p + "/config/" + env[0] + "/")

		e.SetConfigType("yaml")
		e.SetConfigName(name)

		log.Printf("Loading %v config...\n", env[0])
		if err := e.ReadInConfig(); err == nil {
			if err := e.Unmarshal(&config); err != nil {
				log.Printf("Unable to unmarshal %v config %v\n", env, err)
			}
		} else {
			log.Printf("Unable to load %v config: %v\n", env[0], err)
		}

	}

	s := viper.New()
	s.AddConfigPath(p + "/secret")

	s.SetConfigType("yaml")
	s.SetConfigName(name)

	log.Print("Loading secret...")
	if err := s.ReadInConfig(); err == nil {
		if err := s.Unmarshal(&config); err != nil {
			log.Printf("Unable to unmarshal secret %v\n", err)
		}
	} else {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Printf("No secrets configured, bypassing secret unmarshal.\n")
		} else {
			log.Printf("Unable to load secret: %v\n", err)

		}
	}

	return
}
