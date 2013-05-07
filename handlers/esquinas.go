package handlers

import (
	"coleta/modelos"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/esquinas", Esquinas)
}

func Esquinas(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		esquinasGet(w, r)
	} else {
		esquinasPost(w, r)
	}
}

func esquinasGet(w http.ResponseWriter, r *http.Request) {
	gopath := os.Getenv("GOPATH") // NOTA Solução temporária. Apenas para testes

	funcMap := template.FuncMap{"zonas": modelos.ListaDeZonas}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/esquinas.html")
	if err != nil {
		log.Println("Ali:", err)
	}

	var esquina modelos.EsquinaValidada
	err = t.ExecuteTemplate(w, "esquinas.html", &esquina)
	if err != nil {
		log.Println("Aqui:", err)
	}
}

func esquinasPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
	}

	var esquina modelos.Esquina
	esquina.Preencher(r)

	log.Printf("Esquina: %+v\n", esquina)
}
