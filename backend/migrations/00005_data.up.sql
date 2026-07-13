insert into courses (name, description) values
('Math', 'Math course description...'),
('Bio', 'Bio course description...');

insert into lessons (name, content, course_id) values
('Intro', 'Some content...', 1),
('Math basics', 'Some content...', 1),
('Conclusion', 'Some content...', 1),

('Basics', 'Some content...', 2),
('Conclusion', 'Some content...', 2);

insert into lesson_tests (lesson_id) values (2);

insert into questions (text, test_id) values
('12 + 2 = ', 1),
('4 - 3 = ', 1),
('5 / 5 = ', 1),
('1 * 7 = ', 1);

insert into answers (text, is_correct, question_id) values
('10', false, 1),
('14', true, 1),
('12', false, 1),
('2', false, 1),

('3', false, 2),
('0', false, 2),
('7', false, 2),
('1', true, 2),

('19', false, 3),
('5', false, 3),
('1', true, 3),
('6', false, 3),

('1', false, 4),
('7', true, 4),
('6', false, 4),
('8', false, 4);