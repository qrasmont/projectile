# projectile

Project builder

## Usage
```
projectile [-p path] [-a=true] action
```
- Without -p, the path is set to the cwd.
- -a: perform all the config actions sequentially

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
