# projectile

Run your projects commands by calling simple action keyword.

## Usage
```
projectile [OPTIONS] <COMMAND> [<ARGS>]

COMMAND:
  get           List all config actions.
  do <actions>  Execute the actions listed.
  edit          Open the config with $EDITOR.
  add           Add a new action to the project.
  append        Append steps to an existing action.

OPTIONS:
  -p, --path  The project path, the current working directory by default.
  -h, --help  Show this help message
```

## Config file
The config is kept in a json file.
The file's default location is `$HOME/.config/projectile.json`.
This default location can be modified by setting the `PROJECTILE_CONFIG` environment variable.

Example of a projectile config:
```json
{
    "projects": [
        {
            "path": "/path/to/your/project",
            "actions": [
                {
                    "name": "setup",
                    "steps": [
                        "command 1",
                        "command 2"
                    ]
                },
                {
                    "name": "clean",
                    "steps": [
                        "command 1"
                    ]
                },
                {
                    "name": "build",
                    "steps": [
                        "command 1",
                        "command 2",
                        "command 3"
                    ]
                }
            ]
        },
        {
            "path": "/path/to/another/project",
            "actions": [
                {
                    "name": "test",
                    "steps": [
                        "command 1",
                        "command 2"
                    ]
                },
                {
                    "name": "build",
                    "steps": [
                        "command 1"
                    ]
                }
            ]
        },
    ]
}
```
The file lists, for each project, a serie of actions that can be performed.
Each action consist of an array of shell commands to execute.
