package dao

import (
	"coleta/modelos"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrRowsNotAffected = errors.New("Rows not affected.")
)

type EsquinaDAO struct {
	*sql.Tx
	fields string
}

func NewEsquinaDAO(tx *sql.Tx) *EsquinaDAO {
	return &EsquinaDAO{
		Tx:     tx,
		fields: "id, zona_id, cruzamento, localizacao, prioridade",
	}
}

func (dao *EsquinaDAO) Save(esquina *modelos.Esquina) error {
	if esquina.Id == 0 {
		return dao.create(esquina)
	} else {
		return dao.update(esquina)
	}
}

func (dao *EsquinaDAO) create(esquina *modelos.Esquina) error {
	query := fmt.Sprintf("INSERT INTO esquina (%s) VALUES (DEFAULT, ?, ?, ?, ?)",
		dao.fields)
	res, err := dao.Exec(
		query,
		esquina.Zona.Id,
		esquina.Cruzamento,
		esquina.Localização,
		esquina.Prioridade,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	esquina.Id = int(id)

	return nil
}

func (dao *EsquinaDAO) update(esquina *modelos.Esquina) error {
	query := "UPDATE esquina SET cruzamento = ?, localizacao = ?, prioridade = ?, zona_id = ?"
	row, err := dao.Exec(query, esquina.Cruzamento, esquina.Localização, esquina.Prioridade, esquina.Zona.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRowsNotAffected
	}

	return nil
}

func (dao *EsquinaDAO) FindById(id int) (*modelos.Esquina, error) {
	query := fmt.Sprintf("SELECT %s FROM esquina WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	esquina := new(modelos.Esquina)
	esquina.Zona = modelos.Zona{}

	err := row.Scan(
		&esquina.Id,
		&esquina.Zona.Id,
		&esquina.Cruzamento,
		&esquina.Localização,
		&esquina.Prioridade,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return esquina, nil
}

func (dao *EsquinaDAO) BuscaCompletaPorId(id string) (*modelos.Esquina, error) {
	query := fmt.Sprintf("SELECT %s FROM esquina WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	esquina := new(modelos.Esquina)
	esquina.Zona = modelos.Zona{}

	err := row.Scan(
		&esquina.Id,
		&esquina.Zona.Id,
		&esquina.Cruzamento,
		&esquina.Localização,
		&esquina.Prioridade,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	esquina.Participantes = make(map[string]modelos.Participantes)

	líderDAO := NewLiderDAO(dao.Tx)
	líderes, err := líderDAO.BuscaPorEsquina(esquina.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dao.preencherLíderes(esquina.Participantes, líderes)

	voluntárioDAO := NewVoluntarioDAO(dao.Tx)
	voluntários, err := voluntárioDAO.BuscaPorEsquina(esquina.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dao.preencherVoluntários(esquina.Participantes, voluntários)

	return esquina, nil
}

func (dao EsquinaDAO) preencherLíderes(
	mapa map[string]modelos.Participantes,
	líderes []modelos.Líder,
) {
	for _, líder := range líderes {
		for _, turno := range líder.Turnos {
			if _, ok := mapa[turno.Id]; ok {
				participantes := mapa[turno.Id]
				participantes.Líderes = append(participantes.Líderes, líder)
				mapa[turno.Id] = participantes

			} else {
				novoLíder := make([]modelos.Líder, 1)
				novoLíder[0] = líder
				mapa[turno.Id] = modelos.Participantes{Líderes: novoLíder}
			}
		}
	}
}

func (dao EsquinaDAO) preencherVoluntários(
	mapa map[string]modelos.Participantes,
	voluntários []modelos.Voluntário,
) {
	for _, voluntário := range voluntários {
		for _, turno := range voluntário.Turnos {
			if _, ok := mapa[turno.Id]; ok {
				participantes := mapa[turno.Id]
				participantes.Voluntários = append(participantes.Voluntários, voluntário)
				mapa[turno.Id] = participantes

			} else {
				novoVoluntário := make([]modelos.Voluntário, 1)
				novoVoluntário[0] = voluntário
				mapa[turno.Id] = modelos.Participantes{Voluntários: novoVoluntário}
			}
		}
	}
}

func (dao *EsquinaDAO) BuscarPorZona(idDaZona string) ([]modelos.Esquina, error) {
	query := fmt.Sprintf("SELECT %s FROM esquina WHERE zona_id = ?", dao.fields)
	rows, err := dao.Query(query, idDaZona)
	if err != nil {
		return nil, err
	}

	esquinas := make([]modelos.Esquina, 0)
	for rows.Next() {
		var esquina modelos.Esquina
		rows.Scan(
			&esquina.Id,
			&esquina.Zona.Id,
			&esquina.Cruzamento,
			&esquina.Localização,
			&esquina.Prioridade,
		)

		esquinas = append(esquinas, esquina)
	}

	return esquinas, nil
}

func (dao *EsquinaDAO) BuscaCompletaPorZona(idDaZona string) ([]modelos.Esquina, error) {
	query := fmt.Sprintf("SELECT %s FROM esquina WHERE zona_id = ?", dao.fields)
	rows, err := dao.Query(query, idDaZona)
	if err != nil {
		return nil, err
	}

	esquinas := make([]modelos.Esquina, 0)
	for rows.Next() {
		var esquina modelos.Esquina
		rows.Scan(
			&esquina.Id,
			&esquina.Zona.Id,
			&esquina.Cruzamento,
			&esquina.Localização,
			&esquina.Prioridade,
		)

		esquinas = append(esquinas, esquina)
	}

	líderDAO := NewLiderDAO(dao.Tx)
	voluntárioDAO := NewVoluntarioDAO(dao.Tx)
	for i, _ := range esquinas {
		esquinas[i].QtdDeLíderes, err = líderDAO.QtdPorEsquina(esquinas[i].Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		esquinas[i].QtdDeVoluntários, err =
			voluntárioDAO.QtdPorEsquina(esquinas[i].Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return esquinas, nil
}

func (dao *EsquinaDAO) Delete(id int) error {
	query := "DELETE FROM esquina WHERE id = ?"
	res, err := dao.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRowsNotAffected
	}

	return err
}
