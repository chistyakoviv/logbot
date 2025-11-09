package tgbot

import (
	"context"
	"log"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/commands/tgcommand"
	"github.com/chistyakoviv/logbot/internal/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type TgBot struct {
	cfg      *config.Config
	updater  *ext.Updater
	commands tgcommand.TgCommands
}

func New(cfg *config.Config, commands tgcommand.TgCommands) bot.Bot {
	return &TgBot{
		cfg:      cfg,
		commands: commands,
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
			log.Println("an error occurred while handling update:", err.Error())
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

	// Add echo handler to reply to all text messages.
	for name, command := range b.commands {
		dispatcher.AddHandler(handlers.NewCommand(name, command.Handler))
	}
	// Add message handler to handle command stages.
	dispatcher.AddHandler(handlers.NewMessage(noCommands, func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := b.SendMessage(ctx.EffectiveMessage.Chat.Id, "No command received", nil)
		return err
	}))

	// Start the webhook server. We start the server before we set the webhook itself, so that when telegram starts
	// sending updates, the server is already ready.
	webhookOpts := ext.WebhookOpts{
		ListenAddr:  b.cfg.ListenAddr,
		SecretToken: b.cfg.Webhook.Secret, // Setting a webhook secret here allows you to ensure the webhook is set by you (must be set here AND in SetWebhook!).
	}

	// The bot's urlPath can be anything. Here, we use "custom-path/<TOKEN>" as an example.
	// It can be a good idea for the urlPath to contain the bot token, as that makes it very difficult for outside
	// parties to find the update endpoint (which would allow them to inject their own updates).
	err = b.updater.StartWebhook(bot, "tg/"+b.cfg.Token, webhookOpts)
	if err != nil {
		log.Fatalf("failed to start webhook: %s", err.Error())
	}

	err = b.updater.SetAllBotWebhooks(b.cfg.Webhook.Domain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookOpts.SecretToken,
	})
	if err != nil {
		log.Fatalf("failed to set webhook: %s", err.Error())
	}

	logger.Info("Bot has been started...", "bot_username", bot.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	b.updater.Idle()

	return nil
}

func (b *TgBot) Shutdown(ctx context.Context) error {
	return b.updater.Stop()
}
