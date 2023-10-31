# proyectofinalbd
SPA CRUD featuring: HTML templates, web requests and the MVC architecture (code in Spanglish).

Se conecta la app con una base de datos local Oracle para realizar
operaciones CRUD (Create, Read, Update y Delete) sobre una tabla 'alumnos'.\
La base de datos debera estar ejecutándose para que se pueda realizar la
conexión de forma exitosa, para esto, revisar los servicios (WIN + R -> services.msc).\
Para provocar una operación se deberá hacer una petición mediante la GUI web,
esta tiene botones que realizan peticiones web via JavaScript a las URL handler
que ejecutan las operaciones CRUD a la base de datos.\
Si se intenta ingresar a las URL designadas para hacer una operacion CRUD, el usuario
sera redirigido automáticamente de vuelta a la página principal:\
/alumnos/create -> /

---

App establishes a connection with a local OracleSQL database thus enabling
the ability to execute CRUD operations (Create, Read, Update, Delete) over a table 'alumnos'.\
Database's service must be executing to be able to succeed operations, for this, check
Windows Services (WIN + R -> services.msc).\
The triggering of operations is done by sending web requests via JavaScript on the
initialized web server, achieved through buttons on the web GUI.\
If client tries to enter any that URL designed to reply to requests, it'll
be redirected by the server to the homepage.
/alumnos/create -> /

#### Conexión/Connection:
username: proyectofinalbd\
password: 123456\
service: XEPDB1
