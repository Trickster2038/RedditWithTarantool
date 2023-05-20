# RedditWithTarantool

```lua
// Создать vinyl space
post = box.schema.space.create('post', engine: "vinyl")

// Создать пост
box.space.post:auto_increment{"post text"}
box.space.post:insert{0, "post text"}

// Записать логический дамп
local status, error = require('dump').dump('/tmp/dump')

// Выполнить SQL из Lua
box.execute([[SELECT * FROM my_table;]]);

// ===== Spaces =====

/*
	post
	______________________
	id			unsigned 
	content 	string
*/

/*
	comment
	______________________
	id			unsigned 
	content 	string
	ref			unsigned 
*/
```

