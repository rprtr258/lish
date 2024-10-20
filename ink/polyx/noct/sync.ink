` sync operations `

std := load('../vendor/std')
each := std.each

fs := load('fs')

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
