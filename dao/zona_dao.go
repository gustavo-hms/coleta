package dao

import (
	"coleta/modelos"
	"database/sql"
	"fmt"
	"log"
)

const (
	OpçãoNenhuma              = 0
	OpçãoNãoFiltrarBloqueadas = 1 << iota
)

type ZonaDAO struct {
	*sql.Tx
	fields string
}

func NewZonaDAO(tx *sql.Tx) *ZonaDAO {
	return &ZonaDAO{
		Tx:     tx,
		fields: "id, nome",
	}
}

func (dao *ZonaDAO) Save(zona *modelos.Zona) error {
	if zona.Id == 0 {
		return dao.create(zona)
	} else {
		return dao.update(zona)
	}
}

func (dao *ZonaDAO) create(zona *modelos.Zona) error {
	query := fmt.Sprintf("INSERT INTO zona (%s) VALUES (DEFAULT, ?)", dao.fields)
	res, err := dao.Exec(query, zona.Nome)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	zona.Id = int(id)
	return nil
}

func (dao *ZonaDAO) update(zona *modelos.Zona) error {
	query := "UPDATE zona SET nome = ?"
	row, err := dao.Exec(query, zona.Nome)
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

func (dao *ZonaDAO) BuscaCompleta(id string) (*modelos.Zona, error) {
	query := fmt.Sprintf("SELECT %s FROM zona WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	zona := new(modelos.Zona)

	err := row.Scan(
		&zona.Id,
		&zona.Nome,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	esquinaDAO := NewEsquinaDAO(dao.Tx)
	zona.Esquinas, err = esquinaDAO.BuscaCompletaPorZona(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return zona, nil
}

func (dao *ZonaDAO) FindAll() (zonas []*modelos.Zona, err error) {
	return dao.FindAllWithOptions(OpçãoNenhuma)
}

func (dao *ZonaDAO) FindAllWithOptions(opções int) (zonas []*modelos.Zona, err error) {
	query := fmt.Sprintf("SELECT %s FROM zona", dao.fields)
	rows, err := dao.Query(query)
	if err != nil {
		return nil, err
	}

	var filtrarBloqueadas = true
	if opções&OpçãoNãoFiltrarBloqueadas != 0 {
		filtrarBloqueadas = false
	}

	for rows.Next() {
		zona := new(modelos.Zona)
		rows.Scan(&zona.Id, &zona.Nome)

		zonaBloqueada := false

		if filtrarBloqueadas {
			for _, bloqueada := range modelos.ZonasBloqueadas {
				if bloqueada == zona.Nome {
					zonaBloqueada = true
				}
			}
		}

		if !zonaBloqueada {
			zonas = append(zonas, zona)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return zonas, nil
}

func (dao *ZonaDAO) Delete(id int) error {
	query := "DELETE FROM zona WHERE id = ?"
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
