package qframe_handler_influxdb

import (
	"fmt"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/qnib/qframe-types"
)

const (
	version = "0.0.5"
	pluginTyp = "handler"
)

type Plugin struct {
    qtypes.Plugin
	cli client.Client
	metricCount int

}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, name, version),
		metricCount: 0,
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() {
	host := p.CfgStringOr("host", "localhost")
	port, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.port", p.Name), "8086")
	username, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.username", p.Name), "root")
	password, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.password", p.Name), "root")
	addr := fmt.Sprintf("http://%s:%s", host, port)
	cli := client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	}
	var err error
	p.cli, err = client.NewHTTPClient(cli)
	if err != nil {
		p.Log("error", fmt.Sprintf("Error during connection to InfluxDB '%s': %v", addr, err))
	} else {
		p.Log("info", fmt.Sprintf("Established connection to '%s", addr))
	}
}

func (p *Plugin) NewBatchPoints() client.BatchPoints {
	dbName := p.CfgStringOr("database", "qframe")
	dbPrec := p.CfgStringOr("precision", "s")
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: dbPrec,
	})
	if err != nil {
		p.Log("error", fmt.Sprintf("Not able to create BatchPoints: %v", err))
	}
	return bp

}

func (p *Plugin) WriteBatch(points client.BatchPoints) client.BatchPoints {
	err := p.cli.Write(points)
	if err != nil {
		p.Log("error", fmt.Sprintf("Not able to write BatchPoints: %v", err))
	}
	return p.NewBatchPoints()
}

func (p *Plugin) MetricsToBatchPoint(m qtypes.Metric) (pt *client.Point, err error) {
	fields := map[string]interface{}{
		"value": m.Value,
	}
	pt, err = client.NewPoint(m.Name, m.Dimensions, fields, m.Time)
	return
}
// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("info", fmt.Sprintf("Start log handler %sv%s", p.Name, version))
	batchSize := p.CfgIntOr("batch-size", 100)
	tick := p.CfgIntOr("ticker-sec", 1)

	p.Connect()
	bg := p.QChan.Data.Join()
	inputs := p.GetInputs()
	//srcSuccess, err := p.Cfg.BoolOr(fmt.Sprintf("handler.%s.source-success", p.Name), true)
	// Create a new point batch
	bp := p.NewBatchPoints()
	tickChan := time.NewTicker(time.Duration(tick)*time.Second).C
	skipTicker := false
	dims := map[string]string{
		"version": version,
		"plugin": p.Name,
	}
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			case qtypes.Metric:
				m := val.(qtypes.Metric)
				pt, err := p.MetricsToBatchPoint(m)
				if err != nil {
					p.Log("error", fmt.Sprintf("%v", err))
					continue
				}
				bp.AddPoint(pt)
				if ! m.InputsMatch(inputs) {
					continue
				}

				if len(bp.Points()) >= batchSize {
					bLen := len(bp.Points())
					p.Log("debug", fmt.Sprintf("Write batch of %d",bLen))
					p.metricCount += bLen
					p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.size", qtypes.Gauge, float64(bLen), dims, time.Now(), false))
					bp = p.WriteBatch(bp)
					skipTicker = true
				}
			}
		case <- tickChan:
			if ! skipTicker {
				bLen := len(bp.Points())
				p.Log("debug", fmt.Sprintf("Ticker: Write batch of %d",bLen))
				p.metricCount += bLen
				bp = p.WriteBatch(bp)
				p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.size", qtypes.Gauge, float64(bLen), dims, time.Now(), false))
			}
			skipTicker = false
			p.QChan.Data.Send(qtypes.NewExt(p.Name, "influxdb.batch.count", qtypes.Counter, float64(p.metricCount), dims, time.Now(), false))
		}
	}
}
