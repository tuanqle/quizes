package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// printPretty
// Format struct for console output
func printPretty(intf interface{}) string {
	s, _ := json.MarshalIndent(intf, "", "\t")
	return string(s)
}

// getReq
// Manage GET message to server
func getReq(token, path string) ([]byte, error) {
	// Retrieve handle
	url := fmt.Sprintf("https://api.packet.net/%s", path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Set Authentication header
	req.Header.Add("X-Auth-Token", token)

	// Prepare Client
	var resp *http.Response
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	txport := &http.Transport{TLSClientConfig: tlsConfig, DisableCompression: true}
	client := http.Client{Transport: txport}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("[GET] response error %d: body = '%s'", string(body))
	}
	// Return payload
	return body, nil
}

// postReq
// Manage POST message to server
func postReq(token, path string, buf []byte) ([]byte, error) {
	// Retrieve handle
	url := fmt.Sprintf("https://api.packet.net/%s", path)
	reader := bytes.NewReader(buf)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}

	// Add header
	// Set Authentication header
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/json")

	// Prepare client
	// Use insecure client to connect
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	txport := &http.Transport{TLSClientConfig: tlsConfig, DisableCompression: true}
	client := http.Client{Transport: txport}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return body, fmt.Errorf("[POST] response error [%d]: body = '%s'", resp.StatusCode, string(body))
	}

	return body, nil
}

// delReq
// Manage DELETE message to server
func delReq(token, path string, buf []byte) ([]byte, error) {
	url := fmt.Sprintf("https://api.packet.net/%s", path)

	reader := bytes.NewReader(buf)
	req, err := http.NewRequest("DELETE", url, reader)
	if err != nil {
		return nil, err
	}

	// Add header
	// Set Authentication header
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/json")

	// Prepare client
	// Use insecure client to connect
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	txport := &http.Transport{TLSClientConfig: tlsConfig, DisableCompression: true}
	client := http.Client{Transport: txport}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return body, fmt.Errorf("[POST] response error [%d]: body = '%s'", resp.StatusCode, string(body))
	}

	return body, nil
}

// selectOS
// Select an operating system with matching distro
// For this demo, I use distro as a matching criteria
func selectOS(token string, distro string) (*OS, error) {
	if len(distro) == 0 {
		return nil, fmt.Errorf("invalid param: require 'distro' value to select OS")
	}

	fmt.Printf("selectOS(): searching for distro='%s'\n", distro)
	payld, err := getReq(token, "operating-systems")
	if err != nil {
		return nil, fmt.Errorf("error in '/operation-systems' API %v", err)

	}

	// Retreive list of systems
	oses := &OSes{}
	if err := json.Unmarshal(payld, oses); err != nil {
		return nil, fmt.Errorf("internal error unmarshal '/operating-systems' [%v] payload=%s", err, string(payld))
	}

	// Handle errors from Server response
	if oses.Errors != nil {
		return nil, fmt.Errorf("Server error: %v", oses.Errors)
	}

	// Find an available system
	// Provisionable gives information about the plan
	for _, os := range oses.OperatingSystems {
		if strings.Compare(os.Distro, distro) == 0 && len(os.Provisionable) > 0 {
			fmt.Printf("Found: %s\n", printPretty(os))
			return os, nil
		}
	}

	// No system available
	return nil, fmt.Errorf("no available operating-system '%s'", distro)
}

// selectPlan
// Find a plan with matching line
// For this demo, I use plan class as matching criteria
func selectPlan(token, projID, class string) (*Plan, error) {
	if len(class) == 0 {
		return nil, fmt.Errorf("invalid param: require 'class' value to select Plan")
	}

	fmt.Printf("selectPlan(): searching for class='%s'\n", class)
	path := fmt.Sprintf("projects/%s/plans", projID)

	payld, err := getReq(token, path)
	if err != nil {
		return nil, fmt.Errorf("error with '%s' API %v", path, err)
	}

	plans := &Plans{}
	if err := json.Unmarshal(payld, plans); err != nil {
		return nil, fmt.Errorf("internal error unmarshal '%s' [%v] payload=%s", path, err, string(payld))
	}

	// Check for server error response
	if plans.Errors != nil {
		return nil, fmt.Errorf("Server error: %v", plans.Errors)
	}

	// Find a with valid ID
	for _, p := range plans.Plans {
		if strings.Compare(p.Class, class) == 0 {
			fmt.Printf("Found: %s\n", printPretty(p))
			return p, nil
		}
	}

	// No available plan
	return nil, fmt.Errorf("No available plan class='%s'", class)
}

// selectFacility
// Facility selection is base on matching feature
// For this demo, I use single feature as matching criteria
func selectFacility(token, projID string, feature string) (*Facility, error) {
	if len(feature) == 0 {
		return nil, fmt.Errorf("invalid param: require 'feature' value to select Facility")
	}

	fmt.Printf("selectFacility(): searching for feature='%s'\n", feature)
	path := fmt.Sprintf("projects/%s/facilities", projID)

	payld, err := getReq(token, path)
	if err != nil {
		return nil, fmt.Errorf("error with '%s' API %v", path, err)
	}

	facs := &Facilities{}
	if err := json.Unmarshal(payld, facs); err != nil {
		return nil, fmt.Errorf("internal error unmarshal '%s' [%v] payload=%s", path, err, string(payld))
	}

	// Check for server error response
	if facs.Errors != nil {
		return nil, fmt.Errorf("Server error: %v", facs.Errors)
	}

	// Find a facility with valid ID
	for _, f := range facs.Facilities {
		if len(f.Features) == 0 {
			continue
		}
		// Check for matching features
		for _, val := range f.Features {
			if strings.Compare(val, feature) == 0 {
				fmt.Printf("Found: %s\n", printPretty(f))
				return f, nil
			}
		}
	}

	// No available facility
	return nil, fmt.Errorf("No available facility with matching feature '%s'", feature)
}

