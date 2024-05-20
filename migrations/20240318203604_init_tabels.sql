-- +goose Up
-- +goose StatementBegin
create table if not exists tg_users
(
    id                bigserial,
    tg_id             bigint,
    patient_id        bigint,
    name              varchar(255),
    surname           varchar(255),
    username          varchar(255),
    registration_time varchar(16),
    birt_date         varchar(10),
    third_name        varchar(255),
    last_state        varchar(255),
    phone             varchar(30),
    language_code     varchar(5),
    admin             boolean default false
);

create table if not exists message_log
(
    id            bigserial,
    tg_message_id bigint,
    text          text,
    user_tg_id    bigint,
    time          timestamp with time zone
);

create table if not exists translations
(
    id                  bigserial,
    id_in_source_system bigint unique,
    uses                boolean default false,
    slug                text unique,
    ru_text             text,
    eng_text            text,
    pt_Br_text          text
);

create table if not exists appointment
(
    id             bigserial,
    tg_id          bigint,
    speciality_id  bigint,
    appointment_id bigint,
    doctor_id      bigint,
    date           text,
    time_start     text,
    time_end       text,
    draft          boolean
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

set time zone 'UTC';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tg_users;
drop table message_log;
drop table message;
drop table message_condition;
-- +goose StatementEnd
