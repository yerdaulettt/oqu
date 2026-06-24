create table if not exists rating (
    course_id integer references courses(id) on delete cascade,
    lesson_id integer references lessons(id) on delete cascade,
    user_id integer references users(id) on delete cascade,
    completed boolean,
    primary key(lesson_id, user_id)
);