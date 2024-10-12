#someday_maybe #project

# hierachial dependencies / perfect dependency resolution

resolve dependencies perfectly by just getting every version of every needed dependency and reusing same dependencies

the only problem i see now: how to allow dependencies see only its dependencies and not other. in other words if

`a` depends on `x-0.1` and `y`

`b` depends on `x-0.2` and `z`

`a` should see only `x-0.1` and `b` should see only `x-0.2`

as i understand, all(?) modern dependency resolvers would solve this problem by just using `x-0.2` for both `a` and `b`

[FAQ | Documentation | Poetry - Python dependency management and packaging made easy](https://python-poetry.org/docs/faq/#why-is-the-dependency-resolution-process-slow)