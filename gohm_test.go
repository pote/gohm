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
	cleanupConn := Gohm.RedisPool.Get()
	defer cleanupConn.Do("SCRIPT FLUSH")
	defer cleanupConn.Do("FLUSHDB")
	defer cleanupConn.Close()

	user := &User{
		Name: "Marty",
		Email: "marty@mcfly.com",
	}

	err = Gohm.Save(user)
	if err != nil {
		t.Error(err)
	}
}
