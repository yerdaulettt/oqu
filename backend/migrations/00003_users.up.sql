create table if not exists users (
    id serial primary key,
    name varchar(25),
    username varchar(60) unique,
    role varchar(20),
    password text
);