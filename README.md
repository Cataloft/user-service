<h3 align="left">Запуск сервиса</h3>

```sh
$ git clone https://github.com/Cataloft/user-service
$ docker-compose up --build
```
<h3 align="left">Доступные ручки</h3>

* **POST-`localhost:12345/users`:** Получает ФИО пользователя, обогащает из открытых api и кладёт в postgres
* **GET-`localhost:12345/users`:** Отдаёт данные пользователeй с фильтрами и пагинацией
   - **Параметры фильтров:** ageGreater=возраст, ageLower=возраст, nameContain=символыВимени, любоеПолеПользователя=значПоляПользователя
   - **Параметры пагинации:** page=номерСтр, pageSize=размерСтр
* **PATCH-`localhost:12345/update/id`:** Обновляет данные пользователя по id
* **DELETE-`localhost:12345/delete/id`:** Удаляет пользователя по id
