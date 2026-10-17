PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'student',
    password_hash BLOB NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS diary_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATETIME NOT NULL,
    subject TEXT NOT NULL,
    lesson_topic TEXT NOT NULL DEFAULT '',
    is_done BOOLEAN NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_diary_entries_date
    ON diary_entries (date);
CREATE INDEX IF NOT EXISTS idx_diary_entries_user_date
    ON diary_entries (user_id, date);

-- Ten test users. The password is intentionally a dummy value for local tests.
INSERT OR IGNORE INTO users (id, username, email, role, password_hash, created_at)
VALUES
    (1, 'ivan', 'ivan@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (2, 'anna', 'anna@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (3, 'petr', 'petr@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (4, 'maria', 'maria@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (5, 'olga', 'olga@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (6, 'nikita', 'nikita@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (7, 'sofia', 'sofia@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (8, 'maxim', 'maxim@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (9, 'elena', 'elena@example.test', 'student', 'test-hash', '2026-10-01 08:00:00'),
    (10, 'danil', 'danil@example.test', 'student', 'test-hash', '2026-10-01 08:00:00');

-- Exactly 100 diary entries distributed across October 2026.
-- Every user gets 10 entries; dates cover the whole month.
WITH RECURSIVE entries(n) AS (
    SELECT 1
    UNION ALL
    SELECT n + 1 FROM entries WHERE n < 100
), subjects(subject) AS (
    SELECT 'Математика' UNION ALL
    SELECT 'История' UNION ALL
    SELECT 'Физика' UNION ALL
    SELECT 'Информатика' UNION ALL
    SELECT 'Английский'
)
INSERT OR IGNORE INTO diary_entries
    (id, user_id, date, subject, lesson_topic, is_done, notes, created_at, updated_at)
SELECT
    20000 + n,
    ((n - 1) % 10) + 1,
    printf('2026-10-%02d %02d:00:00', ((n - 1) % 31) + 1, 8 + (n % 4)),
    CASE ((n - 1) % 5)
        WHEN 0 THEN 'Математика'
        WHEN 1 THEN 'История'
        WHEN 2 THEN 'Физика'
        WHEN 3 THEN 'Информатика'
        ELSE 'Английский'
    END,
    printf('Тема занятия %02d', n),
    CASE WHEN n % 3 = 0 THEN 1 ELSE 0 END,
    CASE WHEN n % 5 = 0 THEN 'Повторить материал' ELSE '' END,
    '2026-09-30 18:00:00',
    printf('2026-10-%02d 15:00:00', ((n - 1) % 31) + 1)
FROM entries;
