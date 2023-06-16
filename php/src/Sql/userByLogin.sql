SELECT `id`,
       `login`,
       `name`,
       `last_name` as `lastName`,
       `from`,
       `age`
FROM `user`
WHERE login = :login;
