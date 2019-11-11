// Carsten's CESSDA CAFE Coffee Machine
// Copyright CESSDA-ERIC 2019
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	//"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type job struct {
	ID           string `json:"jobId"`
	JobStarted   string `json:"jobStarted"`
	Product      string `json:"product"`
	JobReady     string `json:"jobReady"`
	JobRetrieved string `json:"jobRetrieved"`
}

// Initial coffe jobs
var jobList = []job{
	job{ID: "c1be03bf-d9cc-486b-92af-3d91c27d3ba5", Product: "COFFEE", JobStarted: "2019-02-16T11:31:47+0000", JobReady: "2019-02-16T11:32:34+0000", JobRetrieved: "2019-02-16T10:33:00+0000"},
	job{ID: "90fcb5bd-0f08-4656-8306-4e1efaaea2b0", Product: "CAPPUCCINO", JobStarted: "2019-02-16T10:31:47+0000", JobReady: "2019-02-16T10:32:34+0000", JobRetrieved: "2019-02-16T10:33:00+0000"},
}

// return entire job history
func getAllJobs() []job {
	return jobList
}

// check whether a coffee is still brewing
func systemBrewing() bool {
	for _, o := range jobList {
		readytime, _ := time.Parse(time.RFC3339, o.JobReady)
		if time.Now().Before(readytime) {
			return true
		}
	}
	return false
}

// check whether a coffee needs picking up
func jobWaiting() bool {
	for _, o := range jobList {
		if len(o.JobRetrieved) == 0 {
			return true
		}
	}
	return false
}

// check overall system status
func systemStatus() (int, int, string) {
	var systemStatusCode int
	var systemHTTPStatusCode int
	var systemStatusMessage string

	if systemBrewing() {
		systemStatusCode = 1
		systemHTTPStatusCode = http.StatusConflict
		systemStatusMessage = "System Brewing -- Please wait!"
	} else if jobWaiting() {
		systemStatusCode = 2
		systemHTTPStatusCode = http.StatusUnauthorized
		systemStatusMessage = "Coffee Waiting -- Come and get it!"
	} else {
		systemStatusCode = 0
		systemHTTPStatusCode = http.StatusOK
		systemStatusMessage = "System Ready!"
	}

	return systemStatusCode, systemHTTPStatusCode, systemStatusMessage

}

// set a sepcific job to retrieved if it`s done but still waiting
func retrieveJob(id string) (*job, bool) {
	for index, o := range jobList {
		if o.ID == id {
			// only retrieve when done and only once
			readytime, _ := time.Parse(time.RFC3339, o.JobReady)
			if time.Now().After(readytime) && len(o.JobRetrieved) == 0 {
				o.JobRetrieved = time.Now().Format(time.RFC3339)
				jobList[index].JobRetrieved = o.JobRetrieved
				return &o, true
			}
			return &o, false
		}
	}
	return nil, false
}

// return a job
func getJobbyID(id string) (*job, bool) {
	for _, o := range jobList {
		if o.ID == id {
			return &o, true
		}
	}
	return nil, false
}

// create a new coffee job
func newJob(ID string, Product string) (*job, bool, string) {

	systemStatusCode, _, systemStatusMessage := systemStatus()

	if !(systemStatusCode == 0) {
		return nil, false, systemStatusMessage
	}

	//myjobid, _ := uuid.NewV4()

	var newJob job
	newJob.ID = ID
	newJob.Product = Product
	newJob.JobStarted = time.Now().Format(time.RFC3339)
	newJob.JobReady = time.Now().Add(time.Minute * 1).Format(time.RFC3339)

	jobList = append(jobList, newJob)

	theNewJob, success := getJobbyID(newJob.ID)
	return theNewJob, success, systemStatusMessage

}
