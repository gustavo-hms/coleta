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

type Zona string

func ListaDeZonas() []Zona {
	return []Zona{
		ZonaSantoAmaro,
		ZonaCongonhasBerrini,
		ZonaIbirapuera,
		ZonaVilaMariana,
		ZonaAclimaçãoIpiranga,
		ZonaPaulista,
		ZonaPinheiros,
		ZonaButantã,
		ZonaLapaPompéia,
		ZonaPacaembu,
		ZonaCentro,
		ZonaTatuapé,
		ZonaNorte,
		ZonaSaúde,
	}
}
