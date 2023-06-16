CREATE TABLE `user`
(
    `id`         INT                                   NOT NULL AUTO_INCREMENT,
    `login`      VARCHAR(100)                          NOT NULL COMMENT 'Логин',
    `name`       VARCHAR(100)                          NOT NULL COMMENT 'Имя',
    `last_name`  VARCHAR(100)                          NOT NULL COMMENT 'Фамилия',
    `age`        INT                                   NOT NULL COMMENT 'Возраст',
    `from`       VARCHAR(255)                          NULL     DEFAULT NULL COMMENT 'Откуда',
    `settings`   JSON                                  NULL     DEFAULT NULL COMMENT 'Настройки',
    `created_at` TIMESTAMP                             NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Время создания',
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Время обновления',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;
ALTER TABLE `user` ADD UNIQUE(`login`);
ALTER TABLE `user` ADD INDEX(`name`);
ALTER TABLE `user` ADD INDEX(`age`);
