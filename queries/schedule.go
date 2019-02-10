package queries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/maintenance"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-maintenance/client"
	"strings"
	"time"
)

type Schedule struct {
	PrefixCb func() string
}

func NewSchedule() *Schedule {
	return &Schedule{}
}

func (q *Schedule) GetSchedule(rawJSON bool) error {

	responseBytes := getSchedule()

	if rawJSON {
		client.PrintJSONBytes(responseBytes)
	} else {
		table, err := toScheduleTable(responseBytes)
		if err != nil {
			return err
		}
		client.PrintMessage(table)
	}
	return nil
}

func (q *Schedule) AddSchedule(start time.Time, duration time.Duration, file string) error {
	// TODO get and add ?

	ret := getSchedule()
	resp := master.Response{}
	err := json.Unmarshal(ret, &resp)

	currentSchedule := resp.GetMaintenanceSchedule.Schedule

	mIDs := toMachineIDs(file)

	newWindow := maintenance.Window{
		MachineIDs: mIDs,
		Unavailability: mesos.Unavailability{
			Start: mesos.TimeInfo{
				Nanoseconds: start.UnixNano(),
			},
			Duration: &mesos.DurationInfo{
				Nanoseconds: duration.Nanoseconds(),
			},
		},
	}

	currentSchedule.Windows = append(currentSchedule.Windows, newWindow)

	body := master.Call{
		Type: master.Call_UPDATE_MAINTENANCE_SCHEDULE,
		UpdateMaintenanceSchedule: &master.Call_UpdateMaintenanceSchedule{
			Schedule: currentSchedule,
		},
	}

	err = updateSchedule(body)
	if err != nil {
		return err
	}
	client.PrintMessage("Maintenance schedule updated successfully.")

	return nil
}

func (q *Schedule) RemoveSchedule(file string) error {

	// Get current schedule
	ret := getSchedule()
	resp := master.Response{}
	err := json.Unmarshal(ret, &resp)
	currentSchedule := resp.GetMaintenanceSchedule.Schedule

	// IDs to remove
	mIDs := toMachineIDs(file)
	client.PrintVerbose("Get MachinID list from %s: %d", file, len(mIDs))

	newIDs := []mesos.MachineID{}
	newWindows := []maintenance.Window{}

	// remove ID and set the left IDs back
	for _, window := range currentSchedule.Windows {
		for _, currentID := range window.MachineIDs {
			client.PrintVerbose("Contains %s: %s", currentID.GetIP(), containMID(mIDs, currentID))
			if !containMID(mIDs, currentID) {
				newIDs = append(newIDs, currentID)
			}
		}

		if len(newIDs) != 0 {
			window.MachineIDs = newIDs
			newIDs = nil
			newWindows = append(newWindows, window)
		}
	}

	if len(newWindows) != 0 {
		currentSchedule.Windows = newWindows
	} else {
		currentSchedule.Windows = nil
	}

	body := master.Call{
		Type: master.Call_UPDATE_MAINTENANCE_SCHEDULE,
		UpdateMaintenanceSchedule: &master.Call_UpdateMaintenanceSchedule{
			Schedule: currentSchedule,
		},
	}

	err = updateSchedule(body)
	if err != nil {
		return err
	}

	client.PrintMessage("Maintenance schedule updated successfully.")
	return nil
}

func toScheduleTable(planJSONBytes []byte) (string, error) {
	schedule := master.Response{}
	err := json.Unmarshal(planJSONBytes, &schedule)
	check(err)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintln("Window\tHostname\t\tIP\t\tStart\t\t\t\tDuration"))

	for window, element := range schedule.GetMaintenanceSchedule.Schedule.Windows {

		start := time.Unix(0, element.Unavailability.Start.GetNanoseconds())
		duration := time.Duration(element.Unavailability.GetDuration().GetNanoseconds())

		for _, host := range element.MachineIDs {
			hostname := *host.Hostname
			ip := *host.IP
			buf.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", window, hostname, ip, start, duration))
		}

	}

	return strings.TrimRight(buf.String(), "\n"), nil
}

func getSchedule() []byte {
	body := master.Call{
		Type: master.Call_GET_MAINTENANCE_SCHEDULE,
	}

	requestContent, err := json.Marshal(body)
	check(err)

	responseBytes, err := client.HTTPServicePostJSON("", requestContent)
	check(err)

	return responseBytes
}

func updateSchedule(body master.Call) error {

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client.PrintVerbose("Request JSON: %s", requestContent)
	_, err = client.HTTPServicePostJSON("", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Schedules updated")
	}

	return nil
}
