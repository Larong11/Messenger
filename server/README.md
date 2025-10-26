# документация
## структура проекта с примерами

* ### Domain
  Domain (предметная область) — это сердце приложения.  
  * Содержит сущности, их свойства и бизнес-правила.  
  * Не зависит от базы данных, HTTP, UI или внешних сервисов.  
  * Отвечает только за правильность бизнес-логики.  
  Entity — сущность и правила.  
  Service — бизнес правила не привязанные к одной сущности.
  Repositories — папка с интерфейсами репозиториев.

  #### пример:
  * User — сущность.
  * User имеет функцию create, где есть проверка что поля не пустые, возраст больше минимального и так далее.  
  * User.Service имеет функцию проверки что username и email уникален, вообще такое относящееся к бд можно и в use_case оставить.
  * Repositories.UserRepository — интерфейс репозитория, которые реализован в Infrastructure.  
  
* ### Application
  Application слой — это сценарии работы системы (use-cases).
  * Он координирует действия между Domain и инфраструктурой.  
  * Не содержит бизнес-правил самих сущностей — все правила лежат в Domain.  
  * Знает какие шаги нужно сделать, чтобы выполнить конкретную задачу пользователя.  
  #### пример:
  * UserUseCase.CreateUser — use_case для созданий пользователя.
  * Проверяет уникальность username and email через Domain.User.Service
  * создает domain-entity через User.NewUser
  * обращается к интерфейсу userRepository который лежит в domain
* ### Infrastructure
  Infrastructure — это слой, который обеспечивает конкретные реализации технических деталей, необходимых для работы приложения:

  * работа с базой данных
  * HTTP / WebSocket / gRPC
  * кэширование, очередь сообщений, логирование
  * интеграция с внешними сервисами (email, push, платежи)
  #### пример
    * реализация PostgresUserRepository
    * реализация httpHandler

## Полный пример работы DDD на примере создания пользователя
1. **Сервер получает HTTP-запрос**

2. **Handler (Infrastructure)** обрабатывает запрос и вызывает нужный `UserUseCase.CreateUser`

3. **UserUseCase (Application)** выполняет сценарий:

    1. Вызывает `Domain.UserService` для проверки уникальности `email` и `username`
        - Использует интерфейс `Domain.Repositories.UserRepository`

    2. Создаёт нового пользователя через `Domain.User.CreateNewUser`
        - Проверяется `FirstName`, `LastName`, `age`
        - Интерфейс `UserRepository` реализован в `infrastructure.persistence`

    3. Сохраняет пользователя через `Domain.Repositories.UserRepository`
        - Создаётся код подтверждения
        - Интерфейс `UserRepository` реализован в `infrastructure.persistence`

    4. Вызывает `infrastructure.email.SMTP_Sender` для отправки письма на почту


# Описание запросов к серверу
## Проверка коректности userName
http запрос
* метод: POST
* URL: "/api/users/check-email"
* Content-Type: application/json
* Тело:{ "username": "john_doe" } 
http ответ:
* Успешный запрос, имя доступно
  * HTTP статус: 200 OK
  * JSON: {"available": true}
* Успешный запрос, имя занято
    * HTTP статус: 200 OK
    * JSON: {"available": false}
* Неверный HTTP метод:
  * HTTP статус: 405 Method Not Allowed
  * Тело ответа: method not allowed
* Некорректный JSON в теле запроса
  * HTTP статус: 400 Bad Request
  * Тело ответа: invalid JSON body
* Ошибка валидации username
  * HTTP статус: 400 Bad Request
  * Тело ответа: "username must be 3–20 characters long" или "username may contain only letters, digits, and underscores"
* Внутренняя ошибка сервера
  * HTTP статус: 500 Internal Server Error
  * Тело ответа: <текст ошибки>

## Проверка корректности email

http запрос

* метод: POST

* URL: "/api/users/check-email"

* Content-Type: application/json

* Тело: { "email": "user@example.com" }

http ответ:
* Успешный запрос, email доступен
  * HTTP статус: 200 OK
  * JSON: {"available": true}
* Успешный запрос, email уже занят
  * HTTP статус: 200 OK
  * JSON: {"available": false}
* Неверный HTTP метод
  * HTTP статус: 405 Method Not Allowed
  * Тело ответа: method not allowed
* Некорректный JSON в теле запроса
  * HTTP статус: 400 Bad Request
  * Тело ответа: invalid JSON body
* Ошибка валидации email
  * HTTP статус: 400 Bad Request
  * Тело ответа: "invalid email format"
* Внутренняя ошибка сервера
  * HTTP статус: 500 Internal Server Error
  * Тело ответа: <текст ошибки>

# Поток обработки на сервере
## Check UserName
Client  
&nbsp;&nbsp;&nbsp;&nbsp;→ POST /api/users/check-username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.CheckUserName](infrastructure/http/handlers/user_handler.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [RegisterUserUseCases.CheckUserName](application/use_cases/user/register.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [domain validation](domain/user/validation.go) (username rules)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByUserName](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [SQL запрос в БД](infrastructure/persistence/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← проверка существования  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← результат use case  
&nbsp;&nbsp;&nbsp;&nbsp;← формирование JSON-ответа  
← Client получает {"available": true/false}

## Check Email
Client  
&nbsp;&nbsp;&nbsp;&nbsp;→ POST /api/users/check-email  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.CheckEmail](infrastructure/http/handlers/user_handler.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [RegisterUserUseCases.CheckEmail](application/use_cases/user/register.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [domain validation](domain/user/validation.go) (email rules)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByEmail](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [SQL запрос в БД](infrastructure/persistence/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← проверка существования  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← результат use case  
&nbsp;&nbsp;&nbsp;&nbsp;← формирование JSON-ответа  
← Client получает {"available": true/false}
