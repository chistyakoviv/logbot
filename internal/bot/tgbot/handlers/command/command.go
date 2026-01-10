package command

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	I18n "github.com/chistyakoviv/logbot/internal/i18n"
)

type TgCommandInterface interface {
	ApplyDescription(lang string, i18n I18n.I18nChainInterface)
	GetStageHandler(idx int) (handlers.Response, error)
	GetCallbackHandlers() map[string]handlers.Response
}

type TgCommand struct {
	StageHandlers    []handlers.Response
	CallbackHandlers map[string]handlers.Response
}

func (c *TgCommand) ApplyDescription(lang string, i18n I18n.I18nChainInterface) {
}

func (c *TgCommand) GetStageHandler(idx int) (handlers.Response, error) {
	if idx < 0 || idx >= len(c.StageHandlers) {
		return nil, fmt.Errorf("no handler with idx %d", idx)
	}
	return c.StageHandlers[idx], nil
}

func (c *TgCommand) GetCallbackHandlers() map[string]handlers.Response {
	return c.CallbackHandlers
}

type TgCommands map[string]TgCommandInterface
