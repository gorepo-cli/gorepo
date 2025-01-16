<div align="center">
  <picture>
    <img style="margin-bottom:0;" width="180" src="./assets/gorepo.png" alt="logo">
  </picture>
  <h1 align="center">GOREPO</h1>
</div>
<p align="center">
  A CLI to manage Go monorepos
</p>

<div align="center">
  <a href="https://github.com/gorepo-cli/gorepo/releases">
    <img src="https://img.shields.io/github/release/gorepo-cli/gorepo.svg" alt="Latest Release">
  </a>
  <a href="https://github.com/yourusername/yourrepo/actions">
    <img src="https://github.com/gorepo-cli/gorepo/actions/workflows/cicd.yml/badge.svg" alt="Build Status">
  </a>
  <a href="https://github.com/gorepo-cli/gorepo/commits">
    <img src="https://img.shields.io/github/commit-activity/m/gorepo-cli/gorepo.svg" alt="Commits">
  </a>
  <a href="https://github.com/gorepo-cli/gorepo/blob/main/CONTRIBUTE.md">
    <img src="https://img.shields.io/badge/Contribute-Contributions%20Welcome-brightgreen.svg" alt="Contribute">
  </a>
  <a href="https://discord.gg/dRuqRU7R">
    <img src="https://img.shields.io/discord/1313214880689688626.svg?label=Discord&logo=discord&color=7289DA" alt="Discord">
  </a>
</div>


