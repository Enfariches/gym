CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    name TEXT DEFAULT '',
    surname TEXT DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    departament TEXT DEFAULT '',
    passhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name TEXT DEFAULT '',
    second_name TEXT DEFAULT '',
    surname TEXT DEFAULT '',
    age INTEGER DEFAULT 0,
    sex BOOLEAN DEFAULT FALSE,
    phone TEXT UNIQUE DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    departament TEXT DEFAULT '',
    post TEXT DEFAULT '',
    passhash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    cron_expression TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    video_id INTEGER NOT NULL,
    admin_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);