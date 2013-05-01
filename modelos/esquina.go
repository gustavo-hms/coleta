package modelos

import (
	"net/url"
)

type Esquina struct {
	Cruzamento  string
	Localização url.URL
	Zona        Zona
}
