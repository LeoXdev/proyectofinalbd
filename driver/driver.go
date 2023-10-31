// Package driver se encarga de las operaciones a la base de datos, ya sea
// realizar la conexion a la base de datos o ejecutar queries en caso
// de que algun handler lo requiera.
package driver

import (
	"database/sql"
	"log"
	"errors"
	"proyectofinalbd/model"

	"github.com/sijms/go-ora/v2"
)

var (
	// dbConn guardará un pointer a la conexión de la base de datos.
	// Empezará como null hasta la primera y única llamada a ConexionBaseDatos().
	dbConn *sql.DB
)

// ConexionBaseDatos se encarga de realizar la conexion a la base de datos para posteriormente
// poder realizar las operaciones CRUD.
// La referencia a esta conexión será guardada por la variable de paquete dbConn, que será accesada
// por los métodos que requieran realizar una query: 'dbConn.Query("query")'.
func ConexionBaseDatos() {
	connStr := go_ora.BuildUrl("localhost", 1521, "XEPDB1", "proyectofinalbd", "123456", nil)
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		log.Panicf("error al tratar de conectarse a la base de datos: %v", err)
	}
	err = conn.Ping() // Si Ping() no retorna error, la conexión a la base de datos está garantizada
	if err != nil {
		log.Panicf("error al tratar de mantener encendida base de datos: %v", err)
	}
	log.Printf("La conexión a la base de datos fue realizada exitosamente.")
	dbConn = conn

	// La conexion a la base de datos no se cierra, bajo palabras de la documentación oficial
	// del paquete sql:
	// It is rarely necessary to close a DB.
}

// ObtenerRegistrosAlumnos retorna un struct conteniendo todas las entradas de la tabla 'alumnos',
// será ejecutado por el handler Read.
// Se limita el número de entradas en 50 para evitar una carga excesiva a la plantilla HTML,
// aún se puede buscar un alumno en específico usando la opción de READ dentro de la página web.
func ObtenerRegistrosAlumnos() (model.Alumnos, error) {
	rows, err := dbConn.Query(`SELECT * FROM alumnos WHERE ROWNUM <= 50`)
	if err != nil {
		return model.Alumnos{}, err // Si hubo un error al ejecutar la query se retorna un slice vacio
	}
	defer rows.Close()
	var (
		// Sobre estas variables se guardara el valor de cada registro

		matricula string
		nombre    string
		edad      sql.NullInt16
	)
	// Se guardan en memoria todos los registros obtenidos de la base de datos
	registros := model.Alumnos{}
	for rows.Next() {
		err = rows.Scan(&matricula, &nombre, &edad)
		if err != nil {
			return model.Alumnos{}, err // Si hubo un error al ejecutar la query se retorna un slice vacio
		}

		// Se rellena con 0s las entradas de matricula que podrían no tener 8 dígitos, esto para
		// tener un mejor formato en la renderización de las entradas.
		for len(matricula) != 8 {
			matricula = "0" + matricula
		}
		// Se valida el unico campo que podría llegar nulo desde la base de datos,
		// si es nulo, se le asigna el valor cero.
		if !edad.Valid {
			edad.Int16 = 0
		}
		registros.Alumnos = append(registros.Alumnos, model.Alumno{
			Matricula: matricula,
			Nombre:    nombre,
			Edad:      edad.Int16,
		})
	}

	return registros, nil
}

// CrearRegistro recibe un objeto model.Alumno creado por el handler Create y
// crea un nuevo registro en la base de datos.
func CrearRegistro(nuevoAlumno model.Alumno) error {
	stmt, err := dbConn.Prepare(`INSERT INTO alumnos (matricula, nombre, edad) VALUES (:1, :2, :3)`)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }() // Cierra el statement (~COMMIT...)

	_, err = stmt.Exec(nuevoAlumno.Matricula, nuevoAlumno.Nombre, nuevoAlumno.Edad)
	if err != nil {
		return err
	}

	return nil
}

// ObtenerRegistro recibe una matricula y determina si existe una entrada con dicha matricula,
// si es así, retorna la entrada Alumno identificable por la matricula pasada.
func ObtenerRegistro(matricula string) (model.Alumno, error) {
	removerCerosMatricula(&matricula)

	alumnoBuscado := model.Alumno{}
	// QueryRow difiere el error hasta que se escanea el objeto que retorna
	row := dbConn.QueryRow(`SELECT * FROM alumnos WHERE matricula = :1`, matricula)
	// Concatenar o formatear introduce riesgo de SQL-Injection: https://go.dev/doc/database/sql-injection
	//row = dbConn.QueryRow(`SELECT * FROM alumnos WHERE matricula = `+ matricula)
	err := row.Scan(&alumnoBuscado.Matricula, &alumnoBuscado.Nombre, &alumnoBuscado.Edad)
	if err != nil {
		return alumnoBuscado, err // Si hubo un error al ejecutar la query se retorna un slice vacio
	}
	
	return alumnoBuscado, nil
}

// ActualizarRegistro recibe un objeto model.Alumno y ejecuta una query para tratar de
// actualizar el registro que contenga la matricula del objeto recibido.
// Se retorna error si no se actualizó ningún registro.
func ActualizarRegistro(nuevosDatosAlumno model.Alumno) error {
	stmt, err := dbConn.Prepare(`UPDATE alumnos SET nombre = :2, edad = :3 WHERE matricula = :1`)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }() // Cierra el statement (~COMMIT...)

	result, err := stmt.Exec(nuevosDatosAlumno.Nombre, nuevosDatosAlumno.Edad, nuevosDatosAlumno.Matricula)
	if err != nil {
		return err
	}
	ra, _ := result.RowsAffected()	// Si la query update no afectó (actualizó) ningún registro
	if ra == 0 {
		return errors.New("no se actualizó ningún registro, revisar que la matricula sea correcta")
	}

	return nil
}

// EliminarRegistro recibe un objeto model.Alumno y ejecuta una query para tratar de
// eliminar el registro que contenga la matricula del objeto recibido.
// Se retorna error si no se eliminó ningún registro.
func EliminarRegistro(alumnoAEliminar model.Alumno) error {
	stmt, err := dbConn.Prepare(`DELETE FROM alumnos WHERE matricula = :1`)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }() // Cierra el statement (~COMMIT...)

	result, err := stmt.Exec(alumnoAEliminar.Matricula)
	if err != nil {
		return err
	}
	ra, _ := result.RowsAffected()	// Si la query update no afectó (actualizó) ningún registro
	if ra == 0 {
		return errors.New("no se eliminó ningún registro, revisar que la matricula sea correcta")
	}

	return nil
}

// removerCerosMatricula remueve los primeros n ceros de una matricula recibida.
// Retorna true si se ha modificado la matricula pasada o false en caso contrario.
func removerCerosMatricula(matricula *string) bool {
	// Si la matricula no inicia con 0, no hay necesidad de remover 0s redundantes.
	if string([]rune(*matricula)[0]) != "0" {
		return false
	}
	// Se recorre de izquierda a derecha y se identifica el indice del primer caracter que no sea 0
	var nonZeroIndex int
	for i := 0; i < len(*matricula); i++ {
		// Este identificador posiblemente contenga un caracter que ocupa más de un byte, como un emoji
		caracter := string([]rune(*matricula)[i])
		if caracter != "0" {
			nonZeroIndex = i
		}
	}
	
	// Se quitan caracteres, que supuestamente son 0s
	*matricula = (*matricula) [nonZeroIndex:]
	return true
}
