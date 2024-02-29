-- +goose Up
create table chat (
    id serial primary key,
    usernames text[]
);

create table message (
    id serial primary key,
    "from" text not null,
    text text not null,
    timestamp timestamp
);

-- +goose Down
drop table chat;