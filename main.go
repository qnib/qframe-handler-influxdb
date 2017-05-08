package main

import (
	"log"
	"fmt"
	"os"

	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-handler-influxdb/lib"
	"github.com/qnib/qframe-filter-docker-stats/lib"
	"github.com/qnib/qframe-collector-docker-events/lib"
	"github.com/qnib/qframe-collector-docker-stats/lib"
	"github.com/qnib/qframe-collector-internal/lib"
)

func Run(qChan qtypes.QChan, cfg config.Config, name string) {
	p, _ := qframe_handler_influxdb.New(qChan, cfg, name)
	p.Run()
}

func check_err(pname string, err error) {
	if err != nil {
		log.Printf("[EE] Failed to create %s plugin: %s", pname, err.Error())
		os.Exit(1)
	}
}

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"handler.influxdb.database": "qframe",
		"handler.influxdb.host": "172.17.0.1",
		"handler.influxdb.inputs": "internal,container-stats",
		"handler.influxdb.pattern": "%{INT:number}",
		"handler.influxdb.ticker-msec": "2000",
		"handler.influxdb.batch-size": "500",
		"collector.docker-events.docker-host": "unix:///var/run/docker.sock",
		"filter.container-stats.inputs": "docker-stats",
		"log.level": "info",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	// Start handler
	phi, err := qframe_handler_influxdb.New(qChan, *cfg, "influxdb")
	check_err(phi.Name, err)
	go phi.Run()
	// Start filter
	pfc, err := qframe_filter_docker_stats.New(qChan, *cfg, "container-stats")
	check_err(pfc.Name, err)
	go pfc.Run()
	// start docker-events
	pe, err := qframe_collector_docker_events.New(qChan, *cfg, "docker-events")
	check_err(pe.Name, err)
	go pe.Run()
	// start docker-stats
	pds, err := qframe_collector_docker_stats.New(qChan, *cfg, "docker-stats")
	check_err(pds.Name, err)
	go pds.Run()
	pci, err := qframe_collector_internal.New(qChan, *cfg, "internal")
	check_err(pci.Name, err)
	go pci.Run()
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
