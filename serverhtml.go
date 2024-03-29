//Aplicacion de facturacion
package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*.html"))

func conexionBd() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasena := "emmanuelyk1"
	Nombre := "almacen"
	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1:3305)/"+Nombre)
	if err != nil {
		panic(err.Error())

	}
	return conexion
}

func main() {

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/factura", Factura)
	http.HandleFunc("/scanerfacta", ScanerFac)
	http.HandleFunc("/scanerverifica", ScanerVerifica)
	http.HandleFunc("/guardado", Guadado)

	log.Println("Servidor corriendo........")
	http.ListenAndServe(":8080", nil)
}

func Guadado(w http.ResponseWriter, r *http.Request) {

	plantillas.ExecuteTemplate(w, "guardado", nil)

}

type Guardando struct {
	Ndocumento string
	RutaStru   string
}

func ScanerVerifica(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBd()
	recopilardatos := Guardando{}
	arregloDatos := []Guardando{}

	numfactura := strings.TrimSpace(r.URL.Query().Get("nfactura"))
	numruta := strings.TrimSpace(r.URL.Query().Get("nruta"))
	recopilardatos.Ndocumento = numfactura
	recopilardatos.RutaStru = numruta
	arregloDatos = append(arregloDatos, recopilardatos)

	//verificamos  si existe el numero de factura

	rutaCondicion := "\"" + numruta + "\""
	facturaCondicion := "\"" + numfactura + "\""
	registros, err := conexionEstablecida.Query("SELECT Factura_verifica,Cuenta,Farmacia,Ruta,ver_factura FROM scaner_factura WHERE Ruta=" + rutaCondicion + " AND " + "ver_factura =  0" + " AND " + "Factura_verifica = " + facturaCondicion)

	if err != nil {
		panic(err.Error())

	}
	contenedorV := FacturaScaner{}
	arregloContenedores := []FacturaScaner{}
	for registros.Next() {

		/*var contenedor string
		var ruta string
		*/
		//---------------------
		var FacturaverificaScanerV string
		var CuentaScanerV string
		var RutaScanerV string
		var FarmaciaScanerV string
		var VerfacturaScanerV string

		//----------
		err = registros.Scan(&FacturaverificaScanerV, &CuentaScanerV, &FarmaciaScanerV, &RutaScanerV, &VerfacturaScanerV)
		if err != nil {
			panic(err.Error())

		} else {
			contenedorV.FacturaverificaScaner = FacturaverificaScanerV
			contenedorV.CuentaScaner = CuentaScanerV
			contenedorV.RutaScaner = RutaScanerV
			contenedorV.FarmaciaScaner = FarmaciaScanerV
			contenedorV.VerfacturaScaner = VerfacturaScanerV
			arregloContenedores = append(arregloContenedores, contenedorV)

		}

	}
	//se obtenemos un registro guadamos la informacion
	if len(arregloContenedores) != 0 {
		actualizaRegistro(numfactura)
		//fmt.Println("Se guandar los cambios")
		retornoListadoTabla := listatoRutaScaner(numruta)
		//fmt.Println(retornoListadoTabla)
		plantillas.ExecuteTemplate(w, "scanerfactura", retornoListadoTabla)
	} else {
		plantillas.ExecuteTemplate(w, "guardado", rutaCondicion)

		//http.Error(w, "No esxiste el el cocumento", 500)

	}

	//repintamos la pantalla de verificacion de facturas

	//fmt.Println(arregloContenedores)
	//fmt.Println(numfactura)
	//fmt.Println(numruta)
	//fmt.Println("Cuantos elementos tiene el arreglo")
	//fmt.Println(len(arregloContenedores))
	//fin de verificacion de numero de documento
	//	plantillas.ExecuteTemplate(w, "guardado", nil)
	//fmt.Println(numfactura)
	//fmt.Println(numruta)
	//http.Redirect(w, r, "/scanerfacta", 301)

}

func listatoRutaScaner(rutaCondicion string) []FacturaScaner {
	//fmt.Println("La ruta de condicion es ", rutaCondicion)
	conexionEstablecida := conexionBd()
	//fmt.Println("\"" + rutaPantalla + "\"")
	rutaCondicion = strings.TrimSpace(rutaCondicion)
	rutaCondicion = "\"" + rutaCondicion + "\""
	registros, err := conexionEstablecida.Query("SELECT Factura_verifica,Cuenta,Farmacia,Ruta,ver_factura FROM scaner_factura WHERE Ruta=" + rutaCondicion + " AND " + "ver_factura =  0")

	if err != nil {
		panic(err.Error())

	}
	contenedorV := FacturaScaner{}
	arregloContenedores := []FacturaScaner{}
	for registros.Next() {

		/*var contenedor string
		var ruta string
		*/
		//---------------------
		var FacturaverificaScanerV string
		var CuentaScanerV string
		var RutaScanerV string
		var FarmaciaScanerV string
		var VerfacturaScanerV string

		//----------
		err = registros.Scan(&FacturaverificaScanerV, &CuentaScanerV, &FarmaciaScanerV, &RutaScanerV, &VerfacturaScanerV)
		if err != nil {
			panic(err.Error())

		} else {
			contenedorV.FacturaverificaScaner = FacturaverificaScanerV
			contenedorV.CuentaScaner = CuentaScanerV
			contenedorV.RutaScaner = RutaScanerV
			contenedorV.FarmaciaScaner = FarmaciaScanerV
			contenedorV.VerfacturaScaner = VerfacturaScanerV
			arregloContenedores = append(arregloContenedores, contenedorV)

		}

	}

	//fmt.Println(rutaPantalla)
	//fmt.Println(arregloContenedores)
	return arregloContenedores

	//plantillas.ExecuteTemplate(w, "scanerfactura", arregloContenedores)

}

