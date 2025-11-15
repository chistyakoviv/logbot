-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    chat_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    UNIQUE (token)
);

COMMENT ON TABLE subscriptions IS 'Subscriptions the chat will receive logs from';

COMMENT ON COLUMN subscriptions.id IS 'Subscription unique identifier';
COMMENT ON COLUMN subscriptions.chat_id IS 'Chat ID the subscription belongs to';
COMMENT ON COLUMN subscriptions.token IS 'Subscription token is used to identify the subscription';
COMMENT ON COLUMN subscriptions.created_at IS 'Creation timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
