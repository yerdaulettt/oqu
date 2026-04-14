create table if not exists enrollments (
    course_id int references courses(id),
    user_id int references users(id)
);