# RedditWithTarantool

## Environment example (.env)

```sh
DB_HOST="127.0.0.1:3301"
DB_USER="admin"
DB_PASSWORD="pass"
SERVER_PORT="8085"
```



## Spaces



### post

| Название | post        |
| -------- | ----------- |
| индексы  | primary(id) |
| id       | unsigned    |
| content  | string      |



### comment

| Название | comment                       |
| -------- | :---------------------------- |
| индексы  | primary(id), ref_idx(ref, id) |
| id       | unsigned                      |
| content  | string                        |
| ref      | unsigned                      |



## API

###  [POST] /post  - Создать пост

Запрос

```json
{
    "content": "my text"
}
```

Ответ

```json
{
    "status": "ok"
}
```



### [GET] /posts - Получить все посты

Запрос

```json

```

Ответ

```json
{
    "posts": [
        {
            "id": 1,
            "content": "post text"
        },
    	...
    ]
}
```



### [POST] /comment?id=<post_id> - Создать комментарий

Запрос

```json
{
    "content": "post text",
    "ref": 1
}
```

Ответ

```json
{
    "status": "ok"
}
```



### [GET] /comments?id=<post_id> - Получить все комментарии к посту

Запрос

```json

```

Ответ

```json
{
    "comments": [
        {
            "id": 1,
            "content": "comment text",
            "ref": 1
        },
        ...
    ]
}
```



### [POST] /reset - Очистить spaces

Запрос

```json

```

Ответ

```json
{
    "status": "ok"
}
```



## Lua Cheatsheet

```lua
// Записать/восстановить логический дамп
local status, error = require('dump').dump('/tmp/dump')
local status, error = require('dump').restore('tmp/dump')

// Запустить из docker (для логических дампов)
docker run --name mytarantool -p3301:3301 -e TARANTOOL_USER_NAME=admin -e TARANTOOL_USER_PASSWORD=pass -v "$(pwd)"/dump:/tmp/dump -d tarantool/tarantool

// Создать бинарный дамп
box.backup.start()
box.backup.stop()

// Запустить из docker (для бинарныхх дампов)
docker run --name mytarantool -p3301:3301 -e TARANTOOL_USER_NAME=admin -e TARANTOOL_USER_PASSWORD=pass -v "$(pwd)"/shared:/var/lib/tarantool/ -d tarantool/tarantool

// Подключиться к doсker
docker exec -t -i mytarantool console
    
// Выполнить SQL из Lua
box.execute([[SELECT * FROM my_table;]]);
    
// Создать vinyl space
post = box.schema.space.create('post', {engine= 'vinyl'})

// Работа с сущностями
box.space.post:format({
         {name = 'id', type = 'unsigned'},
         {name = 'band_name', type = 'string'},
         {name = 'year', type = 'unsigned'}
         })
box.space.post:create_index('primary', {
         type = 'tree',
         parts = {'id'}
         })
box.space.post:auto_increment{"post text"}
box.space.post:insert{0, "post text"}
```