# act - Activity Continuous Tracking

## Usage

### `spent`
```bash
act spent 5.67 -i 12345 --activity_id=1 -d 2017-09-15 -m 'Making the world a better place for humans'
```
Parameters/Arguments:
1. time -- ex: 1.5 (for 1 hour 30 minutes)

Options/Flags:
* -i --issue_id -- If not informed it will try retrieve it from the git branch name using the regex config
* --activity_id -- Can be defined a default value at the configuration file
* -d --date -- The default value is the current date
* -m --comment -- If not informed it will try to use the defined text editor

### `log`
```bash
act log -i 12345
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

git:
  path: /usr/bin/git
  regex: '[0-9]*'

editor: vim
```
