create table course
(
    id            bigserial primary key,
    title         varchar(255)             not null,
    description   text,
    owner_user_id bigint                   not null,
    created_at    timestamp with time zone not null default current_timestamp,
    updated_at    timestamp with time zone not null default current_timestamp
);

create table course_listener
(
    course_id bigint not null references course (id) on delete cascade,
    group_id  bigint not null,
    unique (course_id, group_id)
);

create table course_section
(
    id         serial primary key,
    course_id  int          not null references course (id) on delete cascade,
    parent_id  int references course_section (id) on delete cascade,
    title      varchar(255) not null,
    sort_order int default 0
);

create table course_section_item
(
    id           serial primary key,
    section_id   int          not null references course_section (id) on delete cascade,
    item_type    varchar(50)  not null check (item_type in ('lecture', 'assignment', 'test')),
    title        varchar(255) not null,
    sort_order   int     default 0,
    is_published boolean default false
);

create table lectures
(
    item_id              int primary key references course_section_item (id) on delete cascade,
    content              text,
    video_url            varchar(255),
    reading_time_minutes int default 0
);

create table assignments
(
    item_id       int primary key references course_section_item (id) on delete cascade,
    description   text not null,
    max_score     int default 100,
    deadline_days int
);

create table test
(
    id                  bigserial primary key,
    item_id             bigint references course_section_item (id) on delete cascade,
    title               text                     not null,
    description         text                     not null,
    effective_from      timestamp with time zone not null,
    effective_till      timestamp with time zone not null,
    duration_sec        bigint                   not null,
    max_attempts        smallint default 1       not null,
    questions_total     smallint                 not null,
    generation_settings jsonb,
    bank_ids            bigint[]                 not null
);

create table bank
(
    id          bigserial primary key,
    title       text                     not null,
    description text,
    course_id   bigint references course (id) on delete cascade,
    created_at  timestamp with time zone not null default current_timestamp,
    updated_at  timestamp with time zone not null default current_timestamp
);

create table bank_question
(
    id             bigserial primary key,
    bank_id        bigint                   not null references bank (id) on delete cascade,
    question_text  text                     not null,
    question_type  varchar(50)              not null check (question_type in ('single_choice', 'multiple_choice', 'text')),
    default_points int                               default 1,
    created_at     timestamp with time zone not null default current_timestamp
);

create table bank_answers
(
    id          bigserial primary key,
    question_id int     not null references bank_question (id) on delete cascade,
    answer_text text    not null,
    is_correct  boolean not null default false
);
