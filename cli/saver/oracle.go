package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-oci8"
)

type OracleConnector struct {
	db     *sqlx.DB
	stmt   *sqlx.NamedStmt
	config map[string]string
}

func (c *OracleConnector) Init(cfg map[string]string) error {
	var err error
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}
	c.config = cfg
	c.db, err = sqlx.Open("oci8", fmt.Sprintf(`%s/%s@%s:%s/%s?PROTOCAL=TCP`,
		c.config["user"],
		c.config["password"],
		c.config["host"],
		c.config["port"],
		c.config["service"],
	),
	)
	if err != nil {
		return err
	}
	logger.Infof("Соединение с базой %s установлено.", c.config["host"])
	c.stmt, err = c.db.PrepareNamed(`INSERT INTO DISPATCHER.TGPSDATA (DSYSDATA, DDATA, NTIME, CTIME, CID,
		CIP, CLATITUDE, CNS, CLONGTITUDE, CEW, CCURSE, CSPEED, CSATEL, CDATAVALID) VALUES (:DSYSDATA, :DDATA, :NTIME,
		:CTIME, :CID, :CIP, :CLATITUDE, :CNS, :CLONGTITUDE, :CEW, :CCURSE, :CSPEED, :CSATEL, :CDATAVALID)`)
	if err != nil {
		return err
	}
	return c.db.Ping()
}

func (c *OracleConnector) Save(packet *egtsParsePacket) error {
	logger.Debugf("Получена отметка %s от %d (%s)", packet.GUID.String(), packet.Client, packet.ClientIP)
	point, _ := packet.ToDBGpsPoint()
	_, err := c.stmt.Exec(&point)
	if err != nil {
		logger.Errorf("Ошибка при сохранении в базу данных: %v", err)
	}
	logger.Debugf("Сохранена отметка %s от %s (%s)", packet.GUID.String(), point.Cid, point.Cip)
	return err
}

func (c *OracleConnector) Close() error {
	logger.Infof("База данных %s отключена.", c.config["host"])
	return c.db.Close()
}
