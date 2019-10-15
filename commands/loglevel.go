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
}

func HandleLogLevelSection(app *kingpin.Application, q *queries.Loglevel) {
	HandleLogLevelCommands(app.Command("loglevel", "set/get log level").Alias("loglevels"), q)
}

func HandleLogLevelCommands(loglevel *kingpin.CmdClause, q *queries.Loglevel) {
	cmd := &loglevelHandler{q: q}
	set := loglevel.Command("set", "set log level of agent.").Action(cmd.handleSetLogLevel)
	set.Flag("agent-id", "Agent id of mesos-slave.").Required().StringVar(&cmd.agentid)
	set.Flag("duration", "Duration of modified log level. Can use unit h for hours, m for minutes, s for seconds. e.g: 1h.").Required().DurationVar(&cmd.duration)
	set.Flag("level", "Level of log. 0, 1, 2 or 3").Required().IntVar(&cmd.level)

}

func (cmd *loglevelHandler) handleSetLogLevel(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.SetLoglevel(cmd.agentid, cmd.level, cmd.duration)
}
