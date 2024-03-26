-- +goose Up
-- +goose StatementBegin
create table if not exists tg_users
(
    id                bigint,
    tg_id             bigint,
    name              varchar(255),
    surname           varchar(255),
    username          varchar(255),
    registration_time timestamp,
    last_state        varchar(255),
    phone             varchar(11),
    language_code     varchar(5)
);

create table if not exists message_log
(
    id                bigint,
    tg_message_id     bigint,
    system_message_id bigint,
    user_tg_id        bigint,
    time              timestamp
);

create table if not exists message
(
    id         bigint,
    name       text,
    ru_text    text,
    eng_text   text,
    pt_Br_text text

);

create table if not exists message_condition
(
    id bigint
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tg_users;
drop table message_log;
drop table message;
drop table message_condition;
-- +goose StatementEnd
