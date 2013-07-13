update zona set nome="Moema/Berrini" where id = 2;
update zona set nome="Santana/Norte" where id = 13;
update zona set nome="Pacaembu/Barra Funda" where id = 10;

alter table esquina add column prioridade enum('baixa', 'alta');

alter table voluntario drop column turno;
alter table voluntario add column esquina_id INT(11) NOT NULL;
alter table voluntario add column rg varchar(20);
alter table voluntario add column cpf varchar(20);
alter table voluntario add column idade varchar(20);
alter table voluntario add column cadastrado_em DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00';

alter table zona add column bloqueada tinyint(1) not null default 0;

create table turnos_do_voluntario (
	voluntario_id int(11) not null,
	turno varchar(20)
) engine=InnoDB;

