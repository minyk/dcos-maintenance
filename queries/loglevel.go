package queries

import (
	"encoding/json"
	"fmt"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/minyk/dcos-maintenance/client"
	"time"
)

type Loglevel struct {
	PrefixMaster func() string
	PrefixAgent  func(id string) string
}

func NewLoglevel() *Loglevel {
	return &Loglevel{
		PrefixMaster: func() string { return "/mesos/api/v1/" },
		PrefixAgent:  func(id string) string { return fmt.Sprintf("/agent/%s/api/v1", id) },
	}
}

func (q *Loglevel) SetLoglevel(agentid string, level int, duration time.Duration) error {

	body := agent.Call{
		Type: agent.Call_SET_LOGGING_LEVEL,
		SetLoggingLevel: &agent.Call_SetLoggingLevel{
			Level: uint32(level),
			Duration: mesos.DurationInfo{
				Nanoseconds: duration.Nanoseconds(),
			},
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseBytes, err := client.HTTPServicePostJSON(q.PrefixAgent(agentid), requestContent)
	if err != nil {
		return err
	}

	client.PrintJSONBytes(responseBytes)

	return nil
}
