# ИСУ

Сервис имеет следующий API:
- /auth - авторизация при помощи ITMO.ID
- POST /user/users/add - добавить информацию о себе в базу
- GET /user/users/find - найти публичную информацию пользователя по номеру телефона (query argument)
- PUT /user/users/update - обновить свою публичную информацию
- GET /admin/users - получить все данные пользователей
- PUT /admin/users - обновить все данные пользователя
