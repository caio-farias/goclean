package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/robfig/cron"
)

type Scheduler struct {
	cron *cron.Cron
	jobs []CronJob
}

type CronJob struct {
	Id      int    `json:"-"`
	Name    string `json:"name"`
	CronExp string `json:"cron_exp"`
	Command string `json:"command"`
}

func (cj CronJob) Run() {
	log.Printf("should run this at %s the command: %s", cj.CronExp, cj.Command)
	WriteLog(fmt.Sprintf("should run this at %s the command: %s", cj.CronExp, cj.Command))
	cmd := exec.Command("bash", "-c", cj.Command)
	if cmd.Err != nil {
		log.Println(cmd.Err)
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
	return &Scheduler{
		cron: cron.New(),
		jobs: jobs,
	}
}

func (schdl *Scheduler) Init() {
	for _, j := range schdl.jobs {
		schdl.cron.AddJob(j.CronExp, j)
	}
}

func (schdl *Scheduler) add(cron_exp string, callback func()) {
	schdl.cron.AddFunc(cron_exp, callback)
}

func executeCmd(cmd string) {
	log.Println(cmd)
}
