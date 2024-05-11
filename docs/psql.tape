Output docs/psql.gif

Require psql

Set Width 1200
Set Height 650
Set Padding 50
Set TypingSpeed 30ms

# psql cli
Type `psql postgres://rw:password@127.0.0.1:5432/appdb`
Sleep 500ms
Enter

# set tenant config
Type `-- Set 'user-1' as the current session user`
Sleep 500ms
Enter
Type `SELECT set_config('app.current_user', 'user-1', false);`
Sleep 500ms
Enter
Sleep 500ms

# query
Type `-- Query data`
Sleep 500ms
Enter
Type `SELECT id, title, description, completed, user_slug FROM todos;`
Sleep 500ms
Enter
Sleep 5s
