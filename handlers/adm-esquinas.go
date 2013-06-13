package handlers

import (
	"coleta/dao"
	"coleta/modelos"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/adm/esquinas", Esquinas)
}

func Esquinas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		w.WriteHeader(http.StatusNotImplemented)
	case "GET":
		esquinasGet(w, r, new(modelos.EsquinaComErros))
	case "POST":
		esquinasPost(w, r)
	}
}

func esquinasGet(
	w http.ResponseWriter,
	r *http.Request,
	esquina *modelos.EsquinaComErros,
) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Início da transação:", err)
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAll()
	if err != nil {
		zonaDAO.Rollback()
		log.Println(err)
		return
	}

	zonaDAO.Commit()

	funcMap := template.FuncMap{"zonas": func() []zonaComSeleção {
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
		return
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

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	if err := esquinaDAO.Save(&esquina); err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro ao gravar esquina:", err)
		return
	}

	if err := esquinaDAO.Commit(); err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro no commit:", err)
	}

	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/adm-esquinas-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo esquinas-sucesso.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
