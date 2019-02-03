package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type machineHandler struct {
	q    *queries.Machine
	name string
}

func (cmd *machineHandler) handleMachineUp(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.MachineUp(nil)
}

// HandleScheduleSection
func HandleMachineSection(app *kingpin.Application, q *queries.Machine) {
	HandleMachineCommands(app.Command("machine", "Modify machine status.").Alias("machines"), q)
}

// HandleScheduleCommand
func HandleMachineCommands(machine *kingpin.CmdClause, q *queries.Machine) {
	cmd := &machineHandler{q: q}
	up := machine.Command("up", "Stop maintenance for machines").Action(cmd.handleMachineUp)
	up.Flag("machines", "IP or Hostname for machine. Comma separated value can be used.").StringVar(&cmd.name)

	down := machine.Command("down", "Start maintenance for machines").Action(cmd.handleMachineUp)
	down.Flag("machines", "IP or Hostname for machine. Comma separated value can be used.").StringVar(&cmd.name)
}
