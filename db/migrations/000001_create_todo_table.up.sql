CREATE TABLE `todo` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title` TEXT NOT NULL,
    `done` INTEGER NOT NULL,
    `created` DATE DEFAULT (datetime('now'))
);
