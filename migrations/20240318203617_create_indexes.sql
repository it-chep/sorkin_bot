-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists tg_users_id on tg_users (tg_id);
-- create index concurrently if not exists tg_users_id on tg_users (tg_id);
create index concurrently if not exists message_logs_id on message_logs (tg_id);
-- create index concurrently if not exists message_logs_id on message_logs (tg_id);

-- +goose Down
drop index concurrently tg_users_id;
-- drop index concurrently tg_users_id;
