/*
@autor: Jeferson Castro
@fecha: 01/02/2026
@descripcion: Proyecto final. Evaluacion en contacto con el docente
*/

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"Proyectofinal/bd"
	"Proyectofinal/modelos"
)

// INTERFACES
// Las interfaces en Go definen un "contrato" de comportamiento.
// En lugar de depender de un tipo de dato exacto (como un struct Producto),
// exigimos que cualquier objeto que quiera ser "Vendible" debe tener una
// función llamada GenerarEtiqueta(). Esto aporta flexibilidad y polimorfismo al código.
type Vendible interface {
	GenerarEtiqueta() string
}

func main() {
	// Conexion a base de datos
	db := bd.Conectar()
	defer db.Close()

	// 2. Configurar la ruta principal de la web ("/")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Consultar SQL Server
		filas, err := db.Query("SELECT ID, Nombre, Marca, Stock, Precio FROM Productos")
		if err != nil {
			http.Error(w, "Error al leer productos", http.StatusInternalServerError)
			return
		}
		defer filas.Close()

		// SLICES Y ESTRUCTURAS (STRUCTS)
		// Un Slice es un arreglo dinámico en Go. Lo usamos aquí porque
		// no sabemos de antemano cuántos productos hay en la base de datos.
		// El Slice irá creciendo automáticamente (con 'append') a medida que
		// le metamos las estructuras (structs) de cada producto encontrado.
		var catalogo []modelos.Producto

		// Se llena el Slice con los datos de SQL
		for filas.Next() {
			var p modelos.Producto
			filas.Scan(&p.ID, &p.Nombre, &p.Marca, &p.Stock, &p.Precio)
			catalogo = append(catalogo, p)
		}

		// Leer el archivo HTML y enviarle el catálogo
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, catalogo)
	})

	// Este bloque procesa la compra
	http.HandleFunc("/comprar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Leer los datos del formulario web
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		email := r.FormValue("email")
		productoID := r.FormValue("producto_id")

		// CONSULTAS PARAMETRIZADAS (@p1, @p2...) Y SEGURIDAD
		// Nunca debemos concatenar texto del usuario directamente en el SQL.
		// Usamos parámetros (@p1, @p2) propios de SQL Server para que la base de datos
		// limpie los datos. Esto previene ataques de Inyección SQL, asegurando el sistema.
		_, errInsert := db.Exec("INSERT INTO Usuarios (Nombre, Apellido, Email) VALUES (@p1, @p2, @p3)", nombre, apellido, email)
		if errInsert != nil {
			fmt.Println(" Error al insertar usuario:", errInsert)
		}

		//Restar 1 al stock del producto
		_, errUpdate := db.Exec("UPDATE Productos SET Stock = Stock - 1 WHERE ID = @p1", productoID)
		if errUpdate != nil {
			fmt.Println(" Error al actualizar stock:", errUpdate)
		}

		//Obtener el nombre del producto para el mensaje final
		var nombreProducto string

		fmt.Println("ID recibido del formulario HTML:", productoID)

		err := db.QueryRow("SELECT Nombre FROM Productos WHERE ID = @p1", productoID).Scan(&nombreProducto)
		if err != nil {
			fmt.Println(" Error al buscar el nombre del producto en SQL:", err)
		}

		//CONCURRENCIA (GOROUTINES)
		// Enviar un correo electrónico o simularlo puede tardar varios
		// segundos. Al envolver esta acción en una 'go func()', creamos un hilo ligero
		// de ejecución en segundo plano. Así, el usuario no se queda esperando frente
		// a una pantalla de carga; la web responde de inmediato mientras el proceso
		// "pesado" ocurre en segundo plano.

		go func(destinatario, articulo string) {
			fmt.Println("------------------------------------------------")
			fmt.Printf("[PROCESO EN SEGUNDO PLANO] Enviando recibo al correo: %s por el artículo: %s\n", destinatario, articulo)
			fmt.Println("------------------------------------------------")
		}(email, nombreProducto)

		// MAPS
		// Un mapa es una estructura de datos de clave-valor. Lo usamos aquí
		// en lugar de crear un struct nuevo porque solo necesitamos pasar un par de
		// textos temporales a la pantalla de éxito. Es una forma rápida, limpia y
		// eficiente de agrupar datos sin estructurar para enviarlos a la vista HTML.

		datosExito := make(map[string]string)
		datosExito["NombreUsuario"] = nombre + " " + apellido
		datosExito["NombreProducto"] = nombreProducto

		// USO DE LA INTERFAZ: Creamos un producto temporal solo para imprimir su etiqueta
		var articulo Vendible = modelos.Producto{Nombre: nombreProducto, Marca: "MSI"}
		fmt.Println("Generando etiqueta de envío:", articulo.GenerarEtiqueta())

		// Mostrar la página verde de confirmación
		tmpl := template.Must(template.ParseFiles("templates/exito.html"))
		tmpl.Execute(w, datosExito)
	})

	// Servidor HTTP Y LOG.FATAL
	// ListenAndServe mantiene el programa vivo escuchando peticiones.
	// Lo envolvemos en log.Fatal para que, si el puerto 8080 está ocupado o el
	// servidor colapsa, el programa se detenga inmediatamente e imprima el error exacto.
	fmt.Println("------------------------------------------------")
	fmt.Println("Servidor web corriendo en http://localhost:8080")
	fmt.Println("Abre tu navegador para ver la tienda.")
	fmt.Println("------------------------------------------------")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
