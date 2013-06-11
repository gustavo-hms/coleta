package config

import (
	"encoding/json"
	"io/ioutil"
)

var Dados Config

func Ler(caminho string) error {
	confBytes, err := ioutil.ReadFile(caminho)
	if err != nil {
		return err
	}

	err = json.Unmarshal(confBytes, &Dados)
	if err != nil {
		return err
	}

	return nil
}

type Config struct {
	Banco               Banco
	DiretórioDasPáginas string
}

type Banco struct {
	Usuário string
	Senha   string
	Host    string
	Base    string
}
