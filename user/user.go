package user

import (
	"basic-crud/database"
	"encoding/json"
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
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		handlerErrorRequest(w, "Falha ao ler corpo da requisição", http.StatusInternalServerError)
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		handlerErrorRequest(w, "Erro ao converter JSON para struct", http.StatusInternalServerError)
		return
	}

	dbConnection, err := database.Connect()

	if err != nil {
		handlerErrorRequest(w, "Erro ao conectar com o banco de dados", http.StatusInternalServerError)
		return
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare("INSERT INTO usuarios (nome, email) VALUES (?, ?);")

	if err != nil {
		handlerErrorRequest(w, "Erro ao preparar statement", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	insertResult, err := statement.Exec(user.Name, user.Email)

	if err != nil {
		handlerErrorRequest(w, "Erro ao executar statement", http.StatusInternalServerError)
		return
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		handlerErrorRequest(w, "Erro ao obter id do usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": userId})

}

func GetAll(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := database.Connect()

	if err != nil {
		handlerErrorRequest(w, "Erro ao conectar com o banco de dados", http.StatusInternalServerError)
		return
	}

	defer dbConnection.Close()

	queryResult, err := dbConnection.Query("SELECT * FROM usuarios;")

	if err != nil {
		handlerErrorRequest(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}

	defer queryResult.Close()

	var users []user

	for queryResult.Next() {
		var user user

		if err := queryResult.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			handlerErrorRequest(w, "Erro ao scanear o usuário", http.StatusInternalServerError)
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		handlerErrorRequest(w, "Erro ao converter JSON", http.StatusInternalServerError)
		return
	}
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		handlerErrorRequest(w, "Erro ao converter id para uint32", http.StatusInternalServerError)
		return
	}

	dbConnection, err := database.Connect()

	if err != nil {
		handlerErrorRequest(w, "Erro ao conectar com o banco de dados", http.StatusInternalServerError)
		return
	}

	defer dbConnection.Close()

	queryResult, err := dbConnection.Query("SELECT * FROM usuarios WHERE id = ?;", ID)

	if err != nil {
		handlerErrorRequest(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	defer queryResult.Close()

	var user user

	if !queryResult.Next() {
		handlerErrorRequest(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	if err := queryResult.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		handlerErrorRequest(w, "Erro ao scanear o usuário", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		handlerErrorRequest(w, "Erro ao scanear o usuário", http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)

	if err != nil {
		handlerErrorRequest(w, "Erro ao converter parâmetro para inteiro", http.StatusInternalServerError)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		handlerErrorRequest(w, "Falha ao ler corpo da requisição", http.StatusInternalServerError)
		return
	}

	var user user

	if err = json.Unmarshal(requestBody, &user); err != nil {
		handlerErrorRequest(w, "Erro ao converter usuário para struct", http.StatusInternalServerError)
		return
	}

	dbConnection, err := database.Connect()

	if err != nil {
		handlerErrorRequest(w, "Erro ao conectar com o banco de dados", http.StatusInternalServerError)
		return
	}

	defer dbConnection.Close()

	statement, err := dbConnection.Prepare("UPDATE usuarios SET nome = ?, email = ? WHERE id = ?;")

	if err != nil {
		handlerErrorRequest(w, "Erro ao preparar statement", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Email, ID); err != nil {
		handlerErrorRequest(w, "Erro ao executar statement", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func handlerErrorRequest(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}
