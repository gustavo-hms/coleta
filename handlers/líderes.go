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
	t := exibiçãoDoLíder(líder, "líderes.html")
	if t != nil {
		err := t.ExecuteTemplate(w, "líderes.html", líder)
		if err != nil {
			log.Println("Aqui:", err)
		}
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
	erros := líder.Preencher(r.Form)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		líderesGet(w, r, erros)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(&líder); err != nil {
		líderDAO.Rollback()
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

	página, err := ioutil.ReadFile(gopath + "/src/coleta/páginas/líderes-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo esquinas-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}

func exibiçãoDoLíder(líder *modelos.LíderComErros, página string) *template.Template {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Erro ao iniciar transação:", err)
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAll()
	if err != nil {
		zonaDAO.Rollback()
		log.Println("Erro ao buscar zonas:", err)
		return nil
	}

	zonaDAO.Commit()

	funcMap := template.FuncMap{
		"zonas": func() []zonaComSeleção {
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
				if líder != nil {
					for _, t := range líder.Turnos {
						if s.Turno.Id == t.Id {
							s.Selecionado = true
						}
					}
				}

				seleção = append(seleção, s)
			}

			return seleção
		},

		"iguais": func(x, y string) bool {
			return x == y
		},
	}

	t, err := template.New("esquinas").Funcs(funcMap).
		ParseFiles(gopath + "/src/coleta/páginas/" + página)
	if err != nil {
		log.Println("Ali:", err)
		return nil
	}

	return t
}
