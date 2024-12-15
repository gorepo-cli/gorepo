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
If you want to share your use cases and affect the direction it is going, **open an issue** or join [discord](https://discord.gg/dRuqRU7R).

# Disclaimer
- This is not nearly a v1, it provides only basic features
- I code features as I go and as I need them
- Commit before running any command to see exactly what you are doing
- I only test Linux for now, macOS is probably ok, Windows is probably not

# Structure of a command

```
gorepo [global options] <command> [command options]
```

# Pre-requisites

- go: to use the CLI, you need go (used to run go commands)
- git: to build the project, you need git (used to inject the version at build time)

# Getting started

Gorepo is not yet available on any package manager. You will need to build it yourself:
- Clone/download the repository
- Run `make build` to create bin/gorepo
- Add the bin folder to your PATH
- As a result, you can now run `gorepo` from anywhere
- Change code, build, test from anywhere, repeat

Example on Linux:
```
vim ~/.bashrc

# add this:
export PATH="$PATH:/home/my_name/Repositories/gorepo-cli/bin"

# refresh the terminal
source ~/.bashrc
```

# Concepts

- A **monorepo** is a project with a `work.toml` file at the root. Monorepos can not be nested.
- A **module** is a folder containing a `module.toml` file. Currently you have to **create it manually**. It can be an empty file for now. Modules can technically be nested but you should probably avoid it for clarity.
- Modules can have a `scripts` section. They can be executed with `gorepo execute <script_name>` (see reference below).

# Reference

The reference contains information that is relevant to the actual commited version on master. Reference for future development and experimental features should be under [ROADMAP.md](./BRAINSTORM.md).

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

## gorepo list

### Description

List all modules of the monorepo. Formally a module is a folder with a `module.toml` file in it.

### Usage

```
gorepo list
```

## gorepo execute

### Description

Execute a script across all targeted modules.

By default, it runs all the scripts (bash scripts) defined in `module.toml` files that are targeted.
It will not run if the script is missing in one of the targeted module is missing, unless you pass the flag `--allow-missing`.

### Usage

```
gorepo execute [--target] [--exclude] [--allow-missing] [script_name]
```

### Parameters

- `script_name`: the name of the script to execute
- `--target` (optional): comma-separated names of modules to target
- `--exclude` (optional): comma-separated names of modules to exclude
- `--allow-missing` (optional): allows the script to run even if some of the targets does not have the script

### Examples

```
# Will execute 'my_command' script in all modules.
# It will not run if it is missing in one or more modules.
gorepo execute my_command

# Will execute 'my_command' script in all modules that have it
gorepo execute --allow-missing my_command

# Will execute 'my_command' script in modules 1 and 2
gorepo execute --target=mod1,mod2 my_command

# Will execute 'my_command' script in work.toml at the root
gorepo execute --target=root my_command

# Will execute 'my_command' script in all modules except in module X
gorepo execute --exclude=modX my_command
```

## gorepo fmt-ci

### Description

This command is breaking if the code in targeted modules is not formated.
This is primary meant to be used in ci pipelines, it does not modify the code or apply changes.

### Usage

```
gorepo fmt-ci [--target] [--exclude]
```

### Parameters

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
gorepo vet-ci [--target] [--exclude]
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
