package main

import (
	"errors"
	"fmt"
)

type Account interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	Balance() float64
}

type CurrentAccount struct {
	balance float64
}

func (a *CurrentAccount) Deposit(amount float64) error {
	a.balance += amount
	return nil
}

func (a *CurrentAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("negative amount")
	}
	if a.balance < amount {
		return errors.New("not enough funds")
	}
	a.balance -= amount
	return nil
}

func (a *CurrentAccount) Balance() float64 {
	return a.balance
}

type SavingsAccount struct {
	CurrentAccount
}

func (a *SavingsAccount) Withdraw(amount float64) error {
	if a.balance < amount || a.balance <= 500 {
		return errors.New("not enough funds")
	}
	a.balance -= amount
	return nil
}

func ProcessAccount(a Account) {
	err := a.Deposit(500)
	if err != nil {
		fmt.Println("Error depositing funds")
	}
	err = a.Withdraw(200)
	if err != nil {
		fmt.Println("Error withdrawing funds")
	}
	fmt.Printf("Balance: %.2f\n", a.Balance())
}

func main() {
	c := &CurrentAccount{}
	s := &SavingsAccount{}
	ProcessAccount(c)
	ProcessAccount(s)
}
