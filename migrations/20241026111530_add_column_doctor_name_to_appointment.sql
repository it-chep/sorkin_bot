-- +goose Up
-- +goose StatementBegin
alter table appointment add column doctor_name varchar(255), add column type varchar(255);
alter table tg_users add column home_address text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table appointment drop column doctor_name, drop column type;
alter table tg_users drop column home_address;
-- +goose StatementEnd
