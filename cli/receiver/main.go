package main

import (
	"net"
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

	runServer(config.getListenAddress(), store)
}

func runServer(srvAddress string, store Connector) {
	l, err := net.Listen("tcp", srvAddress)
	if err != nil {
		logger.Fatalf("Не удалось открыть соединение: %v", err)
	}
	defer l.Close()

	logger.Infof("Запущен сервер %s...", srvAddress)
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorf("Ошибка соединения: %v", err)
		} else {
			go handleRecvPkg(conn, store)
		}
	}
}
