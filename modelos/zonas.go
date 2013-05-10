package modelos

const (
	ZonaSantoAmaro        = "Santo Amaro"
	ZonaCongonhasBerrini  = "Congonhas/Berrini"
	ZonaIbirapuera        = "Ibirapuera"
	ZonaVilaMariana       = "Vila Mariana"
	ZonaAclimaçãoIpiranga = "Aclimação/Ipiranga"
	ZonaPaulista          = "Paulista"
	ZonaPinheiros         = "Pinheiros"
	ZonaButantã           = "Butantã"
	ZonaLapaPompéia       = "Lapa/Pompéia"
	ZonaPacaembu          = "Pacaembu"
	ZonaCentro            = "Centro"
	ZonaTatuapé           = "Tatuapé"
	ZonaNorte             = "Norte"
	ZonaSaúde             = "Saúde"
)

type Zona struct {
	Id   int
	Nome string
}
