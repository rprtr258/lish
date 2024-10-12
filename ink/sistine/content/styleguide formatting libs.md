#someday_maybe

add to styleguide if actually so
- use [`Wissance/stringFormatter`](https://github.com/Wissance/stringFormatter) over `{html,text}/template` for unstructured texts
- for structured (html, markdow) consider using function constructors like `Render(w io.Writer) error`
    - maybe reuse libs for frontend(html)/markdown(mk)