// createDevice
// Launch device with selected facility and selected plan
func createDevice(token, projID string, os *OS, facility *Facility, plan *Plan) (*Device, error) {
	// Check for nil
	if os == nil || plan == nil {
		return nil, fmt.Errorf("Missing param: minimum param required 'operating-system and plan'")
	}

	path := fmt.Sprintf("projects/%s/devices", projID)
	factName := "any"
	if facility != nil {
		factName = facility.Id
	}

	fmt.Printf("createDevice():\n\tOS='%s' with plan='%s' at facility='%s'\n", os.Distro, plan.Id, factName)
	enroll := &Enroll{
		Facility: factName,
		Plan:     plan.Id,
		OpSystem: os.Slug,
	}

	buf, err := json.Marshal(enroll)
	if err != nil {
		return nil, fmt.Errorf("internal error marshal %v", err)
	}

	// Process response payload
	var payld []byte
	payld, err = postReq(token, path, buf)
	if err != nil {
		return nil, fmt.Errorf("internal error: failed to send path='%s' body='%s': error [%v]", path, string(buf), err)
	}

	// UnMarshal Device
	device := &Device{}
	if err := json.Unmarshal(payld, device); err != nil {
		return nil, fmt.Errorf("internal error unmarshal '%s' [%v] payload=%s", path, err, string(payld))

	}

	fmt.Printf("Created: %s\n", printPretty(device))
	return device, nil
}

// retreiveEvent
// Retreive the latest event for device given
func retreiveEvent(token string, dev *Device) (*Event, error) {
	if dev == nil || len(dev.Id) == 0 {
		return nil, fmt.Errorf("No device given")
	}

	path := fmt.Sprintf("devices/%s/events", dev.Id)
	payld, err := getReq(token, path)
	if err != nil {
		return nil, fmt.Errorf("error with '%s' API %v", path, err)
	}

	events := &Events{}
	if err := json.Unmarshal(payld, events); err != nil {
		return nil, fmt.Errorf("internal error unmarshal '%s' [%v] payload=%s", path, err, string(payld))
	}

	// Check for server response error
	if events.Errors != nil {
		return nil, fmt.Errorf("Server error: %v", events.Errors)

	}

	// No event
	if len(events.Events) == 0 {
		return nil, fmt.Errorf("No event given")
	}

	return events.Events[0], nil
}

// removeDevice
// Decommission device
func removeDevice(token string, device *Device) error {
	if device == nil || len(device.Id) == 0 {
		return fmt.Errorf("no device given")
	}

	// 1 Time use
	type deleteBody struct {
		Force bool `json:"force_delete"`
	}

	force := &deleteBody{Force: true}

	buf, err := json.Marshal(force)
	// Internal error with marshal
	// We ignore payload. Trigger delete without FORCE
	if err != nil {
		fmt.Printf("WARNING: %v", err)
	}

	// Remove device
	path := fmt.Sprintf("devices/%s", device.Id)
	_, err = delReq(token, path, buf)
	if err != nil {
		fmt.Printf("error in DELETE '%s' API data %v", path, string(buf), err)
	}

	// Validate device removal
	_, err = retreiveEvent(token, device)
	// No Event for device indicate device is deprovisioned
	if err != nil {
		fmt.Printf("Device '%s' removed", device.Id)
		return nil
	}
	return fmt.Errorf("something went wrong")
}

// Main func
func main() {
	token := "N5SGxfcjmWuuWhv5zu4hW1sUBDuYWvt5"
	projectID := "ca73364c-6023-4935-9137-2132e73c20b4"

	system, err := selectOS(token, "ubuntu")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	plan := &Plan{}
	plan, err = selectPlan(token, projectID, "c2.medium.x86")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	facility := &Facility{}
	facility, err = selectFacility(token, projectID, "global_ipv4")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	var device *Device
	device, err = createDevice(token, projectID, system, facility, plan)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	// Waiting for provision complete
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		// 30 sec interval
		ticker := time.NewTicker(30 * time.Second)

		defer ticker.Stop()

		fmt.Printf("Provisioning device... wait...\n")
		for {
			select {
			case t := <-ticker.C:
				event, err := retreiveEvent(token, device)
				// Throwing WARNING instead of quit to avoid server hanging
				if err != nil {
					fmt.Printf("WARNING: %v\n", err)
					continue
				}

				fmt.Printf("%s: %s\n", t.Format(time.RFC3339), event.Body)
				if strings.Contains(event.Body, "Provision complete") {
					return
				}
			}
		}
	}()
	wg.Wait()

	// Deprovision
	if err := removeDevice(token, device); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
