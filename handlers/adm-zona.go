package handlers

import (
	"coleta/modelos"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	http.HandleFunc("/adm/zona/", Zona)
}

func Zona(w http.ResponseWriter, r *http.Request) {
	print(idDaZona(r.URL), "\n")
}

func idDaZona(endereço *url.URL) string {
	return strings.SplitN(endereço.Path[10:], "/", 2)[0]
}

func jsonDasEsquinas(esquinas []modelos.Esquina) string {
	objetos := make([]string, len(esquinas))
	for k, esquina := range esquinas {
		objetos[k] = `{"cruzamento": "` + esquina.Cruzamento +
			`", "link": "/adm/esquina/` + fmt.Sprint(esquina.Id) + `"}`
	}

	return `[` + strings.Join(objetos, ",") + `]`
}
