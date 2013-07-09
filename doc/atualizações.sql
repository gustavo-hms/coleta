update zona set nome="Moema/Berrini" where id = 2;
update zona set nome="Santana/Norte" where id = 13;
update zona set nome="Pacaembu/Barra Funda" where id = 10;

alter table esquina add column prioridade enum('baixa', 'alta');
