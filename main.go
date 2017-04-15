package main

import (
	"log"
	"time"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-handler-influxdb/lib"
)

func Run(qChan qtypes.QChan, cfg config.Config, name string) {
	p, _ := qframe_handler_influxdb.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"handler.test.pattern": "%{INT:number}",
		"handler.test.inputs": "test",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p, err := qframe_handler_influxdb.New(qChan, *cfg, "test")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	go p.Run()
	time.Sleep(2*time.Second)
	qm := qtypes.NewQMsg("test", "test")
	qm.Msg = "1"
	log.Println("Send message")
	qChan.Data.Send(qm)
	time.Sleep(1*time.Second)
}
