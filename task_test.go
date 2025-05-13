package main

import (
//	"encoding/hex"
	"testing"
//	"flag"
	"log"
)

/*
type cmdF_T struct {
	ax int64
	bx int64
	as string
	bs string
}
*/

// var cmdF *cmdF_T
var tasks = make([]*task_T, 0, 10)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
//        cmdF = &cmdF_T{}
//        readFs(cmdF)
}

func TestTaskDeploy(t *testing.T) {
	task := task_T{}
	task.instructs = []byte{ins_movsb, 0, 0, 6, 0, 5, 0}
	task.vData = []byte("Hello task!")
	task.deploy()
	for k, task := range tasks {
		t.Log(k)
		t.Log(task.hash())
		t.Log(string(task.vData))
	}
}

func TestTaskCall(t *testing.T) {
	for k, task := range tasks {
		task.excute()
		t.Log(k)
		t.Log(task.hash())
		t.Log(string(task.vData))
	}
}

func TestBoth(t *testing.T) {
	TestTaskDeploy(t)
	TestTaskCall(t)
}

/*
func readFs(cmdF *cmdF_T) {
        flag.IntVar(&cmdF.ax, "ax", 0, "---")
        flag.IntVar(&cmdF.bx, "bx", 0, "---")
        flag.StringVar(&cmdF.as, "as", "", "---")
        flag.StringVar(&cmdF.bs, "bs", "", "---")
}
*/
