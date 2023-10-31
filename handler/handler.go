// Package handler se encarga de inicializar el servidor web donde se invocan las
// operaciones CRUD dirigidas a la base de datos, ademas de contener los handlers que
// se encargan de las operaciones CRUD.
// Un handler es (en pocas palabras), una funcion que se ejecuta al acceder a cierta URL
// dentro del servidor.
package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"proyectofinalbd/driver"
	"proyectofinalbd/model"
)

// CrearServidor crea el servidor web que contiene los handlers que llaman a las operaciones CRUD.
// El usuario lo debe usar accediendo a la página principal, las peticiones web para hacer CRUD
// de registros seran llamadas por la GUI en la página principal, sin necesidad de meterse a URLs.
func CrearServidor() {
	sm := http.DefaultServeMux

	sm.HandleFunc("/", Crud)	//	Single Page Application CRUD

	// FileServer funciona como un handler que sirve archivos estáticos (servira CSS y JS al HTML),
	// el argumento http.Dir representa un directorio presente del lado del servidor.
	fs := http.FileServer(http.Dir("handler/template/static"))
	// Se recorta static/ del request hecho por la plantilla html para acceder al archivo style.css
	// que fue registrado por fs.
    http.Handle("/static/", http.StripPrefix("/static/", fs))

	sm.HandleFunc("/alumnos/create", Create)	//	endpoint para peticiones POST
	sm.HandleFunc("/alumnos/read", Read)		//	endpoint para peticiones GET (usando OPTIONS)
	sm.HandleFunc("/alumnos/update", Update)	//	endpoint para peticiones PUT
	sm.HandleFunc("/alumnos/delete", Delete)	//	endpoint para peticiones DELETE

	log.Printf("La aplicacion se encuentra corriendo en: 'http://localhost:8080/'")

	err := http.ListenAndServe(":8080", sm)
	if err != nil {
		// No se retorna error para no ensuciar main
		log.Panicf("error al tratar de iniciar el servidor web: %v", err)
	}
}

