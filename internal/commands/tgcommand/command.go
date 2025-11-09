package tgcommand

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type TgCommand struct {
	Name    string
	Handler handlers.Response
	Stages  []handlers.Response
}
