
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
    
### Suggestion

* `/plans` API provides hardware architecture [AMD or ARM]  in description.
  Suggest adding a new field: `architecture`
  
* `/projects/<id>/devices` API uses fields, plans, facilites, operating-systems of type strings.
  It is unclear what the associate fields are.
  Suggest replacing it with `plan-ID`, `facility-name`, and `os-slug`
