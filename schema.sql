CREATE TABLE `coin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `symbol` varchar(64) NOT NULL DEFAULT '',
  `name` varchar(64) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_symbol` (`name`,`symbol`)
) ENGINE=InnoDB AUTO_INCREMENT=1691 DEFAULT CHARSET=utf8;

CREATE TABLE `user_coin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(64) NOT NULL DEFAULT '' COMMENT 'user address',
  `symbol` varchar(64) NOT NULL DEFAULT '' COMMENT 'coin symbol',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT 'coin name',
  `cnt` int(11) NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_symbol_name` (`user`,`symbol`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

