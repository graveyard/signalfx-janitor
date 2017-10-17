package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const baseURL = "https://api.signalfx.com/"

var sfxToken = envOrDie("SFX_TOKEN")
var sfxOrgID = envOrDie("SFX_ORG_ID")

func envOrDie(s string) string {
	value := os.Getenv(s)
	if value == "" {
		log.Fatalf("env var %s is required", s)
	}
	return value
}

func main() {
	incidents, err := GetV1Incidents()
	if err != nil {
		log.Fatal("error looking up incidents:", err.Error())
	}

	log.Printf("Found %d incidents\n", len(incidents))

	err = resolveIncidents(incidents)
	if err != nil {
		log.Fatal("error resolving incidents:", err.Error())
	}
}

type SimpleIncident struct {
	Label     string
	ID        string
	CreatedAt time.Time
}

func (si SimpleIncident) String() string {
	timeAgo := time.Now().Sub(si.CreatedAt)
	return fmt.Sprintf("%s (time ago = %s)", si.Label, timeAgo)
}

func GetV1Incidents() ([]SimpleIncident, error) {
	eventTimeSeries, err := listActiveIncidentsV1()
	if err != nil {
		return []SimpleIncident{}, err
	}

	incidents := []SimpleIncident{}
	for _, series := range eventTimeSeries {
		updatedAt := time.Unix(int64(series.UpdatedOnMs/1000), 0)
		label := fmt.Sprint(series.SfDetector, " -- ", series.SfDetectorID)
		incidents = append(incidents, SimpleIncident{
			ID:        series.IncidentID,
			CreatedAt: updatedAt,
			Label:     label,
		})
	}

	return incidents, nil
}

func resolveIncidents(incidents []SimpleIncident) error {
	for _, i := range incidents {
		log.Println("Incident:", i)
		shouldAutoResolve := i.CreatedAt.Before(time.Now().Add(-30 * time.Minute))
		log.Println("Should auto resolve:", shouldAutoResolve)
		if shouldAutoResolve {
			err := clearIncident(i.ID)
			if err != nil {
				return fmt.Errorf("error resolving incident %s: %s ", i.ID, err.Error())
			}
		}
		log.Println("")
	}

	return nil
}

// EventTimeSeries (V1 API)
type EventTimeSeries struct {
	RS []EventTimeSeriesRS `json:"rs"`
}

// EventTimeSeries (V1 API)
type EventTimeSeriesRS struct {
	IncidentID   string  `json:"sf_incidentId"`
	UpdatedOnMs  float64 `json:"sf_updatedOnMs"`
	SfDetector   string  `json:"sf_detector"`
	SfDetectorID string  `json:"sf_detectorId"`
}

func listActiveIncidentsV1() ([]EventTimeSeriesRS, error) {
	url := baseURL + "v1/eventtimeseries"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []EventTimeSeriesRS{}, err
	}

	// Add query params
	q := req.URL.Query()
	q.Add("query", `sf_organizationID:`+sfxOrgID+` AND (NOT sf_archived:true) AND ((((sf_anomalyState:("anomalous" "too high" "too low"))) AND (sf_detector.lowercase:* OR sf_displayName.lowercase:*)))`)
	// TODO: Properly page through results
	q.Add("offset", strconv.Itoa(0))
	q.Add("limit", strconv.Itoa(500))
	q.Add("order_by", `-sf_priority,-sf_anomalyStateUpdateTimestampMs`)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-SF-TOKEN", sfxToken)
	resp, err := client.Do(req)
	if err != nil {
		return []EventTimeSeriesRS{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []EventTimeSeriesRS{}, err
	}
	s := new(EventTimeSeries)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return []EventTimeSeriesRS{}, err
	}
	return s.RS, nil
}

// clearIncident works for V1 and V2 detectors
// https://developers.signalfx.com/v2/reference#incidentidclear
func clearIncident(incidentID string) error {
	url := baseURL + "v2/incident/" + incidentID + "/clear"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-SF-TOKEN", sfxToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Println("error:", string(body))
		return fmt.Errorf("Error clearing incident %s, got StatusCode %d", incidentID, resp.StatusCode)
	}

	return nil
}
