alter table bank_answers rename column text to answer_text;
alter table bank_question rename column points to default_points;
alter table bank_question rename column text to question_text;

alter table bank_question drop constraint if exists bank_question_question_type_check;

alter table bank_question
    alter column question_type type smallint
    using case
        when question_type = 'SINGLE' then 0
        when question_type = 'MULTIPLE' then 1
        when question_type = 'TEXT' then 2
    end;
