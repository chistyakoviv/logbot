-- +goose Up
-- +goose NO TRANSACTION

-- CREATE INDEX CONCURRENTLY cannot run inside a transaction block (SQLSTATE 25001)
CREATE INDEX CONCURRENTLY idx_subscriptions_token_chatid ON subscriptions (token, chat_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
