package queries

import (
	"encoding/json"
	"fmt"
	mesos "github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/mesos/mesos-go/api/v1/lib/master"
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

	_, err = client.HTTPServicePostJSON(q.PrefixAgent(agentid), requestContent)
	if err != nil {
		return err
	}

	client.PrintMessage("Request Accepted: %s", agentid)

	return nil
}

func (q *Loglevel) GetLoglevelAll() error {
	list, err := q.getAgentList()
	if err != nil {
		return err
	}

	for _, id := range list {
		err = q.GetLoglevel(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Loglevel) GetLoglevel(agentid string) error {
	body := agent.Call{
		Type: agent.Call_GET_LOGGING_LEVEL,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseContent, err := client.HTTPServicePostJSON(q.PrefixAgent(agentid), requestContent)
	if err != nil {
		return err
	}

	loglevel := agent.Response{}
	err = json.Unmarshal(responseContent, &loglevel)
	if err != nil {
		return err
	}

	client.PrintMessage("Agent %s Current log level: %d", agentid, loglevel.GetLoggingLevel.Level)
	return nil
}

func (q *Loglevel) getAgentList() ([]string, error) {
	body := master.Call{
		Type: master.Call_GET_AGENTS,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	responseContent, err := client.HTTPServicePostJSON(q.PrefixMaster(), requestContent)
	if err != nil {
		return nil, err
	}

	agents := master.Response{}
	err = json.Unmarshal(responseContent, &agents)
	if err != nil {
		return nil, err
	}

	var list []string

	for _, agent := range agents.GetAgents.Agents {
		list = append(list, agent.AgentInfo.ID.Value)
	}

	return list, nil
}
