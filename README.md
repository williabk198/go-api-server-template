# Go API Server Template

This repo provides the basis of an API server with a dummy database connection.

## Usage

Please note that if you do use this as the basis for your API server, 
then you will need to replace all instances of `github.com/williabk198/go-api-server-template` 
to the name of your project/module, add an actual implementation of the `db.Database` interface
in `db/database.go` (look at the `dummydb` package as an example) 
and update this readme file to properly reflect your project

## Third Party Pacakges

By default, this template uses `go-chi/chi`, `go-chi/cors` and `google/uuid`.
These packages can be updated or removed to better fit your needs at any time. 
