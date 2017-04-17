package qframe_handler_influxdb

import (
	"strings"
	"fmt"
	"strconv"
	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
)

const (
	version = "0.0.0"
)

type Plugin struct {
    qtypes.Plugin
	cli client.Client

}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, name, version),
	}
	return p, err
}

// Connect creates a connection to InfluxDB
func (p *Plugin) Connect() {
	host, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.host", p.Name), "localhost")
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

func (p *Plugin) CreatePoint(qm qtypes.QMsg) (*client.Point, error) {
	// Create a point and add to batch
	mName := qm.KV["name"]
	mValue := qm.KV["value"]
	mTags := map[string]string{}
	if val, ok := qm.KV["tags"] ; ok {
		for _, pair := range strings.Split(val, ",") {
			list := strings.Split(pair, "=")
			if len(list) == 2 {
				mTags[list[0]] = list[1]
			}
		}
	}
	val, _ := strconv.ParseFloat(mValue, 64)
	fields := map[string]interface{}{
		"value": val,
	}
	pt, err := client.NewPoint(mName, mTags, fields, qm.Time)
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	return pt, err
}

func (p *Plugin) NewBatchPoints() client.BatchPoints {
	dbName, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.database", p.Name), "qwatch")
	dbPrec, _ := p.Cfg.StringOr(fmt.Sprintf("handler.%s.precision", p.Name), "s")
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
// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("info", fmt.Sprintf("Start log handler %sv%s", p.Name, version))
	p.Connect()
	bg := p.QChan.Data.Join()
	inStr, err := p.Cfg.String(fmt.Sprintf("handler.%s.inputs", p.Name))
	if err != nil {
		inStr = ""
	}
	inputs := strings.Split(inStr, ",")
	srcSuccess, err := p.Cfg.BoolOr(fmt.Sprintf("handler.%s.source-success", p.Name), true)
	// Create a new point batch
	bp := p.NewBatchPoints()
	for {
		val := bg.Recv()
		qm := val.(qtypes.QMsg)
		if len(inputs) != 0 && ! qutils.IsInput(inputs, qm.Source) {
			continue
		}
		if qm.SourceSuccess != srcSuccess {
			continue
		}
		pt, _ := p.CreatePoint(qm)
		bp.AddPoint(pt)
		bp = p.WriteBatch(bp)
	}
}
