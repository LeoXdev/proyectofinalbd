"use strict";
// -------------------- CREATE --------------------
// Elementos dentro del div create
const createButton = document.getElementById('create-submit')
const createMatricula = document.getElementById('create-matricula')
const createNombre = document.getElementById('create-nombre')
const createEdad = document.getElementById('create-edad')
// ulAlumnos es un elemento <ul> que contiene los alumnos de la base de datos renderizados (en forma de <li>)
const ulAlumnos = document.getElementById('read').querySelector('ul')
const nodoNuevoAlumno = document.createElement('li')

// Este evento manda una petición POST al handler Create.
createButton.addEventListener('click', (e) => {
    // Evitar recarga automatica de la pagina
    e.preventDefault()

    // Se validan datos y se retorna la función sin realizar la operación CREATE en caso de que
    // algún dato este mal.
    // No se eliminan los valores actuales de los input en vaso de que solo haya faltado
    // corregir un pequeño detalle.
    if (!validarLongitudCampo(createMatricula, 8, "Matricula")) return
    if (!validarLongitudCampo(createNombre, 20, "Nombre")) return
    if (!validarLongitudCampo(createEdad, 2, "Edad")) return

    const url = '/alumnos/create';
    const nuevoAlumno = {
        "matricula": createMatricula.value,
        "nombre": createNombre.value,
        "edad": Number(createEdad.value),
    }

    console.log(`Se ha hecho una petición POST con ${nuevoAlumno.matricula}, ${nuevoAlumno.nombre} ${nuevoAlumno.edad}`)
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json', // Se avisa que se mandará un JSON
        },
        body: JSON.stringify(nuevoAlumno),
    })
        .then(response => {
            return response.text();
        })
        .then(data => {
            // Se checa si el mensaje de respuesta contiene un error
            if (data.includes('error')) {
                throw new Error(`El servidor retorno un error: ${data}`);
            }
            console.log(`Response data: ${data}`);

            // Añadir dinamicamente el nuevo elemento creado; sin necesidad de refrescar la página
            // para ver los cambios.
            nodoNuevoAlumno.style.display = 'inline';
            nodoNuevoAlumno.style.minWidth = '70%';
            nodoNuevoAlumno.style.backgroundColor = 'cornsilk';
            nodoNuevoAlumno.textContent = `${nuevoAlumno.matricula}, ${nuevoAlumno.nombre}, ${nuevoAlumno.edad}`
            ulAlumnos.appendChild(nodoNuevoAlumno)


            // Se limpian campos de inputs en caso de que se haya hecho la operación exitosamente
            // Sería frustrante que se borren aún cuando falló y solo faltaba correjir un detalle.
            createMatricula.value = ""
            createNombre.value = ""
            createEdad.value = ""
            createMatricula.focus()
        })
        .catch(error => {
            console.error(`Fetch error: ${error}`);
        });
})

// -------------------- READ/buscar --------------------
// Elementos dentro del div buscar
const readMatricula = document.getElementById('buscar-matricula')
const readButton = document.getElementById('buscar-submit')

// Este evento manda una petición OPTIONS al handler Read, se busca por un alumno en específico por matricula.
// Se usará momentáneamente el método http OPTIONS debido a que el servidor aún no cuenta con lógica
// de autorización para bloquear el paso del cliente a la url /alumnos/read.
readButton.addEventListener('click', (e) => {
    // Evitar recarga automatica de la pagina
    e.preventDefault()

    if (!validarLongitudCampo(readMatricula, 8, "Matricula")) return

    const url = '/alumnos/read'
    const alumnoBuscado = {
        "matricula": readMatricula.value,
        "nombre": null,
        "edad": null,
    }

    console.log(`Se ha hecho una petición OPTIONS con ${alumnoBuscado.matricula}`)
    fetch(url, {
        method: 'OPTIONS',
        headers: {
            'Content-Type': 'application/json', // Se avisa que se mandará un JSON
        },
        body: JSON.stringify(alumnoBuscado),
    })
        .then(response => {
            return response.text();
        })
        .then(data => {
            // Se checa si el mensaje de respuesta contiene un error
            if (data.includes('error')) {
                throw new Error(`El servidor retorno un error: ${data}`);
            }
            console.log(`Response data: ${data}`);

            // Mostrar el elemento encontrado.
            alert(`${data}`)

            // Se limpian campos de inputs en caso de que se haya hecho la operación exitosamente
            // Sería frustrante que se borren aún cuando falló y solo faltaba correjir un detalle.
            readMatricula.value = ""
            readMatricula.focus()
        })
        .catch(error => {
            console.error(`Fetch error: ${error}`);
            alert(`No se encontró el alumno, error: ${error}`)
        })
})

// -------------------- UPDATE --------------------
// Elementos dentro del div update
const updateButton = document.getElementById('update-submit')
const updateMatricula = document.getElementById('update-matricula')
const updateNombre = document.getElementById('update-nombre')
const updateEdad = document.getElementById('update-edad')

