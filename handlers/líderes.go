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
		líderesGet(w, r, modelos.NovoLíderComErros())
	case "POST":
		líderesPost(w, r)
	}
}

func líderesGet(
	w http.ResponseWriter,
	r *http.Request,
	líder *modelos.LíderComErros,
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
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)

	funcMap := template.FuncMap{
		"zonas": func() []zonaComSeleção {
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
		},

		"turnos": func() []turnoComSeleção {
			turnos := modelos.Turnos()
			seleção := make([]turnoComSeleção, 0, len(turnos))
			for _, turno := range turnos {
				s := turnoComSeleção{Turno: turno}
				seleção = append(seleção, s)
			}

			return seleção
		},
	}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/líderes.html")
	if err != nil {
		log.Println("Ali:", err)
		return
	}

	err = t.ExecuteTemplate(w, "líderes.html", líder)
	if err != nil {
		log.Println("Aqui:", err)
		return
	}
}

func líderesPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	var líder modelos.Líder
	validado := líder.Preencher(r.Form)

	if validado != nil {
		w.WriteHeader(http.StatusBadRequest)
		líderesGet(w, r, validado)
		return
	}

	db, err := db.Conn()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(&líder); err != nil {
		log.Println("Erro ao gravar líder:", err)
		erroInterno(w, r)
		return
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println("Erro no commit:", err)
		erroInterno(w, r)
		return
	}
	db.Close()

	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/líderes-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo esquinas-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}
