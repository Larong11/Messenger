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
  * User.Service имеет функцию проверки что username и email уникален. 
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