package handlers

import (
	"coleta/dao"
	"coleta/db"
	"coleta/modelos"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var gopath = os.Getenv("GOPATH") // NOTA Solução temporária. Apenas para testes

func init() {
	http.HandleFunc("/adm/esquinas", Esquinas)
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

func esquinasGet(
	w http.ResponseWriter,
	r *http.Request,
	esquina *modelos.EsquinaValidada,
) {
	// TODO create connection transaction outside
	db, err := db.Conn()
	if err != nil {
		log.Println(err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	zonaDAO := dao.NewZonaDAO(tx)

	funcMap := template.FuncMap{"zonas": func() []zonaComSeleção {
		zonas, err := zonaDAO.FindAll()
		if err != nil {
			log.Println(err)
			return nil
		}

		seleção := make([]zonaComSeleção, 0, len(zonas))
		for _, zona := range zonas {
			s := zonaComSeleção{Zona: *zona}
			if esquina != nil && esquina.Zona.Id == zona.Id {
				s.Selecionado = true
			}
			seleção = append(seleção, s)
		}
		return seleção
	}}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/adm-esquinas.html")
	if err != nil {
		log.Println("Ali:", err)
	}

	err = t.ExecuteTemplate(w, "adm-esquinas.html", esquina)
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
	validada := esquina.Preencher(r.FormValue)

	if validada != nil {
		w.WriteHeader(http.StatusBadRequest)
		esquinasGet(w, r, validada)
		return
	}

	db, err := db.Conn()
	if err != nil {
		log.Println(err)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	if err := esquinaDAO.Save(&esquina); err != nil {
		log.Println("Erro ao gravar esquina:", err)
	}
	if err := esquinaDAO.Commit(); err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro no commit:", err)
	}
	db.Close()

	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/adm-esquinas-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo esquinas-sucesso.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
