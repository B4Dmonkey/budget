Database driver docs:
https://pkg.go.dev/github.com/mattn/go-sqlite3

Debugging
https://github.com/golang/vscode-go/blob/master/docs/debugging.md

vs-code go ext
https://github.com/golang/vscode-go/blob/master/docs/settings.md

# Todo
- For create document Meta, I sent a null field on create. I need a way to enforce that I don't do this. If a field is to be null then it shouldn't be listed in the sqlc query. (i could also make the field non null)