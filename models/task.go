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
	// Scripts is the config file
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

// AddTask adds a task
func AddTask(username string) (int64, error) {
	task := &Task{}

	task.UserName = username
	task.Status = TaskStatusNew
	if id, err := orm.NewOrm().Insert(task); err != nil {
		logs.Error("[AddTask] fail to insert task '%s': %v", username, err)
		return -1, err
	} else {
		return id, nil
	}
}

// UpdateTask updates the output filename or the status
func UpdateTask(task *Task) error {
	t, err := QueryTaskByID(task.ID)
	if err != nil {
		logs.Error("[UpdateTask] fail to querytask %v", err)
		return err
	}

	if t == nil {
		return fmt.Errorf("Cannot find the task %d.", task.ID)
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
	if num, err := orm.NewOrm().Delete(task); err != nil {
		logs.Error("[DeleteTask] %v", err)
		return num, err
	} else {
		return num, nil
	}
}
