# headlessui — spec pro práci na téhle knihovně

Headless (nestylované) UI komponenty nad [templ](https://templ.guide). Tvar API vychází
z [bits-ui](https://bits-ui.com/). Cíl číslo jedna je přístupnost: správný výchozí stav
se vyrenderuje z Go, JavaScript jen udržuje interaktivitu.

Cílový stack konzumenta: **Tailwind v4** na styly, **htmx** na interakci se serverem,
**Alpine** jen tam, kde to komponenta opravdu potřebuje. Nad touhle knihovnou vznikne
stylovaná UI knihovna — v **samostatném repu**, ne tady.

## Příkazy

```bash
make example
templ generate && go test ./...
templ generate -check && gofmt -l . && go vet ./...
cd alpine && npm run build
```

Po každé změně `.templ` souboru nebo attrs funkce pusť `templ generate`. Generované
`*_templ.go` soubory se commitují, aby knihovna šla použít bez templ CLI.

Po každé změně v `alpine/src` pusť `npm run build`. `alpine/dist` se commituje ze stejného
důvodu — Go ho embeduje, takže knihovna jde použít bez Node.

## Pravidla, která se neporušují

- **Jeden plochý Go balíček.** Žádný druhý modul, žádný `go.work`. Podsložky jsou jen dvě:
  [alpine/](alpine) — npm balíček s Alpine pluginy, žádné `.go` soubory — a
  [example/](example) — `package main` se stránkou se všemi komponentami
  (`make example` = build Alpine, `templ generate`, `go run ./example`), který knihovnu
  importuje jako kdokoli jiný. Žádná knihovní logika v něm není. Zatím je v `.gitignore`
  (stejně jako `.claude/`), takže `go.mod` nesmí mít závislost, kterou potřebuje jen on.
- **Dvojice souborů na komponentu:** `komponenta.go` (props a atributy) + `komponenta.templ`
  (tenký markup). Sdílené vnitřnosti bez vlastní komponenty mají vlastní soubor: `base.go`,
  `anchor.go` (pozicování), `options.go` (Listbox a Combobox), `portals.go` (teleport, overlaye
  a výchozí přechody), `pending.go` (skript pro stav běžícího requestu), `checkable.go` (Checkbox, Switch a Radio jako button). Když komponenta potřebuje vlastní JS, přibude plugin
  `alpine/src/komponenta.js` a Go si jeho build vloží přes `//go:embed`.
- **Žádné komentáře v kódu.** Ani doc komentáře, ani v testech. Čitelnost nese pojmenování,
  vysvětlení patří do [README.md](README.md).
- Každá komponenta má vlastní `<Komponenta>Props` s embedovaným `Base`, takže se v literálu
  píše `TextInputProps{ID: "email"}` (embedded field selectors, Go 1.27+).
- Každá komponenta musí jít **nastylovat** (`Class`), **doplnit o libovolný atribut**
  (`Attrs`) a **obohatit o vlastní handler** (`Events`) bez ztráty interního chování.
- `.templ` jen zavolá attrs funkci a vypíše `{ children... }`. Veškerá logika atributů
  je v `.go`, aby šla testovat přímo.
- **Nikdy nepiš `@focus`, `:class` apod. přímo do markupu**, když jde stejný atribut poslat
  přes `Attrs` nebo `Events` — prohlížeč by jeden z duplicitních atributů tiše zahodil.

## Base a merge

`Base` nese `Events`, `ID`, `Class` a `Attrs`. Finální atributy skládá `Base.merge(own)`,
případně `mergeEvents(own, internal)` u komponent s vlastními handlery:

1. `Attrs` od volajícího,
2. vlastní (`own`) atributy komponenty — ty vyhrávají, aby nešlo přepsat a11y wiring,
3. `ID` (nastavuje se úplně nakonec, vyhrává i nad `Attrs`),
4. `class` se **slučuje** přes `MergeClass` (pořadí `Attrs` → `Class`), `style` se
   připojí za vlastní inline styl komponenty (styl z `Attrs` je poslední, takže vyhrává),
5. pojmenované `Events` se **zřetězí za** interní handler stejného typu (`chain`).

**Nic, co volající pošle, se tiše nezahazuje.** Handler v `Attrs` se zřetězí se stejně
pojmenovaným `Events`, třída v `Attrs` se sloučí s `Class`. Interní stav musí doběhnout
dřív než uživatelský výraz, proto je pořadí v `chain` pevné.

`MergeClass` je exportovaná proměnná, ne závislost — stylovaná vrstva si do ní zapojí
tailwind-merge. Jádro zůstává bez závislostí (kromě templ).

## Field

Formulářové komponenty sdílí stav přes `Field`, aby label a input nemohly tvrdit dvě různé
věci. `Field` vloží `FieldProps` do kontextu, ostatní si o ně řeknou přes `fieldFrom(ctx)`.

- `Invalid` je true vždy, když je nastavená `Error`.
- `ID` se odvodí z `Name`, a když chybí obojí, vygeneruje se (`nextID`).
- Wrapper dostane `id="<id>-field"` a `data-field="<name>"` jako cíl pro htmx swap.
- `Description` a `Error` berou text z `FieldProps`, ne z children — input se renderuje dřív
  než oni, takže by jinak `aria-describedby` ukazovalo na neexistující ID.
- Uvnitř `Fieldu` je `Field` jediná pravda: `Name`, `Value`, `Disabled`, `Readonly`,
  `Required` a `Invalid` zadané na vstupu se ignorují.
- **Každá komponenta, která z `Fieldu` čerpá, musí fungovat i samostatně.** Fallback staví
  atributy z vlastních props a u `Description`/`Error` vypíše children.
- Vstupy si stav berou přes `controlFrom(ctx, vlastní)`, zaškrtávací (`Checkbox`, `Switch`,
  `Radio`) přes `checkableFrom`. Obě řeší `Field`, `Fieldset` i samostatný režim najednou.
- Zaškrtávací prvky jsou `<button>` s rolí, skrytým inputem uvnitř (jen když jsou
  zaškrtnuté) a `name`/`value` na samotném buttonu (htmx u buttonu posílá jen jeho vlastní
  hodnotu, nikdy potomky). `Required` řeší validační input hned za buttonem
  (`validationAttrs`, sdílené s Listboxem a Comboboxem).
- `Label` ve `Fieldu` dostane `id="<id>-label"` pro `aria-labelledby` (`Meter`, list
  `Listboxu`).
- `Fieldset` je pro skupinu to, co `Field` pro jeden vstup (`fieldsetFrom(ctx)`). `Field`
  **bez** `Name` uvnitř `Fieldsetu` je volba skupiny (jméno a zaškrtnutí z fieldsetu, ID
  `<fieldset>-<value>`, wrapper `data-option`), `Field` **s** `Name` je samostatné pole, které
  z fieldsetu převezme jen `Disabled`. `Description`/`Error` hledají vlastníka přes
  `ownerFrom(ctx)`: nejdřív `Field`, pak `Fieldset`.
- `Form` na každém vstupu, `Fieldu` i `Fieldsetu` se renderuje jako nativní `form="<id>"`.
  Formulářové prvky musí fungovat i mimo `<form>`, žádný plugin nehledá `closest('form')`.

## Stav ven přes data-*

Stavové atributy se renderují **holé** (`data-invalid`, ne `data-invalid="true"`), protože
na to míří Tailwind v4 varianta `data-invalid:`. Stav visí na wrapperu, labelu, vstupu
i chybě zároveň, takže se skoro nikdy nemusí sahat po `group`. CSS třídy ke komunikaci
stavu nikdy — knihovna sama žádnou třídu negeneruje.

Tlačítka jádra (`Button`, `DropdownItem`, spouštěče, zavírací tlačítka, akce toastu,
záložky) nesou vždy buď `data-enabled`, nebo `data-disabled` (`buttonState` v `button.go`),
aby jedna sada tříd stylovaného tlačítka šla dát na kterékoli z nich.

## JavaScript

- Komponenta má fungovat bez JS tak daleko, jak to jde. Žádný stavový atribut, který by
  po první interakci zastaral (proto `TextInput` bez Alpine nerenderuje `data-empty`).
  Výjimka jsou `Checkbox`, `Switch` a `Radio`: kvůli vlastní ikoně jsou to buttony a bez JS
  jen zobrazí stav.
- `TextInput` má Alpine **volitelný, výchozí vypnutý** (`Alpine: true`). S Tailwindem
  to pokryje `focus-visible:` a `placeholder-shown:`.
- Veškerý JS žije v [alpine/](alpine) jako Alpine pluginy — budoucí veřejný npm balíček
  `@daarxwalker/headlessui-alpine`. Alpine je v něm `peerDependency`, ne závislost.
- Plugin komponenty je `alpine/src/komponenta.js` s `export default function (Alpine)`,
  který registruje přes `Alpine.data(...)`. Přidej ho do `src/index.js` (export i celkový
  plugin) a vytvoř `alpine/builds/komponenta.js`, který ho předá `register`.
- Build dělá tři věci: `dist/module.esm.js` pro `import` a `Alpine.plugin(...)`,
  `dist/cdn.min.js` pro `<script>` se vším a `dist/komponenta.min.js` pro Go.
- **Žádná komponenta nerenderuje `<script>`, nic se nespouští samo.** Konzument načte
  pluginy jednou sám: `Scripts()` (celý `dist/cdn.min.js`, [scripts.go](scripts.go)),
  jednotlivé `<Komponenta>Script()` (`//go:embed alpine/dist/komponenta.min.js`
  + `templ.NewOnceHandle`, jednou na stránku), nebo npm balíček. `register` se umí
  zaregistrovat i když skript dorazí až po startu Alpine.
- V metodách pluginu sahej na kořen komponenty přes `this.$root`, ne `this.$el` — když se
  metoda zavolá z potomka (`<button @click="hide()">`), `$el` je ten potomek.
- **Každé pole instance deklaruj v datech předem** (`layer: null`, `cleanup: null`, …).
  Alpine slučuje scope vnořených komponent, takže `this.layer = …` v `init` vnořeného
  dropdownu by jinak přepsalo `layer` rodičovského modalu.
- Alpine je jen na klientskou interakci (fokus, otevření, zvýraznění, klávesnice, hodnota
  ve skrytém inputu). Filtrování, řazení, výběr řádků a překreslování dělá server přes htmx.
- Overlaye (`Modal`, `Dropdown`, `Tooltip`, `Listbox`, `Combobox`) se registrují do
  sdíleného zásobníku v [alpine/src/layers.js](alpine/src/layers.js): Escape jen nejvyšší
  vrstvě, klik mimo zavře vrstvy nad tou, do které se kliklo, modal nastaví `inert` všemu
  kromě sebe, vrstev nad sebou, `[data-hui-toaster]` a `[data-hui-announcer]`. Zásobník je singleton na `Alpine`,
  protože každý `*.min.js` si modul přibalí zvlášť.
- Položky, které přijdou z htmx partialu, nemají kontext komponenty. Proto `id`,
  `aria-selected`, `data-selected`, `data-highlighted` a indikátor položek listboxu
  maluje plugin (`paint()`), Go renderuje jen počáteční stav.
- Preferuj shorthandy (`@focus`, `:data-focused`) před `x-on:`/`x-bind:`.
- **Ovládací konstanty.** Každá komponenta se stavem otevření exportuje Alpine výrazy
  (`TooltipShow`, `DropdownToggle`, `ModalHide`, …), které jdou dát na libovolný element
  a událost přes `Attrs`/`Events`. Volají metodu s **unikátním jménem pro komponentu**
  (`tooltipToggle($el, $event)`, ne `toggle()`), protože Alpine slučuje scope vnořených
  komponent. Vlastní části (`*Trigger`, `*Close`) se skládají ze stejných konstant.
  Kotvené overlaye (`Tooltip`, `Dropdown`) přenesou `anchor-name` a a11y na
  `$el`, který je otevřel. Konstanty se jmenují `Show`/`Hide`/`Toggle` jako okenní události,
  protože `TooltipClose` a `ModalClose` jsou komponenty.
- **Výchozí přechody.** Části s `x-show` dostanou `x-transition` (overlay modalu
  `x-transition.opacity`, obsah Collapsible `x-collapse`) přes
  `transition(own, attrs, directive)` v `portals.go`. Jakýkoli `x-transition…`/`x-collapse…`
  klíč v `Attrs` výchozí nahradí, `false` ho vypne. `@alpinejs/collapse` je jediná přibalená
  závislost npm balíčku (Alpine zůstává peer), pluginy ho registrují jednou.

## htmx

- **Se serverovým překreslováním a appendováním se počítá vždy.** Každá část musí jít
  vyrenderovat v partialu a přežít výměnu i doplnění. Seznamy mají `With<Komponenta>(ctx,
  props)` (`WithListbox`, `WithCombobox`), aby partial vyrenderoval stejné HTML jako celá
  stránka. ID položek se odvozují z hodnoty v Go i v JS (`idSafe`), stav nikdy nežije
  v DOMu položek (hodnota, skryté inputy, cache popisků) a plugin po každé mutaci dorovná
  `id`, `aria-*` a `data-*`.
- Infinite loading je htmx sentinel (`ListboxLoader`, `ComboboxLoader`, `TableLoader`,
  `data-hui-loader`): šipka za poslední položkou počká na další stránku, list má během
  requestu `aria-busy` a `data-pending`, položky `aria-setsize="-1"`. Tabulka má
  `aria-rowcount`/`aria-rowindex`.

- **Běžící request je `data-pending` + `aria-busy="true"`**, nikdy `data-loading` ani třída.
  [alpine/src/pending.js](alpine/src/pending.js) ho po dobu requestu nastaví zdroji,
  submitteru a `hx-target`, počítá souběžné requesty a serverem vyrenderovaný `data-pending`
  (`Button{Pending}`) nesundá; při `HX-Redirect`/`HX-Refresh` ho nechá do odchodu ze stránky. Konec requestu poslouchá na `document` a páruje přes `ctx`
  (htmx 4 posílá události odpojeného elementu na `document`). Nepotřebuje Alpine, přibalí ho
  `register` do každého pluginu i do `Scripts()`. Bez Alpine ručně `PendingScript()`.
- Loadery nekreslí vlastní spinner, spinner je `children` (headless).
- **Hledání v listboxu je `data-searching`** na search inputu a listu, od prvního znaku do
  odpovědi serveru na tentýž dotaz ([alpine/src/options.js](alpine/src/options.js)).
  Odlišuje přepsání celých dat od donačítání stránky, které dává jen `data-pending`.

- Cílová verze je **htmx 4** (`<hx-partial>`, `hx-include="this"` zahrne vnořené inputy,
  po swapu fokusuje `autofocus`).
- `hx-*` jde přes `Attrs`, `merge` se jich nedotkne (přepisuje jen klíče začínající `@`).
- **Žádné typované HX API** — vědomé rozhodnutí, `Attrs` stačí.
- Stabilní `id` je důležité: htmx na něm drží fokus mezi requesty, míří na něj `hx-target`
  i `hx-swap-oob`. Primární cesta je `Name`; pro zbytek je `WithIDScope(ctx, prefix)`.

## Checklist pro novou komponentu

1. `komponenta.go`: `<Komponenta>Props` embedující `Base` a funkce
   `komponentaAttrs(ctx, props) templ.Attributes` — případně s dalšími návratovými
   hodnotami (`(attrs, text, render bool)` u `Description`/`Error`,
   `(attrs, isLink bool)` u `Button`), když se část markupu má skrýt nebo přepnout.
   `ctx` v signatuře jen tehdy, když ho komponenta opravdu používá.
2. Patří-li do `Field` (label, popis, chyba, libovolný vstup), implementuj čtení
   z `fieldFrom(ctx)` **i** fallback bez něj.
3. Potřebuje-li JS, plugin v `alpine/src`, build entry v `alpine/builds` a `npm run build`.
4. `komponenta.templ`: tenký wrapper, který zavolá attrs funkci a vypíše `{ children... }`.
5. `komponenta_test.go`: jeden `Test<Komponenta>` se subtesty (viz níž).
6. Ověř: jde nastylovat, jde poslat libovolný atribut, jde připojit vlastní handler bez
   ztráty interního chování, funguje s klávesnicí, má správné `aria-*` a `data-*` bez JS.
7. Doplň sekci do [README.md](README.md) — anglicky — a řádky do *Styling reference*.

## Testy

- Jeden soubor na komponentu, v něm **jeden** `Test<Komponenta>` a uvnitř `t.Run` subtesty.
- [testify](https://pkg.go.dev/github.com/stretchr/testify): `assert` běžně, `require` tam,
  kde nemá smysl pokračovat.
- [gofakeit](https://pkg.go.dev/github.com/brianvoe/gofakeit/v7) na data, která mají jen
  protéct skrz (`Word()` pro `Name`, `Email()` pro `Value`, `Sentence()` pro texty). Kde
  na konkrétní hodnotě záleží (`type="submit"`, `role="link"`), piš literál.
- Sdílené pomocníky drž v [tests.go](tests.go) s prefixem `test`: `testRender`,
  `testRenderIn`, `testWithChildren`, `testCount`.
- Attrs funkce se dají testovat přímo (`modalAttrs(ctx, props)`), což je přesnější než
  hledání podřetězce v HTML — dělej to tak všude, kde jde o jednotlivé atributy.
- Názvy subtestů anglicky, ať jdou filtrovat: `go test -run 'TestButton/pending' -v .`

## Rozhodnutí, která už padla

Neotáčej je bez domluvy — každé má důvod:

| Rozhodnutí | Proč |
|---|---|
| Žádné nativní overlaye (`<dialog>`, `popover`), všechno přes `x-teleport` do jednoho `Portals` (`#hui-portals`) | jeden mechanismus pro všechny overlaye: sdílený zásobník vrstev, `x-transition`, ovládání stavu odkudkoli. Pozice dál čistě CSS anchor positioningem (`Placement`, `anchor.go`), bez Floating UI |
| `Modal` je div s `role="dialog"` řízený Alpinem, ne nativní `<dialog>` | stav musí jít ovládat odkudkoli (`x-model` přes `x-modelable`, události, `open`) a `x-show` s `x-transition`; `open` je jediný zdroj pravdy a focus trap, `inert` okolí, Escape, klik mimo, zámek scrollu i vracení fokusu z něj odvozuje plugin. Otevřené modaly tvoří zásobník, reaguje jen ten nejvyšší |
| JS je npm balíček v `alpine/`, Go embeduje jeho `dist` | jeden zdroj pro bundlery, CDN i templ; konzument bez Node dál nic neinstaluje |
| Skripty jsou volitelné, komponenty žádný `<script>` nevkládají | všechno je v Alpine, které si konzument načítá sám (`Scripts()`, `<X>Script()`, npm); partialy a markup tak nikdy nespouští kód samy od sebe |
| `Checkbox`, `Switch` i `Radio` jsou `<button>` s rolí a skrytým inputem, ne nativní input | vlastní ikona (`CheckboxIndicator`, `SwitchThumb`, `RadioIndicator`) a jednotné stylování přes `data-checked:`; cena: bez JS se nedají přepnout. Radio skupina (podle `name` a formuláře) má roving tabindex a šipky |
| Overlay a obsah `Modalu` (a obsah ostatních overlayů) se teleportují do povinného `Portals` v `<body>`, obsah uvnitř overlaye cestuje s ním | žádný předek s `transform`, `filter` nebo `overflow` je neuvězní. Cena: bez JS se obsah nevyrenderuje, `<form>` musí být uvnitř `ModalContent`, htmx se v obsahu zpracuje až při otevření (`htmx.process`) a `z-index` si určuje konzument |
| `Modal` je rozložený na části jako bits-ui Dialog (`Portals`, `Trigger`, `Overlay`, `Content`, `Title`, `Description`, `Close`) | bez `::backdrop` potřebuje overlay vlastní element; `Title`/`Description` dávají `aria-labelledby`/`aria-describedby` z Go zadarmo |
| Alpine v `TextInput` je opt-in a výchozí vypnutý | cílový stack je Tailwind + htmx, `focus-visible:` a `placeholder-shown:` to pokryjí bez JS |
| `MergeClass` je proměnná, ne závislost | jádro zůstane bez závislostí, merger si zapojí stylovaná vrstva (pozor: `tailwind-merge-go` umí jen Tailwind 3.0–3.4) |
| `Field` dává wrapperu `id="<id>-field"` | aby šel jeden formulářový řádek vyměnit htmx swapem bez ručních atributů |
| Tři výběrové komponenty podle Headless UI pro React: `Select`, `Listbox`, `Combobox` | dřív byly `Select` se searchem a `Combobox` totéž; teď je každá na jiný případ a liší se chováním, ne jen jménem |
| `Select` je nativní customizable `<select>` (`<button><selectedcontent>` jako první potomek), bez JS | picker, klávesnici, hodnotu i validaci drží prohlížeč; `appearance: base-select` a `::picker(select)` nastavuje konzument (jádro nerenderuje styl), prohlížeč bez podpory ukáže obyčejný select. Placeholder je `<option value="" data-placeholder>` hned za triggerem. `Multiple`, search ani serverové seznamy nemá, od toho je `Listbox` |
| `Listbox` je vlastní listbox se skrytými inputy | plně stylovatelný a kompozitní, umí `Multiple`. `Required` řeší skrytý validační input bez jména. Vždy `ListboxContent` (popup) > `ListboxList` (`role="listbox"`), search je volitelný (`Searchable` + `ListboxSearch`): server překresluje jen list, takže search nepřijde o fokus; se searchem je trigger obyčejný button a `aria-activedescendant` drží search input |
| `Combobox` má input přímo v poli, popup obsahuje jen list, filtruje server | vede uživatelská přívětivost (tvar Headless UI), ne litera WAI-ARIA. `ComboboxInput` ukazuje popisek hodnoty a nemá `name`; dotaz (jen to, co uživatel napsal) drží plugin, pouští ho přes `hui-combobox:query` a do requestu ho doplní v `htmx:config:request` jako `QueryName`. Zavření vrátí popisek a nefiltrovaný list. `ComboboxTrigger` je volitelný chevron mimo tab order, kotva je kořen. Logiku sdílí s `Listboxem` (`options.go`, `options.js`) |
| Shoda s dotazem je `<mark data-hui-match>` v položkách | maluje ho plugin (`mark()` v `options.js`) po každé mutaci listu i změně dotazu, bez ohledu na velikost písmen a diakritiku, mimo `aria-hidden` části. `<mark>` v DOMu, ne CSS Highlight API, protože to neumí `font-weight`. U Comboboxu se nikdy neznačí zobrazená hodnota, jen napsaný dotaz |
| `Collapsible` přes Alpine a tlačítka, ne `<details>` | tvar bits-ui, `x-model`, `x-collapse`; stav se renderuje z Go. Hlavička umí oba režimy: trigger přes celý řádek, nebo titulek (`*Title`) + ovládací tlačítko vpravo, volitelně `Clickable` (klik kamkoli do řádku kromě interaktivních potomků) |
| `Tooltip` nemá `TooltipTrigger`, jen veřejné konstanty | tooltip se ovládá výhradně programově: konzument dá `TooltipShowDelayed`/`TooltipFocus`/`TooltipToggle`/… na vlastní element a sám si určí hover, focus nebo klik. Výchozí kotva je element uvnitř s `aria-describedby="<id>-content"` |
| `Toast` jen ze serveru přes `<hx-partial>` do `#hui-toaster` | žádné klientské API na toasty, server je jediný zdroj; plugin řeší jen časovač, pauzu a zavření |
| `Table` bez JS | výběr řádků i řazení dělá server přes htmx, komponenta jen renderuje `aria-sort` a `data-selected` |
| Mimo `<form>` jen nativní `form` atribut a dokumentace (`hx-include="this"`) | jádro samo žádný `hx-*` nenastavuje |
| Volba ve `Fieldsetu` (wrapper) nemá vlastní `data-checked` | stav nese control uvnitř, řádek se styluje přes `has-data-checked:` |
| Žádné typované HTMX API | `Attrs` stačí a jádro nebude zaostávat za htmx |
| Stylovaná vrstva do samostatného repa | tahle knihovna zůstane bez Tailwindu |
| `tests.go` **bez** sufixu `_test` | vědomé rozhodnutí uživatele; důsledek je, že testify a gofakeit jsou běžné závislosti v `go.mod` — neopravuj to přejmenováním |
| Žádné E2E testy | byly napsané (server + chromedp) a smazané jako zbytečná nadstavba nad malou knihovnu |
| README je primárně návod pro agenta | agent podle něj staví stylovanou UI knihovnu a showcase appku pro konkrétní projekt; sekce *Building a project UI library*, *Building a showcase app* a *Styling reference* musí odpovídat API a vzoru v `example/` |
| Licence MIT, Go modul i npm balíček | `alpine/LICENSE` je kopie kořenové, npm ji přibalí i přes `files: ["dist"]` |
| `example/` jako jednoduchý `package main` | aby šly komponenty ručně vyzkoušet v prohlížeči (`go run ./example`); jen ukázková stránka a pár endpointů pro htmx, žádné testy ani tooling |

## Co nedělat

- Nezakládat další podsložky, druhé moduly ani `go.work`.
- Nepsat komentáře do kódu.
- Nepřidávat nástroje a strukturu, o které nikdo nežádal.
- Nespouštět žádný git příkaz, který mění stav repa, bez výslovného pokynu.

## Co zbývá

Hotové komponenty: `Field`, `Fieldset`, `Legend`, `Label`, `Description`, `Error`,
`TextInput`, `Textarea`, `Checkbox`, `Switch`, `Radio`, `Range`, `Select`, `Listbox`, `Combobox`,
`FileUpload`, `Button`, `Portals`, `Modal`, `Dropdown`, `Tooltip`, `Tabs`,
`Collapsible`, `Breadcrumb`, `Table`, `Meter`, `Toaster`/`Toast`, `Announcer`, `Alert`,
`SkipLink`.

Nová komponenta jede podle checklistu výš; formulářová se napojí na `Field` i `Fieldset`
přes `controlFrom`/`checkableFrom`, overlay na `Portals`, `Placement` a `layers.js`.
