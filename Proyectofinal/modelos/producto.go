package modelos

import "fmt"

// INTERFAZ: 
// Define que cualquier cosa "Vendible" debe poder generar una etiqueta.
type Vendible interface {
	GenerarEtiqueta() string
}

// Producto define nuestro objeto
type Producto struct {
	ID     int
	Nombre string
	Marca  string
	Stock  int
	Precio float64
}

//  Adjuntamos un método a nuestro struct Producto.
// Con esto, Producto implementa automáticamente la interfaz Vendible.
func (p Producto) GenerarEtiqueta() string {
	return fmt.Sprintf("[Artículo: %s | Marca: %s | Precio: $%.2f]", p.Nombre, p.Marca, p.Precio)
}