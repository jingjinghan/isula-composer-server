# isula-composer-server [![Build Status](https://travis-ci.org/isula/isula-composer-server.svg?branch=master)](https://travis-ci.org/isula/isula-composer-server) [![Go Report Card](https://goreportcard.com/badge/github.com/isula/isula-composer-server)](https://goreportcard.com/report/github.com/isula/isula-composer-server) [![codecov.io](https://codecov.io/github/isula/isula-composer-server/coverage.svg?branch=master)](https://codecov.io/github/isula/isula-composer-server?branch=master)


isula-composer-server: a scheduler to collect build task and run tasks on the server, the user could query the build status and download the output.


## APIs

|Method|Path|Summary|Description|
|------|----|------|-----------|
| POST | `/:user/task` | [Create a task](api.md#task "Create a build task") | create a task, the server will start to run the build server |
| PUT | `/:user/task/:id` | [Operate on a task](api.md#operate "Send a command to a task") | restart or edit and restart a task |
| DELETE | `/:user/task/:id` | [Delete a task](api.md#remove-task "Remove a task") | remove a task, clean all the output |
| GET | `/:user/task` | [List user's task](api.md#task-list "List all the tasks, including status") | list tasks |
| GET  | `/:user/task/:id` | [Task Status](api.md#task-status "Get the task status") | get the task status by id, also tells the output url |
| GET  | `/:user/task/:id/filename` | [Download](api.md#output "Download the output file") | get the output file |
| POST | `/:user/hook` | [Add a hook](api.md#hook "Add a hook to user") | add a hook to a user |
| GET | `/:user/hook` | [List hooks](api.md#hook-list "List hooks") | list hooks |
| GET | `/:user/hook/:id` | [Hook information](api.md#hook-information "Get the hook information") | get the hook information |
| DELETE | `/:user/hook/:id` | [Delete a hook](api.md#remove-hook "Remove a hook") | remove a hook |

## Reused isula library
There are several useful wrappers in `isula/ihub` project, like storage/logger/config, `isula-composer-server` will reuse them.
If there is any bug fix of these libs, please submit PR in `isula/ihub` and fix there. 
`isula-composer-server` will sync with the updates.

