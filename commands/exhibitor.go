package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type exhibitorHandler struct {
	q     *queries.Exhibitor
	znode string
}

func HandleExhibitorSection(app *kingpin.Application, q *queries.Exhibitor) {
	HandleExhibitorCommands(app.Command("exhibitor", "Control zookeeper on DC/OS Master nodes with Exhibitor.").Alias("exhibitors"), q)
}

func HandleExhibitorCommands(exhibitor *kingpin.CmdClause, q *queries.Exhibitor) {
	cmd := &exhibitorHandler{q: q}
	set := exhibitor.Command("delete", "Delete znode on DC/OS Master Exhibitor").Action(cmd.handleDeleteZnode)
	set.Flag("znode", "ZNode path for deletion.").StringVar(&cmd.znode)
	set.Flag("confirm", "Confirm this deletion.").Required().Bool()
}

func (cmd *exhibitorHandler) handleDeleteZnode(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.DeleteZNode(cmd.znode)
}
