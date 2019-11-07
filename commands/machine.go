package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type machineHandler struct {
	q    *queries.Machine
	name string
	file string
}

func (cmd *machineHandler) handleMachineUp(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.MachineUp(cmd.file)
}

func (cmd *machineHandler) handleMachineDown(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.MachineDown(cmd.file)
}

// HandleMachineSection blah blah
func HandleMachineSection(app *kingpin.Application, q *queries.Machine) {
	HandleMachineCommands(app.Command("machine", "Modify machine status.").Alias("machines"), q)
}

// HandleMachineCommands blah blah
func HandleMachineCommands(machine *kingpin.CmdClause, q *queries.Machine) {
	cmd := &machineHandler{q: q}
	up := machine.Command("up", "Stop maintenance for machines").Action(cmd.handleMachineUp)
	up.Flag("list", "Name of a specific file which contains hostname, ip lists for machine up.").StringVar(&cmd.file)

	down := machine.Command("down", "Start maintenance for machines").Action(cmd.handleMachineDown)
	down.Flag("list", "Name of a specific file which contains hostname, ip lists for machine down.").StringVar(&cmd.file)
}
