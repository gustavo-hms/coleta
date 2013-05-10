package dao

import (
	"coleta/modelos"
	"database/sql"
	"fmt"
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

func (dao *ZonaDAO) FindAll() (zonas []*modelos.Zona, err error) {
	query := fmt.Sprintf("SELECT %s FROM zona", dao.fields)
	rows, err := dao.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		zona := new(modelos.Zona)
		rows.Scan(&zona.Id, &zona.Nome)
		zonas = append(zonas, zona)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return zonas, nil
}
