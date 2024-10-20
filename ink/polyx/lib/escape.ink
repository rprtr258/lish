` escaping various formats `

str := import('../vendor/str')
replace := str.replace

html := s => (
  s := replace(s, '&', '&amp;')
  s := replace(s, '<', '&lt;')
)
