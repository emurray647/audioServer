CREATE DATABASE IF NOT EXISTS `audio_db`;

use audio_db;

DROP TABLE IF EXISTS wavs;
CREATE TABLE wavs (
    id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    file_size INTEGER NOT NULL,
    duration FLOAT,
    num_channels INTEGER,
    sample_rate INTEGER,
    audio_format INTEGER,
    avg_bytes_per_second INTEGER,
    file_uri VARCHAR(64)
);

-- almost all lookups are on the name column
CREATE INDEX audio_name_index ON wavs(name);

