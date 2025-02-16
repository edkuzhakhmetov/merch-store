# merch-store

## POST | /api/auth 
Аутентификация и получение JWT-токена. 
При первой аутентификации пользователь создается автоматически.

### Тело запроса (AuthRequest):
- `username` (string, обязательное) - Имя пользователя для аутентификации
- `password` (string, обязательное) - Пароль для аутентификации

### Возможные ответы:

#### 1. Пользователь существует. Логин и пароль верны (200 OK)
Успешный ответ (AuthResponse):
- `token` (string) - JWT-токен для доступа к защищенным ресурсам

#### 2. Пользователь существует. Неправильный пароль (401 Unauthorized)
Ответ с ошибкой (ErrorResponse):
- `errors` (string) - "Неавторизован"

#### 3. Регистрация пользователя (201 Created) (NEW)
Статус изменен, т.к. тут создается новый пользователь
Успешный ответ (AuthResponse):
- `token` (string) - JWT-токен, аналогично п.1

#### 4. Некорректное тело запроса (400 Bad request)
Ответ с ошибкой (ErrorResponse):
- `errors` (string) - "Неверный запрос"

#### 5. Прочие ошибки сервера(500 Internal Server Error)
Ответ с ошибкой (ErrorResponse):
- `errors` (string) - "Внутренняя ошибка сервера"

Каждому новому сотруднику выделяется 1000 монет, которые можно использовать для покупки товаров.

Сотрудников может быть до 100к

## POST | /api/buyItem (NEW)
Купить предмет за монеты. 

Количество монеток не может быть отрицательным, запрещено уходить в минус при операциях с монетками.

**Было GET | /api/buy/{item}.** Требования некорректны, т.к.запрос GET не должен изменять состояние на сервере. 
Наверное, правильнее реализовать:
1) GET | /api/items - получение списка всех предметов (как пользователь заранее узнает какие есть/доступны?)
2) POST | /api/items - добавление предмета
3) GET | /api/items/{item} - получение инфо о конкретном предмете
4) POST |/api/items/{item}/buy - купить конкретный предмет (но тогда надо сделать хотя бы п.3, а зачем делать то, что не требуется? :-) )

Поэтому для покупки решено использовать: **POST | /api/buyItem**

### Тело запроса (buyItem):
- `item` (string, обязательное) - Наименование предмета

### Возможные ответы:
Тело ответа не предложено в требованиях, поэтому минимальная доработка - сообщить об успехе операции.

#### 1. Предмет куплен (201 Created)
Статус изменен
Успешный ответ:
- `message` (string) - "Операция выполнена успешно"

#### 2. Токен отсутствует или неверный (401 Unauthorized)
Ответ с ошибкой:
- `errors` (string) - "Неавторизован"

#### 3. Некорректное тело запроса (400 Bad request)
Ответ с ошибкой:
- `errors` (string) - "Неверный запрос"

#### 4. Недостаточно коинов (402 Payment Required) (NEW)
Новый код больше подходит, чем 400 или 500
Ответ с ошибкой:
- `errors` (string) - "Недостаточно коинов"

#### 5. Несуществующий предмет (404 Not Found) - 404 Not Found (NEW)
Ответ с ошибкой:
- `errors` (string) - "Такой предмет не найден"

#### 6. Прочие ошибки (500 Internal Server Error)
Ответ с ошибкой:
- `errors` (string) - "Внутренняя ошибка сервера"

## POST | /api/sendCoin
Отправить монеты другому пользователю.

### Тело запроса:
- `toUser` (string, обязательное) - Имя пользователя-получателя
- `amount` (integer, обязательное) - Количество монет для отправки

Количество монеток не может быть отрицательным, запрещено уходить в минус при операциях с монетками.

### Возможные ответы:
Тело ответа не предложено в требованиях, поэтому минимальная доработка - сообщить об успехе операции.

#### 1. Коины отправлены (201 Created) (NEW)
Статус изменен
Успешный ответ:
- `message` (string) - "Операция выполнена успешно"

#### 2. Токен отсутствует или неверный (401 Unauthorized)
Ответ с ошибкой:
- `errors` (string) - "Неавторизован"

#### 3. Некорректное тело запроса (400 Bad request)
Ответ с ошибкой:
- `errors` (string) - "Неверный запрос"

#### 4. Недостаточно коинов (402 Payment Required) (NEW)
Ответ с ошибкой:
- `errors` (string) - "Недостаточно коинов"

#### 5. Несуществующий получатель (404 Not Found) (NEW)
Ответ с ошибкой:
- `errors` (string) - "Такой получатель не найден"

#### 6. Прочие ошибки сервера (500 Internal Server Error)
Ответ с ошибкой:
- `errors` (string) - "Внутренняя ошибка сервера"

## GET | /api/info
Получить информацию о монетах, инвентаре и истории транзакций.

Каждый сотрудник должен иметь возможность видеть:
- Сгруппированную информацию о перемещении монеток в его кошельке, включая:  
  - Кто ему передавал монетки и в каком количестве  
  - Кому сотрудник передавал монетки и в каком количестве

### Возможные ответы:
#### 1. Успешный запрос (200 OK)
#### 2. Токен отсутствует или неверный (401 Unauthorized)
#### 3. Некорректный запрос (400 Bad request)
#### 4. Прочие ошибки (500 Internal Server Error)
