// dsysb

package main

import (
	"encoding/hex"
	"net/http"
)

type task_J struct {
	Address string `json:"address"`
	Instructs string `json:"instructs"`
	VData string `json:"vData"`
}

type taskPool_J []*task_J

func tasks2Handler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult2(w, responseResult2_T{false, "dsysb inner error", nil})
		return
	}

	tLength := len(state.tasks)
	tasks := make(taskPool_J, tLength, tLength)
	for k, task := range state.tasks {
		tasks[k] = &task_J {
			task.address,
			hex.EncodeToString(task.instructs[:]),
			hex.EncodeToString(task.vData[:]),
		}
	}

	writeResult2(w, responseResult2_T{true, "ok", tasks})
}

func task2Handler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	values := req.URL.Query()
	taskId := values.Get("id")

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult2(w, responseResult2_T{false, "dsysb inner error", nil})
		return
	}

	for _, task := range state.tasks {
		h := task.hash()
		tId := hex.EncodeToString(h[:])

		if tId == taskId {
			taskJ := task_J {
				task.address,
				hex.EncodeToString(task.instructs[:]),
				hex.EncodeToString(task.vData[:]),
			}
			writeResult2(w, responseResult2_T{true, "ok", taskJ})
			return
		}
	}

	writeResult2(w, responseResult2_T{false, "task " + taskId + " does not exist", nil})
}
