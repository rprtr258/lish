# Inker

Web API to execute and experiment with the [Ink programming language](https://github.com/thesephist/ink).

## Development

Inker is a standard Node.js server. You can start it with `yarn start`.

To work fully, you'll need an Ink interpreter. Inker ships with one at `/bin/ink` compiled for Linux, and running `yarn start` will assume you want to use that binary. If you need to develop on another platform or want to use a different version of Ink, you'll need to build your own binary, which you can learn how to do by going to the Ink repository linked above.

Once you have your ink binary, pass the path to your executable as the `INKPATH` environment variable when you start up Inker, and the eval service will use that instead.

If you have `ink` installed on your system already, running `yarn start:dev` will start the server with `INKPATH` set to `ink` on your `PATH`.
