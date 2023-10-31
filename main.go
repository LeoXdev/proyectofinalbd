// Package main sirve como punto de entrada para la aplicación.
// Se conecta el lenguaje golang con una base de datos local Oracle para realizar
// operaciones CRUD (Create, Read, Update y Delete) sobre una tabla 'alumnos'.
//
// La base de datos debera estar ejecutándose para que se puedan realizar las
// conexiones de forma exitosa, para esto, revisar los servicios (WIN + R -> services.msc).
//
// Para provocar una operación se deberá hacer una petición mediante la GUI web,
// esta tiene botones que realizan peticiones web via JavaScript a las URL handler
// que ejecutan las operaciones CRUD a la base de datos.
// Si se intenta ingresar a las URL designadas para hacer una operacion CRUD, el usuario
// sera redirigido automáticamente de vuelta a la página principal:
// /alumnos/create -> /
//
// Información de la conexión:
//
// username: proyectofinalbd
//
// password: 123456
//
// service: XEPDB1
package main

import (
	"proyectofinalbd/driver"
	"proyectofinalbd/handler"
)

func main() {
	driver.ConexionBaseDatos()
	handler.CrearServidor()
}
