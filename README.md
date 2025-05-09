# recetas-golang-fullstack

## Postman ##
Prueba en Postman: Ahora puedes probar las rutas CRUD para categorías:
GET http://localhost:8080/api/v1/categorias
GET http://localhost:8080/api/v1/categorias/{id} (reemplaza {id} con un ID válido)
POST http://localhost:8080/api/v1/categorias (con body {"nombre": "Nueva Cat"})
PUT http://localhost:8080/api/v1/categorias/{id} (con body {"nombre": "Nombre Actualizado"})
DELETE http://localhost:8080/api/v1/categorias/{id}

Get-Item .\scripts\mysql\init-test-db-complete.sql
Remove-Item .\scripts\mysql\init-test-db-complete.sql
Create-Item .\scripts\mysql\init-test-db-complete.sql
# Crear archivos
New-Item .\scripts\mysql\init-test-db-complete.sql -ItemType File
New-Item .\scripts\mysql\init-test-db.sql -ItemType File

[Mysql]
mysqldump -u root -p --no-data golang-recetas > golang-recetas-schema.sql

# TEST/PowerShell ##
$env:APP_ENV="test"; go test -v ./internal/repository/mysql/... ; Remove-Item Env:\APP_ENV

# UP
docker-compose -f docker-compose.test.yml up -d db_test
# DOWN
docker-compose -f docker-compose.test.yml down -v
# RECREATE
docker-compose -f docker-compose.test.yml up --build --force-recreate

go test -v ./categorias/...

