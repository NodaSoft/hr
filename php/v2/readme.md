# Notification System

## Описание

Этот проект реализует систему уведомлений, которая отправляет уведомления по электронной почте и SMS при изменении статуса [чего-то]. Проект разработан с соблюдением принципов SOLID и использует подходы объектно-ориентированного программирования (ООП).

## Структура проекта

### Контроллеры

- **SendNotificationController**: Обрабатывает HTTP-запросы для отправки уведомлений. Извлекает данные из запроса и передает их в сервис `ReferencesOperation`.

### Сервисы

- **TsReturnOperation**: Реализует логику отправки уведомлений по электронной почте и SMS. Получает необходимые данные из DTO, взаимодействует с репозиториями для получения сущностей и вызывает уведомления.

### DTO (Data Transfer Objects)

- **SendNotificationDTO**: Содержит данные для отправки уведомлений.
- **MessageDifferenceDto**: Содержит данные о различиях статусов.
- **GetNotificationDifferenceDTO**: Содержит данные для получения различий уведомлений.

### Репозитории (Interfaces)

- **ClientRepositoryInterface**: Интерфейс для работы с клиентами.
- **EmployeeRepositoryInterface**: Интерфейс для работы с сотрудниками.
- **SellerRepositoryInterface**: Интерфейс для работы с продавцами.
- **GetDifferencesInterface**: Интерфейс для получения различий уведомлений.

### Сущности (Entities)

- **Client**: Представляет клиента.
- **Employee**: Представляет сотрудника.
- **Seller**: Представляет продавца.

### События и Уведомления

- **EventDispatcher**: Отправляет события.
- **ChangeReturnStatusEvent**: Событие изменения статуса возврата.
- **NotificationDispatcher**: Отправляет уведомления.
- **MailNotificationInterfaceClient**: Уведомление по электронной почте.
- **SmsNotificationInterfaceClient**: Уведомление по SMS.

## Установка

    ```bash
    composer install
    ```
## Запустить Тесты

    ```bash
    vendor/bin/phpunit tests
    ```



