USE `ostodo`;
CREATE TABLE IF NOT EXISTS `tasks` (
    `user` text NOT NULL,
    `uuid` BINARY(16) NOT NULL,
    `Name` text NOT NULL,
    `CompletionTime` datetime NOT NULL,
    `Repetitions` int,
    PRIMARY KEY (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;