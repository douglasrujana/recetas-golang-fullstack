# recetas-golang-fullstack

## Postman ##
Prueba en Postman: Ahora puedes probar las rutas CRUD para categorías:
GET http://localhost:8080/api/v1/categorias
GET http://localhost:8080/api/v1/categorias/{id} (reemplaza {id} con un ID válido)
POST http://localhost:8080/api/v1/categorias (con body {"nombre": "Nueva Cat"})
PUT http://localhost:8080/api/v1/categorias/{id} (con body {"nombre": "Nombre Actualizado"})
DELETE http://localhost:8080/api/v1/categorias/{id}

# TEST
$env:APP_ENV="test"; go test -v ./internal/repository/mysql/... ; Remove-Item Env:\APP_ENV

# UP
docker-compose -f docker-compose.test.yml up -d db_test

# DOWN
docker-compose -f docker-compose.test.yml down -v

# RECREATE
docker-compose -f docker-compose.test.yml up --build --force-recreate


Get-Item .\scripts\mysql\golang-recetas-schema.sql
