package command

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

type TgCommandInterface interface {
	ApplyDescription(lang string, i18n I18n.I18nChainInterface)
	GetStartHandler() handlers.Response
	GetStageHandlers() []handlers.Response
	GetCallbackHandlers() map[string]handlers.Response
}

type TgCommand struct {
	StartHandler     handlers.Response
	StageHandlers    []handlers.Response
	CallbackHandlers map[string]handlers.Response
}

func (c *TgCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
}

func (c *TgCommand) GetStartHandler() handlers.Response {
	return c.StartHandler
}

func (c *TgCommand) GetStageHandlers() []handlers.Response {
	return c.StageHandlers
}

func (c *TgCommand) GetCallbackHandlers() map[string]handlers.Response {
	return c.CallbackHandlers
}

type TgCommands map[string]TgCommandInterface
