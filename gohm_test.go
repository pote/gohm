package gohm

import(
	"testing"
)

type User struct {
	id    string `ohm:"id"`
	Name  string `ohm:"name"`
	Email string `ohm:"email"`
}

func (u *User) ID() string {
	return u.id
}

func TestSave(t *testing.T) {
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


	//cleanupConn := connectionPool.Get()
	//cleanupConn.Do("SCRIPT FLUSH")
	//cleanupConn.Do("FLUSHDB")
	//defer cleanupConn.Close()
}
