package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/davidokk/reminder-manager/internal/commander"
	"gitlab.ozon.dev/davidokk/reminder-manager/internal/storage"
	"gitlab.ozon.dev/davidokk/reminder-manager/utils"
)

const (
	helpCommand           = "help"
	addCommand            = "add"
	listCommand           = "list"
	removeOutdatedCommand = "rmout"
	removeByIdCommand     = "rm"
	editCommand           = "edit"
	todayCommand          = "today"
	forDaysCommand        = "fordays"
)

var description = map[string]string{
	addCommand:            "[dd.mm.yy / today / tomorrow] [text] adds a new reminder",
	listCommand:           "shows all your plans in chronological order",
	removeOutdatedCommand: "removes outdated records",
	removeByIdCommand:     "[id] removes record with given id",
	editCommand:           "[id] [new text] changes the reminder text",
	todayCommand:          "shows today's activities",
	forDaysCommand:        "[count] shows records for next 'count' days",
	helpCommand:           "show this menu",
}

const badArgumentResponse = "Bad argument, try one more time"
const successResponse = "Success! =)"

// AddHandlers registers handlers for given Commander
func AddHandlers(cmd *commander.Commander) {
	cmd.RegisterHandler(listCommand, listFunc)
	cmd.RegisterHandler(addCommand, addFunc)
	cmd.RegisterHandler(removeOutdatedCommand, removeOutdatedFunc)
	cmd.RegisterHandler(removeByIdCommand, removeByIdFunc)
	cmd.RegisterHandler(editCommand, editFunc)
	cmd.RegisterHandler(todayCommand, todayFunc)
	cmd.RegisterHandler(forDaysCommand, forDaysFunc)

	var help string
	for name, desc := range description {
		help += fmt.Sprintf("/%s %s\n", name, desc)
	}

	cmd.RegisterHandler(helpCommand, func(string) string { return help })
}

func forDaysFunc(param string) string {
	cnt, err := strconv.Atoi(param)
	if err != nil || cnt < 1 {
		return badArgumentResponse
	}
	rem := storage.AsStrings(storage.RemindersForDays(cnt))
	if rem == nil {
		return fmt.Sprintf("Nothing to do next %d days =(", cnt)
	}
	return fmt.Sprintf("%d things to do next %d days\n\n%s", len(rem), cnt, strings.Join(rem, "\n"))
}

func todayFunc(string) string {
	rem := storage.AsStrings(storage.RemindersForDays(1))
	if rem == nil {
		return "Nothing to do today =("
	}
	return fmt.Sprintf("%d things to do today\n\n%s", len(rem), strings.Join(rem, "\n"))
}

func editFunc(str string) string {
	params := strings.Split(str, " ")
	id, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil || len(params) < 2 {
		return badArgumentResponse
	}
	if err := storage.Edit(id, strings.Join(params[1:], " ")); err != nil {
		return err.Error()
	}
	return successResponse
}

func removeByIdFunc(params string) string {
	id, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		return badArgumentResponse
	}
	if err := storage.RemoveById(id); err != nil {
		return err.Error()
	}
	return successResponse
}

func removeOutdatedFunc(string) string {
	outdated := storage.RemoveOutdated()
	if outdated == 0 {
		return "There aren't outdated records"
	}
	return fmt.Sprintf("%d records were deleted", outdated)
}

func listFunc(string) string {
	res := storage.AsStrings(storage.Data())
	if len(res) == 0 {
		return "You haven't planned anything yet"
	}
	oldCount := storage.OutdatedCount()
	var outdated string
	var actual string
	if oldCount > 0 {
		outdated += "There are outdated entries on your list\n\n"
		outdated += strings.Join(res[:oldCount], "\n")
	}
	if len(res)-oldCount > 0 {
		actual += "Your actual plans\n\n"
		actual += strings.Join(res[oldCount:], "\n")
	}
	var sep string
	if len(outdated) != 0 && len(actual) != 0 {
		sep = "\n\n"
	}
	return outdated + sep + actual
}

func addFunc(str string) string {
	params := strings.Split(str, " ")
	var date time.Time
	if params[0] == "today" {
		date = utils.UpToDay(time.Now())
	} else if params[0] == "tomorrow" {
		date = utils.UpToDay(time.Now()).Add(time.Hour * 24)
	} else {
		var err error
		date, err = time.Parse("02.01.06", params[0])
		if err != nil || len(params) < 2 {
			return badArgumentResponse
		}
	}
	if err := storage.Add(storage.NewReminder(strings.Join(params[1:], " "), date)); err != nil {
		return err.Error()
	}
	return successResponse
}
