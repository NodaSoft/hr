### Тестовое задание для разработчика на PHP
Мы ожидаем, что Вы найдете все возможные ошибки (синтаксические, проектирования, безопасности и т.д.)

### !
1. Выполнил задание исходя из предположения, что не должно быть зависимостей.
2. Невозможно идентифицировать пользователя без уникального идентификатора, обычно:
    - login
    - email
    - phone
   соответственно добавил поле login.
3. Создал структуру каталогов, сделал похожую на symfony.
4. Пришлось переписать практически полностью оба файла

### Проект можно запустить так:
- git clone git@github.com:onnov/noda_soft_hr.git -b dev-hr
- cd noda_soft_hr/php
- сp -f .env.dist .env
- docker compose up -d
- docker exec -it ns_php composer i

### Запуск тестов:
 docker exec -it ns_php php tests/run.ph
 
#### результат:
```shell
getUserAgeFrom: TRUE
getUserByName: TRUE
getUserByLogin: TRUE
listUsersByName: TRUE
listUsersByLogin: TRUE
addUser: TRUE
addUser (Exception): TRUE
addUsers: TRUE
addUsers (Exception): TRUE
```
7 методов, 9 тестов

### Зпруск фиксеров и статических анализаторов:
 docker exec -it ns_php bin/dev-checks.sh
 