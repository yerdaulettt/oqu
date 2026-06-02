create table if not exists courses (
    id serial primary key,
    name varchar(200) NOT NULL,
    description varchar(500) NOT NULL
);

create table if not exists lessons (
    id serial primary key,
    name varchar(200) NOT NULL,
    content text NOT NULL,
    course_id int references courses(id)
);

create table if not exists users (
    id serial primary key,
    name varchar(25) NOT NULL DEFAULT 'User',
    username varchar(60) unique,
    role varchar(20) NOT NULL,
    password text NOT NULL
);

create table if not exists comments (
    id serial primary key,
    content varchar(1000) NOT NULL,
    lesson_id integer references lessons(id),
    user_id integer references users(id)
);

create table if not exists enrollments (
    course_id int references courses(id),
    user_id int references users(id),
    PRIMARY KEY(course_id, user_id)
);

create table if not exists comment_votes (
    id serial primary key,
    comment_id integer references comments(id),
    user_id integer references users(id),
    voted boolean,
    UNIQUE(comment_id, user_id)
);