# Contributing to Subfnder

First off, thank you for considering contributing to Subfnder! It's people like you that make Subfnder such a great tool for the community.

## How Can I Contribute?

### Reporting Bugs
If you find a bug, please create an issue on GitHub. Include as much detail as possible:
* Your operating system and Go version.
* The exact command you ran.
* The expected output vs the actual output.

### Suggesting Enhancements
We are always looking for new ideas! If you have a suggestion:
* Create an issue on GitHub describing the enhancement.
* Explain why this enhancement would be useful to most users.

### Adding New Free APIs
Subfnder aims to integrate completely free APIs that don't require keys. To add a new source:
1. Fork the repository.
2. Create a new file in `sources/` implementing the `Source` interface (check `core/source.go`).
3. Add the new source to the `sourceList` array in `main.go`.
4. Run the tool to ensure your source works perfectly and doesn't break others.
5. Create a Pull Request (PR) with a description of the API you added.

## Pull Request Process
1. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
2. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
3. You may merge the Pull Request in once you have the sign-off of at least one other developer.

Happy hacking!
