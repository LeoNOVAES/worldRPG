-- --------------------------------------------------------
-- Servidor:                     127.0.0.1
-- Versão do servidor:           10.1.37-MariaDB - mariadb.org binary distribution
-- OS do Servidor:               Win32
-- HeidiSQL Versão:              9.5.0.5196
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Copiando estrutura do banco de dados para world
CREATE DATABASE IF NOT EXISTS `world` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `world`;

-- Copiando estrutura para tabela world.arms
CREATE TABLE IF NOT EXISTS `arms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `army` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.arms: ~5 rows (aproximadamente)
/*!40000 ALTER TABLE `arms` DISABLE KEYS */;
INSERT INTO `arms` (`id`, `army`) VALUES
	(1, 'Machado Grande'),
	(2, 'Machado Pequeno'),
	(3, 'Escudo'),
	(4, 'Adaga'),
	(5, 'Espada Grande'),
	(6, 'Espada Curta');
/*!40000 ALTER TABLE `arms` ENABLE KEYS */;

-- Copiando estrutura para view world.armsview
-- Criando tabela temporária para evitar erros de dependência de VIEW
CREATE TABLE `armsview` (
	`army` VARCHAR(100) NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Copiando estrutura para tabela world.class
CREATE TABLE IF NOT EXISTS `class` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `class` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.class: ~4 rows (aproximadamente)
/*!40000 ALTER TABLE `class` DISABLE KEYS */;
INSERT INTO `class` (`id`, `class`) VALUES
	(1, 'Guerreiro'),
	(2, 'Ladino'),
	(3, 'Paladino'),
	(4, 'Golias');
/*!40000 ALTER TABLE `class` ENABLE KEYS */;

-- Copiando estrutura para tabela world.imgavatar
CREATE TABLE IF NOT EXISTS `imgavatar` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `path` varchar(100) DEFAULT NULL,
  `id_person` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `id_person` (`id_person`),
  CONSTRAINT `imgavatar_ibfk_1` FOREIGN KEY (`id_person`) REFERENCES `persons` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.imgavatar: ~2 rows (aproximadamente)
/*!40000 ALTER TABLE `imgavatar` DISABLE KEYS */;
INSERT INTO `imgavatar` (`id`, `path`, `id_person`) VALUES
	(10, '/images/índice.jpg20190201112742.png', 2),
	(11, '/images/100022.png20190201112905.png', 1),
	(12, '/images/download.jpg20190204092146.png', 4);
/*!40000 ALTER TABLE `imgavatar` ENABLE KEYS */;

-- Copiando estrutura para tabela world.ligation
CREATE TABLE IF NOT EXISTS `ligation` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_person` int(11) NOT NULL,
  `id_raca` int(11) DEFAULT NULL,
  `id_class` int(11) DEFAULT NULL,
  `id_arm` int(11) DEFAULT NULL,
  `id_power` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id_person` (`id_person`),
  KEY `id_raca` (`id_raca`),
  KEY `id_class` (`id_class`),
  KEY `id_arm` (`id_arm`),
  KEY `id_power` (`id_power`),
  CONSTRAINT `ligation_ibfk_1` FOREIGN KEY (`id_person`) REFERENCES `persons` (`id`),
  CONSTRAINT `ligation_ibfk_2` FOREIGN KEY (`id_raca`) REFERENCES `racas` (`id`),
  CONSTRAINT `ligation_ibfk_3` FOREIGN KEY (`id_class`) REFERENCES `class` (`id`),
  CONSTRAINT `ligation_ibfk_4` FOREIGN KEY (`id_arm`) REFERENCES `arms` (`id`),
  CONSTRAINT `ligation_ibfk_5` FOREIGN KEY (`id_power`) REFERENCES `powers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.ligation: ~6 rows (aproximadamente)
/*!40000 ALTER TABLE `ligation` DISABLE KEYS */;
INSERT INTO `ligation` (`id`, `id_person`, `id_raca`, `id_class`, `id_arm`, `id_power`) VALUES
	(1, 1, 6, 1, 2, 1),
	(2, 2, 5, 1, 2, 1),
	(14, 2, NULL, NULL, NULL, NULL),
	(15, 2, NULL, NULL, NULL, 4),
	(16, 2, NULL, NULL, NULL, 3),
	(17, 2, NULL, NULL, NULL, 2);
/*!40000 ALTER TABLE `ligation` ENABLE KEYS */;

-- Copiando estrutura para tabela world.persons
CREATE TABLE IF NOT EXISTS `persons` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `xp` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.persons: ~2 rows (aproximadamente)
/*!40000 ALTER TABLE `persons` DISABLE KEYS */;
INSERT INTO `persons` (`id`, `name`, `xp`) VALUES
	(1, 'Filgofin', 2),
	(2, 'Frodo', 2),
	(3, 'Leandro', 2),
	(4, 'teste', 10);
/*!40000 ALTER TABLE `persons` ENABLE KEYS */;

-- Copiando estrutura para tabela world.powers
CREATE TABLE IF NOT EXISTS `powers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `power` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.powers: ~4 rows (aproximadamente)
/*!40000 ALTER TABLE `powers` DISABLE KEYS */;
INSERT INTO `powers` (`id`, `power`) VALUES
	(1, 'Ataque de furia'),
	(2, 'Cura'),
	(3, 'Furtividade'),
	(4, 'Resistencia Enorme'),
	(5, 'Visao estendida');
/*!40000 ALTER TABLE `powers` ENABLE KEYS */;

-- Copiando estrutura para tabela world.racas
CREATE TABLE IF NOT EXISTS `racas` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `raca` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.racas: ~5 rows (aproximadamente)
/*!40000 ALTER TABLE `racas` DISABLE KEYS */;
INSERT INTO `racas` (`id`, `raca`) VALUES
	(1, 'Humano'),
	(2, 'Elfo'),
	(3, 'Orc'),
	(4, 'Anao'),
	(5, 'Hobbit'),
	(6, 'Anão');
/*!40000 ALTER TABLE `racas` ENABLE KEYS */;

-- Copiando estrutura para tabela world.teste
CREATE TABLE IF NOT EXISTS `teste` (
  `name` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.teste: ~2 rows (aproximadamente)
/*!40000 ALTER TABLE `teste` DISABLE KEYS */;
INSERT INTO `teste` (`name`) VALUES
	('dededexs'),
	('dcdededexs');
/*!40000 ALTER TABLE `teste` ENABLE KEYS */;

-- Copiando estrutura para tabela world.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(100) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

-- Copiando dados para a tabela world.users: ~0 rows (aproximadamente)
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` (`id`, `email`, `password`) VALUES
	(1, 'leandronovaes@hotmail.com', '123456');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;

-- Copiando estrutura para view world.armsview
-- Removendo tabela temporária e criando a estrutura VIEW final
DROP TABLE IF EXISTS `armsview`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `armsview` AS SELECT a.army
FROM persons p
JOIN ligation l
ON l.id_person = p.id
JOIN arms a
ON l.id_arm = a.id ;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
