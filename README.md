## Welcome to BreezyNotes. The place where fantasy becomes reality!
#### Development by [Breezy Innovation RZN](https://about.breezynotes.ru)
![BREEZYNOTES](https://i.ibb.co/PvRh0KvX/favicon.png)
### Technology stack:
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
### Frontend repository:
[GitHub](https://github.com/DaniilaRyadinsky/breezy)

# API ДОКУМЕНТАЦИЯ BREEZYNOTES

## Оглавление

1.  [Общие положения](#общие-положения)
*   [Архитектура системы](#архитектура-системы)
*   [Аутентификация](#аутентификация)
*   [Формат ошибок](#формат-ошибок)
2.  [Сервис аутентификации и работы с пользователями](#сервис-аутентификации-и-работы-с-пользователями)
*   [Аутентификация](#аутентификация-endpoints)
*   [Работа с данными пользователя](#работа-с-данными-пользователя)
3.  [Сервис работы с заметками](#сервис-работы-с-заметками)
*   [Работа с тегами](#работа-с-тегами)
*   [Работа с корзиной](#работа-с-корзиной)
*   [Работа с заметками](#работа-с-заметками)
*   [Работа с блоками](#работа-с-блоками)
4.  [Типы блоков и операции над ними](#типы-блоков-и-операции-над-ними)
*   [Тип: `text`](#тип-text)
*   [Тип: `list`](#тип-list)
*   [Тип: `header`](#тип-header)
*   [Тип: `img`](#тип-img)
*   [Тип: `link`](#тип-link)
*   [Тип: `quote`](#тип-quote)
*   [Тип: `code`](#тип-code)
*   [Тип: `file`](#тип-file)

---

## <a name="общие-положения"></a>1. Общие положения

### <a name="архитектура-системы"></a>Архитектура системы

Система BreezyNotes построена на микросервисной архитектуре.

*   **Сервис аутентификации и пользователей**: Отвечает за регистрацию, вход и управление данными пользователей. Работает с базой данных **PostgreSQL**.
*   **Сервис работы с заметками**: Отвечает за всю логику, связанную с заметками, тегами, корзиной и контентными блоками. Работает с базой данных **MongoDB**.

### <a name="аутентификация"></a>Аутентификация

Аутентификация в системе реализована с помощью `access` и `refresh` JWT-токенов.

1.  После успешного входа (`/api/auth`) или регистрации (`/api/auth/reg`) сервер устанавливает токены в `Cookie` браузера с соответствующим временем жизни (`HttpOnly`, `Secure`).
2.  Для всех последующих запросов, требующих авторизации, идентификатор пользователя (`id_user`) извлекается из `access` токена на стороне шлюза (API Gateway).
3.  Если `access` токен истек, клиент должен выполнить запрос на `/api/auth/token` для его обновления с помощью `refresh` токена.

### <a name="формат-ошибок"></a>Формат ошибок

В случае ошибки сервер возвращает соответствующий HTTP-статус и JSON-объект следующего вида:

    {
        "error": "Текст ошибки"
    }


---

## <a name="сервис-аутентификации-и-работы-с-пользователями"></a>2. Сервис аутентификации и работы с пользователями

### <a name="аутентификация-endpoints"></a>Аутентификация

#### `POST /api/auth`
Аутентификация пользователя. В случае успеха возвращает `access` и `refresh` токены в `Cookie` и в теле ответа.

*   **Возможные статусы и ошибки:**
*   `200 OK` - Успешная аутентификация.
*   `400 Bad Request` - Ошибка в запросе:
    *   `"bad JSON"` - Некорректный формат JSON.
    *   `"email and login is empty"` - Не указан email или логин.
    *   `"pw is empty"` - Не указан пароль.
    *   `"password incorrect"` - Неверный пароль.
*   `404 Not Found` - Пользователь с такими данными не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout` - Внутренние ошибки сервера.

#### `POST /api/auth/reg`
Регистрация нового пользователя. В случае успеха возвращает `access` и `refresh` токены.

*   **Возможные статусы и ошибки:**
*   `200 OK` - Успешная регистрация.
*   `302 Found` - Пользователь с таким email или логином уже существует.
*   `400 Bad Request` - Ошибка в запросе:
    *   `"bad JSON"` - Некорректный формат JSON.
    *   `"password not same"` - Пароли не совпадают.
    *   `"email and login is empty"` - Не указан email или логин.
    *   `"pw is empty"` - Не указан пароль.
    *   `"pw not in policy"` - Пароль не соответствует политике безопасности (длина 5-20 символов, минимум одна заглавная буква, один символ и одна цифра).
*   `502 Bad Gateway` / `504 Gateway Timeout` - Внутренние ошибки сервера.

#### `GET /api/auth/token`
Обновление `access` токена с помощью `refresh` токена из `Cookie`.

*   **Возможные статусы и ошибки:**
*   `200 OK` - Токены валидны.
*   `201 Created` - `access` токен успешно обновлен и установлен в `Cookie`.
*   `400 Bad Request` - Ошибка в запросе:
    *   `"one of tokens is empty"` - Отсутствуют необходимые токены.
*   `401 Unauthorized` - Сессия истекла или токены недействительны:
    *   `"refresh_token cookie missing"` - Отсутствует `refresh` токен.
    *   `"invalid token"` - `refresh` токен недействителен.
    *   `"token expired"` - `refresh` токен истек.
*   `502 Bad Gateway` / `504 Gateway Timeout` - Внутренние ошибки сервера.

### <a name="работа-с-данными-пользователя"></a>Работа с данными пользователя

#### `GET /api/user/data`
Получение данных текущего пользователя. Поле с хешем пароля заменяется на пустую строку.

*   **Возможные статусы и ошибки:**
*   `401 Unauthorized` - Ошибка аутентификации.
*   `404 Not Found` - Пользователь не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `DELETE /api/user`
Удаление текущего пользователя и всех его данных.

*   **Возможные статусы и ошибки:**
*   `400 Bad request` - Возникает при неправильном айди
*   `401 Unauthorized` - Ошибка аутентификации.
*   `404 Not Found` - Пользователь не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/user/about`
Изменение информации "о себе".

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON", "id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/user/email`
Изменение email пользователя.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"email is empty", "id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/user/photo`
Изменение фото профиля (аватара). При отправке пустого поля `photo` устанавливается значение по умолчанию `"images/default.png"`.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON", "id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/user/pw`
Изменение пароля пользователя.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` - Ошибка в запросе:
    *   `"bad JSON"`.
    *   `"password not same"` - Пароли не совпадают.
    *   `"new password not in policy"` - Новый пароль не соответствует политике.
    *   `"id not in uuid"`
*   `401 Unauthorized`.
*   `404 Not Found` - Пользователь не найден или текущий пароль неверен.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

---

## <a name="сервис-работы-с-заметками"></a>3. Сервис работы с заметками

Все запросы к этому сервису требуют валидного `access` токена в `Cookie`. ID пользователя извлекается из токена и используется для авторизации доступа к ресурсам (заметкам, тегам и т.д.).

### <a name="работа-с-тегами"></a>Работа с тегами

#### `POST /api/tag`
Создание нового тега.

*   **Возможные статусы и ошибки:**
*   `201 Created` - Тег успешно создан.
*   `400 Bad Request` (`"bad JSON"`, `"title is empty"`, `"color is empty"`, `"emoji is empty"`).
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/tag/title` | `/api/tag/color` | `/api/tag/emoji`
Изменение названия, цвета или эмодзи тега.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"id not in uuid"`, `"field is empty"` (**title, color, emoji**)).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found` - Тег не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `DELETE /api/tag`
Удаление тега по `id`, переданному в query-параметре.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad param"`, `"id not in uuid"`).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found` - Тег не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/tag/by-user`
Получение всех тегов текущего пользователя.

*   **Возможные статусы и ошибки:**
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

### <a name="работа-с-корзиной"></a>Работа с корзиной

#### `DELETE /api/trash`
Полная очистка корзины (безвозвратное удаление всех заметок в ней).

*   **Возможные статусы и ошибки:**
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PUT /api/trash/to`
Перемещение заметки в корзину. **ВНИМАНИЕ: ЕСЛИ ПОЛЬЗОВАТЕЛЬ БЫЛ НЕ АВТОРОМ, ТО ПРИ ЭТОМ ДЕЙСТВИИ ОН УДАЛЯЕТСЯ ИЗ СПИСКОВ В ЗАМЕТКЕ И ОНА СТАНОВИТСЯ ЕМУ НЕ ДОСТУПНА.**

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad param"`, `"id not in uuid"`).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found` - Заметка не найдена.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PUT /api/trash/from`
Восстановление заметки из корзины.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad param"`, `"id not in uuid"`).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found` - Заметка не найдена.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/trash`
Получение списка заметок, находящихся в корзине.

*   **Возможные статусы и ошибки:**
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

### <a name="работа-с-заметками"></a>Работа с заметками

#### `POST /api/note`
Создание новой заметки.

*   **Возможные статусы и ошибки:**
*   `201 Created` - Заметка успешно создана.
*   `400 Bad Request` (`"title is empty"`, `"bad JSON"`).
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/note`
Получение полной информации о заметке, включая список ее блоков.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad param"`, `"id not in uuid"`).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found` - Заметка не найдена.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/note/all`
Получение списка всех заметок пользователя (с пагинацией).

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` - Ошибки пагинации (`"bad start"`, `"bad end"`, `"start < 0!"`, `"start must be int"`, `"end must be int"`).
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/note/by-tag`
Получение списка заметок по указанному тегу (с пагинацией).

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` - Ошибки пагинации или ID тега.
*   `401 Unauthorized`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/note/title`
Изменение названия заметки.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"id not in uuid"`, `"title is empty"`).
*   `401 Unauthorized` (включая `"you dont have permission"`).
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `POST /api/note/tag`
Добавление тега к заметке.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found` - Заметка или тег не найдены.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `DELETE /api/note/tag`
Удаление тега из заметки.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found` - Заметка или тег не найдены.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/note/share`
Предоставление доступа к заметке другому пользователю (по логину) с указанием роли (`read`, `write`).

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"bad note id"`, `"role is empty"`, `"role undefined"`).
*   `401 Unauthorized` (только владелец может делиться).
*   `404 Not Found` - Заметка или пользователь для шеринга не найден.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

### <a name="работа-с-блоками"></a>Работа с блоками

#### `POST /api/block`
Создание нового блока в заметке.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"bad data"`, `"pos < 0"`, `"type is empty"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `424 Failed Dependency` - Используется незарегистрированный тип блока.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `POST /api/block/op`
Выполнение операции над блоком (редактирование контента, изменение атрибутов).

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"bad data"`, `"op is empty"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `424 Failed Dependency`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/block/type`
Изменение типа существующего блока (например, `text` -> `header`).

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"type is empty"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `424 Failed Dependency`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `DELETE /api/block`
Удаление блока.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"no note id"`, `"no block id"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `PATCH /api/block/order`
Изменение порядка блоков внутри заметки.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"order < 0"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.

#### `GET /api/block`
получение блока. Возможно не используемая.

*   **Возможные статусы и ошибки:**
*   `400 Bad Request` (`"bad JSON"`, `"id not in uuid"`).
*   `401 Unauthorized`.
*   `404 Not Found`.
*   `502 Bad Gateway` / `504 Gateway Timeout`.



---

## <a name="типы-блоков-и-операции-над-ними"></a>4. Типы блоков и операции над ними

### <a name="тип-text"></a>Тип: `text`
Блок для отображения простого текста с возможностью стилизации.

### Создание
`POST /api/block`

```json
{
  "type": "text",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text": [
      {
        "style": "default",
        "string": "Ваш текст"
      }
    ]
  }
}
```

### Операции
`POST /api/block/op`

#### `apply_style`

```json
{
  "op": "apply_style",
  "data": {
    "start": 0,
    "end": 2,
    "style": "bold"
  }
}
```

#### `insert_text`

```json
{
  "op": "insert_text",
  "data": {
    "pos": 3,
    "new_text": "вставляемый текст"
  }
}
```

#### `delete_range`

```json
{
  "op": "delete_range",
  "data": {
    "start": 0,
    "end": 3
  }
}
```

---

### <a name="тип-list"></a>Тип: `list`
Блок для списков (`ordered`, `unordered`, `todo`).
### Создание
`POST /api/block`

```json
{
  "type": "list",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text_data": {
      "text": [
        {
          "style": "default",
          "string": "Элемент списка"
        }
      ]
    },
    "level": 0,
    "type": "ordered",
    "value": 1
  }
}
```

### Операции
`POST /api/block/op`

#### `change_value`

```json
{
  "op": "change_value",
  "data": {
    "new_value": 15
  }
}
```

#### `change_level`

```json
{
  "op": "change_level",
  "data": {
    "new_level": 2
  }
}
```

#### `change_type`

```json
{
  "op": "change_type",
  "data": {
    "new_type": "todo"
  }
}
```
### <a name="тип-header"></a>Тип: `header`
Блок для заголовков (уровни 1-3).

### Создание
`POST /api/block`

```json
{
  "type": "header",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text_data": {
      "text": [
        {
          "style": "default",
          "string": "Текст заголовка"
        }
      ]
    },
    "level": 1
  }
}
```

### Операции

```json
{
  "op": "change_level",
  "data": {
    "new_level": 3
  }
}
```
### <a name="тип-img"></a>Тип: `img`
Блок для изображений.

### Создание

```json
{
  "type": "img",
  "note_id": "...",
  "pos": 0,
  "data": {
    "alt": "Альтернативный текст",
    "src": "путь/к/изображению.png"
  }
}
```

### Операции

```json
{
  "op": "change_src",
  "data": {
    "new_src": "новый/путь/к/img.png"
  }
}
```

```json
{
  "op": "change_alt",
  "data": {
    "new_alt": "Новый альт. текст"
  }
}
```

### <a name="тип-link"></a>Тип: `link`
Блок для гиперссылок.

### Создание

```json
{
  "type": "link",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text": "Текст ссылки",
    "url": "https://example.com"
  }
}
```

### Операции

```json
{
  "op": "change_text",
  "data": {
    "new_text": "Новый текст для ссылки"
  }
}
```

```json
{
  "op": "change_url",
  "data": {
    "new_url": "https://new-example.com"
  }
}
```

### <a name="тип-quote"></a>Тип: `quote`
Блок для цитат.

### Создание

```json
{
  "type": "quote",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text": "Текст цитаты."
  }
}
```

### Операции

```json
{
  "op": "change_text",
  "data": {
    "new_text": "Новый текст цитаты."
  }
}
```


### <a name="тип-code"></a>Тип: `code`
Блок для фрагментов кода.

### Создание

```json
{
  "type": "code",
  "note_id": "...",
  "pos": 0,
  "data": {
    "text": "console.log('Hello, World!');",
    "lang": "javascript"
  }
}
```

### Операции

```json
{
  "op": "change_text",
  "data": {
    "new_text": "print('Hello, Python!')"
  }
}
```

```json
{
  "op": "analyse_lang",
  "data": {}
}
```

### <a name="тип-file"></a>Тип: `file`
Блок для прикрепленных файлов.

### Создание

```json
{
  "type": "file",
  "note_id": "...",
  "pos": 0,
  "data": {
    "src": "путь/к/файлу.pdf"
  }
}
```

### Операции

```json
{
  "op": "change_src",
  "data": {
    "new_src": "новый/путь/к/файлу.docx"
  }
}
```
