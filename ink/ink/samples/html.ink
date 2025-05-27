# html rendering library

{format: f, join} := import('str.ink')
{map} := import('functional.ink')

classes := (classes) => join(classes, ' ')
el := (tag, attrs, children) =>
  f('<{{ tag }}{{ props }}>{{ children }}</{{ tag }}>', {
    tag
    props: len(attrs) :: {
      0 -> ''
      _ -> ' ' + join(map(keys(attrs), (k) => (
        v := attrs.(k)
        f('{{ key }}="{{ value }}"', {key: k, value: v})
      )), ' ')
    }
    children: join(children, '')
  })

bindEl := (tag) => (props, children) => el(tag, props, children)

{
  classes
  el
  title: (s) => el('title', {}, s)
  meta:   bindEl('meta')
  link:   bindEl('link')
  h1:     bindEl('h1')
  h2:     bindEl('h2')
  h3:     bindEl('h3')
  p:      bindEl('p')
  em:     bindEl('em')
  strong: bindEl('strong')
  div:    bindEl('div')
  span:   bindEl('span')
  # simple div helper
  d: (children) => el('div', {}, children)
  # html wrapper helper
  html: (head, body) =>
    '<!doctype html>' +
    el('head', {}, head) +
    el('body', {}, body)
}