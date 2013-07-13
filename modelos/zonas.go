package modelos

const (
	ZonaSantoAmaro         = "Santo Amaro"
	ZonaMoemaBerrini       = "Moema/Berrini"
	ZonaIbirapuera         = "Ibirapuera"
	ZonaVilaMariana        = "Vila Mariana"
	ZonaAclimaçãoIpiranga  = "Aclimação/Ipiranga"
	ZonaPaulista           = "Paulista"
	ZonaPinheiros          = "Pinheiros"
	ZonaButantã            = "Butantã"
	ZonaLapaPompéia        = "Lapa/Pompéia"
	ZonaPacaembuBarraFunda = "Pacaembu/Barra Funda"
	ZonaCentro             = "Centro"
	ZonaTatuapé            = "Tatuapé"
	ZonaSantanaNorte       = "Santana/Norte"
	ZonaSaúde              = "Saúde"
)

var ZonasBloqueadas = []string{
	ZonaPaulista,
}

type Zona struct {
	Id        int
	Nome      string
	Bloqueada bool
	Esquinas  []Esquina
}
