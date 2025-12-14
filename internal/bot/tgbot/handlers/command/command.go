package command

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

type TgCommandInterface interface {
	ApplyDescription(lang string, i18n I18n.I18nChainInterface)
	GetHandler() handlers.Response
	GetStages() []handlers.Response
	GetCallbacks() map[string]handlers.Response
}

type TgCommand struct {
	Handler   handlers.Response
	Stages    []handlers.Response
	Callbacks map[string]handlers.Response
}

func (c *TgCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
}

func (c *TgCommand) GetHandler() handlers.Response {
	return c.Handler
}

func (c *TgCommand) GetStages() []handlers.Response {
	return c.Stages
}

func (c *TgCommand) GetCallbacks() map[string]handlers.Response {
	return c.Callbacks
}

type TgCommands map[string]TgCommandInterface
