package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"coleta/modelos/validação"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	registrarComTransação("/voluntarios", Voluntários{})
}

type Voluntários struct{}

func (v Voluntários) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	voluntário := modelos.NovoVoluntário()
	return v.get(w, r, tx, validação.NovoVoluntárioComErros(voluntário))
}

func (v Voluntários) get(
	w http.ResponseWriter,
	r *http.Request,
	tx *dao.Tx,
	voluntário *validação.VoluntárioComErros,
) error {
	t := exibiçãoDoVoluntário(voluntário, tx, "voluntários.html")

	if t == nil {
		return erroInesperado
	}

	err := t.ExecuteTemplate(w, "voluntários.html", voluntário)
	if err != nil {
		log.Println("Aqui:", err)
		return err
	}

	return nil
}

func (v Voluntários) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return err
	}

	voluntário := modelos.NovoVoluntário()
	voluntário.Preencher(r.Form)
	erros := validação.ValidarVoluntário(voluntário)

	if erros != nil {
		if erros.Líder.Id != 0 {
			líderDAO := dao.NewLiderDAO(tx)
			l, err := líderDAO.FindById(voluntário.Líder.Id)
			if err != nil {
				log.Println(err)
				return err
			}

			erros.Líder = l
		}

		w.WriteHeader(http.StatusBadRequest)
		return v.get(w, r, tx, erros)
	}

	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	if err := voluntárioDAO.Save(voluntário); err != nil {
		log.Println("Erro ao gravar voluntário:", err)
		return err
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/cadastro-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo cadastro-sucesso.html:", err)
		return err
	}

	fmt.Fprintf(w, "%s", página)

	return nil
}

func exibiçãoDoVoluntário(voluntário *validação.VoluntárioComErros, tx *dao.Tx, página string) *template.Template {
	esquinaDAO := dao.NewEsquinaDAO(tx.Tx)
	esquinas, err := esquinaDAO.BuscarPorZona(fmt.Sprintf("%d", voluntário.Zona.Id))
	if err != nil {
		log.Println("Erro ao buscar esquinas:", err)
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx.Tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		log.Println("Erro ao buscar zonas:", err)
		return nil
	}

	funcMap := template.FuncMap{
		"esquinas": func() []esquinaComSeleção {
			seleção := make([]esquinaComSeleção, 0, len(esquinas))
			for _, esquina := range esquinas {
				s := esquinaComSeleção{Esquina: esquina}
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
		log.Println(err)
		return nil
	}

	return t
}
