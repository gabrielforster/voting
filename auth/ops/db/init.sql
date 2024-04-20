CREATE DATABASE IF NOT EXISTS auth_db;

USE auth_db;

CREATE TABLE IF NOT EXISTS `users` (
    id INT PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at DATETIME,
    updated_at DATETIME
)
ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;


INSERT INTO `users` (id, email, password, first_name, last_name, created_at, updated_at)
SELECT 1,'rochafrgabriel@gmail.com', SHA1('1234'), 'Gabriel', 'Rocha', now(), null FROM DUAL
WHERE NOT EXISTS
    (SELECT email FROM `users` WHERE email='rochafrgabriel@gmail.com');
