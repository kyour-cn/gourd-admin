package test

import (
	"app/internal/config"
	"app/internal/global"
	"app/internal/initialize"
	"app/internal/orm/query"
	"encoding/json"
	"testing"
)

func TestDB(t *testing.T) {

	config.SetConfigPath("../configs")

	err := initialize.InitDatabase()
	if err != nil {
		t.Error(err)
		return
	}

	takeUser, err := query.Q.User.Take()
	if err != nil {
		t.Error(err)
		return
	}

	jsonStr, err := json.Marshal(takeUser)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("User %s", jsonStr)

	t.Log(global.GetDb("mysql"))

}
