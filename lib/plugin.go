package qframe_handler_influxdb

import (
	"strings"
	"log"
	"fmt"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
)

const (
	version = "0.0.0"
)

type Plugin struct {
    qtypes.Plugin

}

func New(qChan qtypes.QChan, cfg config.Config, name string) (Plugin, error) {
	var err error
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, name, version),
	}
	return p, err
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	log.Printf("[II] Start log handler %sv%s", p.Name, version)
	bg := p.QChan.Data.Join()
	inStr, err := p.Cfg.String(fmt.Sprintf("handler.%s.inputs", p.Name))
	if err != nil {
		inStr = ""
	}
	inputs := strings.Split(inStr, ",")
	srcSuccess, err := p.Cfg.BoolOr(fmt.Sprintf("handler.%s.source-success", p.Name), true)
	for {
		val := bg.Recv()
		qm := val.(qtypes.QMsg)
		if len(inputs) != 0 && ! qutils.IsLastSource(inputs, qm.Source) {
			continue
		}
		if qm.SourceSuccess != srcSuccess {
			continue
		}
		fmt.Printf("%s %-7s sType:%-6s sName:[%d]%-10s %s\n", qm.TimeString(), qm.LogString(), qm.Type, qm.SourceID, qm.Source, qm.Msg)
	}
}
