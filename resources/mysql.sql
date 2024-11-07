CREATE TABLE `paciente_detalhes` (
`id` int(15) unsigned NOT NULL,
`nome_completo` varchar(100) NOT NULL,
`endereco` varchar(255),
`sexo` varchar(10),
`telefone` int(15) unsigned,
`observacoes` varchar(255),
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `paciente_registros` (
`id` int(15) unsigned NOT NULL,
`token` int(10) unsigned NOT NULL,
PRIMARY KEY (`token`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `inspecao_relatorios` (
`id` int(15) unsigned NOT NULL,
`medicacao` varchar(255),
`exames` varchar(255),
`anotacoes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `inspecao_detalhes` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50) NOT NULL,
`observacoes` varchar(255),
`medicacao` varchar(255),
`exames` varchar(255),
`anotacoes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `medicacao_relatorios` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`dose` varchar(255),
`anotacoes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `exame_relatorios` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`nome_exame` varchar(100),
`resultados` varchar(255),
`situacao` varchar(50),
`anotacoes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `liberacao_relatorios` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`proximo_estado` varchar(50),
`pos_medicacao` varchar(255),
`anotacoes` varchar(255)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE `alta_detalhes` (
`id` int(15) unsigned NOT NULL,
`time` varchar(50),
`estado` varchar(50),
`pos_medicacao` varchar(255),
`anotacoes` varchar(255),
`proxima_visita` varchar(50)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;