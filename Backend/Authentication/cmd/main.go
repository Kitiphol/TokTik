package main

import (
	
    "Auth/internal/database"
    "Auth/internal/routes"
)

// psql -U dev -d dbp2 -c "SELECT * FROM users;"
// psql -U dev -d dbp2 -c "SELECT * FROM videos;"
func main() {

    database.InitDB()
	


    r := routes.Setup()
    r.Run() // default listens on :8080
}
