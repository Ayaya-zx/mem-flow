package inmem

import (
	"testing"

	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

func TestAddUserAndGetUser(t *testing.T) {
	repo := NewInmemUserRepository()

	user := &entity.User{
		Name: "User",
	}

	// Add user
	err := repo.AddUser(user)
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.users) != 1 {
		t.Fatalf("got len(repo.users) = %d; want 1", len(repo.users))
	}

	// Get added user
	user, err = repo.GetUser("User")
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "User" {
		t.Errorf("got user.Name = %s; want \"User\"", user.Name)
	}

	u, err := repo.GetUser("NonExistentUser")
	if err == nil {
		t.Error("got nil; want err")
	}
	if u != nil {
		t.Error("got non-nil; want nil")
	}
}

func TestAddUserWithEmptyName(t *testing.T) {
	repo := NewInmemUserRepository()

	user := &entity.User{
		Name: "",
	}

	// Adding a user with an empty name is prohibited
	err := repo.AddUser(user)
	if err == nil {
		t.Errorf("got nil; want error")
	}
	if len(repo.users) != 0 {
		t.Errorf("got len(repo.users) = %d; want 0", len(repo.users))
	}
}

func TestAddUserWithSameNameTwice(t *testing.T) {
	repo := NewInmemUserRepository()

	u1 := &entity.User{
		Name: "User",
	}
	u2 := &entity.User{
		Name: "User",
	}

	err := repo.AddUser(u1)
	if err != nil {
		t.Fatal(err)
	}

	// We should not be able to add two users with same names
	err = repo.AddUser(u2)
	if err == nil {
		t.Errorf("got nil; want error")
	}
}

func TestRemoveUser(t *testing.T) {
	repo := NewInmemUserRepository()

	u := &entity.User{
		Name: "User",
	}

	err := repo.AddUser(u)
	if err != nil {
		t.Fatal(err)
	}
	err = repo.RemoveUser("User")
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.users) != 0 {
		t.Errorf("got len(repo.users) = %d; want 0", len(repo.users))
	}
}
