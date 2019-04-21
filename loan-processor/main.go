package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	context  *Context
	tasks    map[string]TaskFunc
	workflow map[string][]*Task
)

func init() {
	context = &Context{}

	// Register Task
	tasks = map[string]TaskFunc{
		"basicInfo":  TaskFunc{Handler: basicInfo, Context: context, Kind: "rpc"},
		"refinance":  TaskFunc{Handler: refinance, Context: context, Kind: "rpc"},
		"purchase":   TaskFunc{Handler: purchase, Context: context, Kind: "rpc"},
		"coborrower": TaskFunc{Handler: coBorrower, Context: context, Kind: "rpc"},
		"completion": TaskFunc{Handler: completion, Context: context, Kind: "rpc"},
	}

	// Predefine workflow
	workflow = map[string][]*Task{
		"newAccount": []*Task{
			&Task{Name: "basicInfo", State: "enable"},
			&Task{Name: "refinance", State: "disable"},
			&Task{Name: "purchase", State: "disable"},
			&Task{Name: "coborrower", State: "enable"},
			&Task{Name: "completion", State: "enable"},
		},
	}
}

//
// Run
//
// Execute selected Task
func (t *TaskFunc) Run() {
	if t.Handler == nil {
		log.Print("Task handler is undefined")
		return
	}

	t.wg = &sync.WaitGroup{}
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.Error = t.Handler(t.Context)
	}()

	// Do not wait for parallel task
	// kind: "bg"
	if t.Kind != "bg" {
		t.wg.Wait()
	}
}

//
// RegisterWorkFlow
//
// Tracking task state to dynamically enable/disable next task
//
func (c *Context) RegisterWorkFlow(workName string) error {
	tasks, ok := workflow[workName]
	if !ok {
		return fmt.Errorf("invalid workflow '%s'", workName)
	}
	c.WorkFlow = workName
	c.stateMap = make(map[string]string)

	for _, t := range tasks {
		c.stateMap[t.Name] = t.State
	}
	return nil
}

//
// Execute
//
// Perform workflow tasks for context
//
func (ctx *Context) Execute() {
	for _, task := range workflow[ctx.WorkFlow] {
		t, ok := tasks[task.Name]
		if !ok {
			log.Fatalf("no task '%s' define", task.Name)
		}
		if ctx.stateMap[task.Name] == "enable" {
			t.Run()
			if t.Error != nil {
				log.Fatalf("%v", t.Error)
			}
		}
	}
}

//
// Enum Print
//
func (l loanType) String() string {
	switch l {
	case 1:
		return "purchase"
	case 2:
		return "refinance"
	}
	return "invalid"
}

//
// Client Print
//
func (c *Client) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("  Full name: %s\n", c.Name))
	buff.WriteString(fmt.Sprintf("        Age: %d\n", c.Age))
	return buff.String()
}

//
// Refinance Print
//
func (refi *Refinance) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("\nREFINANCE INFO\n"))
	buff.WriteString(fmt.Sprintf("    Address: %s\n", refi.Addr))
	buff.WriteString(fmt.Sprintf("       City: %s\n", refi.City))
	buff.WriteString(fmt.Sprintf("      State: %s\n", refi.State))
	buff.WriteString(fmt.Sprintf("        Zip: %d\n", refi.ZipCode))
	return buff.String()
}

//
// Context Print
//
func (ctx *Context) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(fmt.Sprintln("\nYou provided the following:\n"))
	buff.WriteString(fmt.Sprintf("YOUR INFORMATION\n"))
	buff.WriteString(fmt.Sprintf("%s", ctx.Client))
	buff.WriteString(fmt.Sprintf("  Loan Type: %s\n", ctx.LoanType))
	switch ctx.LoanType {
	case REFINANCE:
		buff.WriteString(fmt.Sprintf("%s", ctx.Refinance))
	}
	if ctx.CoBorrow != nil {
		buff.WriteString(fmt.Sprintf("\nCO-BORROWER INFO\n"))
		buff.WriteString(fmt.Sprintf("%s", ctx.CoBorrow))
	}
	return buff.String()
}

