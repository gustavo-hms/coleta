CREATE TABLE zona (
	id INT(11) NOT NULL AUTO_INCREMENT,
	nome VARCHAR(255),
	PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE TABLE esquina (
	id INT(11) NOT NULL AUTO_INCREMENT,
	zona_id INT(11) NOT NULL,
	cruzamento VARCHAR(255),
	localizacao VARCHAR(255),
	PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE TABLE lider (
	id INT(11) NOT NULL AUTO_INCREMENT,
	zona_id INT(11) NOT NULL,
	esquina_id INT(11) NOT NULL,
	cadastrado_em DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
	nome_completo VARCHAR(255),
	telefone_residencial VARCHAR(20),
	telefone_celular VARCHAR(20),
	operadora_celular ENUM('Claro', 'Oi', 'Outra', 'TIM', 'Vivo'),
	email VARCHAR(255),
	PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE TABLE turnos_do_lider (
	lider_id INT(11) NOT NULL,
	turno VARCHAR(20)
) ENGINE=InnoDB;

CREATE TABLE voluntario (
	id INT(11) NOT NULL AUTO_INCREMENT,
	zona_id INT(11) NOT NULL,
	lider_id INT(11) NOT NULL,
	nome_completo VARCHAR(255),
	telefone_residencial VARCHAR(20),
	telefone_celular VARCHAR(20),
	operadora_celular ENUM('Claro', 'Oi', 'Outra', 'TIM', 'Vivo'),
	email VARCHAR(255),
	turno VARCHAR(255),
	como_soube_coleta_2013 VARCHAR(255),
	PRIMARY KEY (id)
) ENGINE=InnoDB;

INSERT INTO zona
	(nome)
VALUES
	('Santo Amaro'),
	('Congonhas/Berrini'),
	('Ibirapuera'),
	('Vila Mariana'),
	('Aclimação/Ipiranga'),
	('Paulista'),
	('Pinheiros'),
	('Butantã'),
	('Lapa/Pompéia'),
	('Pacaembu'),
	('Centro'),
	('Tatuapé'),
	('Norte'),
	('Saúde');
