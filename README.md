# SnippetBox using SQLite

This repository is the web application **SnippetBox**, created while working through the book ***Let's Go***,
by *Alex Edwards*.

While the book uses MySQL as its database, this repository uses SQLite. This results in some added code for converting SQLite's datetime strings to Go's time.Time format.

For the session management section of the book, I also use SQLite instead of MySQL. This means the SQLite package for the session store manager, sqlite3store,
was used. Using this package, I was able to follow the book's implementaion of sessions seamlessly.