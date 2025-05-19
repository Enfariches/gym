CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) DEFAULT '',
    surname VARCHAR(64) DEFAULT '',
    email VARCHAR NOT NULL UNIQUE,
    department_id INTEGER NOT NULL,
    passhash TEXT NOT NULL,
    CONSTRAINT department_id_fk FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) DEFAULT '',
    second_name VARCHAR(64) DEFAULT '',
    surname VARCHAR(64) DEFAULT '',
    age INTEGER DEFAULT 18,
    sex BOOLEAN DEFAULT FALSE,
    phone VARCHAR(20) DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    department_id INTEGER DEFAULT 1,
    post VARCHAR DEFAULT '',
    passhash TEXT NOT NULL,
    CONSTRAINT department_id_fk FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE TABLE IF NOT EXISTS schedules (
    id SERIAL PRIMARY KEY,
    cron_expression VARCHAR(11) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    media_id INTEGER NOT NULL,
    admin_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT admin_id_fk FOREIGN KEY (admin_id) REFERENCES admins(id)
);

CREATE TABLE IF NOT EXISTS mediafiles (
    id SERIAL PRIMARY KEY,
    title TEXT,
    admin_id INTEGER NOT NULL,
    department_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT admin_id_fk FOREIGN KEY (admin_id) REFERENCES admins(id)
);

CREATE TABLE IF NOT EXISTS statistics (
    id SERIAL PRIMARY KEY,
    progress VARCHAR(50),
    percentage_view INTEGER,
    employee_id INTEGER NOT NULL,
    department_id INTEGER NOT NULL,
    media_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT employee_id_fk FOREIGN KEY (employee_id) REFERENCES employees(id),
    CONSTRAINT deparment_id_fk FOREIGN KEY (department_id) REFERENCES departments(id)
)