# Brainstorm

This file gives an idea on what's next

## Bug fixes / code quality needed
- review all the logging
- extract errors as constants instead of hardcoding text everywhere
- improve help by adding examples
- improve the methods of Terminal (regarding logging, returning errors etc), needs some investigation for desired behaviour
- add feedback when commands executes successfully
- the testkit is using ToEffects to translate between MockEffects and Effects. We should probably be able to write that better using polymorphism

## Next features
- improve command check to perform a health check (provide flags --fix and --ci to break)
- flag --parallel
- add timeouts
- add version check call to a server (and a flag --no-check)
- add new go commands, test, build (tidy, get, run, some may have --ci)
- add new command `gorepo tree` to display the tree of the filesystem
- add validation on names and stuff

## Longer term

- make the cli generic (not go oriented), make it support various workspaces (go, npm, yarn, pnpm...)
- add modules using @templates, community templates, and support templates-tasks
- support third party plugins
- support caching and incremental builds locally - investigate for server based
- support a more complex dependencies and priorities (and add a command gorepo graph to show it)
- investigate if we should have a nodemon like feature or integrate one as a plugin or not
- generate gitignores
- see if we should somehow support docker in a way or another
- new command remove to remove a module
- `gorepo update` to update the CLI
- `gorepo upgrade` to upgrade the packages to the latest version
- natural language commands
