CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    name TEXT,
    surname TEXT,
    email TEXT NOT NULL UNIQUE,
    departament TEXT,
    passhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name TEXT,
    second_name TEXT,
    surname TEXT,
    age INTEGER,
    sex boolean,
    phone TEXT UNIQUE,
    email TEXT NOT NULL UNIQUE,
    departament TEXT,
    post TEXT,
    passhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    day_of_week TEXT NOT NULL,
    hour TEXT NOT NULL,
    minute TEXT NOT NULL,
    video_id iNTEGER NOT NULL,
    employee_id INTEGER NOT NULL
);