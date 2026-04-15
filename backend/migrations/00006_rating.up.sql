create table if not exists rating (
    course_id int references courses(id),
    lesson_id int references lessons(id),
    lesson_score int,
    user_id int references users(id)
);