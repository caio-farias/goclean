package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron         *cron.Cron
	jobsMap      map[string]CronJob
	cronFilePath string
}

type CronJob struct {
	id      cron.EntryID
	key     string `json:"-"`
	Name    string `json:"name"`
	CronExp string `json:"cronExp"`
	Command string `json:"command"`
}

func (cj CronJob) setEntryId(id cron.EntryID) {
	cj.id = id
}

func (cj CronJob) setKey() {
	cj.key = cj.Name + cj.CronExp + cj.Command
}

func (cj CronJob) Run() {
	var outputBuffer bytes.Buffer

	logJob := fmt.Sprintf("goclean -- Running [%s]", cj.Command)

	log.Println(logJob)
	WriteLog(logJob)

	cmd := exec.Command("bash", "-c", cj.Command)
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	err := cmd.Run()
	if err != nil {
		logJob = fmt.Sprintf("goclean -- Error on job %s runnning [%s]: \n%s", cj.Name, cj.Command, outputBuffer.String())
		log.Println(logJob)
		WriteLog(logJob)
		return
	}
	if len(outputBuffer.String()) > 0 {
		logJob = fmt.Sprintf("goclean -- [%s] output: \n%s", cj.Name, outputBuffer.String())
		log.Println(logJob)
		WriteLog(logJob)
	}
}

func NewScheduler(file_path string) *Scheduler {
	file_bytes, err := os.ReadFile(file_path)
	if err != nil {
		log.Panicln("could not read cronn file", err)
		return nil
	}

	var jobs []CronJob
	json.Unmarshal(file_bytes, &jobs)

	jobsMap := make(map[string]CronJob)
	for _, j := range jobs {
		j.setKey()
		jobsMap[j.key] = j
	}

	return &Scheduler{
		cronFilePath: file_path,
		cron:         cron.New(),
		jobsMap:      jobsMap,
	}
}

func (schdl *Scheduler) getCurrentJobs() (map[string]CronJob, error) {
	file_bytes, err := os.ReadFile(schdl.cronFilePath)
	if err != nil {
		log.Panicln("could not read cronn file", err)
		return nil, err
	}

	var jobs []CronJob
	json.Unmarshal(file_bytes, &jobs)

	jobsMap := make(map[string]CronJob)
	for _, j := range jobs {
		j.setKey()
		jobsMap[j.key] = j
	}
	return jobsMap, nil
}

func (schdl *Scheduler) updateJobs(currentJobsMap map[string]CronJob) {
	for key, cj := range currentJobsMap {
		_, ok := schdl.jobsMap[key]
		if !ok {
			id, err := schdl.cron.AddJob(cj.CronExp, cj)
			if err != nil {
				log.Println("something wird happenede")
			}
			cj.setEntryId(id)
		}
	}

	for k := range schdl.jobsMap {
		j, ok := currentJobsMap[k]
		if !ok {
			schdl.cron.Remove(j.id)
		}
	}
}

func (schdl *Scheduler) Init(wg *sync.WaitGroup) {
	defer wg.Done()

	schdl.updateJobs(schdl.jobsMap)
	schdl.cron.Start()

	for {
		time.Sleep(10 * time.Second)
		cjobs, err := schdl.getCurrentJobs()
		if err != nil {
			log.Fatalf("Something is wrong with .goclean.cron.json")
		}
		schdl.updateJobs(cjobs)
		log.Printf("Jobs scheduled: %d", len(schdl.jobsMap))
	}
}
