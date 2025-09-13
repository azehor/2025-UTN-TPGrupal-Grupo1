# Estructura de /internal/*
- endpoint/
- endpoint/endpoint.go
- endpoint/model/model.go
- endpoint/store/store.go

## endpoint.go
Contiene las respuestas a los requests habilitados, asi como las operaciones necesarias para realizarlas.

## model.go
Contiene la estructura de datos que se maneja en el endpoint, asi como la estructura que se guardara en la base de datos

## store.go
Contiene las operaciones que interactuan con la base de datos para dar respuesta a las operaciones en endpoint.go

