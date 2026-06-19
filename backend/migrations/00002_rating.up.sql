create table if not exists rating (
    course_id integer references courses(id),
    lesson_id integer references lessons(id),
    user_id integer references users(id),
    completed boolean,
    primary key(lesson_id, user_id)
);