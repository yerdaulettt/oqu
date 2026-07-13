create table if not exists lesson_tests (
    id serial primary key,
    lesson_id integer references lessons(id) on delete cascade,
    unique(lesson_id)
);

create table if not exists test_results (
    test_id integer references lesson_tests(id) on delete cascade,
    user_id integer references users(id) on delete cascade,
    completed boolean not null,
    unique(test_id, user_id)
);

create table if not exists questions(
    id serial primary key,
    text varchar(250) not null,
    test_id integer references lesson_tests(id) on delete cascade
);

create table if not exists answers(
    id serial primary key,
    text varchar(250) not null,
    is_correct boolean,
    question_id integer references questions(id) on delete cascade
);

create table if not exists test_submits (
    test_id integer references lesson_tests(id) on delete cascade,
    user_id integer references users(id) on delete cascade,
    question_id integer references questions(id) on delete cascade,
    selected_choice integer references answers(id) on delete cascade,
    unique(user_id, question_id)
);