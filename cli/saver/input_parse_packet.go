package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type dbGpsPoint struct {
	Dsysdata time.Time `db:"DSYSDATA"`
	Ddata    time.Time `db:"DDATA"`
	Ntime    int       `db:"NTIME"`
	Ctime    string    `db:"CTIME"`
	Cid      string    `db:"CID"`
	Cip      string    `db:"CIP"`
	// NNum int
	Clatitude  string `db:"CLATITUDE"`
	Cns        string `db:"CNS"`
	Clongitude string `db:"CLONGTITUDE"`
	Cew        string `db:"CEW"`
	Ccurse     string `db:"CCURSE"`
	Cspeed     string `db:"CSPEED"`
	Csatel     string `db:"CSATEL"`
	Cdatavalid string `db:"CDATAVALID"`
}

type egtsParsePacket struct {
	Client              uint32         `json:"client"`
	ClientIP            string         `json:"client_ip"`
	PacketID            uint32         `json:"packet_id"`
	NavigationTimestamp int64          `json:"navigation_unix_time"`
	ReceivedTimestamp   int64          `json:"received_unix_time"`
	Latitude            float64        `json:"latitude"`
	Longitude           float64        `json:"longitude"`
	Speed               uint16         `json:"speed"`
	Pdop                uint16         `json:"pdop"`
	Hdop                uint16         `json:"hdop"`
	Vdop                uint16         `json:"vdop"`
	Nsat                uint8          `json:"nsat"`
	Ns                  uint16         `json:"ns"`
	Course              uint8          `json:"course"`
	GUID                uuid.UUID      `json:"guid"`
	AnSensors           []anSensor     `json:"an_sensors"`
	LiquidSensors       []liquidSensor `json:"liquid_sensors"`
}

func (eep *egtsParsePacket) ToDBGpsPoint() (*dbGpsPoint, error) {
	ndt := time.Unix(eep.NavigationTimestamp, 0)
	year, month, day := ndt.Date()
	hour, min, sec := ndt.Clock()
	point := &dbGpsPoint{
		Dsysdata:   time.Unix(eep.ReceivedTimestamp, 0),
		Ddata:      time.Date(year, month, day, 0, 0, 0, 0, time.Local),
		Ntime:      hour*60*60 + min*60 + sec,
		Ctime:      ndt.Format("15:04:05"),
		Cid:        strconv.FormatUint(uint64(eep.Client), 10),
		Cip:        eep.ClientIP,
		Clatitude:  fmt.Sprintf("%2.6f", eep.Latitude),
		Cns:        "N",
		Clongitude: fmt.Sprintf("%2.6f", eep.Longitude),
		Cew:        "E",
		Ccurse:     strconv.Itoa(int(eep.Course)),
		Cspeed:     strconv.Itoa(int(eep.Speed)),
		Csatel:     fmt.Sprintf("%d", eep.Nsat),
		Cdatavalid: "V",
	}
	return point, nil
}
func (eep *egtsParsePacket) ToBytes() ([]byte, error) {
	return json.Marshal(eep)
}

func (eep *egtsParsePacket) FromBytes(data []byte) (egtsParsePacket, error) {
	result := egtsParsePacket{}
	err := json.Unmarshal(data, &result)
	return result, err
}

type liquidSensor struct {
	SensorNumber uint8  `json:"sensor_number"`
	ErrorFlag    string `json:"error_flag"`
	ValueMm      uint32 `json:"value_mm"`
	ValueL       uint32 `json:"value_l"`
}

type anSensor struct {
	SensorNumber uint8  `json:"sensor_number"`
	Value        uint32 `json:"value"`
}
