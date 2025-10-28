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

    * Работа с базой данных (PostgreSQL, Redis и т.д.)
    * Интеграция с внешними сервисами (SMTP, push-уведомления, платежи)
    * Взаимодействие с файловой системой, очередями сообщений, логирование
    * Реализация интерфейсов, определённых в Domain или Application

  #### пример:
    * `PostgresUserRepository` — реализация интерфейса `UserRepository`
    * `SMTPMailService` — реализация отправки писем
    * Любые драйверы и адаптеры для внешних сервисов

* ### Interfaces
  Interfaces — слой, через который внешние клиенты взаимодействуют с приложением.  
  Этот слой **не содержит бизнес-логики**, а только обрабатывает входящие запросы и формирует ответы.

    * HTTP хэндлеры (`UserHandler`)
    * gRPC серверы
    * CLI команды

  #### пример:
    * `UserHandler.CheckUserName` — принимает HTTP-запрос, вызывает use-case и возвращает JSON
    * `UserHandler.CheckEmail` — принимает email для проверки, вызывает use-case и формирует ответ
    * `UserHandler.CreateUser` — принимает данные пользователя, вызывает `RegisterUserUseCase`, возвращает результат


## Полный пример работы DDD на примере создания пользователя
1. **Сервер получает HTTP-запрос**

2. **Handler (Interface)** обрабатывает запрос и вызывает нужный `UserUseCase.CreateUser`

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
## Запрос кода подтверждения
http запрос
* метод: POST
* URL: "/api/users/request-verification-code"
* Content-Type: application/json
* Тело:
```json
{
  "email": "john@example.com"
}
```
http ответ:

* Успешная отправка кода

  * HTTP статус: 200 OK

  * Тело ответа: пустое

* Неверный HTTP метод

  * HTTP статус: 405 Method Not Allowed

  * Тело ответа: "method not allowed"

* Некорректный JSON в теле запроса

  * HTTP статус: 400 Bad Request

  * Тело ответа: "invalid JSON body"

* Некорректный формат email

  * HTTP статус: 400 Bad Request

  * Тело ответа: "invalid email format"

* Пользователь с таким email уже существует

  * HTTP статус: 400 Bad Request

  * Тело ответа: "user with such email exists"

* Ошибка генерации кода

  * HTTP статус: 500 Internal Server Error

  * Тело ответа: <текст ошибки>

* Внутренняя ошибка сервера

  * HTTP статус: 500 Internal Server Error

  * Тело ответа: <текст ошибки>
## Создание пользователя
http запрос
* метод: POST
* URL: "/api/users/create-user"
* Content-Type: application/json
* Тело:
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "avatar_url": "https://example.com/avatar.png",
  "verification_code": "123456"
}
```
http ответ:

* Успешное создание пользователя

  * HTTP статус: 200 OK

  * JSON: {"id": 42}

* Неверный HTTP метод

  * HTTP статус: 405 Method Not Allowed

  * Тело ответа: method not allowed

* Некорректный JSON в теле запроса

  * HTTP статус: 400 Bad Request

  * Тело ответа: invalid JSON body

* Ошибка валидации username

  * HTTP статус: 400 Bad Request

  * Тело ответа: "username must be 3–20 characters long" или "username may contain only letters, digits, and underscores"

* Ошибка валидации email

  * HTTP статус: 400 Bad Request

  * Тело ответа: "invalid email format"

* Пользователь с таким username уже существует

  * HTTP статус: 400 Bad Request

  * Тело ответа: "username with such username already exists"

* Пользователь с таким email уже существует

  * HTTP статус: 400 Bad Request

  * Тело ответа: "user with such email already exists"
* Истёк код верификации

  * HTTP статус: 400 Bad Request

  * Тело ответа: "verification code expired"

* Некорректный код верификации

  * HTTP статус: 400 Bad Request

  * Тело ответа: "incorrect verification code"

* Ошибка при генерации хеша пароля

  * HTTP статус: 500 Internal Server Error

  * Тело ответа: "failed to generate password hash"

* Внутренняя ошибка сервера

  * HTTP статус: 500 Internal Server Error

  * Тело ответа: <текст ошибки>
# Поток обработки на сервере
## Check UserName
Client  
&nbsp;&nbsp;&nbsp;&nbsp;→ POST /api/users/check-username  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.CheckUserName](interface/http/handlers/user_handler.go)  
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
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.CheckEmail](interface/http/handlers/user_handler.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [RegisterUserUseCases.CheckEmail](application/use_cases/user/register.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [domain validation](domain/user/validation.go) (email rules)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByEmail](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [SQL запрос в БД](infrastructure/persistence/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← проверка существования  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← результат use case  
&nbsp;&nbsp;&nbsp;&nbsp;← формирование JSON-ответа  
← Client получает {"available": true/false}

## Request Verification Code
Client  
&nbsp;&nbsp;&nbsp;&nbsp;→ POST /api/users/request-verification-code  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.RequestVerificationCode](interface/http/handlers/user_handler.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [RegisterUserUseCases.RequestVerificationCode](application/use_cases/user/register.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [domain validation](domain/user/validation.go) (email rules)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByEmail](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ проверка существования пользователя  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.CreateVerificationCode](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← генерация кода и сохранение в БД  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← результат use case (успех или ошибка)  
&nbsp;&nbsp;&nbsp;&nbsp;← формирование HTTP-ответа  
← Client получает HTTP 200 OK или сообщение об ошибке

## Create User
Client  
&nbsp;&nbsp;&nbsp;&nbsp;→ POST /api/users/create-user  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserHandler.CreateUser](interface/http/handlers/user_handler.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [RegisterUserUseCases.RegisterUser](application/use_cases/user/register.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [domain validation](domain/user/validation.go) (username, email rules)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByUserName](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.FindByEmail](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ проверка существования пользователя  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [packageuser.GeneratePasswordHash](domain/user/generate.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [packageuser.NewUser](domain/user/entity.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.GetVerificationCode](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ проверка кода верификации и TTL  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;→ [UserRepository.CreateUser](domain/repositories/user_repository.go)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← результат use case (успешное создание или ошибка)  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;← формирование HTTP-ответа  
← Client получает {"id": 42} или сообщение об ошибке
