drop index if exists tests_course_section_item_id_idx;

alter table tests
    drop constraint if exists tests_course_section_item_id_fkey;

alter table tests
    rename column course_section_item_id to course_section_id;

alter table tests
    add constraint tests_course_section_id_fkey
        foreign key (course_section_id) references course_section (id) on delete cascade;

create index if not exists tests_course_section_id_idx on tests (course_section_id);
