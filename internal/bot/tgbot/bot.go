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
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type TgBot struct {
	cfg         *config.Config
	updater     *ext.Updater
	commands    command.TgCommands
	cmdstage    handlers.Response
	join        handlers.Response
	webhookPath string
}

type TgBotSpec struct {
	Cfg      *config.Config
	Commands command.TgCommands
	Cmdstage handlers.Response
	Join     handlers.Response
}

func New(spec *TgBotSpec) bot.Bot {
	return &TgBot{
		cfg:         spec.Cfg,
		commands:    spec.Commands,
		cmdstage:    spec.Cmdstage,
		join:        spec.Join,
		webhookPath: "/tg/" + spec.Cfg.Token,
	}
}

func (b *TgBot) Start(ctx context.Context, logger *slog.Logger) error {
	// Create bot from environment value.
	bot, err := gotgbot.NewBot(b.cfg.Token, nil)
	if err != nil {
		log.Fatalf("failed to create new bot: %s", err.Error())
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Error(
				"an error occurred while handling update",
				slog.Attr{Key: "error", Value: slog.StringValue(err.Error())},
			)
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
		// Wait for the release to use slog.Logger
		// Logger:      logger,
	})
	b.updater = ext.NewUpdater(dispatcher, &ext.UpdaterOpts{
		// Wait for the release to use slog.Logger
		// Logger: logger,
	})

	// The order of the handlers is important.
	// If a handler matches the request, the next handler will not be called.

	// Add all command handlers.
	for name, command := range b.commands {
		dispatcher.AddHandler(handlers.NewCommand(name, command.Handler))
	}
	// Add new chat member handler to detect the bot is added to a new chat.
	dispatcher.AddHandler(handlers.NewMessage(message.NewChatMembers, b.join))
	// Add no command handler to handle command stages.
	dispatcher.AddHandler(handlers.NewMessage(noCommand, b.cmdstage))

	// Start the webhook server. We start the server before we set the webhook itself, so that when telegram starts
	// sending updates, the server is already ready.
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  b.cfg.HTTPServer.Address,
		SecretToken: b.cfg.Webhook.Secret, // Setting a webhook secret here allows you to ensure the webhook is set by you (must be set here AND in SetWebhook!).
	}

	// The bot's urlPath can be anything. Here, we use "custom-path/<TOKEN>" as an example.
	// It can be a good idea for the urlPath to contain the bot token, as that makes it very difficult for outside
	// parties to find the update endpoint (which would allow them to inject their own updates).
	// err = b.updater.StartWebhook(bot, b.webhookPath, webhookOpts)
	// if err != nil {
	// 	log.Fatalf("failed to start webhook: %s", err.Error())
	// }
	// The AddWebhook method is called inside the StartWebhook method,
	// so we have to call it manually when using an external server.
	b.updater.AddWebhook(bot, b.webhookPath, &ext.AddWebhookOpts{SecretToken: webhookOpts.SecretToken})

	err = b.updater.SetAllBotWebhooks(b.cfg.Webhook.Domain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
	})
	if err != nil {
		log.Fatalf("failed to set webhook: %s", err.Error())
	}

	// logger.Info("Bot has been started...", "bot_username", bot.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	// b.updater.Idle()

	return nil
}

func (b *TgBot) Shutdown(ctx context.Context) error {
	return b.updater.Stop()
}

func (b *TgBot) HandlerFunc() http.HandlerFunc {
	return b.updater.GetHandlerFunc("/")
}

func (b *TgBot) WebhookPath() string {
	return b.webhookPath
}
