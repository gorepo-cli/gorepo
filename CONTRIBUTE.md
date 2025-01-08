# Contributing

Contributions to the project are more than welcome regardless the form they take.

You can:
- Install and use the cli
- Share how you wish to use such a tool
- Share your use cases/feature requests
- Share bugs and pain points
- Consult [BRAINSTORM.md](BRAINSTORM.md)
- Join the [discord](https://discord.gg/dRuqRU7R)
- Open issues on GitHub
- Open a PR
- Add your wishlist in the file [BRAINSTORM.md](BRAINSTORM.md)

Note this file will grow if needed.

Juniors are also welcome to push code (see how we need to add tests, to add new commands, how error messages are hardcoded etc).

# Build and run from source

- Clone/download the repository
- Run `make build` to create bin/gorepo
- Add the bin folder to your PATH
- As a result, you can now run `gorepo_` from anywhere
- Change code, build, test from anywhere, repeat

Example on Linux to add the bin folder to your PATH:
```
vim ~/.bashrc

# add this:
export PATH="$PATH:/home/my_name/Repositories/gorepo-cli/bin"

# refresh the terminal
source ~/.bashrc
```

Test it is working by typing:
```bash
gorepodev version
```

# How it is built

For now the project is in a single file. I think we should break it only if we start having multiple contributions and we have templates and such significant features. Note that even if the project is in a single file, all side effects are injected to make every function testable.

# Test kit

Note there is a testkit that allows you to write unit tests for your feature.

I will add more on that here when needed.

# Pipeline and homebrew

todo
