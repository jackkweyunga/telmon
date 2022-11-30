package monitor

import (
	log "github.com/jackkweyunga/telmon/logging"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

type Restore struct {
	Exec string `mapstructure:"exec"`
	Cmd  string `mapstructure:"cmd"`
}

type Config struct {
	Port     int     `mapstructure:"port"`
	Addr     string  `mapstructure:"addr"`
	interval string  `mapstructure:"interval"`
	Notify   Notify  `mapstructure:"notify"`
	Restore  Restore `mapstructure:"restore"`
}

func Play() {

	log.Setup()

	Init()

	config, err := LoadConfig()
	if err != nil {
		log.Log.Fatal(err)
	}

	var (
		port = config.Port
		addr = config.Addr
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
				log.Log.Println("Service is going down.")
				log.Log.Printf("Received signal %v \n", s)
				isShuttingDown = true
				c = nil
				continue OUTER
			}
		default:
			if !started {
				log.Log.Println("[Telmon] Monitoring Service started successfully ...")
				started = true
			}

			Monitor(addr, port, config.Notify, config.Restore)
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
