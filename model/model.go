// Package model contiene entidad(es) sobre la(s) que se realizarán operaciones CRUD.
package model

// Alumno es el modelo sobre el que se hacen las operaciones CRUD.
//
// La tabla se describe de la siguiente forma:
//
// matricula	NOT NULL	NUMBER(8)
//
// nombre		NOT NULL	VARCHAR(20)
//
// edad						NUMBER(2)
type Alumno struct {
	Matricula string `json:"matricula"`
	Nombre    string `json:"nombre"`
	Edad      int16  `json:"edad"`
}

// Alumnos guarda en memoria registros de alumnos.
type Alumnos struct {
	Alumnos []Alumno
}
