package main

import (
	"context"
	"encoding/csv"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	filePath      = "users.csv"
	dbString      = "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"
	registerQuery = `
		INSERT INTO user_profile (id, first_name, second_name, birthdate, biography, city, password_hash)
			VALUES ($1, $2, $3, $4, 'bio', $5, $6);`
)

func generateUsers() {
	conn, err := pgx.Connect(context.Background(), dbString)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), 14)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	r := csv.NewReader(f)
	var user []string
	for {
		user, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		id := uuid.New().String()

		name := strings.Fields(user[0])
		firstName := name[0]
		secondName := name[1]

		_, err = conn.Exec(context.Background(), registerQuery, id, firstName, secondName, user[1], user[2], passwordHash)
		if err != nil {
			panic(err)
		}
	}
}
