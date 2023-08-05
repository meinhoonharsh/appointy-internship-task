package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	// "time"
)

type Resource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BusinessHour struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	Quantity   int64  `json:"quantity"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type BlockHour struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type Appointment struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	Quantity   int64  `json:"quantity"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type Duration struct {
	Seconds int64 `json:"seconds"`
}

// endpoint request structs

type ListBusinessHoursRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type ListBlockHoursRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type ListAppointmentRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func TimeToString(tm time.Time) string {
	return tm.Format(time.RFC3339)
}

func StringToTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func main() {

	// - resourceId [Required]: ID of the pitch
	// - date [Required]: date in YYYY-MM-DD format
	// - duration [Required]: time duration in minutes (e.g., 30, 60, 120)
	// - quantity [Required]:  capacity to reserve

	inputParam := map[string]interface{}{
		"resourceId": "res_2",
		"date":       "2023-08-05",
		"duration":   "30",
		"quantity":   "1",
	}

	// Create startTime and EndTime in format YYYY-MM-DDTHH:mm:ss.sssZ from inputParam
	resourceId := inputParam["resourceId"].(string)
	startTime := inputParam["date"].(string) + "T00:00:00Z"
	endTime := inputParam["date"].(string) + "T23:59:00Z"

	// declare payload
	payload := map[string]interface{}{
		"resourceId": resourceId,
		"startTime":  startTime,
		"endTime":    endTime,
	}

	businesshours := apiCall("/business-hours", payload)
	blockhours := apiCall("/block-hours", payload)
	appointment := apiCall("/appointments", payload)

	// change businesshours from string to interface
	var businesshoursInterface interface{}
	json.Unmarshal([]byte(businesshours), &businesshoursInterface)

	var blockhoursInterface interface{}
	json.Unmarshal([]byte(blockhours), &blockhoursInterface)

	var appointmentInterface interface{}
	json.Unmarshal([]byte(appointment), &appointmentInterface)

	// check for availability of resource on given date for given duration
	// if available, create appointment
	// else return error

	// Avaible time slots for given date
	// {start_time: 10:00 am, end_time: 11:00 am, quantity: 5)}

}

// create function for API Call with param endpoint and payload (object)
func apiCall(endpoint string, payload map[string]interface{}) string {

	url := "http://api.internship.appointy.com:8000/v1"
	method := "GET"

	newurl := url + endpoint

	// Add get parameters to url
	if payload != nil {
		newurl = newurl + "?"
		for key, value := range payload {
			newurl = newurl + key + "=" + value.(string) + "&"
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, newurl, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIzLTA4LTEwVDAwOjAwOjAwWiIsInVzZXJfaWQiOjMwMDF9.8pZMhoqZdBLqOKT0V7perD4vkoA347idSHVLaCcdefs")
	req.Header.Add("Content-Type", "application/json")

	// send request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
