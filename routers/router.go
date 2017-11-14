package router

import (
	"errors"
	"fmt"
	"sync"

	"github.com/astaxie/beego"
)

var (
	nssLock sync.Mutex
	nss     = make(map[string]*beego.Namespace, 8)
)

// RegisterRouter regists a router by its name and its namespace
func RegisterRouter(name string, ns *beego.Namespace) error {
	if name == "" {
		return errors.New("cannot register a namespace without a name")
	}

	if ns == nil {
		return errors.New("cannot register a nil namespace")
	}

	nssLock.Lock()
	defer nssLock.Unlock()

	if _, existed := nss[name]; existed {
		return fmt.Errorf("namespace '%s' is already exist", name)
	}
	nss[name] = ns

	return nil
}

// GetNamespaces gets the namespaces of a router
func GetNamespaces() map[string]*beego.Namespace {
	return nss
}