func actualizaRegistro(guadarFactura string) {

	conexionEstablecida := conexionBd()
	guadarFactura = strings.TrimSpace(guadarFactura)
	facturaCondicion := "\"" + guadarFactura + "\""
	actualizaregistro, err := conexionEstablecida.Prepare("UPDATE scaner_factura SET  ver_factura =  1 " + "WHERE " + "Factura_verifica = " + facturaCondicion)
	if err != nil {
		panic(err.Error())

	}
	actualizaregistro.Exec()

	//fmt.Println("se guadala la factura", guadarFactura)

}

type Contenedores struct {
	ContenedorStru string
	RutaStru       string
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, )
	conexionEstablecida := conexionBd()
	registros, err := conexionEstablecida.Query("SELECT Contenedor,Ruta FROM contenedores where ruta = \"COCO\"")
	if err != nil {
		panic(err.Error())

	}
	contenedorV := Contenedores{}
	arregloContenedores := []Contenedores{}
	for registros.Next() {

		var contenedor string
		var ruta string
		err = registros.Scan(&contenedor, &ruta)
		if err != nil {
			panic(err.Error())

		} else {
			contenedorV.ContenedorStru = contenedor
			contenedorV.RutaStru = ruta
			arregloContenedores = append(arregloContenedores, contenedorV)

		}

	}
	//fmt.Println(arregloContenedores)

	plantillas.ExecuteTemplate(w, "inicio", arregloContenedores)

}

type RutaStru struct {
	Rutas string
}

func Factura(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := conexionBd()

	//plantillas.ExecuteTemplate(w, "facturas", nil)
	registros, err := conexionEstablecida.Query("SELECT ruta FROM rutas_facturas ")
	if err != nil {
		panic(err.Error())

	}
	rutaV := RutaStru{}
	arregloRutas := []RutaStru{}
	///----------------
	for registros.Next() {

		var ruta string
		err = registros.Scan(&ruta)
		if err != nil {
			panic(err.Error())

		} else {
			rutaV.Rutas = ruta

			arregloRutas = append(arregloRutas, rutaV)

		}

	}
	//fmt.Println(arregloRutas)
	plantillas.ExecuteTemplate(w, "facturas", arregloRutas)
}

type FacturaScaner struct {
	FacturaverificaScaner string
	CuentaScaner          string
	RutaScaner            string
	FarmaciaScaner        string
	VerfacturaScaner      string
}

func ScanerFac(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBd()
	rutaPantalla := r.URL.Query().Get("browser")
	//fmt.Println("\"" + rutaPantalla + "\"")
	rutaCondicion := "\"" + rutaPantalla + "\""
	registros, err := conexionEstablecida.Query("SELECT Factura_verifica,Cuenta,Farmacia,Ruta,ver_factura FROM scaner_factura WHERE Ruta=" + rutaCondicion + " AND " + "ver_factura =  0")

	if err != nil {
		panic(err.Error())

	}
	contenedorV := FacturaScaner{}
	arregloContenedores := []FacturaScaner{}
	for registros.Next() {

		/*var contenedor string
		var ruta string
		*/
		//---------------------
		var FacturaverificaScanerV string
		var CuentaScanerV string
		var RutaScanerV string
		var FarmaciaScanerV string
		var VerfacturaScanerV string

		//----------
		err = registros.Scan(&FacturaverificaScanerV, &CuentaScanerV, &FarmaciaScanerV, &RutaScanerV, &VerfacturaScanerV)
		if err != nil {
			panic(err.Error())

		} else {
			contenedorV.FacturaverificaScaner = FacturaverificaScanerV
			contenedorV.CuentaScaner = CuentaScanerV
			contenedorV.RutaScaner = RutaScanerV
			contenedorV.FarmaciaScaner = FarmaciaScanerV
			contenedorV.VerfacturaScaner = VerfacturaScanerV
			arregloContenedores = append(arregloContenedores, contenedorV)

		}

	}

	//fmt.Println(rutaPantalla)
	//fmt.Println(arregloContenedores)

	plantillas.ExecuteTemplate(w, "scanerfactura", arregloContenedores)

}
