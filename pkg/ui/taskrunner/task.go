package taskrunner

type Task struct {
	Plugin      Executer
	Description string
	Root        string
	Args        []string
}

type Tasks []Task
