
This folder includes source code that implements `packethost` interview challenge.
* Build a program that uses API, allowing launching a machine and tearing it down**

### Build

`go build`

### Execute

`./packet`

#### What it does?

The program goes through a process of selecting: operating-system, plan, and facility. It uses these
information to launch a new device and waits for this device to provision. As a final step, it removes
the device.

### Code breakdown

#### Files

`types.go` - contains data structure

`main.go` - contains the implementation


#### Methods

There are 4 main methods: `selectOS`, `selectPlan`, `selectFacility`, and `createDevice`.

`func selectOS(token string, distro string)`

    This routine uses `/operating-systems` API to retrieve a list of available operation systems. 
    It filters and select an operating-system by using `distro` field.
    
 `func selectPlan(token, projID string, class string)`
 
    This routine uses `/projects/{projID}/plans` API to retrieve a list of available plans associated 
    with the project `{projID}`.  It filters and selects a plan by using `class` field.
    
 `func selectFacility(token, projID string, feature string)`
 
    This routine uses `/projects/{projID}/facilities` API to retrieve a list of avialable facilities
    belong to project `{projID}`. It filters and selects a facility by using `feature` field. I choose
    to use single feature to select a facility for simplicity. However, it is possible to enhance to
    support multiple features filtering.
    
 `func createDevice(token, projID string, os *OS, facility *Facility, plan *Plan)`
 
    This routine uses `/projects/{projID}/devices` API with `POST` method to create a device. It uses
    data from selected operating-system, plan, and facility to orchestrate the launch.
    
### Improvements

This section describew suggested improvementw to the current program.

#### Automate the selection

Current program uses hardcoded values to demonstrate launching and tearing down a machine. To further improve
the program, instead of using hardcoded values, I would automated the process:

1. Get a list of operating-system which has `provisionable_on` field set.
2. Iterate through the `operating-system` list from step 1 and for each `value` in `provisionable_on`, retrieves a
list of available `plans` with matching `value`
3. Iterate through the `plans` list in step 2 and use `feature` field to retrieve a list of available `facilities`
4. Iterate through the `facilities` list and attempt to create device with `os`,`plans`, and `facility`
5. If unsuccessful, repeat from step 2 until we exhausted all the posibility.

#### Handle device creation failure

Current program can run indefinitely waiting for device provisioning. Scenarios such as: PXE boot failure,
corrupted OS, etc. would triggered the condition. To protect from this scenario, I would add a watchdog.

* Add a watchdog to recover from indefinite provisioning.
* In the case watchdog occurs, rollback the system's state (forced deprovisioned new device)
* Record any log that can help debug the failure

#### Unit Test

The program is lacking unit-test. 

### Suggestions

* `/plans` API provides hardware architecture [AMD or ARM]  in description.
  Suggest adding a new field: `architecture`
  
* `/projects/<id>/devices` API uses fields, plans, facilites, operating-systems of type strings.
  It is unclear what the associate fields are.
  Suggest replacing it with `plan-ID`, `facility-name`, and `os-slug`
