package queries

import (
	"encoding/json"
	"fmt"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	agentcalls "github.com/mesos/mesos-go/api/v1/lib/agent/calls"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	mastercalls "github.com/mesos/mesos-go/api/v1/lib/master/calls"
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

	body := agentcalls.SetLoggingLevel(uint32(level), duration)

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

func (q *Loglevel) SetLoglevelAll(level int, duration time.Duration) error {
	list, err := q.getAgentList()
	if err != nil {
		return err
	}

	for _, id := range list {
		err = q.SetLoglevel(id, level, duration)
		if err != nil {
			return err
		}
	}
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
	body := agentcalls.GetLoggingLevel()

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
	body := mastercalls.GetAgents()

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
