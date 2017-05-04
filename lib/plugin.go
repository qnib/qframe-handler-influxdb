package qframe_handler_influxdb

import (
	"fmt"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/qnib/qframe-types"
)

const (
	version = "0.0.1"
	pluginTyp = "handler"
)

type Plugin struct {
    qtypes.Plugin
	cli client.Client

}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, name, version),
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

/*
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

func (p *Plugin) CreateDockerStatsPoints(cs qtypes.ContainerStats) (pt *client.Point, err error) {
	// Create a point and add to batch
	mTags := map[string]string{
		"image_name": cs.Container.Image,
		"container_id": cs.Container.ID,
		"container_name": strings.TrimPrefix(cs.Container.Names[0], "/"),
		"container_cmd": cs.Container.Command,
	}
	cStats := qtypes.NewCPUStats(cs.Stats)
	fields := map[string]interface{}{
		"user": cStats.UsageInUsermodePercentage,
		"kernel": cStats.UsageInKernelmodePercentage,
		"system": cStats.SystemUsagePercentage,
	}
	pt, err = client.NewPoint("cpu_percent", mTags, fields, cStats.Time)
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	return pt, err
}

func (p *Plugin) CreateDockerStatsMemory(cs qtypes.ContainerStats) (pt *client.Point, err error) {
	// Create a point and add to batch
	mTags := map[string]string{
		"image_name": cs.Container.Image,
		"container_id": cs.Container.ID,
		"container_name": strings.TrimPrefix(cs.Container.Names[0], "/"),
		"container_cmd": cs.Container.Command,
	}
	mStats := qtypes.NewMemoryStats(cs.Stats)
	fields := map[string]interface{}{
		"total_rss": mStats.TotalRss,
		"usage": mStats.Usage,
		"limit": mStats.Limit,
	}
	pt, err = client.NewPoint("memory", mTags, fields, mStats.Time)
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	return
}

func (p *Plugin) CreateDockerStatsMemoryPercent(cs qtypes.ContainerStats) (pt *client.Point, err error) {
	// Create a point and add to batch
	mTags := map[string]string{
		"image_name": cs.Container.Image,
		"container_id": cs.Container.ID,
		"container_name": strings.TrimPrefix(cs.Container.Names[0], "/"),
		"container_cmd": cs.Container.Command,
	}
	mStats := qtypes.NewMemoryStats(cs.Stats)
	if ! (reflect.TypeOf(mStats.UsageP).Kind() == reflect.Float64) || ! (reflect.TypeOf(mStats.TotalRssP).Kind() == reflect.Float64) {
		msg := fmt.Sprintf("Non float64 value: UsageP:%v || TotalRssP:%v", mStats.UsageP, mStats.TotalRssP)
		p.Log("error", msg)
		return pt, errors.New(msg)
	}
	fields := map[string]interface{}{
		"usage": mStats.UsageP,
		"total_rss": mStats.TotalRssP,
	}
	pt, err = client.NewPoint("memory_percent", mTags, fields, mStats.Time)
	if err != nil {
		p.Log("error", fmt.Sprintf("%v", err))
	}
	return
}
*/

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
	for {
		select {
		case val := <-bg.Read:
			switch val.(type) {
			/*case qtypes.QMsg:
				qm := val.(qtypes.QMsg)
				if len(inputs) != 0 && ! qutils.IsInput(inputs, qm.Source) {
					continue
				}
				if qm.SourceSuccess != srcSuccess {
					continue
				}
				switch qm.Data.(type){
				case qtypes.ContainerStats:
					cs := qm.Data.(qtypes.ContainerStats)
					pt, err := p.CreateDockerStatsPoints(cs)
					if err == nil {
						bp.AddPoint(pt)
					}
					pt, err = p.CreateDockerStatsMemoryPercent(cs)
					if err == nil {
						bp.AddPoint(pt)
					}
					pt, err = p.CreateDockerStatsMemory(cs)
					if err == nil {
						bp.AddPoint(pt)
					}
					bp = p.WriteBatch(bp)
				}
			*/
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
					p.Log("info", fmt.Sprintf("Write batch of %d",len(bp.Points())))
					bp = p.WriteBatch(bp)
					skipTicker = true
				}
			}
		case <- tickChan:
			if ! skipTicker {
				p.Log("info", fmt.Sprintf("Ticker: Write batch of %d",len(bp.Points())))
				bp = p.WriteBatch(bp)
			}
			skipTicker = false
		}
	}
}
