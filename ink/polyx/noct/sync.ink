` sync operations `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
each := std.each

fs := import('fs.ink')

actions := {
  up: 0
  down: 1
}

` generate a sync plan from a list of paths and stats `
diff := (local, remote) => (
  ` key: action, where 0 = push up, 1 = pull down `
  plan := {}
  each(keys(local), lpath => remote.(lpath) :: {
    () -> plan.(lpath) := actions.up
    _ -> local.(lpath).hash :: {
      remote.(lpath).hash -> ()
      _ -> plan.(lpath) := local.(lpath).mod > remote.(lpath).mod :: {
        true -> actions.up
        false -> actions.down
      }
    }
  })
  each(keys(remote), rpath => local.(rpath) :: {
    () -> plan.(rpath) := actions.down
  })
  plan
)
