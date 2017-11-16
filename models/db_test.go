package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var dbLock sync.Mutex
var testReady = false

func init() {
	err := InitTestDB()
	if err != nil {
		fmt.Println("Fail to test sql,skip all the sql testing: ", err)
		return
	}

	// FIXME: we need a better testing framework to free data after everything is done.
	// Free anyway before init the data
	FreeTestDBData()
	err = InitTestDBData()
	if err != nil {
		fmt.Println("Fail to init data, mostly you need to check your test sql: ", err)
		return
	}

	testReady = true
}

// InitTestDB
//  conn should be something like root:1234@tcp(localhost:3306)/test?charset=utf8
func InitTestDB() error {
	conn := os.Getenv("TESTCONN")
	if conn == "" {
		return errors.New("Please set the 'TESTCONN' before using db unit test")
	}

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.DefaultTimeLoc = time.UTC

	err := orm.RegisterDataBase("default", "mysql", conn)
	if err == nil {
		orm.RunSyncdb("default", false, true)
		return nil
	}

	return errors.New("Cannot connect to the test database")
}

func RunScript(script string) error {
	scriptsql, err := ioutil.ReadFile(script)
	if err != nil {
		return fmt.Errorf("Failed to init test data: %v", err)
	}

	// Seems orm can only exec one sql at one time.
	for _, sql := range strings.SplitAfter(string(scriptsql), ";") {
		if len(strings.TrimSpace(sql)) == 0 {
			continue
		}
		_, err := orm.NewOrm().Raw(sql).Exec()
		if err != nil {
			fmt.Printf("Fail to exec '%s': %v", sql, err)
			return err
		}
	}

	return nil
}

func InitTestDBData() error {
	err := RunScript("testdata/init.sql")
	if err != nil {
		FreeTestDBData()
		return fmt.Errorf("Failed to exec init sql: %v", err)
	}

	return nil
}

func AlterTable() error {
	return RunScript("testdata/alter_table.sql")
}

func RecoverTable() error {
	return RunScript("testdata/recover_table.sql")
}

func FreeTestDBData() error {
	return RunScript("testdata/free.sql")
}

func TestInitDB(t *testing.T) {
	if !testReady {
		return
	}

	cases := []struct {
		conn     string
		driver   string
		name     string
		expected bool
	}{
		{os.Getenv("TESTCONN"), "mysql", "TestInitDB-case-0", true},
		{os.Getenv("TESTCONN"), "liangsql", "TestInitDB-case-1", false},
		{"localhost:22", "mysql", "TestInitDB-case-2", false},
	}

	for _, c := range cases {
		err := InitDB(c.conn, c.driver, c.name)
		assert.Equal(t, c.expected, err == nil)
	}
}
