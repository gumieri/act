# act - Activity Continuous Tracking

## Usage

### `spent`
```bash
act spent 5.67 -i 12345 --activity_id=1 -d '2017-09-15' -m 'Making the world a better place for humans'
```

* -i --issue_id
* --activity_id
* -d --date
* -m --comment

### `log`
```bash
act log -i 12345
```

* -i --issue_id

### `start`
```bash
act start -i 12345
```

* -i --issue_id

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
act note -i 12345 "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec varius eu purus in aliquam. Ut pellentesque magna purus, eu imperdiet justo convallis ac. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus."
```

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
