create database test_information;
use test_information;
create table hunter(
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	gun VARCHAR(100) NOT NULL,
	kills int,
	PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;
