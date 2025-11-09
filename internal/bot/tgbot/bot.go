package tgbot

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/bot"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/cmdstage"
	"github.com/chistyakoviv/logbot/internal/bot/tgbot/handlers/command"
	"github.com/chistyakoviv/logbot/internal/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type TgBot struct {
	cfg      *config.Config
	updater  *ext.Updater
	commands command.TgCommands
	cmdstage *cmdstage.TgCmdstage
}

type TgBotSpec struct {
	Cfg      *config.Config
	Commands command.TgCommands
	Cmdstage *cmdstage.TgCmdstage
}

func New(spec *TgBotSpec) bot.Bot {
	return &TgBot{
		cfg:      spec.Cfg,
		commands: spec.Commands,
		cmdstage: spec.Cmdstage,
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
	dispatcher.AddHandler(handlers.NewMessage(noCommands, b.cmdstage.Handler))
	// Add new chat member handler to detect the bot is added to a new chat.
	dispatcher.AddHandler(handlers.NewMessage(message.NewChatMembers, func(b *gotgbot.Bot, ctx *ext.Context) error {
		msg := ctx.EffectiveMessage

		for _, member := range msg.NewChatMembers {
			if member.Id == b.Id {
				// The bot itself was added
				_, _ = b.SendMessage(msg.Chat.Id,
					fmt.Sprintf("👋 Hi everyone! I’ve just joined %s.", msg.Chat.Title),
					nil)
				break
			}
		}
		return nil
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
