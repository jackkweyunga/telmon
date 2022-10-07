package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Port      int      `mapstructure:"port"`
	Addr      string   `mapstructure:"addr"`
	Password  string   `mapstructure:"password"`
	Email     string   `mapstructure:"email"`
	Receivers []string `mapstructure:"receivers"`
}

func main() {

	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)

	var (
		port      = config.Port
		addr      = config.Addr
		fromEmail = config.Email
		password  = config.Password
		receivers = config.Receivers
	)

	c := make(chan os.Signal, 1)

	// Notify on any interrupt
	signal.Notify(c, os.Interrupt)

	isShuttingDown := false
	started := false

OUTER:
	for {
		if isShuttingDown {
			break
		}

		select {
		case s, ok := <-c:
			if ok {
				fmt.Println("Service is going down.")
				fmt.Printf("Received signal %v \n", s)
				isShuttingDown = true
				c = nil
				continue OUTER
			}
		default:
			if !started {
				println("Monitoring Service started successfully ...")
				started = true
			}

			t := time.Now()
			fmt.Printf("[%s] running \n", t.Format(time.RFC3339))
			Monitor(addr, port, fromEmail, password, receivers)
			time.Sleep(10000 * time.Millisecond)

		}
	}
}

var vp *viper.Viper

func LoadConfig() (Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName(".telmon-config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
