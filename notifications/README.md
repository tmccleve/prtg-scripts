# PRTG/PagerDuty Notification Integration

## Goals

* Create incidents using version 2 of the PagerDuty API for triggered PRTG alerts.

* Automatically resolve alerts when status returns to normal or paused in PRTG.


## Build & Installation

Install dependencies:

`go get github.com/PagerDuty/go-pagerduty`

`go get gopkg.in/yaml.v2`

Build the package

`go build pagerduty.go`

From an Adminstrator powershell session:
`cp pagerduty.exe "C:\Program Files (x86)\PRTG Network Monitor\Notifications\EXE\"`

## PagerDuty yaml config:
```yaml
apikey: myShineyApiKey
useremail: user@example.com
```

Should exist at "C:\\Program Files (x86)\\PRTG Network Monitor\\Notifications\\EXE\\.pd.yml"

## Configuring notification in PRTG

Create new basic notification. Check "EXECUTE PROGRAM" selecting pagerduty.exe from the Program File dropdown.

Populate the parameter field with the following, substituting the service key with your service integration key

`-probe "%probe" -device "%device" -name "%name" -status "%status" -date "%datetime" -linkdevice %linkdevice -message "%message" -servicekey myShineyServiceKey`
