# html rendering library

log := s => out(s + '\n')
std := load('str')
f := std.format
join := std.join
map := load('functional').map

classes := classes => join(classes, ' ')
el := (tag, attrs, children) =>
  f('<{{ tag }}{{ props }}>{{ children }}</{{ tag }}>', {
    tag: tag
    props: len(attrs) :: {
      0 -> ''
      _ -> ' ' + join(map(keys(attrs), k => (
        v := attrs.(k)
        f('{{ key }}="{{ value }}"', {key: k, value: v})
      )), ' ')
    }
    children: join(children, '')
  })

bindEl := tag => (props, children) => el(tag, props, children)

title := s => el('title', {}, s)
meta := bindEl('meta')
link := bindEl('link')

h1 := bindEl('h1')
h2 := bindEl('h2')
h3 := bindEl('h3')

p := bindEl('p')
em := bindEl('em')
strong := bindEl('strong')

div := bindEl('div')
span := bindEl('span')

# simple div helper
d := children => div({}, children)

# html wrapper helper
html := (head, body) => '<!doctype html>' + el('head', {}, head) + el('body', {}, body)

# example usage
log(
  html(
    title('Test page')
    d([
      h1({class: classes(['title']), itemprop: 'title'}, 'Hello, World!')
      p({class: classes(['body'])}, 'this is a body paragraph')
    ])
  )
)

# export
{
  el: el
  title: title
  meta: meta
  link: link
  h1: h1
  h2: h2
  h3: h3
  p: p
  em: em
  strong: strong
  div: div
  span: span
  d: d
  html: html
}