//
// clientInfo
//
// Collect client's name & age
//
func clientInfo(coborrower bool) (*Client, error) {
	borrower := "your"
	if coborrower == true {
		borrower = "your co-borrower's"
	}
	msgs := []string{
		fmt.Sprintf("  What is %s full name?", borrower),
		fmt.Sprintf("  What is %s age?", borrower),
	}

	client := &Client{}
	scanner := bufio.NewScanner(os.Stdin)

	// Collect Client Name
	fmt.Printf("%s ", msgs[0])
	if scanner.Scan() == false {
		return nil, scanner.Err()
	}
	client.Name = scanner.Text()

	// Get Client Age
	for {
		fmt.Printf("%s ", msgs[1])
		if scanner.Scan() == false {
			return nil, scanner.Err()
		}

		age, err := strconv.Atoi(scanner.Text())
		if err == nil && age > 0 {
			client.Age = age
			break
		}

		fmt.Println("\n    Invalid input... please try again!\n")
	}

	return client, nil
}

//
// loanInfo
//
// Collect loanType: Purchase or Refinance
//
func loanInfo() (loanType, error) {
	msg := "\nIs this loan for:\n" +
		"  1. New purchase\n" +
		"  2. Refinance\n" +
		"Select Option?"

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s ", msg)
		if scanner.Scan() == false {
			return INVALID, scanner.Err()
		}
		selection, _ := strconv.Atoi(scanner.Text())
		if selection > 0 && selection < 3 {
			return loanType(selection), nil
		}
		fmt.Printf("\n    Invalid selection '%d'... please try again!\n\n", selection)
		continue
	}
	return INVALID, nil
}

//
// refinance
//
// Collect information related to refinance
//
func refinance(ctx *Context) error {
	refiMsg := "\nIf you're refinancing your loan, " +
		"please indicate the address of the property on which " +
		"the loan was taken out."
	msgs := []string{
		"  What is the street address?",
		"  What is the city?",
		"  What is the state [i.e: CA]?",
		"  What is the zipcode?",
	}

	fmt.Println(refiMsg)
	refi := &Refinance{}
	scanner := bufio.NewScanner(os.Stdin)

	// Street Addr
	fmt.Printf("%s ", msgs[0])
	if scanner.Scan() == false {
		return scanner.Err()
	}
	refi.Addr = scanner.Text()

	// City
	fmt.Printf("%s ", msgs[1])
	if scanner.Scan() == false {
		return scanner.Err()
	}
	refi.City = scanner.Text()

	// State
	for {
		fmt.Printf("%s ", msgs[2])
		if scanner.Scan() == false {
			return scanner.Err()
		}
		state := scanner.Text()
		if len(state) == 2 {
			refi.State = strings.ToUpper(state)
			break
		}
		fmt.Println("\n    Invalid state code... please try again!\n")
	}

	// Zipcode
	for refi.ZipCode == 0 {
		fmt.Printf("%s ", msgs[3])
		if scanner.Scan() == false {
			return scanner.Err()
		}
		zip, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("\n    Invalid state code... please try again!\n")
			continue
		}
		refi.ZipCode = zip
	}
	ctx.Refinance = refi
	return nil
}

//
// purchase
//
// Collect information for purchase task
//
func purchase(ctx *Context) error {
	return nil
}

//
// coBorrower
//
// Collect coBorrower information
//
func coBorrower(ctx *Context) error {
	msg := []string{
		"  Are you applying with a co-borrower?",
		"\nComplete the following question for your co-borrower.",
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("%s ", msg[0])
	if scanner.Scan() == false {
		return scanner.Err()
	}
	res := strings.ToLower(scanner.Text())
	if res == "yes" || res == "y" {
		fmt.Println(msg[1])
		client, err := clientInfo(true)
		if err != nil {
			return err
		}
		ctx.CoBorrow = client

	}

	return nil
}

//
// basicInfo
//
// Collect client information to open an account
//
func basicInfo(ctx *Context) error {
	msg := "Please answer the following questions:"

	var err error
	fmt.Println(msg)
	if ctx.Client, err = clientInfo(false); err != nil {
		return err
	}

	if ctx.LoanType, err = loanInfo(); err != nil {
		return err
	}

	// Enable next Task based on LoanType
	ctx.stateMap[ctx.LoanType.String()] = "enable"
	return nil
}

func completion(ctx *Context) error {
	fmt.Println("Thank you for your submission.")
	fmt.Printf("%v", ctx)
	return nil
}

//
// main
//
func main() {
	// Welcome Banner
	welcome := "=== Welcome to your loan portal ===\n" +
		"We will collect some basic information about you now " +
		"to get you started in your application.\n\n"

	fmt.Println(welcome)

	// Demonstrate workflow
	myWorkFlow := "newAccount"

	context.RegisterWorkFlow(myWorkFlow)
	context.Execute()
}
