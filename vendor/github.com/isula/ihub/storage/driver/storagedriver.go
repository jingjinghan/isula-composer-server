package driver

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/astaxie/beego/context"
)

// StorageDriver defines methods that a Storage Driver must implement for a
// filesystem-like key/value object storage. Storage Drivers are automatically
// registered via an internal registration mechanism, and generally created
// via the StorageDriverFactory interface (https://godoc.org/github.com/docker/distribution/registry/storage/driver/factory).
// Please see the aforementioned factory package for example code showing how to get an instance
// of a StorageDriver
type StorageDriver interface {
	// Init updates the parameters
	Init(parameters map[string]interface{}) error

	// Valid checks if parameters are sufficient parameters
	Valid(parameters map[string]interface{}) error

	// Name returns the human-readable "name" of the driver, useful in error
	// messages and logging. By convention, this will just be the registration
	// name, but drivers may provide other information here.
	Name() string

	// GetContent retrieves the content stored at "path" as a []byte.
	// This should primarily be used for small objects.
	GetContent(ctx context.Context, path string) ([]byte, error)

	// PutContent stores the []byte content at a location designated by "path".
	// This should primarily be used for small objects.
	PutContent(ctx context.Context, path string, content []byte) error

	// Reader retrieves an io.ReadCloser for the content stored at "path"
	// with a given byte offset.
	// May be used to resume reading a stream by providing a nonzero offset.
	Reader(ctx context.Context, path string, offset int64) (io.ReadCloser, error)

	// Writer returns a FileWriter which will store the content written to it
	// at the location designated by "path" after the call to Commit.
	Writer(ctx context.Context, path string, append bool) (FileWriter, error)

	// Stat retrieves the FileInfo for the given path, including the current
	// size in bytes and the creation time.
	Stat(ctx context.Context, path string) (FileInfo, error)

	// List returns a list of the objects that are direct descendants of the
	//given path.
	List(ctx context.Context, path string) ([]string, error)

	// Move moves an object stored at sourcePath to destPath, removing the
	// original object.
	// Note: This may be no more efficient than a copy followed by a delete for
	// many implementations.
	Move(ctx context.Context, sourcePath string, destPath string) error

	// Delete recursively deletes all objects stored at "path" and its subpaths.
	Delete(ctx context.Context, path string) error
}

// FileWriter provides an abstraction for an opened writable file-like object in
// the storage backend. The FileWriter must flush all content written to it on
// the call to Close, but is only required to make its content readable on a
// call to Commit.
type FileWriter interface {
	io.WriteCloser

	// Size returns the number of bytes written to this FileWriter.
	Size() int64

	// Cancel removes any written content from this FileWriter.
	Cancel() error

	// Commit flushes all content written to this FileWriter and makes it
	// available for future calls to StorageDriver.GetContent and
	// StorageDriver.Reader.
	Commit() error
}

// PathRegexp is the regular expression which each file path must match. A
// file path is absolute, beginning with a slash and containing a positive
// number of path components separated by slashes, where each component is
// restricted to alphanumeric characters or a period, underscore, or
// hyphen.
var PathRegexp = regexp.MustCompile(`^(/[A-Za-z0-9._-]+)+$`)

// ErrUnsupportedMethod may be returned in the case where a StorageDriver implementation does not support an optional method.
type ErrUnsupportedMethod struct {
	DriverName string
}

func (err ErrUnsupportedMethod) Error() string {
	return fmt.Sprintf("%s: unsupported method", err.DriverName)
}

// PathNotFoundError is returned when operating on a nonexistent path.
type PathNotFoundError struct {
	Path       string
	DriverName string
}

func (err PathNotFoundError) Error() string {
	return fmt.Sprintf("%s: Path not found: %s", err.DriverName, err.Path)
}

// InvalidPathError is returned when the provided path is malformed.
type InvalidPathError struct {
	Path       string
	DriverName string
}

func (err InvalidPathError) Error() string {
	return fmt.Sprintf("%s: invalid path: %s", err.DriverName, err.Path)
}

// InvalidOffsetError is returned when attempting to read or write from an
// invalid offset.
type InvalidOffsetError struct {
	Path       string
	Offset     int64
	DriverName string
}

func (err InvalidOffsetError) Error() string {
	return fmt.Sprintf("%s: invalid offset: %d for path: %s", err.DriverName, err.Offset, err.Path)
}

// Error is a catch-all error type which captures an error string and
// the driver type on which it occurred.
type Error struct {
	DriverName string
	Enclosed   error
}

func (err Error) Error() string {
	return fmt.Sprintf("%s: %s", err.DriverName, err.Enclosed)
}

var (
	sdsLock sync.Mutex
	sds     = make(map[string]StorageDriver, 8)
)

// Register regists a storage driver to the system driver map
func Register(name string, driver StorageDriver) error {
	if name == "" {
		return errors.New("Could not register a Storage with an empty name")
	}
	if driver == nil {
		return errors.New("Could not register a nil Storage")
	}

	sdsLock.Lock()
	defer sdsLock.Unlock()

	if _, exists := sds[name]; exists {
		return fmt.Errorf("Storage type '%s' is already registered", name)
	}

	sds[name] = driver
	return nil
}

// FindDriver returns the storage driver by its name and parameters
func FindDriver(name string, params map[string]interface{}) (StorageDriver, error) {
	if d, ok := sds[name]; ok {
		err := d.Valid(params)
		return d, err
	}
	return nil, fmt.Errorf("Cannot find '%s' storage driver", name)
}
