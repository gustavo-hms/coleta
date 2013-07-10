package dao

import (
	"coleta/modelos"
	"database/sql"
	"fmt"
)

type VoluntarioDAO struct {
	*sql.Tx
	fields string
}

func NewVoluntarioDAO(tx *sql.Tx) *VoluntarioDAO {
	return &VoluntarioDAO{
		Tx: tx,
		fields: "id, zona_id, lider_id, nome_completo, telefone_residencial, " +
			"telefone_celular, operadora_celular, email, turno, " +
			"como_soube_coleta_2013",
	}
}

func (dao *VoluntarioDAO) Save(voluntario *modelos.Voluntário) error {
	if voluntario.Id == 0 {
		return dao.create(voluntario)
	} else {
		return dao.update(voluntario)
	}
}

func (dao *VoluntarioDAO) create(voluntario *modelos.Voluntário) error {
	query := fmt.Sprintf("INSERT INTO voluntario (%s) VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dao.fields)
	res, err := dao.Exec(query,
		voluntario.Zona.Id,
		voluntario.Lider.Id,
		voluntario.Nome,
		voluntario.TelefoneResidencial,
		voluntario.TelefoneCelular,
		voluntario.OperadoraCelular,
		voluntario.Email,
		voluntario.Turno,
		voluntario.ComoSoubeColeta2013)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	voluntario.Id = int(id)
	return nil

}

func (dao *VoluntarioDAO) update(voluntario *modelos.Voluntário) error {
	query := "UPDATE voluntario SET zona_id = ?, lider_id = ?, " +
		"nome_completo = ?, telefone_residencial = ?, telefone_celular = ?, " +
		"operadora_celular = ?, email = ?, turno = ?, como_soube_coleta_2013 = ?"
	row, err := dao.Exec(query,
		voluntario.Zona.Id,
		voluntario.Lider.Id,
		voluntario.Nome,
		voluntario.TelefoneResidencial,
		voluntario.TelefoneCelular,
		voluntario.OperadoraCelular,
		voluntario.Email,
		voluntario.Turno,
		voluntario.ComoSoubeColeta2013)

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

func (dao *VoluntarioDAO) FindById(id int) (*modelos.Voluntário, error) {
	query := fmt.Sprintf("SELECT %s FROM voluntario WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	voluntario := new(modelos.Voluntário)
	voluntario.Zona = new(modelos.Zona)
	voluntario.Lider = new(modelos.Líder)

	err := row.Scan(&voluntario.Id,
		&voluntario.Zona.Id,
		&voluntario.Lider.Id,
		&voluntario.Nome,
		&voluntario.TelefoneResidencial,
		&voluntario.TelefoneCelular,
		&voluntario.OperadoraCelular,
		&voluntario.Email,
		&voluntario.Turno,
		&voluntario.ComoSoubeColeta2013)

	if err != nil {
		return nil, err
	}

	return voluntario, nil
}

func (dao *VoluntarioDAO) Delete(id int) error {
	query := "DELETE FROM voluntario WHERE id = ?"
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
