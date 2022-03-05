package user

import (
	"basic-crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.Write([]byte("Falha ao ler corpo da requisição"))
		return
	}
	var user user
	if error = json.Unmarshal(requestBody, &user); error != nil {
		w.Write([]byte("Erro ao converter JSON para struct"))
		return
	}
	dbConnection, error := database.Connect()
	if error != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}
	defer dbConnection.Close()
	statement, error := dbConnection.Prepare("INSERT INTO usuarios (nome, email) VALUES (?, ?);")
	if error != nil {
		w.Write([]byte("Erro ao preparar statement"))
		return
	}
	defer statement.Close()
	insertResult, error := statement.Exec(user.Name, user.Email)
	if error != nil {
		w.Write([]byte("Erro ao executar statement"))
		return
	}
	userId, error := insertResult.LastInsertId()
	if error != nil {
		w.Write([]byte("Erro ao obter id do usuário"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário criado com sucesso: %d", userId)))
}
