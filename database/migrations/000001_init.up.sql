create table if not exists users (
    id serial primary key,
    name text not null
);

insert into users (name) values ('Bill');