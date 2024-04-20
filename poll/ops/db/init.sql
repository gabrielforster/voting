CREATE DATABASE IF NOT EXISTS poll_db;

USE poll_db;

CREATE TABLE IF NOT EXISTS `polls` (
    id INT AUTO INCREMENT,
    title VARCHAR(255) NOT NULL,
    hash CHAR(32) NOT NULL,
    decription TEXT,
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT now(),
    updated_at DATETIME DEFAULT now() ON UPDATE now(),

    PRIMARY KEY (`id`),
    INDEX created_by_index (created_by)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `options` (
    id INT AUTO INCREMENT,
    title VARCHAR(255) NOT NULL,
    pool_id INT NOT NULL,
    created_at DATETIME DEFAULT now(),

    PRIMARY KEY (`id`),
    FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS votes (
    id INT AUTO INCREMENT,
    user_id INT NOT NULL,
    pool_id INT NOT NULL,
    option_id INT NOT NULL,
    created_at DATETIME DEFAULT now(),

    PRIMARY KEY (`id`),
    FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES `options`(id) ON DELETE CASCADE,

    INDEX pool_option_index (poll_id, option_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
