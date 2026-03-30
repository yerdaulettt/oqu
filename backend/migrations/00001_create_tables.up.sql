create table if not exists courses (
    id serial primary key,
    name varchar(200),
    description varchar(500)
);

create table if not exists lessons (
    id serial primary key,
    name varchar(200),
    content text,
    course_id int references courses(id)
);