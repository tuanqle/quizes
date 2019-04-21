package main

import (
	"sync"
)

const (
	INVALID loanType = iota
	PURCHASE
	REFINANCE
)

type loanType int
type taskHandler func(context *Context) error

type TaskFunc struct {
	Kind    string
	wg      *sync.WaitGroup
	Handler taskHandler
	Context *Context
	Error   error
}

type Task struct {
	Name  string
	State string
}

type Context struct {
	Client    *Client    `json:"client"`
	LoanType  loanType   `json:"loan-type"`
	Refinance *Refinance `json:"refinance"`
	CoBorrow  *Client    `json:"co-borrower,omitempty"`
	stateMap  map[string]string
	WorkFlow  string `json:"work-flow"`
}

type Client struct {
	Name string `json:"full-name"`
	Age  int    `json:"age"`
}

type Refinance struct {
	Addr    string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode int    `json:"zipcode"`
}
