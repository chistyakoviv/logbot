package command

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type TgCommand struct {
	Handler   handlers.Response
	Stages    []handlers.Response
	Callbacks map[string]handlers.Response
}

type TgCommands map[string]*TgCommand
