package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"gopkg.in/yaml.v2"
)

type pdConfig struct {
	ApiKey    string `yaml:"apikey"`
	UserEmail string `yaml:"useremail"`
}

type prtgEvent struct {
	Probe       string
	Device      string
	Name        string
	Status      string
	Date        string
	Link        string
	Message     string
	ServiceKey  string
	IncidentKey string
}

const configPath = "C:\\Program Files (x86)\\PRTG Network Monitor\\Notifications\\EXE\\.pd.yml"

var config pdConfig

func init() {
	var probe = flag.String("probe", "local", "The PRTG probe name")
	var device = flag.String("device", "device", "The PRTG device name")
	var name = flag.String("name", "name", "The PRTG sensor name for the device")
	var status = flag.String("status", "status", "The current status for the event")
	var date = flag.String("date", "date", "The date time for the triggered event")
	var link = flag.String("linkdevice", "http://localhost", "The link to the triggering sensor")
	var message = flag.String("message", "message", "The PRTG message for the alert")
	var serviceKey = flag.String("servicekey", "myServiceKey", "The PagerDuty v2 service integration key")
}

func main() {
	config = config.getConf(configPath)
	flag.Parse()

	pd := &prtgEvent{
		Probe:       *probe,
		Device:      *device,
		Name:        *name,
		Status:      *status,
		Date:        *date,
		Link:        *link,
		Message:     *message,
		ServiceKey:  *serviceKey,
		IncidentKey: *probe + "-" + *device + "-" + *name,
	}

	if strings.Contains(pd.Status, "Up") || strings.Contains(pd.Status, "ended") {
		ctx := pagerduty.NewClient(config.ApiKey)
		incident, _ := getTriggeredIncidents(ctx, pd.IncidentKey)
		resolveEvent(ctx, incident.Incidents)
	} else {
		event, err := triggerEvent(pd)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(event)
	}
}

func triggerEvent(prtg *prtgEvent) (*pagerduty.EventResponse, error) {
	event := &pagerduty.Event{
		Type:        "trigger",
		IncidentKey: prtg.IncidentKey,
		ServiceKey:  prtg.ServiceKey,
		Description: prtg.IncidentKey,
		Details: "Link: " + prtg.Link +
			"\nIncidentKey: " + prtg.IncidentKey +
			"\nStatus: " + prtg.Status +
			"\nDate: " + prtg.Date +
			"\nMessage: " + prtg.Message,
		ClientURL: prtg.Link,
	}
	res, err := pagerduty.CreateEvent(*event)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func resolveEvent(ctx *pagerduty.Client, incidents []pagerduty.Incident) {
	var incidentSlice []pagerduty.Incident
	for _, incident := range incidents {
		resolve := &pagerduty.Incident{
			APIObject: pagerduty.APIObject{Type: "incident_reference",
				ID: incident.ID},
			Service:            incident.Service,
			LastStatusChangeAt: incident.LastStatusChangeAt,
			LastStatusChangeBy: incident.LastStatusChangeBy,
			EscalationPolicy:   incident.EscalationPolicy,
			Status:             "resolved",
		}
		incidentSlice = append(incidentSlice, *resolve)
	}
	e := ctx.ManageIncidents(config.UserEmail, incidentSlice)
	if e != nil {
		log.Fatalln(e)
	}
}

func getTriggeredIncidents(ctx *pagerduty.Client, incidentKey string) (*pagerduty.ListIncidentsResponse, error) {
	var statuses []string
	incidentsOptions := &pagerduty.ListIncidentsOptions{
		APIListObject: pagerduty.APIListObject{Limit: 100},
		IncidentKey:   incidentKey,
		Statuses:      append(statuses, "triggered"),
	}
	incidents, err := ctx.ListIncidents(*incidentsOptions)
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

func (c pdConfig) getConf(yamlPath string) pdConfig {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}
