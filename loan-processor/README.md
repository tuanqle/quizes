This is a command-line program demonstrating a `workflow` engine for a loan application.
A `workflow` consists of a series of tasks. These tasks are executed in a particular order.
The next `task` is based on client's responses.

### Prerequisite
This program is written using GOlang. Download the lastest `GO` compiler using the link below
    
    https://golang.org/dl/
    

### Installation
Untar the archive to current working directory
   
    tar xvzf loan-processor.tar.gz
  

Source code can also be found here
 
    https://github.com/tuanqle/quizes/loan-processor


### Compiling

    cd loan-processor
    go build


### Execution

    ./loan-processor


#### What it does?

As the program starts, it initializes a pre-defined set of tasks: `basicInfo`, `refinance`,
`purchase`, `co-borrower`, and `completion`. It also initializes pre-defined work-flow,
`newAccount`, which consists of an orderly set of `task` for execution. A `context` is initialized
with the selected work-flow, `newAccount`. `context.Execute()` begins to execute the work-flow.
As the program progresses, it enables or disables a follow-up `task` based on client's responses.

### Contents

    `types.go`  - contains data structure definition
    `main.go`   - program implementation
    `README.md` - This README file

### Code breakdown

#### Key Components
    `tasks      map[string]TaskFunc` 
    This map uses `task`'s name as a key. Each entry of the map holds `TaskFunc` struct.
    The purpose of this map is to lookup `task` handler and its setting to execute.
    
    `type TaskFunc struct {}`
    This holds settings to how to execute the method.  This is used by the `task` executor engine.
        `Handler` - set to the task func/method
        `Context` - set to the current working context
        `Kind`    - set to the execution mode. It consists of 2 modes: `bg` and `rpc`.
                    `bg` mode is to launch the `task` and return to caller immediately.
                    `rpc` mode is to wait until `task` is completed before return to caller.

    `workflow   map[string][]*Task`
    `workflow` map holds a pre-defined work-flow using its name as the key.  Its purpose is to
    define a series of `task` in the order of execution. 

    `type Task struct {}`
    This holds `task`'s name and `state`. This struct allows tuning `task`'s state according to 
    the work-flow
        `Name` - `task`'s name'
        `State` - consists of 2 states: `enable` or `disable`

    `type Context struct{}`
    This holds client's data, `workflow`'s state', and `task`'s state at run-time.
        `Client`    - store client's information: Name and Age
        `LoanType`  - type of loan: `refinance` or `purchase`
        `Refinance` - `refinance`information: Address, City, and State
        `CoBorrow`  - store co-borrower's information (if any): Name and Age
        `stateMap`  - map of `task` state according to the current run-time.
                      This allows dynamically tuning the state of a next `task` based
                      on client's response.
        `Workflow`  - name of the work-flow that the `context` is executing

#### Methods
This section describes `funct` or `methods`.

    `type taskHandler func(context *Context) error`
    This type defines handler's function syntax

    `func (t *TaskFunc) Run()`
    This method executes the `task` by calling the associate `func` pointed by `Handler`.
    It references `Kind` to determine how to execute the `Handler`

    `func (c *Context) RegisterWorkFlow()`
    This method initializes `context.stateMap` with the pre-defined `task`'s state.
    It also saves the `workflow`'s name with this context.

    `func (c *Context) Execute()`
    This method lookups the `workflow` and executes a series of `task` in the order defined
    by the `workflow`. It launches `task` of which state is `enable`.

    `func clientInfo()`
    This method prompts to collect client's data: Name and Age. It also uses to
    collect co-borrower's data.

    `func loanInfo()`
    This method prompts to collect the type of loan.

    `func refinance()`
    Task's handler to collect `refinance` data such as:
    address, city, and state

    `func purchase()`
    Task's handler to perform `purchase` loan-type. This is currently emptied.

    `func coBorrower()`
    Task's handler to collect co-borrower's data: Name and Age

    `func basicInfo()`
    Task's handler `basicInfo` to start a loan application.  It prompts client
    for their information and to select a loan type.

    `func completion()`
    Task's handler to summarize the loan application and print out a thank you message.

    Methods for pretty-print
    `func (l loanType) String() string`
    `func (c *Client) String() string`
    `func (refi *Refinance) String() string`
    `func (ctx *Context) String() string`

### Improvement
This section describes possible enhancements that can be done for the program.

    `allow 2 or more co-borrowers`
    An improvement to this program is to allow more than 1 co-borrower. This can
    be done by changing `context.CoBorrow *Client` to `context.CoBorrow []*Client`

    `Save and resume workflow`
    An improvement to this program is to allow client to resume their workflow. This
    can be achieved by storing the current `context` in JSON format to a file. This
    can then easily restored to the `context` for continuation.

    `live service API endpoint`
    The program is designed with `context` which can be extended to allow concurrency
    with API endpoints.  Each `context` would be independent for concurrency access.
    It also allows different work-flow performed by different clients through API endpoint.

    `Task executor improvement`
    A mechanism to notify caller of task completion in 'bg' mode.
