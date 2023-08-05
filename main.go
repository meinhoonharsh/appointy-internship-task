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

type Slot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
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
		"quantity":   4,
	}

	// Create startTime and EndTime in format YYYY-MM-DDTHH:mm:ss.sssZ from inputParam
	resourceId := inputParam["resourceId"].(string)
	startTime := inputParam["date"].(string) + "T00:00:00Z"
	endTime := inputParam["date"].(string) + "T23:59:00Z"
	// quantity is number of slots to be booked
	quantity := int64(inputParam["quantity"].(int))

	// declare payload
	payload := map[string]interface{}{
		"resourceId": resourceId,
		"startTime":  startTime,
		"endTime":    endTime,
	}

	businesshours := apiCall("/business-hours", payload)
	blockhours := apiCall("/block-hours", payload)
	appointment := apiCall("/appointments", payload)

	// change businesshours from string to Maps of BusinessHour struct
	var businesshoursMap []BusinessHour
	json.Unmarshal([]byte(businesshours), &businesshoursMap)

	var blockhoursMap []BlockHour
	json.Unmarshal([]byte(blockhours), &blockhoursMap)

	var appointmentMap []Appointment
	json.Unmarshal([]byte(appointment), &appointmentMap)

	var availableSlots []Slot

	// check for availability of resource on given date for given duration
	for i := 0; i < len(businesshoursMap); i++ {
		// fmt.Println(businesshoursMap[i].StartTime)
		// fmt.Println(businesshoursMap[i].EndTime)

		// convert string to time
		startTime, _ := StringToTime(businesshoursMap[i].StartTime)
		endTime, _ := StringToTime(businesshoursMap[i].EndTime)

		// convert duration from string to int64
		duration, _ := time.ParseDuration(inputParam["duration"].(string) + "m")

		// fmt.Println("Business Hours: ", i+1)
		// fmt.Println("Start Time: ", startTime, "End Time: ", endTime)
		// fmt.Println("Duration: ", duration)

		// check all available slots in businesshours

		for j := startTime; j.Before(endTime); j = j.Add(duration) {
			var available bool = true

			// fmt.Println("Slot: ", j, "to", j.Add(duration))

			// check if j is in blockhours
			for k := 0; k < len(blockhoursMap); k++ {
				// convert string to time
				blockStartTime, _ := StringToTime(blockhoursMap[k].StartTime)
				blockEndTime, _ := StringToTime(blockhoursMap[k].EndTime)

				if (j.Equal(blockStartTime) || j.After(blockStartTime)) && (j.Before(blockEndTime) || j.Before(blockEndTime)) {

					available = false
					break
				}
			}

			// check if j is in appointment hours and quantity is available
			for l := 0; l < len(appointmentMap); l++ {
				// convert string to time
				appointmentStartTime, _ := StringToTime(appointmentMap[l].StartTime)
				appointmentEndTime, _ := StringToTime(appointmentMap[l].EndTime)

				if (j.Equal(appointmentStartTime) || j.After(appointmentStartTime)) && (j.Before(appointmentEndTime) || j.Before(appointmentEndTime)) {

					businessHourQuantity := businesshoursMap[i].Quantity
					appointmentQuantity := appointmentMap[l].Quantity

					// change into int64 and check if quantity is available
					if businessHourQuantity-appointmentQuantity < quantity {
						available = false
						break
					}

				}
			}

			if available {
				fmt.Println("available")
				availableSlots = append(availableSlots, Slot{TimeToString(j), TimeToString(j.Add(duration))})
			} else {
				fmt.Println("blocked")
			}

		}

		fmt.Println("")

	}

	fmt.Println("Available Slots: ", availableSlots)
	// Convert availableSlots to json
	availableSlotsJson, _ := json.Marshal(availableSlots)
	fmt.Println("Available Slots Json: ", string(availableSlotsJson))

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

	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIzLTA4LTEwVDAwOjAwOjAwWiIsInVzZXJfaWQiOjc1fQ.2m_F_k0qqZiHlFduAZbkBX4i2DbngyROw1rYKqvXOwg")
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
