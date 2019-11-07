package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
	"time"
)

type loglevelHandler struct {
	q        *queries.Loglevel
	agentid  string
	level    int
	duration time.Duration
	all      bool
}

func HandleLogLevelSection(app *kingpin.Application, q *queries.Loglevel) {
	HandleLogLevelCommands(app.Command("loglevel", "set/get log level").Alias("loglevels"), q)
}

func HandleLogLevelCommands(loglevel *kingpin.CmdClause, q *queries.Loglevel) {
	cmd := &loglevelHandler{q: q, agentid: ""}
	set := loglevel.Command("set", "set log level of agent.").Action(cmd.handleSetLogLevel)
	set.Flag("agent-id", "Agent id of the target mesos slave node. If missing, all slave nodes are targeted.").StringVar(&cmd.agentid)
	set.Flag("duration", "Duration of modified log level. Can use unit h for hours, m for minutes, s for seconds. e.g: 1h.").Required().DurationVar(&cmd.duration)
	set.Flag("level", "Level of log. 0 to 6.").Required().IntVar(&cmd.level)

	get := loglevel.Command("get", "Get logging level of an agent.").Action(cmd.handleGetLogLevel)
	get.Flag("agent-id", "Agent ID of the target mesos slave node. If missing, all slave nodes are targeted.").StringVar(&cmd.agentid)
}

func (cmd *loglevelHandler) handleSetLogLevel(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	if cmd.agentid != "" {
		return cmd.q.SetLoglevel(cmd.agentid, cmd.level, cmd.duration)
	}
	return cmd.q.SetLoglevelAll(cmd.level, cmd.duration)
}

func (cmd *loglevelHandler) handleGetLogLevel(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	if cmd.agentid != "" {
		return cmd.q.GetLoglevel(cmd.agentid)
	}
	return cmd.q.GetLoglevelAll()
}
