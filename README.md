# act - Activity Continuous Tracking

## Configuration file

The configuration file need to be named `.act.yaml` and placed on the user's home directory.

On Windows it is the `%userprofile%`, on any Unix like it is the `~` path.

An example of configuration:

```
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
