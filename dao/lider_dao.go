package dao

import (
	"coleta/modelos"
	"database/sql"
	"errors"
	"fmt"
)

type LiderDAO struct {
	*sql.Tx
	fields string
}

func NewEsquinaDAO(tx *sql.Tx) *EsquinaDAO {
	return &LiderDAO{
		Tx: tx,
		fields: "id, zona_id, esquina_id, cadastrado_em, nome_completo, " +
			"telefone_residencial, telefone_celular, operadora_celular, " +
			"email, turno",
	}
}

func (dao *LiderDAO) Save(lider *modelos.Líder) error {
	if lider.Id == 0 {
		dao.create(lider)
	} else {
		dao.update(lider)
	}
}

func (dao *LiderDAO) create(lider *modelos.Líder) error {
	query := fmt.Sprintf("INSERT INTO lider (%s) VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?",
		dao.fields)
	res, err := dao.Exec(query,
		lider.Zona.Id,
		lider.Esquina.Id,
		lider.CadastradoEm,
		lider.Nome,
		lider.TelefoneResidencial,
		lider.TelefoneCelular,
		lider.Operadora,
		lider.Email,
		lider.Turnos)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	lider.Id = int(id)
	return nil
}

func (dao *LiderDAO) update(lider *modelos.Líder) error {
	query := "UPDATE lider SET zona_id = ?, esquina_id = ?, cadastrado_em = ?, " +
		"nome_completo,  = ?, telefone_residencial = ?, telefone_celular = ?, " +
		"operadora_celular = ?, email = ?, turno = ?"
	row, err := dao.Exec(query,
		lider.Zona.Id,
		lider.Esquina.Id,
		lider.CadastradoEm,
		lider.Nome,
		lider.TelefoneResidencial,
		lider.TelefoneCelular,
		lider.Operadora,
		lider.Email,
		lider.Turnos)
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

func (dao *LiderDAO) findById(id int) (*modelos.Zona, error) {
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
		&lider.Email,
		&lider.Turnos)

	if err != nil {
		return nil, err
	}

	return lider, nil
}

func (dao *LiderDAO) Delete(id int) error {
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
