package queries

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	body := master.Call{
		Type: master.Call_GET_MAINTENANCE_SCHEDULE,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseBytes, err := client.HTTPServicePostData("", requestContent, "application/json")
	if err != nil {
		return err
	}
	if rawJSON {
		client.PrintJSONBytes(responseBytes)
	} else {
		tree, err := toScheduleTable(responseBytes)
		if err != nil {
			return err
		}
		client.PrintMessage(tree)
	}
	return nil
}

func (q *Schedule) GetStatus(rawJSON bool) error {
	body := master.Call{
		Type: master.Call_GET_MAINTENANCE_STATUS,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseBytes, err := client.HTTPServicePostJSON("", requestContent)
	if err != nil {
		return err
	}

	if rawJSON {
		client.PrintJSONBytes(responseBytes)
	} else {
		// TODO print human friendly result
		list, err := toStatusTable(responseBytes)
		if err != nil {
			return err
		}
		client.PrintMessage(list)
	}

	return nil
}

func (q *Schedule) UpdateSchedule(start time.Time, duration time.Duration, file string) error {
	body := master.Call{
		Type: master.Call_UPDATE_MAINTENANCE_SCHEDULE,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Schedules updated")
	}

	return nil
}

func toScheduleTable(planJSONBytes []byte) (string, error) {
	schedule := master.Response{}
	err := json.Unmarshal(planJSONBytes, &schedule)
	if err != nil {
		return "", err
	}

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

func toStatusTable(planJSONBytes []byte) (string, error) {
	status := master.Response{}
	err := json.Unmarshal(planJSONBytes, &status)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintln("Status\t\tHostname\t\tIP\t\tOffer"))

	for _, element := range status.GetMaintenanceStatus.Status.DrainingMachines {
		hostname := *element.ID.Hostname
		ip := *element.ID.IP
		if len(element.Statuses) != 0 {
			offer := element.Statuses[0].Status
			buf.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", "Draining", hostname, ip, offer))
		} else {
			buf.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", "Draining", hostname, ip, "none"))
		}
	}

	for _, element := range status.GetMaintenanceStatus.Status.DownMachines {
		hostname := *element.Hostname
		ip := *element.IP
		buf.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", "Downed", hostname, ip, "none"))
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}
