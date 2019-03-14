
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
    In order to pick an operating-system to use, I chose `distro` as the filter.
    
 `func selectPlan(token, projID string, class string)`
 
    This routine uses `/projects/<projID>/plans` API to retrieve a list of available plans associated 
    with the project.  To pick which plan to use, I chose `class` as the filter.
    
 `func selectFacility(token, projID string, feature string)`
 
    This routine uses `/projects/<projID>/facilities` API to retrieve a list of avialable facilities
    associated with the project with ID: `projID`. To pick facilities, I chose `feature` as the filter.
    
 `func createDevice(token, projID string, os *OS, facility *Facility, plan *Plan)`
 
    This routine uses `/projects/<projID>/devices` API with `POST` method to create a device. It uses
    fields from selected operating-system, plan, and facility data as message input.
    
### Suggestion

* `/plans` API provides hardware architecture [AMD or ARM]  in description.
  Suggest adding a new field: `architecture`
  
* `/projects/<id>/devices` API uses fields, plans, facilites, operating-systems of type strings.
  It is unclear what the associate fields are.
  Suggest replacing it with `plan-ID`, `facility-name`, and `os-slug`
