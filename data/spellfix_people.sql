.load ./extensions/spellfix

drop table if exists spellfix_people;
create virtual table spellfix_people using spellfix1;

insert into spellfix_people(word) select first_name from people;
