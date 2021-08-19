package main

/*
Плагин для работы с NATS.
Плагин отправляет пакет в топик NATS messaging system.

Раздел настроек, которые должны отвечають в конфиге для подключения плагина:

[store]
plugin = "nats.so"
servers = "nats://localhost:1222, nats://localhost:1223, nats://localhost:1224"
topic = "receiver"
*/

import (
	"fmt"

	natsLib "github.com/nats-io/nats.go"
)

type NatsConnector struct {
	connection  *natsLib.Conn
	econnection *natsLib.EncodedConn
	config      map[string]string
}

func (c *NatsConnector) Init(cfg map[string]string) error {
	var (
		err error
	)
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}
	c.config = cfg

	var options = make([]natsLib.Option, 3)

	options = append(options, natsLib.Name(fmt.Sprintf("Saver handler, topic: %s", c.config["topic"])))

	if user, uOk := c.config["user"]; uOk {
		if password, pOk := c.config["password"]; pOk {
			options = append(options, natsLib.UserInfo(user, password))
		}
	}

	if c.connection, err = natsLib.Connect(c.config["servers"], options...); err != nil {
		return fmt.Errorf("Ошибка подключения к nats шине: %v", err)
	}
	if c.econnection, err = natsLib.NewEncodedConn(c.connection, natsLib.JSON_ENCODER); err != nil {
		return fmt.Errorf("Ошибка подключения NewEncodedConn к nats шине: %v", err)
	}
	return err
}

func (c *NatsConnector) Bind(recvCh chan *egtsParsePacket) (*natsLib.Subscription, error) {
	return c.econnection.BindRecvQueueChan(c.config["topic"], "save_orcl", recvCh)
}

func (c *NatsConnector) Close() error {
	c.connection.Close()
	return nil
}
