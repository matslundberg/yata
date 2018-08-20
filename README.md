# Yet Another Todo App

> Super duper notepad thingy.

**NOTE!!** THIS IS PRE-ALPHA QUALITY. Please do not use yet!

## What is it?

A way to interact with notes in text-files stored somewhere. You can:
 - View individual notes *(coming soon)*
 - Search all notes *(coming soon)*
 - Archive a note *(coming soon)*
 - Style notes using markdown *(coming soon)*
 - Embed tasks directly in notes using `[ ] An important todo`
 - List all tasks in notes
 - Tag notes using `+mytag`
 - Do @-mention of projects, persons
 - Add parseable dates using `//yyyy-mm-dd`
 - Embed file or image using `!image.jpg` *(coming later)*

Yata has these interfaces:
 - Command Line Interface: `yata`
  -- e.g `yata list tasks due today`
 - Single Page Application *(coming soon)*

Yata can read notes from the following places:
 - Filesystem
 - Dropbox-folder *(coming later)*
 - Git repository, e.g. Github, Gitlab *(coming later)*

## How to use
### Command Line Interface -- `yata`

```
yata help
yata list tasks
yata list tasks due today
yata list tasks status is completed
yata list tasks status is ongoing
yata list tasks status is rejected
yata list tasks status is open
yata list tasks due tomorrow and status is pending
yata list tasks due next 7 days
yata list tasks due last 7 days
yata list tasks @project
yata list tasks +tag
yata list tasks +tag @project due 7 days open
yata list tags
yata list projects
yata complete task 1afe2bb6
yata complete tasks in file filename
yata complete tasks in 
yata list tasks from jira due in 7 days
yata close these
```

## Build Setup

``` bash
```

## Prereqs

```
go get -u github.com/logrusorgru/aurora
```
