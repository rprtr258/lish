#project

https://github.com/rprtr258/beedrill

[Дизайн документ](Дизайн%20документ)
[Backend REST API](Backend%20REST%20API)

analyser

no analyser — можно смотреть только на фенотипы

basic analyser — показывает генотипы (в виде букв/кода)

advanced analyser — раздекоденный генотип

superadvanced analyser — точное представление генотипов

codominance — числовые гены, выражающиеся в фенотипе средним арифметическим обеих аллелей, в продвинутом analyser показывается только отрезок

fertility (1, 2, 3, 4)

lifetime (short, medium, long...)

gene extraction - дорогая операция по вытаскиванию гена для культивации в несколько этапов:

вытащить гены одного организма

вытащить из них полезный trait

вытащить гены из таргетного организма

вытащить из них соответствующий ген

вставить полезный trait

имплантировать полученные гены обратно во второй организм

marketplace

trade bees, products, (money)

auction?

automatic auctioning?

bee house

bee house/apiary → bee house(apiary) lvl 1, 2, 3, ..?

bees keep dying even not in bee houses

energy supply?

building graph of bee houses, centrifuges, storages, etc.

bee houses coloring, labels

label for stable bees production

label for mutation bee houses / mutation for one specimen

automatically generate all wiki pages

admin page

all users history

breeding colors

colors gene

animation effect gene

energy supply

wires

engines

resource → energy consumption

interface sketch

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled.png)

