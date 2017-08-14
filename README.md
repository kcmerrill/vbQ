[![Build Status](https://travis-ci.org/kcmerrill/vbQ.svg?branch=master)](https://travis-ci.org/kcmerrill/vbQ)[![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/vbQ)](https://goreportcard.com/report/github.com/kcmerrill/vbQ)
![vbQ](assets/vbQ.png "vbQ")

At it's core vbQ is a simple flat file queue system. Queues are just folders, tasks are just files inside of said folders. Workers are bash commands so you can be as creative as you need to be. Tasks will continue to be processed so long as tasks are available.

# Why?

If you have a folder with a bunch of bash scripts that need to be manually run or perhaps you have a wiki with instructions you have to follow then you'll understand why. Any tasks that can be automated should be automated, and a lot of the times you get to do it simply because you have the permissions to do so. Need to setup that new github repo? Need to add a new chef user? What about adding those public ssh keys to test/stage/prod boxes? Automate it!

# The vision

Everything is automated and everything is automated through a `pull request`.

## vbQ Configuration file(default: '.vbQ', optional)

vbQ main configuration file is just a yaml file. It's completely optional if you don't wish to use VCS. The only 2 keys available are `startup`, `shutdown`. 

1. `startup` is a bash command that gets run upon startup and before tasks get run
1. `shutdown` is a bash command that gets run upon shutdown and after tasks get run

An example `.vbQ` file.

```yaml
startup: |
    git reset HEAD --hard
    git clean -fd
    git pull

shutdown: |
    git add .
    git commit -m "Finished"
    git push origin master
```

Of course this is just a trivial example but you can see how one might use, say github, as the queue backup.

## Queue Configuration files(default: '.q')

Queues are just folders. Inside each folder, there is a `.q` file(configurable via cli arguments). Here is a fully baked `.q` file with comments. 

```yaml
queue:
    # when a task gets completed succesfully(0 exit code), the folder to send it to 
    completed: ".completed"
    # when a task gets fails(exit code != 0) send to this folder. By default it stays put and is the current directory
    failed: ".failed" 
workers:
    # should we display stdout/stderr's for each task? Default is false
    verbose: true 
    # How many current workers will be running? Default it will be 1 worker
    count: 100 
    # bash command to be run.  Remember it's yaml, so you can use `|` or `>` if need be.
    # {{ .Name }} is the name of the task(aka the filename)
    # {{ task .Args "key" }} Inside the file you can have a key: value sets inside. 
    # {{ uniqId }} generates a uniq id on the fly for you
    # That is how you retrieve those values
    command:
        ssh user@someplace.com mysuperawesomecommand {{ .Name }} {{ task .Args "key"}}
    # ^^^ if you had 100 of these tasks, it would run 100 at the same time due to the `count` key
```

### Messages

Messages are just simply files inside a queue folder. The name of the file can be accessed via `{{ .Name }}` and inside the file you can have key/value pairs. Or not ... if you choose to use the file contents for something else, you can access it via `{{ .Contents }}` To access the keys in your `command` in the `.q` file, use: `{{ task .Args "yourkeyhere"}}`. Here is what a task file might look like. Lets say the task file was named `kcmerrill@gmail.com`. 

```yaml
first_name: kc
last_name: merrill
something: >
    elsewouldgohere
```

### ProTips

1. vbQ should be automated itself! With Jenkins/Build system based off of PR's!
1. Your build system should only allow for 1 build at a time and _NOT_ setup for multiple builds concurrently.
1. Be creative. Use the key=value pairs, or use the whole file contents for stdin for a custom application.
1. To rerun the messages again, simply copy the correct files back to the top level of the queue
1. Trigger sequential steps by creating task files upon completion. 
1. dotfiles/README.md/folders/.template.yml are all ignored and not considered tasks. Only files. Even zips. Be creative.
1. Speaking of readme files, place a README.md on the same level as your `.q` file and put a basic message template in place so those creating the PR's can just copy paste without leaving github. Or create a `file.template.yml`  as all `*.template.yml` files are ignored
