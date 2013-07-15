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
	registrarSeguro("/adm/lider/", AdmLíder{})
}

type AdmLíder struct{}

func idDoLíder(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (l AdmLíder) Get(w http.ResponseWriter, r *http.Request) {
	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro: %s", stringDoId, err)
		erroInterno(w, r)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	líder, err := líderDAO.FindById(id)
	if err != nil {
		líderDAO.Rollback()
		log.Printf("Erro ao carregar líder com id %d: %s", id, err)
		erroInterno(w, r)
		return
	}

	líderDAO.Commit()

	l.get(w, r, &validação.LíderComErros{Líder: *líder})
}

func (l AdmLíder) get(
	w http.ResponseWriter,
	r *http.Request,
	líder *validação.LíderComErros,
) {
	t := exibiçãoDoLíderAdm(líder, "adm-líder.html")
	if t != nil {
		err := t.ExecuteTemplate(w, "adm-líder.html", líder)
		if err != nil {
			log.Println("Aqui:", err)
			erroInterno(w, r)
			return
		}
	}
}

func (l AdmLíder) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		erroInterno(w, r)
		return
	}

	stringDoId := idDoLíder(r.URL)
	id, err := strconv.Atoi(stringDoId)
	if err != nil {
		log.Printf("Não foi possível converter %s para um inteiro:", stringDoId, err)
		erroInterno(w, r)
		return
	}

	líder := modelos.NovoLíder()
	líder.Id = id
	líder.Preencher(r.Form)
	erros := validação.ValidarLíder(líder)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		l.get(w, r, erros)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		erroInterno(w, r)
		return
	}

	líderDAO := dao.NewLiderDAO(tx)
	if err := líderDAO.Save(líder); err != nil {
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

	página, err := ioutil.ReadFile(config.Dados.DiretórioDasPáginas + "/atualização-sucesso.html")
	if err != nil {
		log.Println("Erro ao abrir o arquivo líderes-sucesso.html:", err)
		erroInterno(w, r)
		return
	}

	fmt.Fprintf(w, "%s", página)
}

func exibiçãoDoLíderAdm(líder *validação.LíderComErros, página string) *template.Template {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Erro ao iniciar transação:", err)
		return nil
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquinas, err := esquinaDAO.BuscarPorZona(fmt.Sprintf("%d", líder.Zona.Id))
	if err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro ao buscar esquinas:", err)
		return nil
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
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
		log.Println("Ali:", err)
		return nil
	}

	return t
}
