package main

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"fmt"
	"log"
	"os"
)

const (
	MáximoDeLíderes                    = 1
	MáximoDeVoluntáriosAltaPrioridade  = 17
	MáximoDeVoluntáriosBaixaPrioridade = 13
)

func init() {
	f, err := os.OpenFile("alocação.log", os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.Dados.Banco.Base = "base"
	config.Dados.Banco.Usuário = "usuário"
	config.Dados.Banco.Host = "127.0.0.1"
	config.Dados.Banco.Senha = ""

	if err := dao.Conn(); err != nil {
		fmt.Println("Não foi possível conectar-se ao banco. Erro:", err)
		os.Exit(1)
	}
}

func uso() string {
	return "Uso: " + os.Args[0] + " <arquivo de configuração>"
}

func main() {
	transação, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	tx := &dao.Tx{transação}

	líderes := carregarLíderes(tx)
	for _, líder := range líderes {
		esquina := esquinaLivreParaLíder(tx, líder)
		if esquina != nil {
			esquina.AcrescentarLíder(líder)
			líder.Esquina = esquina
			gravarLíder(tx, líder)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		os.Exit(1)
	}

	transação, err = dao.DB.Begin()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	tx = &dao.Tx{transação}

	voluntários := carregarVoluntários(tx)
	for _, voluntário := range voluntários {
		esquina := esquinaLivreParaVoluntário(tx, voluntário)
		if esquina != nil {
			esquina.AcrescentarVoluntário(voluntário)
			voluntário.Esquina = esquina
			gravarVoluntário(tx, voluntário)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		tx.Rollback()
		os.Exit(1)
	}
}

func carregarLíderes(tx *dao.Tx) []modelos.Líder {
	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.BuscaPorEsquina(0)
	if err != nil {
		log.Println(err)
	}

	return líderes
}

func esquinaLivreParaLíder(tx *dao.Tx, líder modelos.Líder) *modelos.Esquina {
	esquinaDAO := dao.NewEsquinaDAO(tx)
	id := fmt.Sprint(líder.Zona.Id)

	esquinas, err := esquinaDAO.BuscaCompletaPorZona(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	for i, _ := range esquinas {
		livre := true
		for _, turno := range líder.Turnos {
			livre = livre && cabeLíder(esquinas[i], turno)
		}

		if livre {
			return &esquinas[i]
		}
	}

	return nil
}

func gravarLíder(tx *dao.Tx, líder modelos.Líder) {
	líderDAO := dao.NewLiderDAO(tx)
	err := líderDAO.Save(&líder)
	if err != nil {
		log.Println(err)
	}
}

func carregarVoluntários(tx *dao.Tx) []modelos.Voluntário {
	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.BuscaPorEsquina(0)
	if err != nil {
		log.Println(err)
	}

	return voluntários
}

func esquinaLivreParaVoluntário(tx *dao.Tx, voluntário modelos.Voluntário) *modelos.Esquina {
	zona := voluntário.Zona

	if voluntário.Líder.Id != 0 {
		if voluntário.Líder.Esquina == nil {
			líderDAO := dao.NewLiderDAO(tx)
			líder, err := líderDAO.FindById(voluntário.Líder.Id)
			if err != nil {
				log.Println(err)
				return nil
			}

			voluntário.Líder = líder
		}

		if voluntário.Líder.Esquina.Id == 0 {
			// O líder ainda não foi alocado; esperar mais uma hora
			return nil
		}

		// Tentar primeiro alocá-lo na esquina do líder
		esquinaDAO := dao.NewEsquinaDAO(tx)
		id := fmt.Sprint(voluntário.Líder.Esquina.Id)
		esquina, err := esquinaDAO.BuscaCompletaPorId(id)
		if err != nil {
			log.Println(err)
			return nil
		}

		livre := true
		for _, turno := range voluntário.Turnos {
			livre = livre && cabeVoluntário(*esquina, turno)
		}

		if livre {
			return esquina
		}

		// A esquina não estava livre, alocar então pela zona do líder
		zona = voluntário.Líder.Zona
	}

	esquinaDAO := dao.NewEsquinaDAO(tx)
	id := fmt.Sprint(zona.Id)

	esquinas, err := esquinaDAO.BuscaCompletaPorZona(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	for i, _ := range esquinas {
		livre := true
		for _, turno := range voluntário.Turnos {
			livre = livre && cabeVoluntário(esquinas[i], turno)
		}

		if livre {
			return &esquinas[i]
		}
	}

	return nil
}

func gravarVoluntário(tx *dao.Tx, voluntário modelos.Voluntário) {
	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	err := voluntárioDAO.Save(&voluntário)
	if err != nil {
		log.Println(err)
	}
}

func cabeLíder(esquina modelos.Esquina, turno modelos.Turno) bool {
	return len(esquina.Participantes[turno.Id].Líderes) < MáximoDeLíderes
}

func cabeVoluntário(esquina modelos.Esquina, turno modelos.Turno) bool {
	if esquina.Prioridade == modelos.AltaPrioridade {
		return len(esquina.Participantes[turno.Id].Voluntários) <
			MáximoDeVoluntáriosAltaPrioridade
	}

	return len(esquina.Participantes[turno.Id].Voluntários) <
		MáximoDeVoluntáriosBaixaPrioridade
}
