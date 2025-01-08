<div align="center">
    <picture>
        <img style="margin-bottom:0;" width="130" src="./assets/gorepo.png" alt="logo">
    </picture>
    <h1 align="center" >GOREPO</h1>
</div>

<p align="center">
    A CLI to manage Go monorepos.
</p>

- Discord: [https://discord.gg/dRuqRU7R](https://discord.gg/dRuqRU7R)
- Contribute: [CONTRIBUTE.md](./CONTRIBUTE.md)

# Philosophy

The CLI should:
- be dumb to use
- allow running all commands from anywhere since having to cd is just annoying
- allow running CI/CD commands (test, lint, build, etc.) for all modules at once
- be transparent to the user regarding what it does behind the hood

If you want to know more about the direction the project is taking, see [BRAINSTORM.md](./BRAINSTORM.md).
If you want to share your use cases and affect the direction it is going, [open an issue](https://github.com/gorepo-cli/gorepo/issues) or join [discord](https://discord.gg/dRuqRU7R).

# Disclaimer
- This is not nearly a v1, it provides only basic features
- I code features as I go and as I need them
- Commit before running any command to see exactly what you are doing
- I only test Linux for now, macOS is probably ok, Windows is probably not

# Homebrew

## Pre-requisites
You must have go and git installed.

## Install via homebrew
```bash
brew tap gorepo-cli/gorepo
brew install gorepo
```

Test it is working with:
```bash
gorepo version
```

Note in some rare cases, if gorepo is not recognized, you may have to add the folder to the path

## Update via homebrew
```bash
brew upgrade gorepo
```

Note you may have to kill local cache:
```bash
brew update
brew uninstall gorepo
brew cleanup
brew install gorepo
```

## Build from source

To learn how build from source, visit [CONTRIBUTE.md](./CONTRIBUTE.md)

# Pre-requisites

- go: to use the CLI, you need go (used to run go commands)
- git: to build the project, you need git (used to inject the version at build time)

# Concepts

- A **monorepo** is a project with a `work.toml` file at the root. Monorepos can not be nested.
- A **module** is a folder containing a `module.toml` file. Currently you have to **create it manually**. It can be an empty file for now. Modules can technically be nested but you should probably avoid it for clarity.
- Modules can have a `scripts` section. They can be executed with `gorepo execute <script_name>` (see reference below).

# Reference

The reference contains information that is relevant to the actual commited version on master. Reference for future development and experimental features should be under [ROADMAP.md](./BRAINSTORM.md).

**Structure of a command:**

```
gorepo [global options] <optional_module> <command> [command options]
```

- <module_name>: Optional. Specifies the module to execute commands within. If omitted, commands will target the root.
- <command>: The CLI command to execute.
[options]: Command-specific flags and options.

## gorepo init

### Description

Initialize a new monorepo at the working directory.

This command creates two primary files:
- `work.toml` at the work directory
- `go.work` file if the strategy is set as 'workspace' and one does not exist yet. This runs `go work init` behind the hood


### Usage

```
gorepo init
```

### Examples

```
# The most basic way to start:
gorepo init

# You can also pass a name to name your monorepo
gorepo init some_name
```

## gorepo add

### Description

Add a new module to the monorepo.

This command creates a new folder with a `module.toml` and a `go.mod` file in it.
If the strategy used is a workspace, it will also add the module to the workspace.
Please note it will add the module at the directory provided from the root of the monorepo,
not from the current directory.

### Usage

```
gorepo add [module_name]
```

### Parameters

No parameters yet.

### Examples

```
# The most basic way to add a module
gorepo add my_module

# You can also pass a path to add the module at a specific location
gorepo add some_folder/my_module
```

## gorepo list

### Description

List all modules of the monorepo. Formally a module is a folder with a `module.toml` file in it.

### Usage

```
gorepo list
```

## gorepo exec

### Description

Executes a script at the root of the monorepo, or in one, many or all modules.

### Usage

```
gorepo [module_name] exec [--target] [--exclude] [--allow-missing] [script_name]
```

### Parameters

- `module_name`: the name of the module to execute the script in (put none to execute at the root)
- `script_name`: the name of the script to execute
- `--target` (optional): comma-separated names of modules to target, or all
- `--exclude` (optional): comma-separated names of modules to exclude
- `--allow-missing` (optional): allows the script to run even if some of the targets does not have the script

### Examples

```
# Executes 'my_command' script at the root
gorepo exec my_command

# Executes 'my_command' script in module named mod1
gorepo mod1 exec my_command

# Executes 'my_command' accross all modules
# Will fail if the script is missing in some modules
gorepo exec --target=all my_command

# Executes 'my_command' script in all modules that have it
gorepo exec --target=all --allow-missing my_command

# Executes 'my_command' script in modules 1 and 2
gorepo exec --target=mod1,mod2 my_command

# Executes 'my_command' script in all modules except in module X
gorepo exec --target=all --exclude=modX my_command
```

## gorepo fmt-ci

### Description

This command is breaking if the code in targeted modules is not formated.
This is primary meant to be used in ci pipelines, it does not modify the code or apply changes.

### Usage

```
gorepo [module_name] fmt-ci [--target] [--exclude]
```

### Parameters

- `module_name`: the name of the module to execute the script in (put none to execute at the root)
- `--target` (optional): comma-separated names of modules to target
- `--exclude` (optional): comma-separated names of modules to exclude

### Exemples

For the usage of the flags, refer to the reference of `gorepo execute`

## gorepo vet-ci

### Description

This command is breaking if `go vet` returns an error in one of the targeted modules.
This is primary meant to be used in ci pipelines, it does not modify the code or apply changes.

### Usage

```
gorepo [module_name] vet-ci [--target] [--exclude]
```

### Parameters

- `--target` (optional): comma-separated names of modules to target
- `--exclude` (optional): comma-separated names of modules to exclude

### Exemples

For the usage of the flags, refer to the reference of `gorepo execute`

## gorepo version

### Description

Print the version of the CLI

### Usage

```
gorepo version
```

## Contributing

Contributions are welcome, check out [CONTRIBUTE.md](./CONTRIBUTE.md)

## Releases

Check out [RELEASES.md](./RELEASES.md)

## License

This project is licensed under the MIT License.
