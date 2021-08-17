package main

import "github.com/nats-io/nats.go"

//Connector интерфейс для подключения внешних хранилищ
type Connector interface {
	// установка соединения с хранилищем
	Init(map[string]string) error

	// Подписаться на пакеты
	Bind(recvCh chan *egtsParsePacket) (*nats.Subscription, error)

	//закрытие соединения с хранилищем
	Close() error
}
