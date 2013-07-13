package dao

import (
	"coleta/modelos"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type VoluntarioDAO struct {
	*sql.Tx
	fields string
}

func NewVoluntarioDAO(tx *sql.Tx) *VoluntarioDAO {
	return &VoluntarioDAO{
		Tx: tx,
		fields: "id, zona_id, lider_id, esquina_id, nome_completo, telefone_residencial, " +
			"telefone_celular, operadora_celular, email, rg, cpf, idade, " +
			"como_soube_coleta_2013, cadastrado_em",
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
	var idDaZona int
	if voluntario.Zona != nil {
		idDaZona = voluntario.Zona.Id
	}

	var idDoLíder int
	if voluntario.Líder != nil {
		idDoLíder = voluntario.Líder.Id
	}

	var idDaEsquina int
	if voluntario.Esquina != nil {
		idDaEsquina = voluntario.Esquina.Id
	}

	query := fmt.Sprintf("INSERT INTO voluntario (%s) VALUES (DEFAULT, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dao.fields)
	res, err := dao.Exec(
		query,
		idDaZona,
		idDoLíder,
		idDaEsquina,
		voluntario.Nome,
		voluntario.TelefoneResidencial,
		voluntario.TelefoneCelular,
		voluntario.Operadora,
		voluntario.Email,
		voluntario.RG,
		voluntario.CPF,
		voluntario.Idade,
		voluntario.ComoSoube,
		voluntario.CadastradoEm,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	voluntario.Id = int(id)
	return dao.createTurnos(voluntario)
}

func (dao *VoluntarioDAO) createTurnos(voluntário *modelos.Voluntário) error {
	query := "INSERT INTO turnos_do_voluntario (voluntario_id, turno) VALUES (?, ?)"
	for _, turno := range voluntário.Turnos {
		_, err := dao.Exec(query, voluntário.Id, turno.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dao *VoluntarioDAO) update(voluntario *modelos.Voluntário) error {
	query := "UPDATE voluntario SET zona_id = ?, lider_id = ?, esquina_id = ?, " +
		"nome_completo = ?, telefone_residencial = ?, telefone_celular = ?, " +
		"operadora_celular = ?, email = ?, rg = ?, cpf = ?, idade = ?, " +
		"como_soube_coleta_2013 = ?, cadastrado_em = ? WHERE id = ?"
	row, err := dao.Exec(
		query,
		voluntario.Zona.Id,
		voluntario.Líder.Id,
		voluntario.Esquina.Id,
		voluntario.Nome,
		voluntario.TelefoneResidencial,
		voluntario.TelefoneCelular,
		voluntario.Operadora,
		voluntario.Email,
		voluntario.RG,
		voluntario.CPF,
		voluntario.Idade,
		voluntario.ComoSoube,
		voluntario.CadastradoEm,
		voluntario.Id,
	)

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

	if err := dao.deleteTurnos(voluntario.Id); err != nil {
		return err
	}

	return dao.createTurnos(voluntario)
}

func (dao *VoluntarioDAO) deleteTurnos(id int) error {
	query := "DELETE FROM turnos_do_voluntario WHERE voluntario_id = ?"
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

func (dao *VoluntarioDAO) FindByEmail(email string) (*modelos.Voluntário, error) {
	query := fmt.Sprintf("SELECT %s FROM voluntario WHERE email = ?", dao.fields)
	row := dao.QueryRow(query, email)

	voluntario := new(modelos.Voluntário)
	voluntario.Zona = new(modelos.Zona)
	voluntario.Líder = new(modelos.Líder)
	voluntario.Esquina = new(modelos.Esquina)
	var cadastradoEm string

	err := row.Scan(
		&voluntario.Id,
		&voluntario.Zona.Id,
		&voluntario.Líder.Id,
		&voluntario.Esquina.Id,
		&voluntario.Nome,
		&voluntario.TelefoneResidencial,
		&voluntario.TelefoneCelular,
		&voluntario.Operadora,
		&voluntario.Email,
		&voluntario.RG,
		&voluntario.CPF,
		&voluntario.Idade,
		&voluntario.ComoSoube,
		cadastradoEm,
	)

	if err != nil {
		return nil, err
	}

	if voluntario.CadastradoEm, err = time.Parse("2006-01-02 15:04:05", cadastradoEm); err != nil {
		return nil, err
	}

	voluntario.Turnos, err = dao.loadTurnos(voluntario.Id)
	if err != nil {
		return nil, err
	}

	return voluntario, nil
}

func (dao *VoluntarioDAO) FindById(id int) (*modelos.Voluntário, error) {
	query := fmt.Sprintf("SELECT %s FROM voluntario WHERE id = ?", dao.fields)
	row := dao.QueryRow(query, id)

	voluntario := new(modelos.Voluntário)
	voluntario.Zona = new(modelos.Zona)
	voluntario.Líder = new(modelos.Líder)
	voluntario.Esquina = new(modelos.Esquina)
	var cadastradoEm string

	err := row.Scan(
		&voluntario.Id,
		&voluntario.Zona.Id,
		&voluntario.Líder.Id,
		&voluntario.Esquina.Id,
		&voluntario.Nome,
		&voluntario.TelefoneResidencial,
		&voluntario.TelefoneCelular,
		&voluntario.Operadora,
		&voluntario.Email,
		&voluntario.RG,
		&voluntario.CPF,
		&voluntario.Idade,
		&voluntario.ComoSoube,
		&cadastradoEm,
	)

	if err != nil {
		return nil, err
	}

	if voluntario.CadastradoEm, err = time.Parse("2006-01-02 15:04:05", cadastradoEm); err != nil {
		return nil, err
	}

	esquinaDAO := NewEsquinaDAO(dao.Tx)
	esquina, _ := esquinaDAO.FindById(voluntario.Esquina.Id)
	if esquina != nil {
		voluntario.Esquina = esquina
	}

	voluntario.Turnos, err = dao.loadTurnos(id)
	if err != nil {
		return nil, err
	}

	return voluntario, nil
}

func (dao *VoluntarioDAO) BuscaPorEsquina(idDaEsquina int) (voluntários []modelos.Voluntário, err error) {
	query := fmt.Sprintf("SELECT %s FROM voluntario WHERE esquina_id = ? ORDER BY lider_id DESC", dao.fields)
	rows, err := dao.Query(query, idDaEsquina)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var voluntário modelos.Voluntário
		voluntário.Zona = new(modelos.Zona)
		voluntário.Líder = new(modelos.Líder)
		voluntário.Esquina = new(modelos.Esquina)
		var cadastradoEm string

		err := rows.Scan(
			&voluntário.Id,
			&voluntário.Zona.Id,
			&voluntário.Líder.Id,
			&voluntário.Esquina.Id,
			&voluntário.Nome,
			&voluntário.TelefoneResidencial,
			&voluntário.TelefoneCelular,
			&voluntário.Operadora,
			&voluntário.Email,
			&voluntário.RG,
			&voluntário.CPF,
			&voluntário.Idade,
			&voluntário.ComoSoube,
			&cadastradoEm,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		if voluntário.CadastradoEm, err = time.Parse("2006-01-02 15:04:05", cadastradoEm); err != nil {
			log.Println(err)
			return nil, err
		}

		voluntários = append(voluntários, voluntário)
	}

	for i, _ := range voluntários {
		voluntários[i].Turnos, err = dao.loadTurnos(voluntários[i].Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return voluntários, nil
}

func (dao *VoluntarioDAO) loadTurnos(idDoVoluntario int) ([]modelos.Turno, error) {
	query := "SELECT turno FROM turnos_do_voluntario WHERE voluntario_id = ?"
	rows, err := dao.Query(query, idDoVoluntario)
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

func (dao *VoluntarioDAO) QtdPorEsquina(id int) (qtd int, err error) {
	query := "SELECT count(*) FROM voluntario WHERE esquina_id = ?"
	row := dao.QueryRow(query, id)

	err = row.Scan(&qtd)
	return
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
