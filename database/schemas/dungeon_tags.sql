PRAGMA foreign_keys=ON;

DROP TABLE IF EXISTS `tag_taxonomies`;
CREATE TABLE IF NOT EXISTS `tag_taxonomies` (
    `uuid` TEXT PRIMARY KEY,
    `name` TEXT NOT NULL,
    `internal` INTEGER NOT NULL,
    `cluster_domain` TEXT NOT NULL
);

DROP TABLE IF EXISTS `dungeon_tags`;
CREATE TABLE IF NOT EXISTS `dungeon_tags` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` TEXT NOT NULL,
    `taxonomy` TEXT NOT NULL,
    `name_taxonomy` TEXT NOT NULL UNIQUE,
    FOREIGN KEY(`taxonomy`) REFERENCES `tag_taxonomies`(`uuid`) ON DELETE CASCADE
);

DROP TABLE IF EXISTS `taggings`;
CREATE TABLE IF NOT EXISTS `taggings` (
    `tagging_id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `tag` INTEGER NOT NULL,
    `taggable_id` TEXT NOT NULL,
    FOREIGN KEY(`tag`) REFERENCES `dungeon_tags`(`id`) ON DELETE CASCADE,
    UNIQUE(`tag`, `taggable_id`)
);