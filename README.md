# projectile

Project builder

## Usage
```
projectile [OPTIONS] <COMMAND> [<ARGS>]

COMMAND:
  get           List all config actions.
  do <actions>  Execute the actions listed.
  do all        Execute all the actions in the config sequentially

OPTIONS:
  -p, --path  The project path, the current working directory by default.
  -h, --help  Show this help message
```

## Config file

Example of a project config:
```json
{
    "actions":[
        {
            "name": "setup",
            "steps": [
                "git submodule update",
                "git submodule init"
            ]
        },
        {
            "name": "clean",
            "steps": [
                "make clean"
            ]
        },
        {
            "name": "build",
            "steps": [
                "make -j6"
            ]
        }
    ]
}
```
The file lists a serie of actions that can be performed.
Each action consist of an array of shell commands to execute.
