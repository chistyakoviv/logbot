package tgbot

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type TgBot struct {
	cfg         *config.Config
	updater     *ext.Updater
	commands    command.TgCommands
	cmdstage    handlers.Response
	join        handlers.Response
	webhookPath string
	logger      *slog.Logger
	bot         *gotgbot.Bot
}

type TgBotSpec struct {
	Cfg      *config.Config
	Commands command.TgCommands
	Cmdstage handlers.Response
	Join     handlers.Response
	Logger   *slog.Logger
}

func New(spec *TgBotSpec) bot.Bot {
	b := &TgBot{
		cfg:         spec.Cfg,
		commands:    spec.Commands,
		cmdstage:    spec.Cmdstage,
		join:        spec.Join,
		webhookPath: "/tg/",
		logger:      spec.Logger,
	}

	// Split initialization and registeriing webhook
	// so that it can be possible to start the server first.
	b.init()

	return b
}

func (tgb *TgBot) init() {
	var err error

	tgb.bot, err = gotgbot.NewBot(tgb.cfg.Token, nil)
	if err != nil {
		log.Fatalf("failed to create new bot: %s", err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			tgb.logger.Error(
				"an error occurred while handling update",
				slog.Attr{Key: "error", Value: slog.StringValue(err.Error())},
			)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		// Wait for the release to use slog.Logger
		// Logger:      logger,
	})
	tgb.updater = ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		// Wait for the release to use slog.Logger
		// Logger: logger,
	})

	// The order of the handlers is important.
	// If a handler matches the request, the next handler will not be called.

	// Add all command handlers.
	for name, command := range tgb.commands {
		dispatcher.AddHandler(handlers.NewCommand(name, command.GetHandler()))
		// Register all the callbacks the command provides
		for cbName, cb := range command.GetCallbacks() {
			dispatcher.AddHandler(handlers.NewCallback(callbackquery.Prefix(cbName), cb))
		}
	}
	// Add new chat member handler to detect the bot is added to a new chat.
	dispatcher.AddHandler(handlers.NewMessage(message.NewChatMembers, tgb.join))
	// Add no command handler to handle command stages.
	dispatcher.AddHandler(handlers.NewMessage(noCommand, tgb.cmdstage))

	// The AddWebhook method is called inside the StartWebhook method,
	// so we have to call it manually when using an external server.
	// Url is what the route must look like after trimming the prefix specified in the GetHandlerFunc method.
	// So register each webhook by token to be able to use many bots at the same time.
	err = tgb.updater.AddWebhook(tgb.bot, tgb.cfg.Token, &ext.AddWebhookOpts{SecretToken: tgb.cfg.Webhook.Secret})
	if err != nil {
		log.Fatalf("failed to add webhook: %s", err.Error())
	}
}

func (tgb *TgBot) Start(ctx context.Context) {
	// At this point the server must be started.
	err := tgb.updater.SetAllBotWebhooks(tgb.cfg.Webhook.Domain+tgb.webhookPath, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        tgb.cfg.Webhook.Secret,
	})
	if err != nil {
		log.Fatalf("failed to set webhook: %s", err.Error())
	}

	tgb.logger.Info("Telegram bot is ready to handle requests", "bot_username", tgb.bot.Username)
}

func (tgb *TgBot) Shutdown(ctx context.Context) error {
	return tgb.updater.Stop()
}

func (tgb *TgBot) HandlerFunc() http.HandlerFunc {
	return tgb.updater.GetHandlerFunc(tgb.webhookPath)
}

func (tgb *TgBot) WebhookPath() string {
	return tgb.webhookPath
}

func (tgb *TgBot) SendMessage(chatId int64, text string, opts interface{}) error {
	options, ok := opts.(*gotgbot.SendMessageOpts)
	if !ok {
		options = &gotgbot.SendMessageOpts{}
	}
	_, err := tgb.bot.SendMessage(chatId, text, options)
	return err
}
