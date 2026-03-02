package bd

import (
	"database/sql"
	"fmt"
	"log"

	
	_ "github.com/microsoft/go-mssqldb"
)

// Conectar establece la conexión con SQL Server y devuelve el objeto de la base de datos
func Conectar() *sql.DB {
	servidor := "localhost\\SQLEXPRESS" // Si usaste SQLEXPRESS, podría ser "localhost\\SQLEXPRESS"
	puerto := 1433
	usuario := "sa"            
	password := "ProyectoF320" 
	baseDatos := "ecommerce"

	// Cadena de conexión usando SQL Authentication (Usuario y Clave)
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		servidor, usuario, password, puerto, baseDatos)

	

	// Abrimos el pool de conexión
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error configurando la conexión: ", err.Error())
	}

	// Ping verifica que realmente hay comunicación con la base de datos
	err = db.Ping()
	if err != nil {
		log.Fatal("Error conectando a la base de datos (revisa credenciales o TCP/IP): ", err.Error())
	}

	fmt.Println("¡Conexión exitosa a la base de datos ecommerce!")
	return db
}
