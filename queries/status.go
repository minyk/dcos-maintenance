package queries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-maintenance/client"
	"strings"
)

type Status struct {
	PrefixCb func() string
}

func NewStatus() *Status {
	return &Status{
		PrefixCb: func() string { return "/mesos/api/v1/" },
	}
}

func (q *Status) GetStatus(rawJSON bool) error {
	body := master.Call{
		Type: master.Call_GET_MAINTENANCE_STATUS,
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	responseBytes, err := client.HTTPServicePostJSON(q.PrefixCb(), requestContent)
	if err != nil {
		return err
	}

	if rawJSON {
		client.PrintJSONBytes(responseBytes)
	} else {
		list, err := toStatusTable(responseBytes)
		if err != nil {
			return err
		}
		client.PrintMessage(list)
	}

	return nil
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
