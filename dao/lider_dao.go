package dao

import (
	"coleta/modelos"
	"database/sql"
	"fmt"
	"log"
)

type LiderDAO struct {
	*sql.Tx
	fields string
}

func NewLiderDAO(tx *sql.Tx) *LiderDAO {
	return &LiderDAO{
		Tx: tx,
		fields: "id, zona_id, esquina_id, cadastrado_em, nome_completo, " +
			"telefone_residencial, telefone_celular, operadora_celular, " +
			"email",
	}
}

func (dao *LiderDAO) Save(lider *modelos.Líder) error {
	if lider.Id == 0 {
		return dao.create(lider)
	} else {
		return dao.update(lider)
	}
}

func (dao *LiderDAO) create(lider *modelos.Líder) error {
	var idDaZona int
	if lider.Zona != nil {
		idDaZona = lider.Zona.Id
	}

	var idDaEsquina int
	if lider.Esquina != nil {
		idDaEsquina = lider.Esquina.Id
	}

	query := fmt.Sprintf("INSERT INTO lider (%s) VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?)",
		dao.fields)
	res, err := dao.Exec(query,
		idDaZona,
		idDaEsquina,
		lider.CadastradoEm,
		lider.Nome,
		lider.TelefoneResidencial,
		lider.TelefoneCelular,
		lider.Operadora,
		lider.Email)
	if err != nil {
		log.Printf("%s, %v, %v, %v, %v, %v, %v, %v, %v\n", query, idDaZona,
			idDaEsquina,
			lider.CadastradoEm,
			lider.Nome,
			lider.TelefoneResidencial,
			lider.TelefoneCelular,
			lider.Operadora,
			lider.Email)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	lider.Id = int(id)

	return dao.createTurnos(lider)
}

func (dao *LiderDAO) createTurnos(lider *modelos.Líder) error {
	query := "INSERT INTO turnos_do_lider (lider_id, turno) VALUES (?, ?)"
	for _, turno := range lider.Turnos {
		_, err := dao.Exec(query, lider.Id, turno.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *LiderDAO) update(lider *modelos.Líder) error {
	query := "UPDATE lider SET zona_id = ?, esquina_id = ?, cadastrado_em = ?, " +
		"nome_completo,  = ?, telefone_residencial = ?, telefone_celular = ?, " +
		"operadora_celular = ?, email = ?"
	row, err := dao.Exec(query,
		lider.Zona.Id,
		lider.Esquina.Id,
		lider.CadastradoEm,
		lider.Nome,
		lider.TelefoneResidencial,
		lider.TelefoneCelular,
		lider.Operadora,
		lider.Email)
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

	if err := dao.deleteTurnos(lider.Id); err != nil {
		return err
	}

	return dao.createTurnos(lider)
}

func (dao *LiderDAO) deleteTurnos(id int) error {
	query := "DELETE FROM turnos_do_lider WHERE lider_id = ?"
	res, err := dao.Exec(query, id)

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRowsNotAffected
	}

	return err
}

func (dao *LiderDAO) FindById(id int) (*modelos.Líder, error) {
	query := fmt.Sprintf("SELECT %s FROM lider WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	lider := new(modelos.Líder)
	lider.Zona = new(modelos.Zona)
	lider.Esquina = new(modelos.Esquina)

	err := row.Scan(&lider.Id,
		&lider.Zona.Id,
		&lider.Esquina.Id,
		&lider.CadastradoEm,
		&lider.Nome,
		&lider.TelefoneResidencial,
		&lider.TelefoneCelular,
		&lider.Operadora,
		&lider.Email)

	if err != nil {
		return nil, err
	}

	lider.Turnos, err = dao.loadTurnos(id)
	if err != nil {
		return nil, err
	}

	return lider, nil
}

func (dao *LiderDAO) loadTurnos(liderId int) ([]modelos.Turno, error) {
	query := "SELECT turno FROM turnos_do_lider WHERE lider_id = ?"
	rows, err := dao.Query(query)
	if err != nil {
		return nil, err
	}

	turnos := make([]modelos.Turno, 0)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		turnos = append(turnos, modelos.TurnoComId(id))
	}

	return turnos, nil
}

func (dao *LiderDAO) Delete(id int) error {
	if err := dao.deleteTurnos(id); err != nil {
		return err
	}

	query := "DELETE FROM lider WHERE id = ?"
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
