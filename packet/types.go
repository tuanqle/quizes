package main

// Hour
type HrCharge struct {
	Price float64 `json:"price"`
	Multi string  `json:"multiplier"`
}

// Price
type Price struct {
	Hour *HrCharge `json:"hour"`
}

// OS
type OS struct {
	Id            string   `json:"id"`
	Slug          string   `json:"slug"`
	Name          string   `json:"name"`
	Distro        string   `json:"distro"`
	Provisionable []string `json:"provisionable_on"`
	Version       string   `json:"version"`
	Pricing       *Price   `json:"pricing"`
	Licensed      bool     `json:"licensed"`
}

// /operating-systems response
type OSes struct {
	OperatingSystems []*OS    `json:"operating_systems,omitempty"`
	Errors           []string `json:"errors,omitempty"`
}

// Actions
type Action struct {
	Type string `json:"type"`
}

// Plan
type Plan struct {
	Id          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Line        string `json:"line"`
	Class       string `json:"class"`
}

// /plans response
type Plans struct {
	Plans  []*Plan  `json:"plans,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Facility
type Facility struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	Features []string `json:"features"`
	IPRanges []string `json:"ip_ranges"`
}

// /facilities response
type Facilities struct {
	Facilities []*Facility `json:"facilities,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

// Event
type Event struct {
	Body         string `json:"body"`
	CreateAt     string `json:"created_at"`
	Id           string `json:"id"`
	Interpolated string `json:"interpolated"`
	IP           string `json:"ip"`
	Type         string `json:"type"`
}

// /events response
type Events struct {
	Events []*Event `json:"events,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Device enrollment
type Enroll struct {
	Facility      string   `json:"facility"`
	Plan          string   `json:"plan"`
	Hostname      string   `json:"hostname,omitempty"`
	Description   string   `json:"description,omitempty"`
	BillCycle     string   `json:"billing_cycle,omitempty"`
	OpSystem      string   `json:"operating_system"`
	PXE           bool     `json:"alwyas_pxe,omitempty"`
	IPXEScriptUrl string   `json:"ipxe_script_url,omitempty"`
	UserData      string   `json:"userdata,omitempty"`
	Locked        bool     `json:"locked,omitempty"`
	CustomData    string   `json:"customdata,omitempty"`
	HwRevId       string   `json:"hardware_reservation_id,omitempty"`
	SpotInstance  bool     `json:"spot_instance,omitempty"`
	SpotPriceMax  int      `json:"spot_price_max,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	projSSH       []string `json:"project_ssh_keys,omitempty"`
	userSSH       []string `json:"user_ssh_keys,omitempty"`
	features      []string `json:"features,omitempty"`
	subnetSize    int      `json:"public_ipv4_subnet_size,omitempty"`
}

// Provisioned device
type Device struct {
	Id          string    `json:"id"`
	ShortId     string    `json:"short_id"`
	Hostname    string    `json:"hostname"`
	Description string    `json:"description"`
	User        string    `json:"user"`
	IQN         string    `json:"iqn"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	Os          *OS       `json:"operating_system"`
	Facility    *Facility `json:"facility"`
	Plan        *Plan     `json:"plan"`
}
