package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// User describes a customer of the bank.
type User struct {
	name     string
	password string
	balance  float64
}

// Deposit increases balance by amount.
// Amount must be positive.
func (user *User) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("Amount less than 0")
	}
	user.balance += amount
	return nil
}

// Withdraw increases balance by amount.
// Amount must be positive.
func (user *User) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("Amount less than 0")
	}
	if user.balance-amount < 0 {
		return errors.New("Not enogh balance to withdraw amount")
	}
	user.balance -= amount
	return nil
}

// SPrintBalance returns a string of the users balance,
// formatted to be user facing.
func (user *User) SPrintBalance() string {
	return fmt.Sprintf("%v kr.", user.balance)
}

// UserList is a list of Users.
type UserList []User

func main() {
	customers := UserList{
		{"Rick Mann", "Password1", 301000.50},
		{"Peter Pan", "neverland", 2750.00},
	}

	found, user := customers.login(askCredentials())
	if !found {
		fmt.Println("Fel användarnamn eller lösenord.")
		return
	}
	fmt.Printf("Hello, %v!\n", user.name)
	for {
		fmt.Println("What would you like to do?")
		fmt.Println("1. Check balance")
		fmt.Println("2. Deposit money")
		fmt.Println("3. Withdraw money")
		fmt.Println("Or exit")
		choice := userChoice([]string{"1", "1.", "2", "2.", "3", "3.", "exit"})
		switch strings.Trim(choice, ".") {
		case "1":
			fmt.Println(user.SPrintBalance())
		case "2":
			fmt.Println("How much would you like to deposit?")
			n := askForNumber()
			err := user.Deposit(n)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("New balance ", user.SPrintBalance())
			}
		case "3":
			fmt.Println("How much would you like to withdraw?")
			n := askForNumber()
			err := user.Withdraw(n)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("New balance ", user.SPrintBalance())
			}
		case "exit":
			return
		}
	}
}

// login checks the given credentials agains UserList and returns
// true and the User object if it finds a match.
// Otherwis false and an empty User{}
func (list UserList) login(username, password string) (found bool, user User) {
	for _, v := range list {
		if (matchIgnoreCase(username, v.name)) && (password == v.password) {
			user = v
			break
		}
	}
	if (user == User{}) {
		// Om user är tom efter loopen betyder det att ingen
		// användare med det namn och lösenordet hittades.
		return false, user
	}
	return true, user
}

// askCredentials asks the user for a username and password.
// The password field is now echoed when entered.
func askCredentials() (username string, password string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Name: ")
	username, _ = reader.ReadString('\n')
	username = strings.Trim(username, " \n\r")

	fmt.Print("Password: ")
	pwbytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println("") // Move down one line
	if err != nil {
		panic(err)
	}
	return username, strings.Trim(string(pwbytes), " \r")
}

// matchIgnoreCase returns if the two strings are equal,
// ignoring upper or lower case.
func matchIgnoreCase(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

// userChoice takes userinput that must match one of the given
// options and returns first valid input.
func userChoice(options []string) string {
	lookup := make(map[string]bool)
	for _, v := range options {
		lookup[strings.ToLower(v)] = true
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		s, _ := reader.ReadString('\n')
		// Trim trailing whitespace
		s = strings.Trim(s, "\n\r ")
		if lookup[strings.ToLower(s)] {
			return s
		}
		fmt.Println("Not a valid option")
	}
}

// askForNumber takes user input and converts it to float64.
// If the user types an invalid numer they get to try again
// until they enter a valid number or exit.
// If user types exit, returns 0
func askForNumber() float64 {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, _ := reader.ReadString('\n')
		s = strings.Trim(s, " \n\r")
		n, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return n
		} else if s == "exit" {
			return 0
		}
		fmt.Print("Not a valid number. Try again: ")
	}
}
