package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	// TaskStatusNew represents a task is new
	TaskStatusNew = "new"
	// TaskStatusRunning represents a task is running
	TaskStatusRunning = "running"
	// TaskStatusFinish represents a task is finish succesfully
	TaskStatusFinish = "finish"
	// TaskStatusFailed represents a task is failed
	TaskStatusFailed = "failed"
)

// Task defines the task information
type Task struct {
	ID int64 `orm:"column(id);auto"`

	// A user table will be created and UserName will be unique.
	UserName string `orm:"column(username);size(255);null"`
	Status   string `orm:"column(status);size(32);null"`

	// OutputFile is the filename of the building iso/image
	OutputFile string `orm:"column(output_file);size(255);null"`
	// ConfigFile is the yml file of the building configuration
	// FIXME: is 4096 big enough? I prefer NOT to save it on a storage.
	ConfigFile string `orm:"column(config_file);size(4096);null`
	// Scripts is the commands
	// It should be empty if the builder server is just build iso/image,
	// we provide default build command in that case.
	// So, either a user provides `outputFile and configFile`
	// or he/she provides a `scripts`.
	Scripts string `orm:"column(scripts);size(4096);null"`
}

var taskModels = []interface{}{
	new(Task),
}

func init() {
	orm.RegisterModel(taskModels...)
}

const (
	queryTaskListByUser = `select * from task where username=?`
	queryTaskByID       = `select * from task where id=?`
)

// QueryTaskListByUser returns the task list by 'username'
func QueryTaskListByUser(username string) ([]Task, error) {
	var tasks []Task
	_, err := orm.NewOrm().Raw(queryTaskListByUser, username).QueryRows(&tasks)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryTaskListByUser] %s", err)
		return nil, err
	}

	return tasks, nil
}

// QueryTaskByID returns a task by its id
func QueryTaskByID(id int64) (*Task, error) {
	var tasks []Task

	_, err := orm.NewOrm().Raw(queryTaskByID, id).QueryRows(&tasks)

	if err != nil && err != orm.ErrNoRows {
		logs.Error("[QueryTaskByID] %v", err)
		return nil, err
	}

	if len(tasks) == 0 {
		logs.Debug("[QueryTaskByID] cannot find the row.")
		return nil, nil
	}

	return &tasks[0], nil
}

// AddTaskFull adds a task
func AddTaskFull(username, output, config, script string) (int64, error) {
	task := &Task{}
	task.UserName = username
	task.Status = TaskStatusNew
	task.OutputFile = output
	task.ConfigFile = config
	task.Scripts = script

	id, err := orm.NewOrm().Insert(task)
	if err != nil {
		logs.Error("[AddTask] fail to insert task '%s': %v", username, err)
		return -1, err
	}

	return id, nil
}

// UpdateTask updates the output filename or the status
func UpdateTask(task *Task) error {
	t, err := QueryTaskByID(task.ID)
	if err != nil {
		logs.Error("[UpdateTask] fail to querytask %v", err)
		return err
	}

	if t == nil {
		return fmt.Errorf("cannot find the task %d", task.ID)
	}

	if _, err := orm.NewOrm().Update(task); err != nil {
		logs.Error("[UpdateTask] %v", err)
		return err
	}
	return nil
}

// RemoveTask removes a task from db
func RemoveTask(id int64) (int64, error) {
	task := &Task{}
	task.ID = id
	num, err := orm.NewOrm().Delete(task)
	if err != nil {
		logs.Error("[DeleteTask] %v", err)
		return num, err
	}
	return num, nil
}
