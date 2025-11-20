package main

import (
	
)

func main(){
	r := router.SetupRouter()
// start server on port 8080
r.Run(":8080")
}