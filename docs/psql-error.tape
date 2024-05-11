Output docs/psql-error.gif

Require psql

Set Width 1200
Set Height 350
Set Padding 50
Set TypingSpeed 30ms

# psql cli
Type `psql postgres://rw:password@127.0.0.1:5432/appdb`
Sleep 500ms
Enter

# query
Type `-- Attempt to query data without setting the session user`
Sleep 500ms
Enter
Type `SELECT id, title, description, completed, user_slug FROM todos;`
Sleep 500ms
Enter
Sleep 5s
