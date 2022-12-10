CREATE DATABASE IF NOT EXISTS `audio_db`;

use audio_db;

DROP TABLE IF EXISTS wavs;
CREATE TABLE wavs (
    id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64),
    length_seconds INTEGER,
    file_url VARCHAR(64)
);


