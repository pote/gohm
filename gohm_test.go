package gohm

import(
	`github.com/pote/redisurl`
	`testing`
)

type User struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	UUID  string `ohm:"name unique"`
	Email string `ohm:"email index"`
}

func dbCleanup() {
	conn, _ := redisurl.Connect()

	conn.Do(`SCRIPT FLUSH`)
	conn.Do(`FLUSHDB`)
	conn.Close()
}

func TestSaveLoadsID(t *testing.T) {
	dbCleanup()
	defer dbCleanup()
	Gohm, err := NewDefaultGohm()
	if err != nil {
		t.Error(err)
	}

	user := &User{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	}

	err = Gohm.Save(user)
	if err != nil {
		t.Error(err)
	}

	if user.ID != `1` {
		t.Errorf(`id is not set (expected "1", got "%v")`, user.ID)
	}
}

func TestFindStruct(t *testing.T) {
	dbCleanup()
	defer dbCleanup()
	Gohm, _ := NewDefaultGohm()
	Gohm.Save(&User{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	})

	var user User
	err := Gohm.Find(`1`, &user)
	if err != nil {
		t.Error(err)
	}

	if user.ID != `1` {
		t.Errorf(`id not correctly set in model (expected "1", was "%v")`, user.ID)
	}

	if user.Name != "Marty" {
		t.Errorf(`incorrect Name set (expected "Marty", got "%v")`, user.Name)
	}
}
