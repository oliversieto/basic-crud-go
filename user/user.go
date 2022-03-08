package user

import (
	"basic-crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Create(w http.ResponseWriter, r *http.Request) {
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

func GetAll(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := database.Connect()

	if err != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	defer dbConnection.Close()

	queryResult, err := dbConnection.Query("SELECT * FROM usuarios;")

	if err != nil {
		w.Write([]byte("Erro ao buscar usuários"))
		return
	}

	defer queryResult.Close()

	var users []user

	for queryResult.Next() {
		var user user

		if err := queryResult.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Erro ao scanear o usuário"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Erro ao converter JSON"))
		return
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		w.Write([]byte("Erro ao converter id para uint32"))
		return
	}

	dbConnection, err := database.Connect()

	if err != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados"))
		return
	}

	defer dbConnection.Close()

	queryResult, err := dbConnection.Query("SELECT * FROM usuarios WHERE id = ?;", ID)
	if err != nil {
		w.Write([]byte("Erro ao buscar usuário"))
		return
	}

	defer queryResult.Close()

	var user user

	if !queryResult.Next() {
		w.Write([]byte("Usuário não encontrado"))
		return
	}

	if err := queryResult.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		w.Write([]byte("Erro ao scanear o usuário"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Erro ao converter JSON"))
		return
	}
}
