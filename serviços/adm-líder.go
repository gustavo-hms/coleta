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
	"net/url"
	"strconv"
	"strings"
)

func init() {
	registrarSeguroComTransação("/adm/lider/", AdmLíder{})
}

type AdmLíder struct{}

func idDoLíder(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (l AdmLíder) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro: %s", stringDoId, err)
		return err
	}

	líderDAO := dao.NewLiderDAO(tx)
	líder, err := líderDAO.FindById(id)
	if err != nil {
		log.Printf("Erro ao carregar líder com id %d: %s", id, err)
		return err
	}

	return l.get(w, r, tx, &validação.LíderComErros{Líder: *líder})
}

func (l AdmLíder) get(
	w http.ResponseWriter,
	r *http.Request,
	tx *dao.Tx,
	líder *validação.LíderComErros,
) error {
	t := exibiçãoDoLíderAdm(líder, tx, "adm-líder.html")
	if t == nil {
		return erroInesperado
	}

	err := t.ExecuteTemplate(w, "adm-líder.html", líder)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (l AdmLíder) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return erroInesperado
	}

	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro:", stringDoId, err)
		return erroInesperado
	}

	líder := modelos.NovoLíder()
	líder.Id = id
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

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/atualização-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo líderes-sucesso.html:", err)
		return erroInesperado
	}

	fmt.Fprintf(w, "%s", página)

	return nil
}

func exibiçãoDoLíderAdm(líder *validação.LíderComErros, tx *dao.Tx, página string) *template.Template {
	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(fmt.Sprintf("%d", líder.Zona.Id))
	if err != nil {
		log.Println("Erro ao buscar esquinas:", err)
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx)
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
