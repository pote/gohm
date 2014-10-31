package gohm

import(
	`github.com/pote/redisurl`
	`testing`
)

type user struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email index"`
	UUID  string `ohm:"uuid unique"`
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

	u := &user{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	}

	err = Gohm.Save(u)
	if err != nil {
		t.Error(err)
	}

	if u.ID != `1` {
		t.Errorf(`id is not set (expected "1", got "%v")`, u.ID)
	}
}

func TestFindStruct(t *testing.T) {
	dbCleanup()
	defer dbCleanup()
	Gohm, _ := NewDefaultGohm()
	Gohm.Save(&user{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	})

	var u user
	err := Gohm.Find(`1`, &u)
	if err != nil {
		t.Error(err)
	}

	if u.ID != `1` {
		t.Errorf(`id not correctly set in model (expected "1", was "%v")`, u.ID)
	}

	if u.Name != "Marty" {
		t.Errorf(`incorrect Name set (expected "Marty", got "%v")`, u.Name)
	}
}
