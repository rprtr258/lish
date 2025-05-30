# a ray is a vector out from a position in space

vec3 := import('vec3')
vneg := vec3.neg
vadd := vec3.add
vsub := vec3.sub
vmul := vec3.multiply

create := (pos, dir) => {
  pos: pos
  dir: dir
}

Zero := create(vec3.Zero, vec3.Zero)

# march ray to time t
at := (ray, t) => vadd(ray.pos, vmul(ray.dir, t))

{
  create: create
  Zero: Zero
  at: at
}
