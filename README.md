![vbQ](assets/vbQ.png "vbQ")

# vbQ

At it's core vbQ is a simple flat file queue system. Queues are just folders, tasks are just files inside of said folders. Workers are bash commands(this means you can program workers to be whatever you'd like them to be. Python. Go. PHP. Whatever.) that can be run 1 at a time or as many as you need and your box can handle. Once the tasks are completed, depending on whether they pass or fail they will get sent to the appropriate folders(configurable) and then merged back into your favorite *VCS. 

# Why?

If you have a folder with a bunch of bash scripts that need to be manually run or perhaps you have a wiki with instructions you have to follow then you'll understand why. Any tasks that can be automated should be automated, and a lot of the times you get to do it simply because you have the permissions to do so. Need to setup that new github repo? Need to add a new chef user? What about adding those public ssh keys to test/stage/prod boxes? Automate it!

# The vision

Instead of sending in a ticket to your team to create a new chef user, provision a new ldap user, or XYZ, that person instead sends in a PR to github. The PR is just a simple file with a filename, and whatever you want inside. By default, we'll parse out the yaml so you can have configurable arguments to pass into your application/script. You or your team gets that PR, merges it in, github then chats with a Jenkins box and vbQ will take that file or task, run the accompanying commands and move it to completed. Rinse and Repeat. Because it's just a shell task runner, you can write full applications or simple bash scripts. Send out slack messages, send email, encrypt and email keys, etc ... It's really up to you. 

The idea is to automate everything. 


### .q file

Queues are just folders. Inside each folder, there is a `.q` file. Here is a fully baked `.q` file with comments. 

```yaml
queue:
    # when a task gets completed succesfully(0 exit code), the folder to send it to 
    completed: ".completed"
    # when a task gets fails(exit code != 0) send to this folder. By default it stays put.
    failed: ".failed" 
workers:
    # How many current workers will be running? Default it will be 1 worker
    count: 100 
    # bash command to be run.  Remember it's yaml, so you can use `|` or `>` if need be.
    # {{ .Name }} is the name of the task(aka the filename)
    # {{ task .Args "key" }} Inside the file you can have a key: value sets inside. 
    # That is how you retrieve those values
    cmd:
        ssh user@someplace.com mysuperawesomecommand {{ .Name }} {{ task .Args "key"}}
    # ^^^ if you had 100 of these tasks, it would run 100 at the same time due to the `count` key
```

### Messages

Messages are just simply files inside a queue folder. The name of the file can be accessed via `{{ .Name }}` and inside the file you can have key/value pairs. Or not ... if you choose to use the file contents for something else, you can access it via `{{ .Contents }}` To access the keys in your `cmd` in the `.q` file, use: `{{ task .Args "yourkeyhere"}}`. Here is what a task file might look like. Lets say the task file was named `kcmerrill@gmail.com`. 

```yaml
first_name: kc
last_name: merrill
something: >
    elsewouldgohere
```

# ProTips

Just remember, this is a proof of concept. Having said that, here are some protips:

1. vbQ should be automated itself! With Jenkins/Build system based off of PR's!
1. Your build system should only allow for 1 build at a time and _NOT_ setup for multiple builds concurrently.
1. Be creative. Use the key=value pairs, or use the whole file contents for stdin for a custom application.
1. To rerun the messages again, simply copy the correct files back to the top level of the queue
1. dotfiles/README.md/folders are all ignored and not considered tasks. Only files. Even zips. Be creative.
1. Speaking of readme files, place a README.md on the same level as your `.q` file and put a basic message template in place so those creating the PR's can just copy paste without leaving github.
1. Take a peek at [vbq-tasks](https://github.com/kcmerrill/vbq-tasks) for a quick example of what it might look like.
1. vbQ WILL COMPLETELY RESET your git directory ... be careful. Don't be silly.


#TODO/Roadmap

1. Refactor VCS, meaning, allow 100% customizations. Update/Commit/Or not ... use it as a flat file Queue.
1. Alloow for retry/attempt logic
1. Run the jobs again depending if anything new was added to the queue. mid flight.


*Currently only git is supported, but maybe more are coming. It is a POC afterall.
