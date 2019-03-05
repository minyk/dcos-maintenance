package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-maintenance/client"
)

type Machine struct {
	PrefixCb func() string
}

func NewMachine() *Machine {
	return &Machine{}
}

func (q *Machine) MachineDown(file string) error {

	// IDs to down
	mIDs := toMachineIDs(file)
	_, err := client.PrintVerbose("Get MachinIDs list from %s: %d", file, len(mIDs))
	check(err)

	body := master.Call{
		Type: master.Call_START_MAINTENANCE,
		StartMaintenance: &master.Call_StartMaintenance{
			Machines: mIDs,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("", requestContent)
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

	body := master.Call{
		Type: master.Call_STOP_MAINTENANCE,
		StopMaintenance: &master.Call_StopMaintenance{
			Machines: mIDs,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("", requestContent)
	if err != nil {
		return err
	}

	for _, mId := range mIDs {
		client.PrintMessage("Maintenance stopped for: %s", mId.GetHostname())
	}

	return nil
}