[Forestry (Minecraft Mod)](https://github.com/ForestryMC)

[Meet Google Drive - One place for all your files](https://drive.google.com/drive/folders/1JW-tvwYYt_aN-llnmIIGrNkdUB-uKt-e?usp=sharing)

[Tutorial:Bee Basics (Apiculture)](https://ftbwiki.org/Tutorial:Bee_Basics_(Apiculture))

[Bee Breeding](https://feed-the-beast.fandom.com/wiki/Bee_Breeding)

[Bee Breeding Branches](https://feed-the-beast.fandom.com/wiki/Bee_Breeding_Branches)

[Master Apiarist Database v2.0](https://docs.google.com/spreadsheets/d/1_moZHLnL35_u-bJ7kFDxWDxY9OuMWK_4l0EB4wIx0_s/edit?type=view&gid=0&f=true&colid0=1&filterstr0=Lapis&sortcolid=6&sortasc=true&rowsperpage=250#gid=3)

[Minecraft Bees](https://web.archive.org/web/20150323071019/http://mc.nessirojgaming.eu/bees/index.php?title=Main_Page)

[Thesixler's Half-Assed Guide to How To Bees](https://imgur.com/a/F1fXf)

punnett square

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%201.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%202.png)

pedigree chart

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%203.png)

[Gendustry](https://minecraft-ru.gamepedia.com/Gendustry)

[Twitch Emotes - Bringing a little Kappa to you everyday](https://twitchemotes.com/search?query=bee)

![](data/static/old/someday_maybe/to_buy/Untitled%204.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%205.png)

## Текстуры

[ForestryMC/ForestryMC](https://github.com/ForestryMC/ForestryMC/tree/mc-1.12/src/main/resources/assets/forestry/textures/items)

## Мануал

[ForestryMC/ForestryMC](https://github.com/ForestryMC/ForestryMC/tree/b92a1a10e96b216e8eb72cc0ee2a61a2065bc486/src/main/resources/assets/forestry/manual/en_us)

<aside>
💡 Сделать возможность менять кучу дронов на принцессу, чтобы было куда девать ненужных дронов и была возможность получать принцесс.

</aside>

<aside>
💡 Авторежим - пчелы в улье автоматически продолжают работать после смерти.

</aside>

Класс генома пчелы

[ForestryMC/ForestryAPI](https://github.com/ForestryMC/ForestryAPI/blob/master/forestry/api/apiculture/IBeeGenome.java)

Enum хромосомы

[ForestryMC/ForestryAPI](https://github.com/ForestryMC/ForestryAPI/blob/master/forestry/api/apiculture/EnumBeeChromosome.java)

Виды пчел

[ForestryMC/ForestryMC](https://github.com/ForestryMC/ForestryMC/blob/926645f2a4df313f692b453a713cb558ae6a7be1/src/main/java/forestry/apiculture/genetics/BeeDefinition.java)

Стили текста, фонов и тд

forestry.css
```css
.gui {
    title: #404040;
    screen: #ffffff;
    book: #000000;
}
.gui.table {
    header: #ababab;
    row: #ababab;
}
.gui.beealyzer {
    binomial: #14d50b;
    recessive: #3687ec;
    dominant: #ec3661;
}
.gui.mail {
    lettertext: #cfa738;
    text: #6689dc;
}
.gui.greenhouse {
    temperature-header: #404040;
    humidity-header: #404040;
    modifiers-subheader: #aaafb8;
}

.ledger.error {
    background: #ff3535;
    header: #e1c92f;
    text: #000000;
}
.ledger.hint {
    background: #ea38ff;
    header: #e1c92f;
    text: #000000;
}
.ledger.owner {
    background: #ffffff;
    header: #e1c92f;
    subheader: #aaafb8;
    text: #000000;
}
.ledger.climate {
    background: #35a4ff;
    header: #e1c92f;
    subheader: #aaafb8;
    text: #000000;
}
.ledger.habitatformer {
    background: #2ca22c;
    header: #e1c92f;
    subheader: #aaafb8;
    text: #000000;
}
.ledger.power {
    background: #d46c1f;
    header: #e1c92f;
    subheader: #aaafb8;
    text: #000000;
}
.ledger.farm {
    background: #2ca22c;
    header: #e1c92f;
    subheader: #aaafb8;
    text: #000000;
}
.item.circuit.basic {
    primary: #191919;
    secondary: #6dcff6;
}
.item.circuit.enhanced {
    primary: #191919;
    secondary: #cb7c32;
}
.item.circuit.refined {
    primary: #191919;
    secondary: #c9c9c9;
}
.item.circuit.intricate {
    primary: #191919;
    secondary: #e2cb6b;
}
```

[addons-forestry](https://www.curseforge.com/minecraft/mc-mods/mc-addons/addons-forestry)

[Bee Genetics: Modifying Bees to have Desirable Traits](http://modjominecraft.blogspot.com/2014/12/bee-genetics-modifying-bees-to-have.html)

minigames

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%206.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%207.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%208.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%209.png)

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2010.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2011.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2012.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2013.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2014.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2015.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2016.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2017.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2018.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2019.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2020.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2021.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2022.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2023.png)

[Bees](https://ftb.gamepedia.com/Category:Bees?pageuntil=Oblivion+Bee#mw-pages)

[Bee Species](https://ftbwiki.org/Bee_Species)

draw this as planar graph

![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2024.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2025.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2026.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2027.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2028.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2029.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2030.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2031.png)
![](data/static/old/someday_maybe/programming_projects/Bee%20breeder/_static/Untitled%2032.png)

[Swarm Simulator](https://www.swarmsim.com/#/tab/meat/unit/hive)

[Critter Mound](https://crittermound.com/)

[Numbers Getting Bigger - Envato Tuts+ Game Development Tutorials](https://gamedevelopment.tutsplus.com/series/numbers-getting-bigger--cms-847)

[The Math of Idle Games, Part I](https://blog.kongregate.com/the-math-of-idle-games-part-i/)

[](https://pixl.nmsu.edu/files/2018/02/2018-chi-idle.pdf)

[bees in games - Google Search](https://www.google.com/search?q=bees+in+games)

[bee breeding games - Google Search](https://www.google.com/search?q=bee+breeding+games)

[Honey Combs](https://binnie.mods.wiki/wiki/Honey_Combs)

улей / теплица с параметрами окружения (свет / погода / влажность  / температура) - docker like containers

[Master Apiarist Database v2.0](https://docs.google.com/spreadsheets/d/1_moZHLnL35_u-bJ7kFDxWDxY9OuMWK_4l0EB4wIx0_s/edit?usp=sharing)

[Modest Bee](https://ftb.gamepedia.com/Modest_Bee)

[Modest Bee](https://ftbwiki.org/Modest_Bee)

[Nonsanity's Bee Breeding Charts](http://nonsanity.com/bees/#)

[Дизайн и математика игр-кликеров](https://habr.com/ru/post/335754/)

[Guide to genetics](https://tgstation13.org/wiki/Guide_to_genetics)

[Guide to xenobiology](https://tgstation13.org/wiki/Guide_to_xenobiology)

[Reagent: Minimalistic React for ClojureScript](https://reagent-project.github.io/)

[Comparing Elm to React/Redux](https://dev.to/rametta/comparing-elm-to-react-redux-2emo)

[Clojure - Concurrent Programming](https://www.tutorialspoint.com/clojure/clojure_concurrent_programming.htm)
https://habr.com/en/post/657603/