[//]: # ([![Code Size]&#40;https://img.shields.io/github/languages/code-size/gorepo-cli/gorepo.svg&#41;]&#40;https://github.com/gorepo-cli/gorepo&#41;)

[//]: # ([![Stars]&#40;https://img.shields.io/github/stars/gorepo-cli/gorepo.svg?style=social&#41;]&#40;https://github.com/gorepo-cli/gorepo/stargazers&#41;)

## Highlights & philosophy

- Use module names as commands
- Use tasks names as commands
- Run tasks at root or across modules
- Define queues of tasks if multiple steps are needed
- Define priorities when order matters
- Break during the CI with the flag `--ci`, ex: `gorepo fmt --ci`
- Run commands from anywhere

If you want to know more about the future of the project, see [BRAINSTORM.md](./BRAINSTORM.md).
If you want to influence it, [open an issue](https://github.com/gorepo-cli/gorepo/issues) or join the [discord](https://discord.gg/dRuqRU7R).

## Disclaimer

- I code features as I need them, make requests [here](https://github.com/gorepo-cli/gorepo/issues)
- This is not a battle-tested tool yet, also it only provides basic features for now
- Features and syntax may still change all of a sudden
- I only test Linux for now, macOS is probably ok, Windows is probably not
- It's ok if you have non-go packages, works fine

## Cheat sheet

This is a cheat sheet of the commands. Refer to the Reference bellow for more details.

| Command                           | Description                                                                                                      | Compatible with module context |
|-----------------------------------|------------------------------------------------------------------------------------------------------------------|--------------------------------|
| [gorepo init](#gorepo-init)       | Initialize a new monorepo. Creates `work.toml` and optionally `go.work`.<br/>exanmple: `gorepo init my_monorepo` | No                             |
| [gorepo add](#gorepo-add)         | Add a new module. Creates `module.toml` and `go.mod`, adds to workspace.                                         | No                             |
| [gorepo list](#gorepo-list)       | List all modules in the monorepo.                                                                                | No                             |
| [gorepo exec](#gorepo-exec)       | Execute a task at the root or in specific modules. Note the command 'exec' is optional. You can just drop it.    | Yes                            |
| [gorepo fmt](#gorepo-fmt)         | Run `go fmt` on targeted modules, supports `--ci` for CI environments.                                           | Yes                            |
| [gorepo vet](#gorepo-vet)         | Run `go vet` on targeted modules, supports `--ci` for CI environments.                                           | Yes                            |
| [gorepo check](#gorepo-check)     | Analyze and provide information about the monorepo configuration.                                                | No                             |
| [gorepo version](#gorepo-version) | Print the CLI version.                                                                                           | No                             |
| [gorepo help](#gorepo-help)       | Display help for GOREPO commands.                                                                                | No                             |

## Installation via Homebrew (recommended)

### Install via homebrew

```bash
brew tap gorepo-cli/gorepo
brew install gorepo
```

Test it is working with:

```bash
gorepo version
```

Note in some rare cases, if gorepo is not recognized, you may have to add the folder to the path

### Update via homebrew

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

## Pre-requisites

You need go installed.

## Concepts

- A **monorepo** is a project with a `work.toml` file at the root. Monorepos can not be nested.
- A **module** is a folder containing a `module.toml` file. Modules can technically be nested, but you should probably
  avoid it.
- A module can be of the following types:
    - **executable**: the code can be built and executed
    - **library**: the code can be built
    - **script**: the code can be executed
    - **static**: the code is not meant to be built, just imported from other modules
- Modules can have a `scripts` section. They can be executed with `gorepo exec <script_name>` (see reference below).
- Some commands can be run at the root, like `gorepo exec start`, but also in the concept of a module, like
  `gorepo module1 exec start`. Note the module is here a command.
- THe monorepo structure relies on go workspaces since the industry finally acknowledges it is ok to commit go.work
  files for that

## Reference

<div>
  <picture>
    <img src="assets/banner.webp" alt="banner" />
  </picture>
</div>

**Structure of a command:**

```
gorepo [--global_options] [module] <command> [--command-options] [args]
```

### gorepo init

```
NAME:
   gorepo init - Initialize a new monorepo

USAGE:
   gorepo [--global_options] init [--command_options] [monorepo_name]

DESCRIPTION:
   Initialize a new monorepo at the working directory.

   This command creates two primary files:
   - 'work.toml' at the work directory
   - 'go.work' file if the strategy is set as 'workspace' and one does not exist yet. This runs 'go work init' behind the hood

OPTIONS:
   --help, -h  show help
```

#### Examples

```
# The most basic way to start:
gorepo init

# You can also pass a name to name your monorepo
gorepo init some_name
```

### gorepo add

```
NAME:
   gorepo add - Add a module

USAGE:
   gorepo [--global_options] add [--command_options] <module_name>

DESCRIPTION:
   Add a new module to the monorepo.

   This command creates a new folder with 2 file, 'module.toml' and 'go.mod'. It also adds the module to the go workspace. You can pass a path ending with the module name.

OPTIONS:
   --help, -h  show help
```

#### Examples

```
# The most basic way to add a module
gorepo add my_module

# You can also pass a path to add the module at a specific location
gorepo add some_folder/my_module
```

### gorepo list

```
NAME:
   gorepo list - List modules

USAGE:
   gorepo [--global_options] list [--command_options]

DESCRIPTION:
   List all modules of the monorepo. Formally a module is a folder with a module.toml file in it, regardless of the language it uses.

OPTIONS:
   --help, -h  show help
```

### gorepo exec

```
NAME:
   gorepo exec - Execute a script

USAGE:
   gorepo [global_options] [module_name] exec [command_options] <script_name>

DESCRIPTION:
   Compatible with module syntax.

   Execute a script at the root of the monorepo, or in one, many or all modules. Scripts are declared in the files 'work.toml' and 'module.toml'.

OPTIONS:
   --target value   Target specific modules or root (comma separated) (default: "root")
   --exclude value  Exclude specific modules (comma separated)
   --allow-missing  Allow executing the scripts, even if some module don't have it (default: false)
   --help, -h       show help
```

#### Examples

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

### gorepo fmt

```
NAME:
   gorepo fmt - Run go fmt, break with --ci (module syntax compatible)

USAGE:
   gorepo [global_options] [module_name] fmt [command_options]

DESCRIPTION:
   Compatible with module syntax.

   This command runs fmt in all targeted modules.
   It breaks without formating the files if you pass --ci.

OPTIONS:
   --target value   Target specific modules or root (comma separated) (default: "root")
   --exclude value  Exclude specific modules (comma separated)
   --ci             Enable mode CI (default: false)
   --help, -h       show help
```

#### Examples

```
# Format all modules
gorepo fmt

# Breaks if modules are not formated
gorepo fmt --ci
```

### gorepo vet

```
NAME:
   gorepo vet - Run go vet, breaks if needed (module syntax compatible)

USAGE:
   gorepo [global_options] [module_name] vet [command_options] <script_name>

DESCRIPTION:
   Compatible with module syntax.

   This command runs vet in all targeted modules and breaks if vet breaks.

OPTIONS:
   --target value   Target specific modules or root (comma separated) (default: "root")
   --exclude value  Exclude specific modules (comma separated)
   --help, -h       show help
```

### gorepo check

```
NAME:
   gorepo check - Check the configuration

USAGE:
   gorepo [global_options] check [command_options]

DESCRIPTION:
   Gives information about the configuration.
   In the future it will also analyse the configuration

OPTIONS:
   --help, -h  show help
```

### gorepo version

```
NAME:
   gorepo version - Print version

USAGE:
   gorepo [global_options] version [command_options]

OPTIONS:
   --help, -h  show help
```

### gorepo help

```
gorepo help
gorepo command help
gorepo module help
```

### Contributing

Contributions are welcome, check out [CONTRIBUTE.md](./CONTRIBUTE.md)

### License

This project is licensed under the MIT License.
