# Cat Houses API Documentation

## Overview
API для управления домиками для кошек. Домики имеют название и вместимость. Кошек можно помещать в домики, если есть свободные места. Один кот может находиться только в одном домике.

## Models

### CatHouse
```json
{
  "id": 1,
  "name": "Большой дом",
  "capacity": 5,
  "current_occupancy": 0,
  "available_spaces": 5,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### CatHouseCreateRequest
```json
{
  "name": "Новый домик",
  "capacity": 3
}
```

### CatHouseUpdateRequest
```json
{
  "name": "Обновленное название",
  "capacity": 4
}
```

### CatHouseOccupancyRequest
```json
{
  "cat_house_id": 1
}
```

### CatHouseOccupancyResponse
```json
{
  "id": 1,
  "cat_id": 1,
  "cat_house_id": 1,
  "cat_name": "Мурзик",
  "cat_house_name": "Большой дом",
  "occupied_at": "2024-01-01T00:00:00Z"
}
```

## Public Endpoints

### GET /api/v1/cat-houses
Получить список всех домиков.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Большой дом",
    "capacity": 5,
    "current_occupancy": 0,
    "available_spaces": 5,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /api/v1/cat-houses/{id}
Получить информацию о конкретном домике.

**Response:** CatHouse

### GET /api/v1/cat-houses/{id}/cats
Получить список котов в домике.

**Response:**
```json
[
  {
    "id": 1,
    "cat_id": 1,
    "cat_house_id": 1,
    "cat_name": "Мурзик",
    "cat_house_name": "Большой дом",
    "occupied_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /api/v1/cat-houses/{id}/with-cats
Получить домик со списком котов в нем.

**Response:**
```json
{
  "house": {
    "id": 1,
    "name": "Большой дом",
    "capacity": 5,
    "current_occupancy": 2,
    "available_spaces": 3,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "cats": [
    {
      "id": 1,
      "name": "Мурзик",
      "age": 3,
      "description": "Ласковый и игривый кот",
      "user_id": 1,
      "cat_breed_id": null,
      "created_at": "2024-01-01T00:00:00Z",
      "can_edit": false
    }
  ]
}
```

## Protected Endpoints (Require Authentication)

### POST /api/v1/cat-houses
Создать новый домик.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request:** CatHouseCreateRequest

**Response:** CatHouse (201 Created)

### PUT /api/v1/cat-houses/{id}
Обновить информацию о домике.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request:** CatHouseUpdateRequest

**Response:**
```json
{
  "message": "Cat house updated successfully"
}
```

**Errors:**
- `400 Bad Request`: Если новая вместимость меньше текущего количества котов

### DELETE /api/v1/cat-houses/{id}
Удалить домик.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Cat house deleted successfully"
}
```

**Note:** При удалении домика все коты автоматически удаляются из него.

### GET /api/v1/cat-houses/cats/{catId}/occupancy
Получить информацию о том, в каком домике находится кот.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:** CatHouseOccupancyResponse

**Errors:**
- `404 Not Found`: Если кот не находится в домике

### POST /api/v1/cat-houses/cats/{catId}/add
Поместить кота в домик.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request:** CatHouseOccupancyRequest

**Response:** CatHouseOccupancyResponse (201 Created)

**Errors:**
- `400 Bad Request`: Если домик заполнен
- `400 Bad Request`: Если кот уже находится в домике
- `403 Forbidden`: Если пользователь не является владельцем кота
- `404 Not Found`: Если кот или домик не найдены

### DELETE /api/v1/cat-houses/cats/{catId}/remove
Удалить кота из домика.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "message": "Cat removed from house successfully"
}
```

**Errors:**
- `400 Bad Request`: Если кот не находится в домике
- `403 Forbidden`: Если пользователь не является владельцем кота
- `404 Not Found`: Если кот не найден

## Business Rules

1. **Емкость домика:** Домик имеет фиксированную вместимость (capacity)
2. **Текущая занятость:** Текущее количество котов в домике (current_occupancy)
3. **Свободные места:** available_spaces = capacity - current_occupancy
4. **Один кот - один домик:** Кот может находиться только в одном домике одновременно
5. **Проверка владельца:** Только владелец кота (или администратор) может помещать/удалять кота из домика
6. **Обновление вместимости:** Нельзя уменьшить вместимость домика ниже текущего количества котов
7. **Удаление домика:** При удалении домика все связи с котами удаляются автоматически

## Error Responses

Все ошибки возвращаются в формате:
```json
{
  "error": "Error message"
}
```

### Common Error Codes:
- `400 Bad Request`: Неверные данные запроса
- `401 Unauthorized`: Требуется аутентификация
- `403 Forbidden`: Нет прав доступа
- `404 Not Found`: Ресурс не найден
- `500 Internal Server Error`: Внутренняя ошибка сервера