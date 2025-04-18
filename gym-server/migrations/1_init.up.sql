CREATE TABLE IF NOT EXISTS admins_new (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    departament TEXT NOT NULL,
    passhash TEXT NOT NULL
);
