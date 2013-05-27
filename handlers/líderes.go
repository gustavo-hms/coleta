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
)

func init() {
	http.HandleFunc("/líderes", Líderes)
}

func Líderes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		w.WriteHeader(http.StatusNotImplemented)
	case "GET":
		líderesGet(w, r, new(modelos.LíderValidado))
	case "POST":
		//		líderesPost(w, r)
	}
}

func líderesGet(
	w http.ResponseWriter,
	r *http.Request,
	líder *modelos.LíderValidado,
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
			if líder != nil && líder.Zona.Id == zona.Id {
				s.Selecionado = true
			}
			seleção = append(seleção, s)
		}
		return seleção
	}}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/líderes.html")
	if err != nil {
		log.Println("Ali:", err)
	}

	err = t.ExecuteTemplate(w, "líderes.html", líder)
	if err != nil {
		log.Println("Aqui:", err)
	}
}

func líderesPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
	}

	var líder modelos.Líder
	validado := líder.Preencher(r.FormValue)

	if validado != nil {
		w.WriteHeader(http.StatusBadRequest)
		líderesGet(w, r, validado)
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

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(&líder); err != nil {
		log.Println("Erro ao gravar líder:", err)
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println("Erro no commit:", err)
	}
	db.Close()

	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/adm-esquinas-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo esquinas-sucesso.html:", err)
	}

	fmt.Fprintf(w, "%s", página)
}
