package main

import "testing"

func TestValidation(t *testing.T) {
	customers := UserList{
		{"Rick Mann", "Password1", 30100.50},
		{"Peter Pan", "neverland", 2750.00},
	}

	testData := []struct {
		a, b string
		want bool
	}{
		{"Rick Mann", "Password1", true},
		{"Random", "wrong", false},
		{"rIcK mAnN", "Password1", true},
	}
	for _, v := range testData {
		found, _ := customers.login(v.a, v.b)
		if found != v.want {
			t.Fail()
		}
	}
}

func TestMatchIgnoreCase(t *testing.T) {
	testData := []struct {
		a, b string
		want bool
	}{
		{"Rick Mann", "rick mann", true},
		{"Rick Mann", "Rick Mann", true},
		{"RickiMann", "Rick Mann", false},
	}
	for _, v := range testData {
		b := matchIgnoreCase(v.a, v.b)
		if b != v.want {
			t.Fail()
		}
	}
}
