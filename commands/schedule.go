package commands

import (
	"github.com/araddon/dateparse"
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
	"time"
)

type scheduleHandler struct {
	q        *queries.Schedule
	name     string
	start    string
	duration time.Duration
	filename string
	rawJSON  bool
}

func (cmd *scheduleHandler) handleList(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.GetSchedule(cmd.rawJSON)
}

func (cmd *scheduleHandler) handleUpdate(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	t, err := dateparse.ParseAny(cmd.start)
	if err != nil {
		panic(err)
	}
	return cmd.q.AddSchedule(t, cmd.duration, cmd.filename)
}

func (cmd *scheduleHandler) handleRemove(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.RemoveSchedule(cmd.filename)
}

// HandleScheduleSection
func HandleScheduleSection(app *kingpin.Application, q *queries.Schedule) {
	HandleScheduleCommands(app.Command("schedule", "View current maintenance schedule").Alias("schedules"), q)
}

// HandleScheduleCommand
func HandleScheduleCommands(schedule *kingpin.CmdClause, q *queries.Schedule) {
	cmd := &scheduleHandler{q: q}
	view := schedule.Command("view", "Display current maintenance schedule").Action(cmd.handleList)
	view.Flag("json", "Show raw JSON response instead of user-friendly tree").BoolVar(&cmd.rawJSON)

	add := schedule.Command("add", "Add maintenance schedule").Action(cmd.handleUpdate)
	add.Flag("start-at", "Start time of this maintenance schedule.").Required().StringVar(&cmd.start)
	add.Flag("duration", "Duration of maintenance schedule. Can use unit h for hours, m for minutes, s for seconds. e.g: 1h.").Required().DurationVar(&cmd.duration)
	add.Flag("file", "Name of a specific file to update").Required().StringVar(&cmd.filename)

	remove := schedule.Command("remove", "Remove maintenance schedule").Action(cmd.handleRemove)
	remove.Flag("file", "Name of a specific file to update").Required().StringVar(&cmd.filename)
}
