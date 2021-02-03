package main

import (
	"duomly.com/go-bank-backend/api"
	"duomly.com/go-bank-backend/migrations"
)

//todo : 실제 db 구동해서 테스트해보기
func main() {
	migrations.Migrate()
	api.StartApi()
}
