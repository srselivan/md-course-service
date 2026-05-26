alter table bank rename column user_id to course_id;

alter table bank
    add constraint bank_course_id_fkey
        foreign key (course_id) references course (id) on delete cascade;
