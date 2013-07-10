package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"coleta/modelos/validação"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrar("/voluntarios", Voluntários{})
}

type Voluntários struct{}

func (v Voluntários) Get(w http.ResponseWriter, r *http.Request) {
	voluntário := modelos.NovoVoluntário()
	v.get(w, r, validação.NovoVoluntárioComErros(voluntário))
}

func (v Voluntários) get(
	w http.ResponseWriter,
	r *http.Request,
	voluntário *validação.VoluntárioComErros,
) {
	t := exibiçãoDoVoluntário(voluntário, "voluntários.html")
	if t != nil {
		err := t.ExecuteTemplate(w, "voluntários.html", voluntário)
		if err != nil {
			log.Println("Aqui:", err)
		}
	}
}

/*
func (v Voluntários) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	voluntário := modelos.NovoVoluntário()
	voluntário.Preencher(r.Form)
	erros := validação.ValidarVoluntário(voluntário)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		v.get(w, r, erros)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(voluntário); err != nil {
		líderDAO.Rollback()
		log.Println("Erro ao gravar voluntário:", err)
		erroInterno(w, r)
		return
	}
	if err := líderDAO.Commit(); err != nil {
		líderDAO.Rollback()
		log.Println("Erro no commit:", err)
		erroInterno(w, r)
		return
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/cadastro-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo cadastro-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}
*/
func exibiçãoDoVoluntário(voluntário *validação.VoluntárioComErros, página string) *template.Template {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Erro ao iniciar transação:", err)
		return nil
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(fmt.Sprintf("%d", voluntário.Zona.Id))
	if err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro ao buscar esquinas:", err)
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
		"esquinas": func() []esquinaComSeleção {
			seleção := make([]esquinaComSeleção, 0, len(esquinas))
			for _, esquina := range esquinas {
				s := esquinaComSeleção{Esquina: *esquina}
				if voluntário != nil && voluntário.Esquina.Id == esquina.Id {
					s.Selecionado = true
				}
				seleção = append(seleção, s)
			}
			return seleção
		},

		"zonas": func() []zonaComSeleção {
			seleção := make([]zonaComSeleção, 0, len(zonas))
			for _, zona := range zonas {
				s := zonaComSeleção{Zona: *zona}
				if voluntário != nil && voluntário.Zona.Id == zona.Id {
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
				if voluntário != nil {
					for _, t := range voluntário.Turnos {
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
		ParseFiles(config.Dados.DiretórioDasPáginas + "/" + página)
	if err != nil {
		log.Println("Ali:", err)
		return nil
	}

	return t
}
