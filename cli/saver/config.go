package main

/*
Описание конфигурационного файла
*/

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
)

type settings struct {
	Store  map[string]string
	Oracle map[string]string
	Log    logSection
}

func (c *settings) Load(confPath string) error {
	if _, err := toml.DecodeFile(confPath, c); err != nil {
		return fmt.Errorf("Ошибка разбора файла настроек: %v", err)
	}

	return nil
}

func (c *settings) getLogLevel() log.Lvl {
	return c.Log.getLevel()
}

type store struct {
	Host         string `toml:"host"`
	Port         string `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Exchange     string `toml:"exchange"`
	DeliveryMode string `toml:"delivery_mode"`
	Queue        string `toml:"queue"`
}

type oracle struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	host     string `toml:"host"`
	port     string `toml:"port"`
	service  string `toml:"service"`
}

type logSection struct {
	Level string
}

func (l *logSection) getLevel() log.Lvl {
	var lvl log.Lvl

	switch l.Level {
	case "DEBUG":
		lvl = log.DEBUG
		break
	case "INFO":
		lvl = log.INFO
		break
	case "WARN":
		lvl = log.WARN
		break
	case "ERROR":
		lvl = log.ERROR
		break
	default:
		lvl = log.INFO
	}
	return lvl
}
