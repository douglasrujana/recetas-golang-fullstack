-- /scripts/mysql/init-test-db.sql

-- El usuario 'testuser' y la database 'recetas_test_db' DEBEN ser creados
-- por las variables de entorno MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE
-- antes de que este script se ejecute.

-- Este script asegura que 'testuser' pueda conectarse desde cualquier host ('%').
-- Usar CREATE USER IF NOT EXISTS y GRANT es más seguro que ALTER USER.

-- Crear el usuario si por alguna razón las ENV VARS no lo hicieron correctamente
-- o para asegurar el host '%'. Reemplaza la contraseña si ya existe.
CREATE USER IF NOT EXISTS 'testuser'@'%' IDENTIFIED BY 'testpassword';

-- Otorgar todos los privilegios necesarios en la base de datos de test
GRANT ALL PRIVILEGES ON `recetas_test_db`.* TO 'testuser'@'%';

-- Recargar privilegios para asegurar que los cambios se apliquen
FLUSH PRIVILEGES;

-- Opcional: Crear tablas si no usas migraciones automáticas en el test
-- USE `recetas_test_db`;
-- CREATE TABLE IF NOT EXISTS categorias ( ... );