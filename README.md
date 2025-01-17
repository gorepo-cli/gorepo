<div align="center">
  <picture>
    <img style="margin-bottom:0;" width="150" src="https://raw.githubusercontent.com/gorepo-cli/gorepo/refs/heads/master/assets/gorepo.png" alt="logo">
  </picture>
  <h1>GOREPO</h1>
  <p>
    <a href="https://github.com/gorepo-cli/gorepo/releases" style="text-decoration: none;">
      <img src="https://img.shields.io/github/release/gorepo-cli/gorepo.svg" alt="Latest Release">
    </a>
    <a href="https://github.com/yourusername/yourrepo/actions" style="text-decoration: none;">
      <img src="https://github.com/gorepo-cli/gorepo/actions/workflows/cicd.yml/badge.svg" alt="Build Status">
    </a>
    <a href="https://github.com/gorepo-cli/gorepo/commits" style="text-decoration: none;">
      <img src="https://img.shields.io/github/commit-activity/m/gorepo-cli/gorepo.svg" alt="Commits">
    </a>
  </p>

  <p>
    A CLI to manage Go monorepos
  </p>

  <p>
<a href="./REFERENCE.md">
        Reference
      </a> |
    <a href="./CONTRIBUTE.md">
        Contribute
      </a> |
      <a href="https://discord.gg/dRuqRU7R">
        Discord
      </a>
  </p>
</div>

