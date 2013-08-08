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
	registrarSeguroComTransação("/adm/esquinas", Esquinas{})
}

type Esquinas struct{}

func (e Esquinas) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	esquina := new(modelos.Esquina)
	return e.get(w, r, tx, validação.NovaEsquinaComErros(esquina))
}

func (e Esquinas) get(
	w http.ResponseWriter,
	r *http.Request,
	tx *dao.Tx,
	esquina *validação.EsquinaComErros,
) error {
	zonaDAO := dao.NewZonaDAO(tx)
	zonas, err := zonaDAO.FindAllWithOptions(dao.OpçãoNãoFiltrarBloqueadas)
	if err != nil {
		log.Println(err)
		return err
	}

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
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-esquinas.html", esquina)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e Esquinas) Post(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("Erro ao analisar formulário:", err)
		return err
	}

	esquina := new(modelos.Esquina)
	esquina.Preencher(r.FormValue)
	erros := validação.ValidarEsquina(esquina)

	if erros != nil {
		w.WriteHeader(http.StatusBadRequest)
		return e.get(w, r, tx, erros)
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	if err := esquinaDAO.Save(esquina); err != nil {
		log.Println("Erro ao gravar esquina:", err)
		return err
	}

	e.Get(w, r, tx)
	return nil
}
