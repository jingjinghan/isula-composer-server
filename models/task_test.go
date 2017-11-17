package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryTaskListByUser(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		username string
		len      int
		expected bool
	}{
		{"isula", 2, true},
		{"non-exist", 0, true},
	}

	for _, c := range cases {
		tasks, err := QueryTaskListByUser(c.username)
		assert.Equal(t, c.expected, err == nil)
		assert.Equal(t, c.len, len(tasks), "fail to match list length")
	}
}

func TestQueryTaskByID(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		id       int64
		exist    bool
		username string
		expected bool
	}{
		{1, true, "isula", true},
		{10000, false, "", true},
	}

	for _, c := range cases {
		task, err := QueryTaskByID(c.id)
		assert.Equal(t, c.exist, task != nil)
		assert.Equal(t, c.expected, err == nil)
		if c.exist {
			assert.Equal(t, task.UserName, c.username, "fail to get correct username")
		}
	}
}

func TestAddTaskFull(t *testing.T) {
	if !testReady {
		return
	}

	username := "testaddtask"

	id, err := AddTaskFull(username, "", "", "")
	assert.Nil(t, err)

	task, err := QueryTaskByID(id)
	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, username, task.UserName)
	assert.Equal(t, TaskStatusNew, task.Status)
}

func TestUpdateTask(t *testing.T) {
	task := &Task{}

	task.ID = 1
	task.Status = TaskStatusFinish
	err := UpdateTask(task)
	assert.Nil(t, err)
	t1, err := QueryTaskByID(task.ID)
	assert.Nil(t, err)
	assert.NotNil(t, t1)
	assert.Equal(t, task.Status, t1.Status)

	task.ID = 10000
	err = UpdateTask(task)
	assert.NotNil(t, err)
}

func TestRemoveTask(t *testing.T) {
	testID := 2
	num, err := RemoveTask(int64(testID))
	assert.Nil(t, err)
	assert.Equal(t, int64(1), num)

	num, err = RemoveTask(int64(testID))
	assert.Nil(t, err)
	assert.Equal(t, int64(0), num)
}
