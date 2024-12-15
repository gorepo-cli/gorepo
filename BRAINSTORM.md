
# Brainstorm

- allow building modules based on templates and community templates
- the creator of a template should be allowed to define template-scripts
- the bash and go functions should return the console output and/or errors
- add tests
- write a custom 'help' command with some ascii art
- The CLI could also handle incremental builds, given the user configures a storage
- generate a gitignore file for go repos
- see how we could handle docker
- see how we could handle pipelines
- implement some nodemon feature (`gorepo wath` or `gorepo run --watch`)

## New Commands

- add:    to add a module `gorepo add new_mod`
- health: to check the health of the modules (or check), with --fix
- remove: to remove a module
- fmt
- vet
- test
- get
- build   (check how to set priority)
- run     (check how to know the path + priority)
- tidy
- `gorepo check` flag `--fix` (or health)
- `gorepo tree` to display the tree of dependencies of the monorepo
- `gorepo update` to update the CLI
- `gorepo upgrade` to upgrade the packages to the latest version
- `gorepo start` (call what was built) option `--watch` (runs dev, if docker), option `--no-docker` (runs dev, without docker)

## New flags
- [executionFlags] parallel: to run the commands in parallel
- [global]         dry-run:  to show what would be done 

```
acceptable names has:
ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.@!#$%^&()[]{}'+,;=~
```

//{
//	Name:   "add",
//	Usage:  "Add a new module to the monorepo",
//	Action: commands.Add,
//	Flags: []cli.Flag{
//		&cli.BoolFlag{
//			Name:  "verbose",
//			Usage: "Enable verbose output",
//		},
//		&cli.StringFlag{
//			Name:  "template",
//			Usage: "Choose a template (not implemented)",
//		},
//	},
//},
//{}, // sanitize / lint / health / check


Add Context Support for Cancellation

Issue: Long-running operations cannot be cancelled by the user.

Recommendation: Pass context.Context to functions to handle cancellation and timeouts.

go
Copy code
func (cmd *Commands) Run(c *cli.Context) error {
ctx := c.Context
// Pass ctx to functions and check for cancellation
}

Provide Execution Summaries

Issue: Users don't receive a summary of the executed commands.

Recommendation: Collect and display a summary at the end of script execution.

go
Copy code
var failedModules []string
// ... (during execution)
if err != nil {
failedModules = append(failedModules, module.Name)
}
// After execution
if len(failedModules) > 0 {
cmd.SystemUtils.Logger.Warning("Scripts failed in modules: " + strings.Join(failedModules, ", "))
} else {
cmd.SystemUtils.Logger.Success("All scripts ran successfully.")
}