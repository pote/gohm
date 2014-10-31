package gohm

import(
	"github.com/pote/redisurl"
	"testing"
)

type User struct {
	ID    string
	Name  string `ohm:"name"`
	Email string `ohm:"email"`
}

func dbCleanup() {
	conn, _ := redisurl.Connect()

	conn.Do("SCRIPT FLUSH")
	conn.Do("FLUSHDB")
	conn.Close()
}


func TestSave(t *testing.T) {
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
