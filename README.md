# projectile

Project builder

## Usage
```
projectile [-p path] action
```
- Without -p, the path is set to the cwd.
- Reserved action "all" will run all action squencially

## Config file

Example of a project config:
```json
{
  "setup": [
    "git submodule init"
    "git submodule update",
  ],
  "clean": [
    "make clean"
  ],
  "build": [
    "make -j6"
  ]
}
```
The file lists a serie of actions that can be performed.
Each action consist of an array of shell commands to execute.
