CREATE DATABASE IF NOT EXISTS homeworkdb;

CREATE USER 'homework'@'%' IDENTIFIED BY 'homework';
GRANT ALL PRIVILEGES ON homeworkdb.* TO 'homework'@'%';
FLUSH PRIVILEGES;


