const {StyledComponent} = Torus;

const EXAMPLES = [
  {
    slug: 'helloworld',
    description: 'Basic hello world example.',
  },
  {
    slug: 'fizzbuzz',
    description: 'FizzBuzz implemented in Ink.',
  },
  {
    slug: 'fibonacci',
    description: 'Naive and optimized (memoized) fibonacci sequence generators.',
  },
  {
    slug: 'list',
    description: 'Various list-related functions implemented in Ink.',
  },
  {
    slug: 'pi',
    description: 'Estimating Pi with a Monte Carlo (statistical) simulation.',
  },
  {
    slug: 'prime',
    description: 'A prime sieve (for computing all primes under N) written by composing functions.',
  },
  {
    slug: 'quicksort',
    description: 'An efficient in-place quicksort with Hoare partitioning.',
  },
  {
    slug: 'newton',
    description: 'Newton\'s root finding algorithm applied to compute square roots quickly.',
  },
  {
    slug: 'std',
    description: 'The Ink standard library, a set of useful functions and utilities for composing programs',
  },
  {
    slug: 'str',
    description: 'The Ink standard string library, including idiomatic implementations of tail-recursive string manipulation functions',
  },
  {
    slug: 'quine',
    description: 'A quine in Ink -- a program whose output is its own source code',
  },
];

async function evalInk(inkSource) {
  const resp = await fetch('/eval', {
    method: 'POST',
    headers: {
      'Content-Type': 'text/plain',
    },
    body: inkSource,
    cache: 'no-cache',
    credentials: 'same-origin',
  });
  return await resp.json();
}

async function getExample(slug) {
  const resp = await fetch(`/static/ex/${slug}.ink`, {
    method: 'GET',
    cache: 'no-cache',
  });
  return await resp.text();
}

function debounce(fn, duration) {
  let timer;
  const dfn = (...args) => {
    clearTimeout(timer);
    timer = setTimeout(() => {
      fn(...args);
    }, duration);
  }
  return dfn;
}

function ExampleItem(item, insertCb) {
  return jdom`<li>
    <button class="exampleItem block" onclick="${async () => {
      insertCb('\` loading example... \`')
      insertCb(await getExample(item.slug));
    }}" tabindex="0">
      <div class="filename">${item.slug}.ink</div>
      <div class="description">${item.description}</div>
    </button>
  </li>`;
}

function ExampleList(insertCb, closeCb) {
  return jdom`<div class="backdrop" onclick="${closeCb}">
    <div class="examples" onclick="${evt => evt.stopPropagation()}">
      <h2>Examples</h2>
      <p>Tap on an example to try it.</p>
      <p>You can also find <a href="https://github.com/thesephist/ink/tree/master/samples" target="_blank" rel="noopener noreferer">more sample Ink programs on GitHub</a>, though many of them won't be able to run in this browser context today, if they take some user input or generate files as output.</p>
      <ul>
        ${EXAMPLES.map(it => ExampleItem(it, insertCb))}
      </ul>
    </div>
  </div>`;
}

class IOBox extends StyledComponent {
  init() {
    this.waiting = false;
    this.showSaved = false;
    this.verticalSplit = window.innerWidth > 600;
    this.flipState = 0;

    const lastSave = window.localStorage.getItem('ink_eval_save');
    this.stdin = lastSave || `out('Hello, World!')`;
    this.result = {
      exit: 0,
      error: null,
      output: 'click Run to see output...',
      duration: 0,
    };

    window.addEventListener('beforeunload', evt => {
      this.persistInput(false);
    });

    this.persistInput = this.persistInput.bind(this);
    this.debouncedPersist = debounce(this.persistInput, 5000);
    this.handleInput = this.handleInput.bind(this);
    this.handleKeydown = this.handleKeydown.bind(this);
    this.switchSplit = this.switchSplit.bind(this);
    this.requestEval = this.requestEval.bind(this);
  }

  insert(src) {
    this.stdin = src;
    this.render();
  }

