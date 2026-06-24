insert into courses (name, description) values
('Test course 1', 'Some description...'),
('Test course 2', 'Some description...'),
('Test course 3', 'Some description...');

insert into lessons (name, content, course_id) values
('Test lesson 1', 'Some content...', 1),
('Test lesson 2', 'Some content...', 1),
('Test lesson 3', 'Some content...', 1),
('Test lesson 4', 'Some content...', 2),
('Test lesson 5', 'Some content...', 2),
('Test lesson 6', 'Some content...', 3),
('Test lesson 7', 'Some content...', 3),
('Test lesson 8', 'Some content...', 3);

insert into questions (text, lesson_id) values
('Test question 1', 2),
('Test question 2', 2),
('Test question 3', 2),
('Test question 4', 2);

insert into answers (text, is_correct, question_id) values
('Test answer 1', false, 1),
('Test answer 2', true, 1),
('Test answer 3', false, 1),
('Test answer 4', false, 1),
('Test answer 1', false, 2),
('Test answer 2', false, 2),
('Test answer 3', false, 2),
('Test answer 4', true, 2),
('Test answer 1', false, 3),
('Test answer 2', false, 3),
('Test answer 3', true, 3),
('Test answer 4', false, 3),
('Test answer 1', true, 4),
('Test answer 2', true, 4),
('Test answer 3', false, 4),
('Test answer 4', false, 4);