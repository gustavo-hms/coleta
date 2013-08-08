package serviços

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	registrarSeguroComTransação("/adm/esquina/", new(AdmEsquina))
}

type AdmEsquina struct {
	turno   string
	esquina *modelos.Esquina
}

func idDaEsquina(endereço *url.URL) string {
	idx := strings.LastIndex(endereço.Path, "/")
	return endereço.Path[idx+1:]
}

func (e *AdmEsquina) Get(w http.ResponseWriter, r *http.Request, tx *dao.Tx) error {
	id := idDaEsquina(r.URL)
	esquinaDAO := dao.NewEsquinaDAO(tx)
	esquina, err := esquinaDAO.BuscaCompletaPorId(id)
	if err != nil {
		log.Println(err)
		return err
	}

	e.esquina = esquina
	turno := r.FormValue("turno")
	if turno != "" {
		e.turno = turno
		return e.exibirTurno(w, r)
	} else {
		return e.exibiçãoGeral(w, r)
	}
}

func (e AdmEsquina) exibirTurno(w http.ResponseWriter, r *http.Request) error {
	funcMap := template.FuncMap{
		"turno": func() string {
			return e.turno
		},
	}

	t, err := template.New("esquina").Funcs(funcMap).
		ParseFiles(
		config.Dados.DiretórioDasPáginas+"/adm-esquina-turno.html",
		config.Dados.DiretórioDasPáginas+"/adm-esquina-líder.html",
		config.Dados.DiretórioDasPáginas+"/adm-esquina-voluntário.html",
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-esquina-turno.html", e.esquina)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e AdmEsquina) exibiçãoGeral(w http.ResponseWriter, r *http.Request) error {
	funcMap := template.FuncMap{
		"plural": func(tamanho int) bool {
			return tamanho != 1
		},
	}

	t, err := template.New("esquina").Funcs(funcMap).
		ParseFiles(config.Dados.DiretórioDasPáginas + "/adm-esquina.html")
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, "adm-esquina.html", e.esquina)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
