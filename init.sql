CREATE DATABASE IF NOT EXISTS `audio_db`;

use audio_db;

DROP TABLE IF EXISTS wavs;
CREATE TABLE wavs (
    id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    file_size INTEGER NOT NULL,
    length_seconds FLOAT,
    num_channels INTEGER,
    sample_rate INTEGER,
    audio_format INTEGER,
    avg_bytes_per_sec INTEGER,
    file_uri VARCHAR(64)
);


