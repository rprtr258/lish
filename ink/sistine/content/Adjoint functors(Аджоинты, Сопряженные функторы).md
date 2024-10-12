#article

Изначально я хотел включить эту тему в статью про функторы, но тут рассказать можно **очень** много и придется много разбираться, поэтому я выделил отдельную статью.

# **Prerequisits**

Важно знать понятие функтора и естественного преобразования, поэтому советую прочитать вторую и третью(достаточно начала, где описываются естественные преобразования) статьи по теории категорий, чтобы хотя бы иметь шанс разобраться, что здесь происходит.

# **???**

Для начала начнем с примера, который я называю каррирование и который даже иногда используется в разных разделах математики в том или ином виде. Рассмотрим два функтора: **F**(A) = AxB и **G**(C) = C^B. Существует биекция между следующими множествами:

\sigma: Hom(F(A), C) ~= Hom(A, G(C))или, что то же самое:\sigma: Hom(AxB, C) ~= Hom(A, C^B)

Как работает эта биекция: она принимает функцию AxB→C, о которой можно думать, как о функции двух переменных: одна типа A, вторая типа B, и возвращает функцию, которая принимает один аргумент типа A и возвращает функцию, принимающую переменную типа B. Звучит сложно, потому что это функция, принимающая функцию и возвращающая функцию, но все на деле элементарно. Разберемся, как она работает на каком-нибудь примере:

Пусть A = B = Int — тип целых чисел и мы хотим узнать, чему равно \sigma(+), где + — функция суммы, а именно +(x, y) = x + y. \sigma(+) даст эту функцию, но в виде двух функций от одной переменной: \sigma(+)(x) = \y → x + y, где \y → x + y это результат \sigma(+)(x). Результат как бы запоминает первую переменную и, принимая вторую, вычисляет исходную функцию. В общем случае если есть функция f:AxB→C, то \sigma(f)(x) = \y → f(x, y).

Не будем здесь углубляться в то, для чего нужна это соответствие, нам сейчас важно другое. Несложно проверить, что \sigma — биекция. Таким образом у нас появляется связь двух функторов F и G с помощью соответствия \sigma. Оказывается, что даже более слабая версия этого соответствия бывает довольна полезна.

[Connection between categorical notion of adjunction and dual space/adjoint in vector spaces](https://math.stackexchange.com/questions/1009615/connection-between-categorical-notion-of-adjunction-and-dual-space-adjoint-in-ve)

[What is an intuitive view of adjoints? (version 2: functional analysis)](https://mathoverflow.net/questions/6552/what-is-an-intuitive-view-of-adjoints-version-2-functional-analysis)

[What is an intuitive view of adjoints? (version 1: category theory)](https://mathoverflow.net/questions/6551/what-is-an-intuitive-view-of-adjoints-version-1-category-theory)

[Can adjoint linear transformations be naturally realized as adjoint functors?](https://mathoverflow.net/questions/476/can-adjoint-linear-transformations-be-naturally-realized-as-adjoint-functors)

[Adjoint functors](https://en.wikipedia.org/wiki/Adjoint_functors)