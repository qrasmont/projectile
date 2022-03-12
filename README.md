# projectile

Run your projects commands by calling simple action keyword.

## Usage
```
projectile [OPTIONS] <COMMAND> [<ARGS>]

COMMAND:
  add <name> [steps]     Add a new action to the project.
  append  <name> [steps] Append steps to an existing action.
  do [actions]           Execute the actions listed.
  edit                   Open the config with $EDITOR.
  get                    List all config actions.
  rm [actions]           Remove the actions listed.

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
    "Projects": [
        {
            "Path": "/path/to/your/project",
            "Actions": [
                {
                    "Name": "setup",
                    "Steps": [
                        "command 1",
                        "command 2"
                    ],
                    "SubDir": ""
                },
                {
                    "Name": "clean",
                    "Steps": [
                        "command 1"
                    ],
                    "SubDir": ""
                },
                {
                    "Name": "build",
                    "Steps": [
                        "command 1",
                        "command 2",
                        "command 3"
                    ],
                    "SubDir": ""
                }
            ]
        },
        {
            "Path": "/path/to/another/project",
            "Actions": [
                {
                    "Name": "test",
                    "Steps": [
                        "command 1",
                        "command 2"
                    ],
                    "SubDir": ""
                },
                {
                    "Name": "build",
                    "Steps": [
                        "command 1"
                    ],
                    "SubDir": ""
                }
            ]
        },
    ]
}
```
The file lists, for each project, a serie of actions that can be performed.
Each action consist of an array of shell commands to execute.
