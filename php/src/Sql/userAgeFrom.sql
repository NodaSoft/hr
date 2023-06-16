SELECT `id`,
       `login`,
       `name`,
       `last_name`                       as `lastName`,
       `from`,
       `age`,
       JSON_EXTRACT(`settings`, '$.key') as `key`
FROM `user`
WHERE `age` > :age;
