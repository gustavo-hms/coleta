package dao

import (
	"coleta/modelos"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrRowsNotAffected = errors.New("Rows not affected.")
)

type EsquinaDAO struct {
	Tx     *sql.Tx
	fields string
}

func NewEsquinaDAO(tx *sql.Tx) *EsquinaDAO {
	return &EsquinaDAO{
		Tx:     tx,
		fields: "id, zona_id, cruzamento, localizacao",
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
	query := fmt.Sprintf("INSERT INTO esquina (%s) VALUES (DEFAULT, ?, ?, ?)",
		dao.fields)
	res, err := dao.Tx.Exec(query,
		esquina.Zona.Id,
		esquina.Cruzamento,
		esquina.Localização)
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
	query := "UPDATE esquina SET cruzamento = ?, localizacao = ?"
	row, err := dao.Tx.Exec(query, esquina.Cruzamento, esquina.Localização)
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

func (dao *EsquinaDAO) findById(id int) (*modelos.Esquina, error) {
	query := fmt.Sprintf("SELECT %s FROM esquina WHERE id = ?", dao.fields)
	row := dao.Tx.QueryRow(query, id)

	esquina := new(modelos.Esquina)
	esquina.Zona = modelos.Zona{}

	err := row.Scan(&esquina.Id,
		&esquina.Zona.Id,
		&esquina.Cruzamento,
		&esquina.Localização)

	if err != nil {
		return nil, err
	}

	return esquina, nil
}

func (dao *EsquinaDAO) Delete(id int) error {
	query := "DELETE FROM esquina WHERE id = ?"
	res, err := dao.Tx.Exec(query, id)
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
