-- +goose Up
-- +goose StatementBegin
create table if not exists tg_users
(
    id                bigserial,
    tg_id             bigint,
    name              varchar(255),
    surname           varchar(255),
    username          varchar(255),
    registration_time timestamp,
    last_state        varchar(255),
    phone             varchar(30),
    language_code     varchar(5)
);

create table if not exists message_log
(
    id            bigserial,
    tg_message_id bigint,
    text          text,
    user_tg_id    bigint,
    time          timestamp
);

create table if not exists translations
(
    id bigserial,
    slug text unique,
    ru_text text,
    eng_text text,
    pt_Br_text text
);

create table if not exists message
(
    id         bigserial,
    name       text,
    ru_text    text,
    eng_text   text,
    pt_Br_text text

);

create table if not exists message_condition
(
    id bigserial
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tg_users;
drop table message_log;
drop table message;
drop table message_condition;
-- +goose StatementEnd
