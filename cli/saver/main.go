package main

import (
	"os"

	"github.com/labstack/gommon/log"
)

var (
	config settings
	logger *log.Logger
)

func main() {
	var (
		store Connector
	)
	logger = log.New("-")
	logger.SetHeader("${time_rfc3339_nano} ${short_file}:${line} ${level} -${message}")

	if len(os.Args) == 2 {
		if err := config.Load(os.Args[1]); err != nil {
			logger.Fatalf("Ошибка парсинга конфига: %v", err)
		}
	} else {
		logger.Fatalf("Не задан путь до конфига")
	}
	logger.SetLevel(config.getLogLevel())

	store = &NatsConnector{}

	if err := store.Init(config.Store); err != nil {
		logger.Fatal(err)
	}
	defer store.Close()

	db := &OracleConnector{}
	if err := db.Init(config.Oracle); err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	runServer(store, db)
}

func runServer(store Connector, db *OracleConnector) {
	recvCh := make(chan *egtsParsePacket)
	store.Bind(recvCh)

	for {
		point := <-recvCh
		go db.Save(point)
	}
}
