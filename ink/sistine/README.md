# Sistine 🏰
**Sistine** is a simple, flexible, productive static site generator written in [Ink](https://dotink.co/) and built on [Merlot](https://github.com/thesephist/merlot)’s Markdown engine. You can see a live demo of a Sistine site [on the Sistine docs website](https://sistine.vercel.app/).

![A screenshot of the Sistine docs site, built with Sistine](/static/img/sistine-screenshot.png)

## Documentation
Sistine's documentation lives on its own website at **[sistine.vercel.app](https://sistine.vercel.app)**. There, you'll find information on how to install and use Sistine, as well as a detailed reference for its templating language.

## Development
This repository technically contains two projects. First, the Ink source code for the Sistine static site generator; and second, the documentation site for Sistine, which is generated by Sistine itself from assets in this repo. Both parts of this repository uses a Makefile to manage common build commands.

### Sistine, the tool
Sistine's source code mostly lives in `./src`, with vendored dependencies copied into `./vendor`. Tests for Sistine utilities are in `./test`.

- `make check` or `make t` runs all tests in the repository.
- `make fmt`  or `make f` formats all Ink files (including tests) in the repository, if you have [`inkfmt`](https://github.com/thesephist/inkfmt) installed.

### Sistine, the website
The Sistine documentation website is a normal Sistine project, living in this repository.  The repository is set up with Vercel so that contents of `./public` auto-deploys on every commit to `main`. Other parts of the website, like the content pages and templates, are set up exactly as a normal Sistine project, in `./content` and `./tpl` respectively.

- `make` will run the in-repository copy of Sistine to build the documentation site to `./public`.

### Contributing & reporting issues
Given that it's currently quite slow and written in Ink, you probably shouldn't use it for anything important. But if you are interested, and want to ask questions about how it works or what's coming next, feel free to [reach out on Twitter](https://twitter.com/thesephist) or file a [GitHub issue](https://github.com/thesephist/sistine/issues).
