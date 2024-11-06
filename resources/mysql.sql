CREATE TABLE `paciente_detalhes` (
`id` int(15) unsigned NOT NULL,
`nome_completo` varchar(100) NOT NULL,
`endereco` varchar(255),
`sexo` varchar(10),
`telefone` int(15) unsigned,
`observacoes` varchar(255),
PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;