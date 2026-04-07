create table if not exists comments (
    id serial primary key,
    content varchar(1000),
    lesson_id integer references lessons(id)
);