alter table bank_question
    alter column question_type type text
    using case
        when question_type = 0 then 'SINGLE'
        when question_type = 1 then 'MULTIPLE'
        when question_type = 2 then 'TEXT'
    end;

alter table bank_question
    add constraint bank_question_question_type_check
        check (question_type in ('SINGLE', 'MULTIPLE', 'TEXT'));

alter table bank_question rename column question_text to text;
alter table bank_question rename column default_points to points;
alter table bank_answers rename column answer_text to text;
