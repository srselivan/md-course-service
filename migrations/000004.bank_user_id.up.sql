alter table bank drop constraint if exists bank_course_id_fkey;

alter table bank rename column course_id to user_id;
