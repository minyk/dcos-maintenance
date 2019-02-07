package queries

import (
	"bufio"
	"encoding/csv"
	"github.com/mesos/mesos-go/api/v1/lib"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func toMachineIDs(filename string) []mesos.MachineID {
	f, err := os.Open(filename)
	check(err)

	rdr := csv.NewReader(bufio.NewReader(f))
	rdr.Comment = '#'
	lists, err := rdr.ReadAll()
	check(err)

	var mids []mesos.MachineID

	for _, line := range lists {
		mid := mesos.MachineID{
			IP:       &line[1],
			Hostname: &line[0],
		}
		mids = append(mids, mid)
	}

	return mids
}

func containMID(s []mesos.MachineID, e mesos.MachineID) bool {

	for _, a := range s {
		if a.Equal(e) {
			return true
		}
	}
	return false
}
