package queries

import (
	"encoding/json"
	//mesos "github.com/mesos/mesos-go/api/v1/lib"
	//"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/mesos/mesos-go/api/v1/lib/master/calls"
	"github.com/minyk/dcos-maintenance/client"
)

type Machine struct {
	PrefixCb func() string
}

func NewMachine() *Machine {
	return &Machine{
		PrefixCb: func() string { return "/mesos/api/v1/" },
	}
}

func (q *Machine) MachineDown(file string) error {

	// IDs to down
	mIDs := toMachineIDs(file)
	_, err := client.PrintVerbose("Get MachinIDs list from %s: %d", file, len(mIDs))
	check(err)

	body := calls.StartMaintenance(mIDs...)

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON(q.PrefixCb(), requestContent)
	if err != nil {
		return err
	}

	for _, mId := range mIDs {
		client.PrintMessage("Maintenance started for: %s", mId.GetHostname())
	}

	return nil
}

func (q *Machine) MachineUp(file string) error {
	mIDs := toMachineIDs(file)
	_, err := client.PrintVerbose("Get MachinIDs list from %s: %d", file, len(mIDs))
	check(err)

	body := calls.StopMaintenance(mIDs...)

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON(q.PrefixCb(), requestContent)
	if err != nil {
		return err
	}

	for _, mId := range mIDs {
		client.PrintMessage("Maintenance stopped for: %s", mId.GetHostname())
	}

	return nil
}