  persistInput(showBadge) {
    window.localStorage.setItem('ink_eval_save', this.stdin);
    if (showBadge === false) {
      return;
    }

    this.showSaved = true;
    this.render();
    setTimeout(() => {
      this.showSaved = false;
      this.render();
    }, 1200);
  }

  handleInput(evt) {
    this.stdin = evt.target.value;
    this.debouncedPersist();
  }

  handleKeydown(evt) {
    const ctrlCmd = evt.ctrlKey || evt.metaKey;
    if (evt.key === 'Enter' && ctrlCmd) {
      this.requestEval();
    } else if (evt.key.toLowerCase() === 's' && ctrlCmd) {
      evt.preventDefault();
      this.persistInput(true);
    } else if (evt.key === 'Tab' && !evt.shiftKey) {
      evt.preventDefault();
      const idx = evt.target.selectionStart;
      if (idx !== null) {
        const front = evt.target.value.substr(0, idx);
        const back = evt.target.value.substr(idx);
        evt.target.value = front + '\t' + back;
        evt.target.setSelectionRange(idx + 1, idx + 1);
      }
    }
  }

  switchSplit() {
    this.verticalSplit = !this.verticalSplit;
    this.flipState ++;
    this.render();
  }

  async requestEval() {
    if (!this.waiting) {
      this.waiting = true;
      this.render();

      try {
        this.result = await evalInk(this.stdin);
        this.render();
      } catch (e) {
        this.result.output = 'Client error: ' + e.toString();
      }
      this.waiting = false;
      this.render();
    }
  }

  buttons() {
    return jdom`<div>
      <button
        class="block mobile-hidden"
        title="Flip / rotate panes"
        onclick=${this.switchSplit}>Flip</button>
      <button
        class="block"
        title="Save to your browser"
        onclick=${this.persistInput}>Save</button>
      <button
        class="accent block"
        title="Run Ink program & show output"
        onclick=${this.requestEval}>Run</button>
    </div>`;
  }

  styles() {
    return css`
    display: flex;
    flex-direction: ${(
      this.verticalSplit ? 'row' : 'column'
    ) + (
      this.flipState % 4 < 2 ? '' : '-reverse'
    )};
    justify-content: space-between;
    height: 0;
    flex-grow: 1;
    flex-shrink: 1;

    .inkInputBox,
    .inkOutputBox {
      width: 50%;
      box-sizing: border-box;
      position: relative;
      flex-shrink: 1;
    }
    .divider {
      background: #222;
      width: 3px;
      height: 100%;
      flex-shrink: 0;
    }
    &.flipped {
      .inkInputBox,
      .inkOutputBox {
        height: 50%;
        width: 100%;
      }
      .divider {
        width: 100%;
        height: 3px;
      }
      .savedBadge {
        top: 14px;
      }
    }

    textarea,
    .output {
      font-family: 'Dank Mono', 'Menlo', 'Monaco', monospace;
      font-size: 1em;
      box-sizing: border-box;
      padding: 12px 14px;
      background: transparent;
      resize: none;
      width: 100%;
      height: 100%;
      overflow-y: auto;
      -webkit-overflow-scrolling: touch;
      border: 0;
      &:focus {
        outline-color: var(--block-accent-color);
      }
    }
    .output {
      background: #eee;
    }
    .elapsedTime {
      color: #777;
    }
    code {
      font-family: 'Dank Mono', 'Menlo', 'Monaco', monospace;
      display: block;
      font-size: 1em;
      line-height: 1.4em;
      margin-bottom: .4em;
    }
    .savedBadge {
      position: absolute;
      top: 8px;
      right: 8px;
      color: #555;
      transition: opacity .2s;
      padding: 4px 6px;
      background: rgba(255, 255, 255, 0.9);
      border-radius: 3px;

      &.hidden {
        opacity: 0;
      }
    }
    `;
  }

