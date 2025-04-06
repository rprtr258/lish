print := (t, v) => out(t+': '+string(v)+'\n')
print('print 1', 1)

{sqrt} := import('math.ink')
print('sqrt 36', sqrt(36))

pyth := (x, y) => sqrt(x*x + y*y)
print('pyth(3, 4)', pyth(3, 4))

kmul := (x, y, k) => k(x * y)
kplus := (x, y, k) => k(x + y)
ksqrt := (x, k) => k(sqrt(x))
kpyth := (x, y, k) => kmul(x, x, x2 => kmul(y, y, y2 => kplus(x2, y2, x2py2 => ksqrt(x2py2, k))))
kpyth(3, 4, n => print('kpyth 3, 4', n))

mreturn := x => k => k(x)
mbind := (x, mf) => x(mf)
cont := (x, y, mf) => x(xval => y(yval => mf(xval, yval)))
mpyth := (x, y) => (
  x2 := mpow2(x)
  y2 := mpow2(y)
  x2py2 := cont(x2, y2, madd)
  mbind(x2py2, msqrt)
)
mpow2 := x => mreturn(x*x)
madd := (x, y) => mreturn(x+y)
msqrt := x => mreturn(sqrt(x))
mpyth(3, 4)(n => print('mpyth 3, 4', n))

loaded := n => (_, f, _) => f(n)
state := loaded(42)
print('state loaded', state(_ => 'loading', _ => 'loaded', _ => 'error'))
print('state loaded param', state(_ => 'loading', n => 'loaded '+string(n), _=>'error'))

keq := (x, y, k) => k(x == y)
kminus := (x, y, k) => k(x - y)

kfactorial := (n, k) =>
  keq(n, 0,
    b => true :: {
      b -> k(1)
      _ -> kminus(n, 1,
      n1 => kfactorial(n1,
        fn1 => kmul(n, fn1, k)))
    })

kfactorial(5, n => print('kfactorial 5', n))

keq2 := (x, y, ktrue, kfalse) => (x :: {
  y -> ktrue
  _ -> kfalse
})()

kfactorial2 := (n, k) =>
  keq2(n, 0,
    () => k(1),
    () => kminus(n, 1,
      n1 => kfactorial2(n1,
        fn1 => kmul(n, fn1, k))))

kfactorial2(5, n => print('kfactorial2 5', n))

kdiv := (x, y, k) => k(x / y)

kdivSafe := (x, y, kok, kerr) =>
  keq2(y, 0,
    () => kerr(string(x)+'/'+string(y)+': div by zero!'),
    () => kdiv(x, y, kok))

kdivSafe(5, 2, ok => print('ok', ok), err => print('err', err))
kdivSafe(5, 0, ok => print('ok', ok), err => print('err', err))

mdiv := (x, y) => mreturn(x / y)
mdiv(5, 2)(ok => print('mdiv', ok))

mdivSafe := (x, y, kerr) =>
  k => keq2(y, 0,
    () => kerr(string(x)+'/'+string(y)+': div by zero!'),
    () => kdiv(x, y, k))

mdivSafe(5, 2, err => print('err', err))(ok => print('ok', ok))
mdivSafe(5, 0, err => print('err', err))(ok => print('ok', ok))

mtok := f => (x, k) => f(x)(k)
ktom := f => x => k => f(x, k)