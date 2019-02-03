package commands

import (
	"github.com/minyk/dcos-maintenance/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
	"time"
)

type scheduleHandler struct {
	q        *queries.Schedule
	name     string
	startt   string
	start    time.Time
	duration time.Duration
	filename string
	rawJSON  bool
}

func (cmd *scheduleHandler) handleList(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.GetSchedule(cmd.rawJSON)
}

func (cmd *scheduleHandler) handleUpdate(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.GetSchedule(cmd.rawJSON)
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

	update := schedule.Command("update", "Update maintenance schedule").Action(cmd.handleUpdate)
	update.Arg("start", "Start time of this maintenance schedule").Required().StringVar(&cmd.startt)
	update.Arg("duration", "Duration of maintenance schedule").Required().DurationVar(&cmd.duration)
	update.Arg("file", "Name of a specific file to update").Required().StringVar(&cmd.filename)
	update.Flag("json", "Show raw JSON response instead of user-friendly tree").BoolVar(&cmd.rawJSON)
}