[//]: # ([![Code Size]&#40;https://img.shields.io/github/languages/code-size/gorepo-cli/gorepo.svg&#41;]&#40;https://github.com/gorepo-cli/gorepo&#41;)
[//]: # ([![Stars]&#40;https://img.shields.io/github/stars/gorepo-cli/gorepo.svg?style=social&#41;]&#40;https://github.com/gorepo-cli/gorepo/stargazers&#41;)

## Highlights & philosophy

- Use module names as commands to use their context
- Use tasks names as commands to execute them
- Run tasks at root or across modules from anywhere
- Define queues of tasks if multiple steps are needed
- Define priorities when order matters
- Break CI with the flag `--ci`, like `gorepo fmt --ci`

If you want to know more about the future of the project, see [BRAINSTORM.md](./BRAINSTORM.md).
If you want to influence it, [open an issue](https://github.com/gorepo-cli/gorepo/issues) or join the [discord](https://discord.gg/dRuqRU7R).

## Disclaimer

- I code features as I need them, make requests [here](https://github.com/gorepo-cli/gorepo/issues)!
- This is not a battle-tested tool yet, and it only provides basic features for now
- Features and syntax may change all of a sudden (follow semantic versioning to know)
- I only test Linux for now, macOS is probably ok, Windows is probably not
- It's ok if you have non-go packages, it works fine

## Getting started

### Install via homebrew

```bash
# use the brew tap 
brew tap gorepo-cli/gorepo

# install gorepo
brew install gorepo

# test installation
gorepo version
```

You should update frequently with `brew upgrade gorepo`

### Build from sources

Follow instructions in [contribution](./CONTRIBUTE.md).

## Cheat sheet

Check the full reference [by clicking here](./REFERENCE.md) if needed.

| Command                                   | Description                                                                                                      | Compatible with module context |
|-------------------------------------------|------------------------------------------------------------------------------------------------------------------|--------------------------------|
| [gorepo init](./REFERENCE.md#gorepo-init) | Initialize a new monorepo. Creates `work.toml` and optionally `go.work`.<br/>exanmple: `gorepo init my_monorepo` | No                             |
| [gorepo add](./REFERENCE.md#gorepo-add)                 | Add a new module. Creates `module.toml` and `go.mod`, adds to workspace.                                         | No                             |
| [gorepo list](./REFERENCE.md#gorepo-list)               | List all modules in the monorepo.                                                                                | No                             |
| [gorepo exec](./REFERENCE.md#gorepo-exec)               | Execute a task at the root or in specific modules. Note the command 'exec' is optional. You can just drop it.    | Yes                            |
| [gorepo fmt](./REFERENCE.md#gorepo-fmt)                 | Run `go fmt` on targeted modules, supports `--ci` for CI environments.                                           | Yes                            |
| [gorepo vet](./REFERENCE.md#gorepo-vet)                 | Run `go vet` on targeted modules, supports `--ci` for CI environments.                                           | Yes                            |
| [gorepo check](./REFERENCE.md#gorepo-check)             | Analyze and provide information about the monorepo configuration.                                                | No                             |
| [gorepo version](./REFERENCE.md#gorepo-version)         | Print the CLI version.                                                                                           | No                             |
| [gorepo help](./REFERENCE.md#gorepo-help)               | Display help for GOREPO commands.                                                                                | No                             |

## Reference

[//]: # (<div>)
[//]: # (  <picture>)
[//]: # (    <img src="assets/banner.webp" alt="banner" />)
[//]: # (  </picture>)
[//]: # (</div>)

Check the full reference [by clicking here](./REFERENCE.md). Note you can also use the `help` command (even on module names).

## Deeper on concepts

### What does it mean that a module can be used as a command ?

It means you can use a module name directly to execute command in their context. For example, run `gorepo mod1 test` to run tests in the module mod1. Most other cli will provide a command for that but I like writing less.
Note this syntax is syntactic sugar, you can prefer `gorepo exec --target=mod1 test`.


### What does it mean that a task can be used as a command ?

In your `work.toml` and `module.toml`, you can define tasks (like 'scripts' in a package.json). To run a task, you can use the exec keyword, like `gorepo exec mytask`. You can also drop the command exec and simply write `gorepo mytask`.
Note this syntax is syntactic sugar, you can prefer `gorepo exec mytask`.

### Can they be combined ?

The two concepts described above can be combined to execute a task in a module. You can simply write `gorepo mod1 test` to run tests in the module mod1, instead of `gorepo exec --target=mod1 test`.

### How about conflicts between commands, module names and task names ?

Since module names and task names are used as commands, they should not conflict between them or with commands. If you have conflicting them, just know commands will win over the others, and modules will win over tasks. You can skip the syntactic sugar and use the full version of the commands.

### What is a monorepo ?

In the context of this CLI, a monorepo is a repository containing a `work.toml` file at the root (generated by `gorepo init`).
It is important to note that the CLI will generate and maintain a go workspace for you. Same thing applies if you have javascript modules, you should use a workspace to manage dependencies.

Structure of a work.toml:
```toml
name = 'my_super_project'
version = '0.1.0'
vendor = true

[scripts]
start = "echo 'start something here'"
```

### What is a module ?

A **module** is a folder containing a `module.toml` file. They have the following structure:
```toml
template = '@default'
type = 'executable'
language = 'go'
priority = 0

[scripts]
some_task = "echo 'do something here'"
```

- The module **name** is the name of the folder
- Modules can have the following **types**:
  - **executable**: the code can be built and executed
  - **library**: the code can be built
  - **script**: the code can be executed
  - **static**: the code is just meant to be imported from other modules
- The **language** field is used to know if it can use default Go commands.
- The **template** field is unused yet
- The **priority** is a number to order the modules when anything is executed. The higher goes first.

### How does priority work ?

For now priority is handled in a very straightforward way.
Each module can have a priority set in their `module.toml`. The modules with the higher priority goes first in all commands.
Type `gorepo list` to see all priorities.

I will probably provide a much more robust dependency graph and pipelines, but simple priorities with tasks queues already cover most cases.

### How to define queues of tasks ?

When you define a task at root or in a module, you can define it as a single task (a string) or a queue of tasks (an array of strings).

In this example, task1 is a single task, task2 is a queue of tasks:
```toml
[tasks]
task1 = "echo 'this is a simple task'"
task2 = [
  "echo 'step 1",
  "echo 'step 2'"
]
```

I will probably provide a much more robust dependency graph and pipelines, but simple priorities with tasks queues already cover most cases.

### Can I override go commands ?

The tool allow you to run go commands in all modules at once, like `gorepo fmt` (refer to the reference to know which are implemented).
You can override this behaviour using the following rules:
- If a task of the same name exist at the root, it will run instead (the command will not be run in each module)
- If a task of the same name exist in a go module, it will run instead for this module
- If a task of the same name exist in a non-go module, it will run (otherwise nothing will run for this module)

## Contributing

Contributions are welcome, check out [CONTRIBUTE.md](./CONTRIBUTE.md)

## License

This project is licensed under the MIT License.
