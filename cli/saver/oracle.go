package main

import (
	"fmt"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

type OracleConnector struct {
	db     *sqlx.DB
	stmt   *sqlx.Stmt
	config map[string]string
}

func (c *OracleConnector) Init(cfg map[string]string) error {
	var err error
	if cfg == nil {
		return fmt.Errorf("Не корректная ссылка на конфигурацию")
	}
	c.config = cfg
	c.db, err = sqlx.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`,
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
	c.stmt, err = c.db.Preparex(`INSERT INTO DISPATCHER.TGPSDATA (DSYSDATA, DDATA, NTIME, CTIME, CID,
		CIP, CLATITUDE, CNS, CLONGTITUDE, CEW, CCURSE, CSPEED, CSATEL) VALUES (:DSYSDATA, :DDATA, :NTIME,
		:CTIME, :CID, :CIP, :CLATITUDE, :CNS, :CLONGTITUDE, :CEW, :CCURSE, :CSPEED, :CSATEL);`)
	if err != nil {
		return err
	}
	return c.db.Ping()
}

func (c *OracleConnector) Save(packet *egtsParsePacket) error {
	point, _ := packet.ToDBGpsPoint()
	_, err := c.stmt.Exec(&point)
	if err != nil {
		logger.Errorf("Ошибка при сохранении в базу данных: %v", err)
	}
	return err
}

func (c *OracleConnector) Close() error {
	return c.db.Close()
}
