package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-maintenance/client"
)

type Machine struct {
	PrefixCb func() string
}

func NewMachine() *Machine {
	return &Machine{}
}

func (q *Machine) MachineDown(mids []mesos.MachineID) error {
	body := master.Call{
		Type: master.Call_START_MAINTENANCE,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseBytes, err := client.HTTPServicePostData("", requestContent, "application/json")
	if err != nil {
		return err
	}
	client.PrintJSONBytes(responseBytes)
	return nil
}
