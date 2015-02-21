package gohm

import(
	`github.com/pote/redisurl`
	`testing`
)

type user struct {
	ID      string `ohm:"id"`
	Name    string `ohm:"name"`
	Email   string `ohm:"email index"`
	UUID    string `ohm:"uuid  unique"`
	//Friends []user `ohm:"collection"`
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
	gohm, err := NewGohm()
	if err != nil {
		t.Error(err)
	}

	u := &user{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	}

	err = gohm.Save(u)
	if err != nil {
		t.Error(err)
	}

	if u.ID != `1` {
		t.Errorf(`id is not set (expected "1", got "%v")`, u.ID)
	}
}

func TestLoad(t *testing.T) {
	dbCleanup()
	defer dbCleanup()
	gohm, _ := NewGohm()
	gohm.Save(&user{
		Name: `Marty`,
		Email: `marty@mcfly.com`,
	})

	u := &user{ID: `1`}
	err := gohm.Load(u)
	if err != nil {
		t.Error(err)
	}

	if u.ID != `1` {
		t.Errorf(`id not correctly set in model (expected "1", was "%v")`, u.ID)
	}

	if u.Name != `Marty` {
		t.Errorf(`incorrect Name set (expected "Marty", got "%v")`, u.Name)
	}

	u2 := &user{}
	if err = gohm.Load(u2); err == nil {
		t.Error(`Load should return an error when loading model without a set id`)
	}
}

func TestLoadInvalidID(t *testing.T) {
	dbCleanup()
	defer dbCleanup()

	u := &user{ID: `1000000`}

	gohm, _ := NewGohm()

	err := gohm.Load(u)
	if err == nil {
		t.Error(`did not return an error when fetching an invalid ID`)
	}

	if err.Error() != `Couldn't find "user:1000000" in redis` {
		t.Error(`did not return the expected error message, returned "` + err.Error() + `"`)
	}
}
