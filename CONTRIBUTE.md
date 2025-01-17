# Contributing

Contributions to the project are more than welcome regardless the form they take. Juniors are also more than welcome 🎉

I would recommend you start installing and using the CLI to get a sense of what is happening there.

## Ways to contribute

- Join [discord](https://discord.gg/dRuqRU7R)
- Open an issue to discuss a new feature
- Share ideas about your dream tool
- Share bugs and pain points
- Consult [BRAINSTORM.md](BRAINSTORM.md) to see what you contribute with
- Fork and open a PR

## Build from source and Run

- Clone the repository
- Run `make build`, it will create bin/gorepo_
- Add the bin folder to your PATH
```
# How to add the folder to path
vim ~/.bashrc

# Add this:
export PATH="$PATH:/home/my_name/Repositories/gorepo-cli/bin"

# Refresh the terminal:
source ~/.bashrc
```
- Now, you can now run `gorepo_` from anywhere (I added the underscore so you can also install the production one in parallel)
- Change code, build, test locally, repeat

# Project's structure

More to come on that

# Testkit

Note there is a testkit that allows you to write unit tests for your feature (it allows TDD).

More to come on that

# Release

More to come on that
