-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT NOT NULL,
    lang INTEGER DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id)
);

COMMENT ON TABLE user_settings IS 'User settings';

COMMENT ON COLUMN user_settings.user_id IS 'User ID';
COMMENT ON COLUMN user_settings.lang IS 'Current language of the user';
COMMENT ON COLUMN user_settings.updated_at IS 'Last update timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
