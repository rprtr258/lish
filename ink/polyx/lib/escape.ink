` escaping various formats `

str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
replace := str.replace

{
  html: s => (
    s := replace(s, '&', '&amp;')
    s := replace(s, '<', '&lt;')
  )
}