// Crud es un handler que renderiza la plantilla HTML que contiene los elementos que provocan
// las operaciones CRUD mediante peticiones web.
// No retorna error para no modificar la firma necesaria para ser handler, en su lugar
// causa panic.
func Crud(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("handler/template/crud.html")
	if err != nil {
		log.Panicf("error al escanear los archivos: %v", err)
	}

	// Se cargan en memoria los datos de la base de datos.
	readData, err := driver.ObtenerRegistrosAlumnos()
	if err != nil {
		log.Panicf("error al obtener los registros de alumnos: %v", err)
	}

	// La plantilla será escrita en w (ya con los valores recibidos de la base de datos)
	err = tpl.Execute(w, readData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
// Create es el handler que crea un nuevo registro en la base de datos, se ejecuta al recibir
// una peticion POST mandada por JavaScript desde un botón en el HTML.
func Create(w http.ResponseWriter, r *http.Request) {
	log.Printf("Se ha recibido una petición %v para Create.", r.Method)

	// Redireccionar en caso de que la petición no sea POST, se retorna para evitar un
	// error superfluous response.WriteHeader call.
	if redireccionMetodoInvalido(w, r, http.MethodPost, "http://localhost:8080/alumnos/create") {
		return
	}
	
	var nuevoAlumno model.Alumno
	err := json.NewDecoder(r.Body).Decode(&nuevoAlumno)
	if err != nil {
		logResponse(w, "error al tratar de leer el nuevo objeto alumno (JSON):", err, http.StatusBadRequest)
		return
	}
	
	err = driver.CrearRegistro(nuevoAlumno)
	if err != nil {
		logResponse(w, "error al crear el nuevo registro alumno:", err, http.StatusInternalServerError)
		return
	}

	logResponse(w, "Se ha procesado la petición POST exitosamente.", nil, http.StatusOK)
}

// Read es el handler que se encarga de realizar peticiones de búsqueda (de un único alumno)
// Se procesa con el método http OPTIONS debido a que el servidor aún no cuenta con lógica
// de autorización para bloquear el paso del cliente a la url /alumnos/read.
func Read(w http.ResponseWriter, r *http.Request) {
	log.Printf("Se ha recibido una petición %v para Read.", r.Method)

	// Redireccionar en caso de que la petición no sea OPTIONS, se retorna para evitar un
	// error superfluous response.WriteHeader call.
	if redireccionMetodoInvalido(w, r, http.MethodOptions, "http://localhost:8080/alumnos/read") {
		return
	}

	// Se utiliza la misma clase de model.Alumno para buscar un registro por matricula para evitar
	// tener que definir otra struct que será utilizada una única vez.
	// El único campo no null sera el de matricula.
	var alumnoBuscado model.Alumno
	err := json.NewDecoder(r.Body).Decode(&alumnoBuscado)
	if err != nil {
		logResponse(w, "error al tratar de leer el objeto alumnoBuscado del request body (JSON):", err, http.StatusBadRequest)
		return
	}
	
	registroAlumno, err := driver.ObtenerRegistro(alumnoBuscado.Matricula)
	if err != nil {
		logResponse(w, "error al buscar el registro alumno:", err, http.StatusInternalServerError)
		return
	}
	
	// El frontend encontrará el registro dentro de los corchetes
	logResponse(w,
		"Se ha procesado la petición GET (OPTIONS) exitosamente: [" +
			registroAlumno.Matricula + "," +
			registroAlumno.Nombre + "," +
			strconv.Itoa(int(registroAlumno.Edad)) + "]", nil, http.StatusOK)
}

// Update es el handler usado para actualizar por completo un registro de un alumno ya
// existente.
func Update(w http.ResponseWriter, r *http.Request) {
	log.Printf("Se ha recibido una petición %v para Update.", r.Method)

	// Redireccionar en caso de que la petición no sea PUT, se retorna para evitar un
	// error superfluous response.WriteHeader call.
	if redireccionMetodoInvalido(w, r, http.MethodPut, "http://localhost:8080/alumnos/update") {
		return
	}

	var nuevosDatosAlumno model.Alumno
	err := json.NewDecoder(r.Body).Decode(&nuevosDatosAlumno)
	if err != nil {
		logResponse(w, "error al tratar de leer el objeto alumno a actualizar (JSON):", err, http.StatusBadRequest)
		return
	}

	err = driver.ActualizarRegistro(nuevosDatosAlumno)
	if err != nil {
		logResponse(w, "error al tratar de actualizar el registro alumno:", err, http.StatusInternalServerError)
		return
	}
	
	logResponse(w, "Se ha procesado la petición PUT exitosamente.", nil, http.StatusOK)
}

// Delete es el handler que elimina un registro de la base de datos, siempre y cuando exista.
func Delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("Se ha recibido una petición %v para Delete.", r.Method)

	// Redireccionar en caso de que la petición no sea DELETE, se retorna para evitar un
	// error superfluous response.WriteHeader call.
	if redireccionMetodoInvalido(w, r, http.MethodDelete, "http://localhost:8080/alumnos/delete") {
		return
	}

	var alumnoAEliminar model.Alumno
	err := json.NewDecoder(r.Body).Decode(&alumnoAEliminar)
	if err != nil {
		logResponse(w, "error al tratar de leer el objeto alumno a eliminar (JSON):", err, http.StatusBadRequest)
		return
	}

	err = driver.EliminarRegistro(alumnoAEliminar)
	if err != nil {
		logResponse(w, "error al tratar de eliminar el registro alumno:", err, http.StatusInternalServerError)
		return
	}
	
	logResponse(w, "Se ha procesado la petición DELETE exitosamente.", nil, http.StatusOK)
}

// logResponse es una función generalizada para logear respuestas tanto exitosas o erróneas
// según sea necesario hacia las consola de Go y HTML.
// Se usa cuando se encuentra un error en una petición mandada por el frontend o cuando se
// logra ejecutar la operación de la petición exitosamente.
func logResponse(w http.ResponseWriter, mensaje string, err error, statusCode int) {	
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	// Si el mensaje de respuesta es erróneo
	if err != nil {
		w.Write([]byte(mensaje + " " + err.Error()))
		log.Printf("%v %v", mensaje, err)
		return
	}
	// Si el mensaje de respuesta es exitoso, no se tiene que acceder a err
	w.Write([]byte(mensaje))
	log.Printf("%v", mensaje)
	return
}
// redireccion es una funcion que se utiliza para redireccionar al cliente a
// la página principal en caso de que se haga una petición con un método no permitido.
// Se retorna un bool verdadero en caso que se haya direccionado al usuario, esto es necesario
// para que el método llamante pueda retornar y cancelar el handler llamado con metodo incorrecto.
func redireccionMetodoInvalido(w http.ResponseWriter, r *http.Request, method, urlRequestada string) bool {
	// Redireccionar en caso de que la petición no sea la deseada
	if r.Method != method {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Printf("Se ha redirigido al usuario desde '%v' hacia 'http://localhost:8080/'", urlRequestada)
		return true
	}
	return false
}