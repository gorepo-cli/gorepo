
## Reference

<div>
  <picture>
    <img src="https://raw.githubusercontent.com/gorepo-cli/gorepo/refs/heads/master/assets/banner.webp" alt="banner" />
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
