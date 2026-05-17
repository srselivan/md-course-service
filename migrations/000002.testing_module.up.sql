create table tests
(
    id                  bigserial primary key,
    course_id           bigint                   not null references course (id) on delete cascade,
    course_section_id   bigint                   not null references course_section (id) on delete cascade,
    title               text                     not null,
    description         text                     not null,
    available_from      timestamp with time zone not null,
    available_to        timestamp with time zone not null,
    duration_seconds    integer                  not null check (duration_seconds > 0),
    max_attempts        integer                  not null default 1 check (max_attempts > 0),
    max_score           integer                  not null default 0 check (max_score >= 0),
    questions_count     integer                  not null check (questions_count > 0),
    generation_settings jsonb                    not null default '{}'::jsonb,
    created_at          timestamp with time zone not null default current_timestamp,
    bank_ids            bigint[]                 not null
);
create index tests_course_id_idx on tests (course_id);
create index tests_course_section_id_idx on tests (course_section_id);

create table test_attempts
(
    id             bigserial primary key,
    test_id        bigint                   not null references tests (id) on delete cascade,
    user_id        bigint                   not null,
    attempt_number integer                  not null check (attempt_number > 0),
    status         text                     not null check (status in ('IN_PROGRESS', 'COMPLETED')),
    grade_status   text                     not null default 'NONE' check (grade_status in ('NONE', 'READY', 'PRE_GRADED')),
    started_at     timestamp with time zone not null default current_timestamp,
    expires_at     timestamp with time zone not null,
    completed_at   timestamp with time zone,
    total_score    integer                  not null default 0 check (total_score >= 0),
    unique (test_id, user_id, attempt_number)
);
create index test_attempts_test_user_idx on test_attempts (test_id, user_id);
create index test_attempts_status_idx on test_attempts (status);

create table attempt_questions
(
    id                 bigserial primary key,
    attempt_id         bigint  not null references test_attempts (id) on delete cascade,
    question_id        bigint  not null references bank_question (id) on delete restrict,
    order_index        integer not null,
    question_type      text    not null check (question_type in ('SINGLE', 'MULTIPLE', 'TEXT')),
    question_text      text    not null,
    points             integer not null default 1 check (points >= 0),
    score_awarded      integer not null default 0 check (score_awarded >= 0),
    instructor_comment text,
    unique (attempt_id, question_id),
    unique (attempt_id, order_index)
);
create index attempt_questions_attempt_id_idx on attempt_questions (attempt_id);

create table attempt_question_answers
(
    id                  bigserial primary key,
    attempt_question_id bigint  not null references attempt_questions (id) on delete cascade,
    source_answer_id    bigint references bank_answers (id) on delete restrict,
    answer_text         text    not null,
    is_correct          boolean not null default false
);
create index attempt_question_answers_attempt_question_id_idx on attempt_question_answers (attempt_question_id);

create table attempt_answers
(
    id                  bigserial primary key,
    attempt_question_id bigint not null references attempt_questions (id) on delete cascade,
    selected_answer_id  bigint references attempt_question_answers (id) on delete restrict,
    text_response       text,
    check (selected_answer_id is not null or text_response is not null)
);
create unique index attempt_answers_selected_answer_idx
    on attempt_answers (attempt_question_id, selected_answer_id)
    where selected_answer_id is not null;
create unique index attempt_answers_text_response_idx
    on attempt_answers (attempt_question_id)
    where selected_answer_id is null;
