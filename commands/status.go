package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type statusHandler struct {
	q       *queries.Status
	name    string
	rawJSON bool
}

// HandleScheduleSection
func HandleStatusSection(app *kingpin.Application, q *queries.Status) {
	HandleStatusCommands(app.Command("status", "Display current maintenance status").Alias("statuses"), q)
}

func (cmd *statusHandler) handleStatus(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.GetStatus(cmd.rawJSON)
}

// HandleScheduleCommand
func HandleStatusCommands(status *kingpin.CmdClause, q *queries.Status) {
	cmd := &statusHandler{q: q}

	view := status.Action(cmd.handleStatus)
	view.Flag("json", "Show raw JSON response instead of user-friendly tree").BoolVar(&cmd.rawJSON)
}
