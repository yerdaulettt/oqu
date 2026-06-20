create table if not exists questions(
    id serial primary key,
    text varchar(250) not null,
    lesson_id integer references lessons(id)
);

create table if not exists answers(
    id serial primary key,
    text varchar(250) not null,
    is_correct boolean,
    question_id integer references questions(id)
);