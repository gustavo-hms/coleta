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
	registrarComTransação("/lideres", Líderes{})
}

type Líderes struct{}

func (l Líderes) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	líder := modelos.NovoLíder()
	return l.get(w, r, tx, validação.NovoLíderComErros(líder))
}

func (l Líderes) get(
	w http.ResponseWriter,
	r *http.Request,
	tx *dao.Tx,
	líder *validação.LíderComErros,
) error {
	t := exibiçãoDoLíder(líder, tx, "líderes.html")

	if t == nil {
		return erroInesperado
	}

	err := t.ExecuteTemplate(w, "líderes.html", líder)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (l Líderes) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return erroInesperado
	}

	líder := modelos.NovoLíder()
	líder.Preencher(r.Form)
	erros, falha := validação.ValidarLíder(líder, tx)

	if falha != nil {
		return falha
	}

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		return l.get(w, r, tx, erros)
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(líder); err != nil {
		log.Println(err)
		return err
	}

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/cadastro-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo cadastro-sucesso.html:", err)
		return erroInesperado
	}

	fmt.Fprintf(w, "%s", página)

	return nil
}

func exibiçãoDoLíder(líder *validação.LíderComErros, tx *dao.Tx, página string) *template.Template {
	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(fmt.Sprintf("%d", líder.Zona.Id))
	if err != nil {
		log.Println("Erro ao buscar esquinas:", err)
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAll()
	if err != nil {
		log.Println("Erro ao buscar zonas:", err)
		return nil
	}

	funcMap := template.FuncMap{
		"esquinas": func() []esquinaComSeleção {
			seleção := make([]esquinaComSeleção, 0, len(esquinas))
			for _, esquina := range esquinas {
				s := esquinaComSeleção{Esquina: esquina}
				if líder != nil && líder.Esquina.Id == esquina.Id {
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
				if líder != nil && líder.Zona.Id == zona.Id {
					s.Selecionado = true
				}
				seleção = append(seleção, s)
			}
			return seleção
		},

		"turnos": func() []turnoComSeleção {
			turnos := modelos.Turnos()
			turnos = turnos[:2]
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
		ParseFiles(config.Dados.DiretórioDasPáginas + "/" + página)
	if err != nil {
		log.Println(err)
		return nil
	}

	return t
}
