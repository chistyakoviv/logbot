-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_settings (
    chat_id BIGINT NOT NULL,
    collapse_period INTERVAL NOT NULL DEFAULT INTERVAL '0',
    mute_until TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    silence_until TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (chat_id)
);

COMMENT ON TABLE chat_settings IS 'Chat settings';

COMMENT ON COLUMN chat_settings.chat_id IS 'Chat ID';
COMMENT ON COLUMN chat_settings.collapse_period IS 'Number of seconds to suppress repeated error notifications';
COMMENT ON COLUMN chat_settings.mute_until IS 'Time until which notifications are not sent';
COMMENT ON COLUMN chat_settings.silence_until IS 'Time until which notifications are sent silently';
COMMENT ON COLUMN chat_settings.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
