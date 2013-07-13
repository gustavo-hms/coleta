package main

import (
	"coleta/config"
	"coleta/dao"
	"coleta/modelos"
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	MáximoDeLíderes                    = 1
	MáximoDeVoluntáriosAltaPrioridade  = 1
	MáximoDeVoluntáriosBaixaPrioridade = 13
)

func init() {
	f, _ := os.Open("alocação.log")
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 2 {
		fmt.Println(uso())
		os.Exit(1)
	}

	if err := config.Ler(os.Args[1]); err != nil {
		fmt.Printf("Não foi possível ler o arquivo de configuração %s. Erro: %s", os.Args[1], err)
		os.Exit(1)
	}

	if err := dao.Conn(); err != nil {
		fmt.Println("Não foi possível conectar-se ao banco. Erro:", err)
		os.Exit(1)
	}
}

func uso() string {
	return "Uso: " + os.Args[0] + " <arquivo de configuração>"
}

func main() {
	tx, err := dao.DB.Begin()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	líderes := carregarLíderes(tx)
	for _, líder := range líderes {
		esquina := esquinaLivreParaLíder(tx, líder)
		if esquina != nil {
			print("Colocando líder ", líder.Id, " na esquina ", esquina.Id, "\n")
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

	tx, err = dao.DB.Begin()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	voluntários := carregarVoluntários(tx)
	for _, voluntário := range voluntários {
		esquina := esquinaLivreParaVoluntário(tx, voluntário)
		if esquina != nil {
			print("Colocando voluntário ", voluntário.Id, " na esquina ", esquina.Id, "\n")
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

func carregarLíderes(tx *sql.Tx) []modelos.Líder {
	líderDAO := dao.NewLiderDAO(tx)
	líderes, err := líderDAO.BuscaPorEsquina(0)
	if err != nil {
		log.Println(err)
	}

	return líderes
}

func esquinaLivreParaLíder(tx *sql.Tx, líder modelos.Líder) *modelos.Esquina {
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

func gravarLíder(tx *sql.Tx, líder modelos.Líder) {
	líderDAO := dao.NewLiderDAO(tx)
	err := líderDAO.Save(&líder)
	if err != nil {
		log.Println(err)
	}
}

func carregarVoluntários(tx *sql.Tx) []modelos.Voluntário {
	voluntárioDAO := dao.NewVoluntarioDAO(tx)
	voluntários, err := voluntárioDAO.BuscaPorEsquina(0)
	if err != nil {
		log.Println(err)
	}

	return voluntários
}

func esquinaLivreParaVoluntário(tx *sql.Tx, voluntário modelos.Voluntário) *modelos.Esquina {
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

func gravarVoluntário(tx *sql.Tx, voluntário modelos.Voluntário) {
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
