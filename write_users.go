package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	registerQueryShort = `
		INSERT INTO user_profile (id, first_name)
			VALUES ($1, $2);`
)

func writeUsersInfinite() {
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

	var i int
	for {
		i++

		id := uuid.New().String()
		firstName := strconv.Itoa(i)

		_, err = conn.Exec(context.Background(), registerQueryShort, id, firstName)
		if err != nil {
			panic(err)
		}

		fmt.Println(i)
	}
}
