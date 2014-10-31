package gohm

import(
	"github.com/pote/redisurl"
	"testing"
)

type User struct {
	ID    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email"`
}

func dbCleanup() {
	conn, _ := redisurl.Connect()

	conn.Do("SCRIPT FLUSH")
	conn.Do("FLUSHDB")
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
		Name: "Marty",
		Email: "marty@mcfly.com",
	}

	err = Gohm.Save(user)
	if err != nil {
		t.Error(err)
	}

	if user.ID != "1" {
		t.Errorf("id is not set: %v", user.ID)
	}
}

func TestFindStruct(t *testing.T) {
	dbCleanup()
	defer dbCleanup()
	Gohm, _ := NewDefaultGohm()
	Gohm.Save(&User{
		Name: "Marty",
		Email: "marty@mcfly.com",
	})


	var user User
	Gohm.Find("1", user)

	if user.Name != "Marty" {
		t.Errorf("incorrect Name set: %v", user.Name)
	}
}
