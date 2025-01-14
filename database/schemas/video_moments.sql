PRAGMA foreign_keys=ON;

DROP TABLE IF EXISTS `videos`;
CREATE TABLE IF NOT EXISTS `videos` (
    `uuid` TEXT PRIMARY KEY,
    `cluster_uuid` TEXT NOT NULL,
);

DROP TABLE IF EXISTS `video_moments`;
CREATE TABLE IF NOT EXISTS `video_moments` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `video_uuid` TEXT NOT NULL,
    `moment_time` INTEGER NOT NULL,
    FOREIGN KEY(`video_uuid`) REFERENCES `videos`(`uuid`) ON DELETE CASCADE
);