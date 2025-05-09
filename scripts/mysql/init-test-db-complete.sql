-- /scripts/mysql/init-test-db-complete.sql
-- Script ÚNICO y Completo para inicializar la base de datos de prueba

-- 1. Crear Base de Datos (si no existe) y seleccionarla
CREATE DATABASE IF NOT EXISTS `recetas_test_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `recetas_test_db`;

-- --------------------------------------------------------
--      SECCIÓN DE CREACIÓN DE TABLAS
--      ¡¡¡ REEMPLAZA ESTO CON TU SCHEMA REAL !!!
-- --------------------------------------------------------

-- Tabla: categorias
CREATE TABLE IF NOT EXISTS `categorias` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_categorias_slug` (`slug`),     -- Nombre de constraint único sugerido
  UNIQUE KEY `uk_categorias_nombre` (`nombre`) -- Nombre de constraint único sugerido
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL, -- Hash de la contraseña
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: recetas
CREATE TABLE IF NOT EXISTS `recetas` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `categoria_id` int unsigned DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL, -- Asumiendo relación con users
  `nombre` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(180) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tiempo_preparacion` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `descripcion` text COLLATE utf8mb4_unicode_ci,
  `foto` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_recetas_slug` (`slug`),
  KEY `idx_recetas_categoria` (`categoria_id`), -- Nombre de índice sugerido
  KEY `idx_recetas_user` (`user_id`),       -- Nombre de índice sugerido
  CONSTRAINT `fk_recetas_categorias` FOREIGN KEY (`categoria_id`) REFERENCES `categorias` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_recetas_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE -- O ON DELETE CASCADE si prefieres borrar recetas si se borra el usuario
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: ingredientes (Ejemplo, ajusta según tu diseño)
CREATE TABLE IF NOT EXISTS `ingredientes` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
     `created_at` datetime(3) DEFAULT NULL,
     `updated_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_ingredientes_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: receta_ingredientes (Tabla de unión Muchos-a-Muchos, ejemplo)
CREATE TABLE IF NOT EXISTS `receta_ingredientes` (
    `receta_id` int unsigned NOT NULL,
    `ingrediente_id` int unsigned NOT NULL,
    `cantidad` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- Ej: "2 tazas", "100 gr"
    PRIMARY KEY (`receta_id`, `ingrediente_id`), -- Clave primaria compuesta
    KEY `idx_recetaing_ingrediente` (`ingrediente_id`),
    CONSTRAINT `fk_recetaing_recetas` FOREIGN KEY (`receta_id`) REFERENCES `recetas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE, -- Si borro receta, borro la relación
    CONSTRAINT `fk_recetaing_ingredientes` FOREIGN KEY (`ingrediente_id`) REFERENCES `ingredientes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE -- Si borro ingrediente, borro la relación
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ... (AÑADE AQUÍ CUALQUIER OTRA TABLA QUE NECESITES) ...


-- --------------------------------------------------------
--      SECCIÓN DE USUARIO Y PERMISOS PARA TESTS
-- --------------------------------------------------------

-- Crear el usuario de prueba si no existe, permitiendo conexión desde cualquier host
-- Asegúrate de que la contraseña coincida con tu config.test.yaml y docker-compose.test.yml
CREATE USER IF NOT EXISTS 'testuser'@'%' IDENTIFIED BY 'testpassword';

-- Otorgar todos los privilegios necesarios en la base de datos de test a ese usuario
GRANT ALL PRIVILEGES ON `recetas_test_db`.* TO 'testuser'@'%';

-- Aplicar los cambios de privilegios inmediatamente
FLUSH PRIVILEGES;

-- Mensaje final (opcional, se ve en logs de MySQL si falla algo después)
SELECT 'Script init-test-db-complete.sql ejecutado.' AS status;