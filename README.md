# act - Activity Continuous Tracking

A tool for timekeeping and tracking. For now only for Redmine.

## Code status

[![Go Report Card](https://goreportcard.com/badge/github.com/gumieri/act)](https://goreportcard.com/report/github.com/gumieri/act) [![Build Status](https://travis-ci.org/gumieri/huexe.svg?branch=master)](https://travis-ci.org/gumieri/huexe)

## Configuration file

The configuration file need to be named `.act.yaml` and placed on the user's home directory.

On Windows it is the `%userprofile%`, on any Unix like it is the `~` path.

An example of configuration:

```yaml
redmine:
  url: 192.168.3.41
  access_key: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'

default:
  activity_id: 24

activities:
  dev: 24
  cr: 25

git:
  path: /usr/bin/git
  regex: '[0-9]*'

editor: vim
```

## How to track an activity

The commands to manage the tracking are `start`, `stop`, `status`, `push` and `rm`.
A very simple workflow example:

```bash
act start
```
```
2017/10/19 20:17:01 Activity 1234 started.
```
```bash
act status
```
```
        Issue   Started At      Stopped At      Spent           Comment
{0}     #1234   8:17PM          -               5m40.367443079s ""
```
```bash
act stop
```
```
2017/10/19 20:22:47 Activity 1234 stopped. Time elapsed 0.10 (5m45.347366369s)
```
```bash
act status
```
```
        Issue   Started At      Stopped At      Spent           Comment
{0}     #1234   8:17PM          8:22PM          5m45.347366369s ""
```
```bash
act push
```
```
Added 0.10 hour(s) to the Issue #1234.
```

## Commands

### `spent`
```bash
act spent 5.67 -i 12345 -a dev -d 2017-09-15 -m 'Making the world a better place for humans'
```
Parameters/Arguments:
1. time -- it can be informed as:
    * `1.5` -- (1 hour 30 minutes) As a fraction of hour
    * `1:45` -- (1 hour 45 minutes) As time with hour and minute

Options/Flags:
* -i --issue_id -- If not informed it will try retrieve it from the git branch name using the regex config
* -d --date -- It can be informed as:
    * `2017-09-22` -- Complete date
    * `09-22` -- Only the month and day. The year will be the current one
    * `22` -- Only the day. The year and month will be the current ones
    * `-1` -- Informing how many days back from the current date
    * And if not informed, it will use the current date
* -m --comment -- If not informed it will try to use the defined text editor
* -a --activity -- The name of an activity ID defined on the configuration file under activities
* --activity_id -- The activity ID. It Can be defined a default value to be used at the configuration file

### `log`
```bash
act log -i 12345
```

Options/Flags:
* -i --issue_id

### `link`
```bash
act link -i 12345
```

Options/Flags:
* -i --issue_id

### `start`
```bash
act start -i 12345
```

Options/Flags:
* -i --issue_id -- If not informed it will try retrieve it from the git branch name using the regex config

### `stop`
```bash
act stop
```

### `push`
```bash
act push
```

### `status`
```bash
act status
```

### `rm`
```bash
act rm 0
```

Parameters/Arguments:
1. activity index on `status` list (without the curly brackets).

### `note`
```bash
act note -i 12345 "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec varius eu purus in aliquam. Ut pellentesque magna purus, eu imperdiet justo convallis ac."
```

```bash
act note -i 12345 -t mr
```

Parameters/Arguments:
1. note -- The string note. If not informed it will try open the configured editor.

Options/Flags:
* -i --issue_id -- If not informed it will try retrieve it from the git branch name using the regex config
* -t --template -- The name of a template file to be loaded to be edited. These template files need to be stored in `templates` inside `.act` in the home directory (`~/.act/templates/text_file`)

## License

Act is released under the [MIT License](http://www.opensource.org/licenses/MIT).