// Este evento manda una petición PUT al handler Update.
updateButton.addEventListener('click', (e) => {
    // Evitar recarga automatica de la pagina
    e.preventDefault()

    // Se validan datos y se retorna la función sin realizar la operación UPDATE en caso de que
    // algún dato este mal.
    // No se eliminan los valores actuales de los input en vaso de que solo haya faltado
    // corregir un pequeño detalle.
    if (!validarLongitudCampo(updateMatricula, 8, "Matricula")) return
    if (!validarLongitudCampo(updateNombre, 20, "Nombre")) return
    if (!validarLongitudCampo(updateEdad, 2, "Edad")) return

    const url = '/alumnos/update';
    const nuevosDatosAlumno = {
        "matricula": updateMatricula.value,
        "nombre": updateNombre.value,
        "edad": Number(updateEdad.value),
    }

    console.log(`Se ha hecho una petición PUT con ${nuevosDatosAlumno.matricula}, ${nuevosDatosAlumno.nombre}, ${nuevosDatosAlumno.edad}`)
    fetch(url, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json', // Se avisa que se mandará un JSON
        },
        body: JSON.stringify(nuevosDatosAlumno),
    })
        .then(response => {
            return response.text();
        })
        .then(data => {
            // Se checa si el mensaje de respuesta contiene un error
            if (data.includes('error')) {
                throw new Error(`El servidor retorno un error: ${data}`);
            }
            console.log(`Response data: ${data}`);

            for (const child of Array.from(ulAlumnos.children)) {
                // datosEntrada contiene el string[n] de cada entrada renderizada:
                //  [0] = matricula + ",",
                //  [1] = nombre
                //  [2] = edad
                let datosEntrada = child.textContent.trim().split(' ')

                // Añadir 0s a la matricula en caso que no tenga 8 digitos.
                while (nuevosDatosAlumno.matricula.length < 8)
                    nuevosDatosAlumno.matricula = "0" + nuevosDatosAlumno.matricula;
                    
                if (datosEntrada[0].slice(0, -1) == nuevosDatosAlumno.matricula) {
                    let nuevoTextContent = `${nuevosDatosAlumno.matricula}, ${nuevosDatosAlumno.nombre} ${nuevosDatosAlumno.edad}`
                    child.textContent = nuevoTextContent
                    break
                }
            }

            // Se limpian campos de inputs en caso de que se haya hecho la operación exitosamente
            // Sería frustrante que se borren aún cuando falló y solo faltaba correjir un detalle.
            updateMatricula.value = ""
            updateMatricula.focus()
        })
        .catch(error => {
            console.error(`Fetch error: ${error}`);
            alert(`No se encontró el alumno, error: ${error}`)
        })
})
// -------------------- DELETE --------------------
// Elementos dentro del div delete
const deleteButton = document.getElementById('delete-submit')
const deleteMatricula = document.getElementById('delete-matricula')

// Este evento manda una petición DELETE al handler Update.
deleteButton.addEventListener('click', (e) => {
    // Evitar recarga automatica de la pagina
    e.preventDefault()

    // Se validan datos y se retorna la función sin realizar la operación DELETE en caso de que
    // algún dato este mal.
    // No se eliminan los valores actuales de los input en vaso de que solo haya faltado
    // corregir un pequeño detalle.
    if (!validarLongitudCampo(deleteMatricula, 8, "Matricula")) return

    const url = '/alumnos/delete';
    const alumnoAEliminar = {
        "matricula": deleteMatricula.value,
        "nombre": null,
        "edad": null,
    }

    console.log(`Se ha hecho una petición DELETE con ${alumnoAEliminar.matricula}`)
    fetch(url, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json', // Se avisa que se mandará un JSON
        },
        body: JSON.stringify(alumnoAEliminar),
    })
        .then(response => {
            return response.text()
        })
        .then(data => {
            // Se checa si el mensaje de respuesta contiene un error
            if (data.includes('error')) {
                throw new Error(`El servidor retorno un error: ${data}`)
            }
            console.log(`Response data: ${data}`)

            // Obtener todos los nodos con la etiqueta <li> en la lista de alumnos.
            let nodos = document.querySelectorAll('#read ul li')

            let matriculaAlumnoEliminado = `${alumnoAEliminar.matricula}`

            // Loop through the <li> nodos and find the one containing the string
            for (let i = 0; i < nodos.length; i++) {
                // matriculaEntradaI contiene la matricula de la entrada renderizada i.
                let matriculaEntradaI = nodos[i].textContent.trim().split(' ')[0].slice(0, -1)

                // Añadir 0s a la matricula en caso que no tenga 8 digitos.
                while (alumnoAEliminar.matricula.length < 8)
                    alumnoAEliminar.matricula = "0" + alumnoAEliminar.matricula;

                if (matriculaEntradaI == alumnoAEliminar.matricula) {
                    nodos[i].parentElement.removeChild(nodos[i])

                    break
                }
            }

            // Se limpian campos de inputs en caso de que se haya hecho la operación exitosamente
            // Sería frustrante que se borren aún cuando falló y solo faltaba correjir un detalle.
            deleteMatricula.value = ""
            deleteMatricula.focus()
        })
        .catch(error => {
            console.error(`Fetch error: ${error}`);
            alert(`No se encontró el alumno, error: ${error}`)
        })
})
// -----------------------------------------------------

// validarLongitudCampo revisa que el parametro campo tenga una longitud igual o menor a len.
// Si es así, retorna true, false para caso contrario.
// Se encapsula esta funcionalidad debido a que se usa en múltiples formularios a lo largo
// de la pagina HTML.
const validarLongitudCampo = (campo, len, nombre) => {
    if (campo.value.length > len) {
        alert(`El campo ${nombre} no puede exceder los ${len} caracteres.\nNo se realizó la operación create.`)
        return false
    }
    return true
}