package main

import (
	"log"
	"fmt"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-handler-influxdb/lib"
	"github.com/qnib/qframe-filter-docker-stats/lib"
	"github.com/qnib/qframe-collector-docker-events/lib"
	"github.com/qnib/qframe-collector-docker-stats/lib"
)

func Run(qChan qtypes.QChan, cfg config.Config, name string) {
	p, _ := qframe_handler_influxdb.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"handler.influxdb.database": "qframe",
		"handler.influxdb.host": "172.17.0.1",
		"handler.influxdb.inputs": "container-stats",
		"handler.influxdb.pattern": "%{INT:number}",
		"handler.influxdb.ticker-sec": "5",
		"handler.influxdb.batch-size": "100",
		"collector.docker-events.docker-host": "unix:///var/run/docker.sock",
		"filter.container-stats.inputs": "docker-stats",
		"log.level": "info",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	phi, err := qframe_handler_influxdb.New(qChan, *cfg, "influxdb")
	if err != nil {
		log.Printf("[EE] Failed to create filter: %v", err)
		return
	}
	go phi.Run()
	// Start filter
	pfc, err := qframe_filter_docker_stats.New(qChan, *cfg, "container-stats")
	if err != nil {
		log.Printf("[EE] Failed to docker-stats filter: %v", err)
		return
	}
	go pfc.Run()
	// start docker-events
	pe, err := qframe_collector_docker_events.New(qChan, *cfg, "docker-events")
	if err != nil {
		log.Printf("[EE] Failed to docker-event collector: %v", err)
		return
	}
	go pe.Run()
	// start docker-stats
	pds, err := qframe_collector_docker_stats.New(qChan, *cfg, "docker-stats")
	if err != nil {
		log.Printf("[EE] Failed to docker-stats collector: %v", err)
		return
	}
	go pds.Run()
	dc := qChan.Data.Join()
	for {
		select {
		case msg := <-dc.Read:
			switch msg.(type) {
			case qtypes.Metric:
				qm := msg.(qtypes.Metric)
				if qm.IsLastSource("container-stats") {
					msg := fmt.Sprintf("%s Metric %s: %v %s\n", qm.GetTimeRFC(), qm.Name, qm.Value, qm.GetDimensionList())
					if false {
						fmt.Printf(msg)
					}
				}
			}
		}
	}}