  compose() {
    return jdom`<div class="ioBox ${this.verticalSplit ? '' : 'flipped'}">
      <div class="inkInputBox">
        <textarea
          class="inputBox"
          oninput="${this.handleInput}"
          onkeydown="${this.handleKeydown}"
          placeholder="out('Hello, Ink!')"
          autocomplete="off" autocorrect="off" autocapitalize="off" spellcheck="${false}"
          value="${this.stdin}"/>
        <div class="savedBadge ${this.showSaved ? '' : 'hidden'}">
          Saved to this device
        </div>
      </div>
      <div class="divider"></div>
      <div class="inkOutputBox">
        <div class="output">
          <code class="elapsedTime">elapsed: ${this.result.duration}s</code>
          ${this.waiting ?  'running...' :
            this.result.output.split('\n')
              .map(line => jdom`<code>${line}</code>`)}
        </div>
      </div>
    </div>`;
  }
}

class App extends StyledComponent {
  init() {
    this.showExamples = false;

    this.ioBox = new IOBox();

    this.insertExample = this.insertExample.bind(this);
    this.openExamples = this.openExamples.bind(this);
    this.closeExamples = this.closeExamples.bind(this);
  }

  insertExample(src) {
    this.ioBox.insert(src);
    this.closeExamples();
  }

  openExamples() {
    this.showExamples = true;
    this.render();

    // eck. cheap trick but it works
    document.querySelector('.exampleItem').focus();
  }

  closeExamples() {
    this.showExamples = false;
    this.render();
  }

  styles() {
    return css`
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100%;

    header {
      display: block;
      width: 100%;
      height: 36px;
      background: #222;
      color: #fff;
      padding: 8px 16px;
      justify-content: space-between;
      box-sizing: border-box;
      flex-grow: 0;

      &, .left, .right {
        display: flex;
        flex-direction: row;
        align-items: center;
      }
    }
    h1 {
      font-weight: normal;
      font-size: 1em;
      margin: 0;
    }
    .block {
      font-size: .9em;
      font-weight: normal;
      text-decoration: none;
      position: relative;
      top: 8px;
      padding: 3px 6px;
    }
    .backdrop {
      position: fixed;
      z-index: 10;
      background: rgba(0, 0, 0, .3);
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
    }
    .examples {
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: #f8f8f8;
      box-shadow: 0 3px 8px -1px rgba(0, 0, 0, .3);
      border-radius: 3px;
      overflow: hidden;

      min-width: 300px;
      max-width: 700px;
      width: 50%;

      h2 {
        font-weight: normal;
        font-size: 1em;
        margin: 0;
        display: block;
        width: 100%;
        height: 36px;
        background: #222;
        color: #fff;
        padding: 8px 16px;
        justify-content: space-between;
        box-sizing: border-box;
      }
      p {
        line-height: 1.5em;
      }
      a {
        color: var(--block-accent-color);
      }
      p, ul {
        margin: 12px;
      }
      ul {
        padding-left: 0;
        overflow-y: auto;
        min-height: 200px;
        max-height: 60vh;
        height: 50vh;
      }
      li {
        width: 100%;
      }
    }
    .exampleItem {
      cursor: pointer;
      list-style: none;
      padding: 12px;
      margin-bottom: 16px;
      text-align: left;
      width: calc(100% - 12px);
      &:focus {
        outline: initial;
      }
      .filename {
        font-weight: bold;
        margin-bottom: 6px;
      }
    }
    @media only screen and (max-width: 700px) {
      .examples {
        top: 52px;
        transform: translate(-50%, 0);
        min-width: 92%;
      }
    }
    `;
  }

  compose() {
    return jdom`<main>
      <header>
        <div class="left">
          <h1>Inker</h1>
        </div>
        <div class="right">
          <a href="https://github.com/thesephist/ink"
            rel="noopener noreferrer"
            target="_blank"
            class="block tiny-hidden"
            title="About the Ink programming language"
            >About</a>
          <button
            class="block"
            title="Show example snippets"
            onclick=${this.openExamples}>Examples</button>
          ${this.ioBox.buttons().children}
        </div>
      </header>
      ${this.ioBox.node}
      ${this.showExamples ?
        ExampleList(this.insertExample, this.closeExamples) : null}
    </main>`;
  }
}

const app = new App();
document.getElementById('root').appendChild(app.node);
