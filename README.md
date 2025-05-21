# recetas-golang-fullstack
// go:build ignore
// +build ignore

# para correr el api
go run ./cmd/app/main.go

# Limpiado de dependencias
go mod tidy

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

[Lint]
# para correr en root/
golangci-lint run
# para correr en backend/para detectar problemas estáticos.
golangci-lint run ./... 

[Format]
# para correr en root/
gofmt -d .
#para correr en backend/
gofmt -d ./...

[Mysql]
mysqldump -u root -p --no-data golang-recetas > golang-recetas-schema.sql

# TEST/PowerShell ##
#(corriendo desde backend/ si el yml está en la raíz del proyecto)
docker-compose -f ../docker-compose.test.yml up -d db_test

docker-compose -f docker-compose.test.yml down -v

#(corriendo desde root/)
$env:APP_ENV="test"; go test -v ./internal/repository/mysql/... ; Remove-Item Env:\APP_ENV
# EN POWERSHELL:
$env:APP_ENV="test"; go test -v ./recetas/... ; Remove-Item Env:\APP_ENV
$env:APP_ENV="test"; go test -v ./contactos/... ; Remove-Item Env:\APP_ENV


# UP
docker-compose -f docker-compose.test.yml up -d db_test
# DOWN
docker-compose -f docker-compose.test.yml down -v
# RECREATE
docker-compose -f docker-compose.test.yml up --build --force-recreate

go test -v ./categorias/...
go test -v ./recetas/...
go test -v ./contactos/...
go test -v ./seguridad/...
go test -v ./usuarios/...
go test -v ./shared/...

#Para correr TODOS los tests de integración que encuentren (categorías y recetas)
go test -v ./...


#para correr solo la suite de tests unitarios del servicio.
go test -v -run TestCategoriaServiceTestSuite ./categorias/... 
# para correr solo la suite de tests de integración del repositorio.
go test -v -run TestCategoriaRepositoryIntegrationTestSuite ./categorias/... 

bfb5acae6b29ef2caa261d48f533f26fff57746f

[Seeder]
# para correr el seeder.
go run cmd/seeder/main.go

[Docs]
1. godox : encuentra partes del código no documentado y genera un reporte
go install github.com/nikolaydubina/godox@latest
godox ./...

2. godoxgen: Genera plantillas de comentarios automáticamente
go install github.com/nikolaydubina/godoxgen@latest
godoxgen -dir=./backend

3. staticcheck: Usa staticcheck para auditar documentación
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
staticcheck -checks=SA
staticcheck ./... > godoc-report.txt


## Design-First (o API-First): ##
oapi-codegen

# mas usado, completo pero rquiere más onfiguración al iniciar
#go-swagger 

# Idiomático en Go: Usar oapi-codegen es muy común. Generas los tipos y las interfaces de los handlers, y luego implementas la lógica de # # # # negocio en esos handlers.
# # Generar tipos y stubs de servidor para chi router
# oapi-codegen -generate types,chi-server -package api_generated -o api/generated.go openapi.yaml

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

## Code-First (Generar OpenAPI desde el Código Go): ##
librerias para instrpeccioanr el código
swag 

go install github.com/swaggo/swag/cmd/swag@latest

go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files // Necesario para servir los archivos de Swagger UI

# para generar la documentación
swag init -g ./cmd/app/main.go --output ./docs --parseDependency --parseInternal -d ./,./shared,./categorias,./recetas
# ruta de documentación
http://localhost:8080/swagger/index.html#/

TO DO
=====
Endpoint de Health Check (si lo tuviéramos):
validar funciones de fecha: mimetype
validar el nombre:(vacio) de la receta desde le postform 
eliminar caregoria(validar que exista/ cascada) en recetas 
Por ahora, lo he simplificado a c.Error()

[Módulo/Utils]
- seguridad.go: libreria para generar token/google-uuid y hash (bcript), costo, se usa para campo token de usuario
- validaciones.go: regexp(validaciones de email)
- 

[Módulo/Recetas]
Api:v1/recetas
Api:v1/recetas/{id}
Api:v1/recetas/{id}/ingredientes
Api:v1/recetas/{id}/ingredientes/{ingrediente_id}
Api:v1/recetas/{id}/comentarios
Api:v1/recetas/{id}/comentarios/{comentario_id}

recetas_helper.go
rutas_helper.go
-Receta_Helper_Usuario(): funcion para listar todas las recetas credas por un susuario 
-Receta_helper_home(): funcion para listar las recetas mas vistas
-Receta_helper_busqueda(): funcion para listar las recetas por nombre
-Receta_helper_ordenamiento(): funcion para listar las recetas por ordenamiento
lisr receta por slug : mostar el detalle de la receta
- usuario que creo la receta: usuario_id (relacion con usuario) fk usuario
-busqueda de recetas: combobox
-ordenamiento de recetas
- Editar foto de receta

[Módulo/Contactanos]
dtos: requeridos, formato de email, si espacios en blanco
handlers
models
repositories
services
validators
crud:fechas

[Requisitos Funcionales]
- al registrar los datos se debe enviar un correo al admin 
- enviar un mensaje de confirmación al usuario
- Crear un utilitario para usar mailtrap: mailtrap.io: SMTP

[Requisitos No Funcionales]

[Modulo/Seguridad]
# Desde backend/
go get golang.org/x/crypto/bcrypt

Api:v1/seguridad/verificar-token
Api:v1/seguridad/registro
Api:v1/seguridad/login
Api:v1/seguridad/confirmar
Api:v1/seguridad/olvidar
Api:v1/seguridad/cambiar
Api:v1/seguridad/actualizar
Api:v1/seguridad/eliminar
Api:v1/seguridad/validar
validacion de roles
validacion de permisos
autenticacion
verificar
dtos: login, registro, confirmar, olvidar, cambiar, actualizar, eliminar, validar
-Valida: correo, contraseña, usuario, token
pkg: jwt
handlers
models:
- estado
- usuario
repositories
services
validators
crud

POST:seguridad/registro
GET:seguridad/verificacion

pkg: seguridad
Proteger rutas
custom middlaware: verificar-token

[Módulo/Usuarios]
Api:v1/usuarios/registro
Api:v1/usuarios/login
Api:v1/usuarios/confirmar
Api:v1/usuarios/olvidar
Api:v1/usuarios/cambiar
Api:v1/usuarios/actualizar
Api:v1/usuarios/eliminar
Ruta_FRONTEND: http://localhost:5173/usuarios

dtos: campos obligatorios, formato de email, si espacios en blanco, validaciones de contraseña
handlers
models: relacion con metadata datos personales(tlf, nombre, direccion, etc) roles, estado, token
repositories
services: validaciones de contraseña, encriptacion de contraseña, generacion de token, validacion de token, envio de correo, eliminacion de usuario, actualizacion de usuario, registro de usuario, login de usuario, confirmacion de usuario, olvido de contraseña, cambio de contraseña
validators
crud
utils: seguridad.go: libreria para generar token y hash (bcript), costo,


[CORS]

configurar cors / middleware


## Requistos Fcunonales ##
- Validar que el correo no se repita anets de guardar
- Validar que la contraseña tenga al menos 8 caracteres
- Validar que el usuario tenga al menos 8 caracteres
- Validar que el usuario tenga al menos 8 caracteres


[Vue/Compositon API]

