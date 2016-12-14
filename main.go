package main

import ()

func main() {
	psql := PostgresDB{}
	psql.Open()
	// psql.write()
	psql.Query(1)
}
