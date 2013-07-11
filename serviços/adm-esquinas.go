package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"coleta/modelos/validação"
	"html/template"
	"log"
	"net/http"
)

func init() {
	registrarSeguro("/adm/esquinas", Esquinas{})
}

type Esquinas struct{}

func (e Esquinas) Get(w http.ResponseWriter, r *http.Request) {
	esquina := new(modelos.Esquina)
	e.get(w, r, validação.NovaEsquinaComErros(esquina))
}

func (e Esquinas) get(
	w http.ResponseWriter,
	r *http.Request,
	esquina *validação.EsquinaComErros,
) {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println("Início da transação:", err)
		return
	}

	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		zonaDAO.Rollback()
		log.Println(err)
		return
	}

	zonaDAO.Commit()

	funcMap := template.FuncMap{
		"zonas": func() []zonaComSeleção {
			seleção := make([]zonaComSeleção, 0, len(zonas))
			for _, zona := range zonas {
				s := zonaComSeleção{Zona: *zona}
				if esquina != nil && esquina.Zona.Id == zona.Id {
					s.Selecionado = true
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
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-esquinas.html")
	if err != nil {
		log.Println("Ali:", err)
		return
	}

	err = t.ExecuteTemplate(w, "adm-esquinas.html", esquina)
	if err != nil {
		log.Println("Aqui:", err)
	}
}

func (e Esquinas) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
	}

	esquina := new(modelos.Esquina)
	esquina.Preencher(r.FormValue)
	erros := validação.ValidarEsquina(esquina)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		e.get(w, r, erros)
		return
	}

	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	if err := esquinaDAO.Save(esquina); err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro ao gravar esquina:", err)
		return
	}

	if err := esquinaDAO.Commit(); err != nil {
		esquinaDAO.Rollback()
		log.Println("Erro no commit:", err)
	}

	e.Get(w, r)
}
