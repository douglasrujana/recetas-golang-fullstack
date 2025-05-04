-- Este script se ejecuta automáticamente al iniciar el contenedor MySQL
-- cuando se monta en /docker-entrypoint-initdb.d/

-- Asegúrate de que se está usando la base de datos correcta
CREATE DATABASE IF NOT EXISTS recetas_test_db;
USE recetas_test_db;

-- Crear el usuario de prueba si no fue creado por las variables de entorno
CREATE USER IF NOT EXISTS 'testuser'@'%' IDENTIFIED BY 'testpassword';

-- Otorgar privilegios completos en la base de datos de test
GRANT ALL PRIVILEGES ON recetas_test_db.* TO 'testuser'@'%';

-- Aplicar los cambios de privilegios
FLUSH PRIVILEGES;

-- ✅ Puedes agregar aquí tablas mínimas si no usas migraciones automáticas:
-- CREATE TABLE IF NOT EXISTS categorias (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     nombre VARCHAR(100) NOT NULL
-- );

-- Puedes imprimir algo en logs de MySQL para depuración (no se ve en consola docker)
-- SELECT 'init-test-db.sql ejecutado correctamente';
