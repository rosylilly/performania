CREATE DATABASE IF NOT EXISTS `performania`;
CREATE USER performania IDENTIFIED BY 'performania';
GRANT ALL PRIVILEGES ON performania.* TO 'performania'@'%';

SET PERSIST local_infile=1;
