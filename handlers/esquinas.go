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
	switch r.Method {
	default:
		w.WriteHeader(http.StatusNotImplemented)
	case "GET":
		esquinasGet(w, r, new(modelos.EsquinaValidada))
	case "POST":
		esquinasPost(w, r)
	}
}

type zonaComSeleção struct {
	Zona        modelos.Zona
	Selecionado bool
}

func esquinasGet(
	w http.ResponseWriter,
	r *http.Request,
	esquina *modelos.EsquinaValidada,
) {
	gopath := os.Getenv("GOPATH") // NOTA Solução temporária. Apenas para testes

	funcMap := template.FuncMap{"zonas": func() []zonaComSeleção {
		zonas := modelos.ListaDeZonas()
		seleção := make([]zonaComSeleção, 0, len(zonas))
		for _, zona := range zonas {
			s := zonaComSeleção{Zona: zona}
			if esquina != nil && esquina.Zona == zona {
				s.Selecionado = true
			}
			seleção = append(seleção, s)
		}
		return seleção
	}}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/esquinas.html")
	if err != nil {
		log.Println("Ali:", err)
	}

	err = t.ExecuteTemplate(w, "esquinas.html", esquina)
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
	validada := esquina.Preencher(r)

	if validada != nil {
		w.WriteHeader(http.StatusBadRequest)
		esquinasGet(w, r, validada)
		return
	}

	log.Printf("Esquina: %+v\n", esquina)
}
