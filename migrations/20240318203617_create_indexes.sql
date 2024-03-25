-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists tg_users on tg_id (tg_id);
-- create index concurrently if not exists tg_users on tg_id (tg_id);
create index concurrently if not exists message_logs on referals (tg_id);
-- create index concurrently if not exists message_logs on referal (tg_id);
create index concurrently if not exists admins_login on admins (login);
-- create index concurrently if not exists admins_login on admin (login);

-- +goose Down
drop index concurrently tg_users;
-- drop index concurrently users_tg_id;
drop index concurrently referals_tg_id;
-- drop index concurrently referals_tg_id;
drop index concurrently admins_login;
-- drop index concurrently admins_login;
