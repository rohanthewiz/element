// highlight.js — tiny, dependency-free syntax highlighters for the element
// playground: Go source and generated HTML (with CSS inside <style> blocks).
//
//   eleHi.go(src)            -> HTML string (token <span>s)
//   eleHi.html(out, swatch)  -> HTML string (optional color swatches in CSS)
//   eleHi.escape(s)          -> HTML-escaped string
//   eleHi.editor(ta, code)   -> wires a <textarea> to its overlay <code>;
//                               returns a repaint function
//
// Invariant shared by every tokenizer: the text content of the returned HTML
// is character-identical to the input, so the overlay editor's textarea and
// highlight layer stay column-aligned. Color swatches (an extra inline <i>)
// are therefore only used in the output pane, never in an editor.
(() => {
'use strict';

const ESC = { '&': '&amp;', '<': '&lt;', '>': '&gt;' };
const esc = s => s.replace(/[&<>]/g, c => ESC[c]);
const span = (cls, s) => s ? '<span class="' + cls + '">' + esc(s) + '</span>' : '';

// --- Go ----------------------------------------------------------------------

const GO_KW = /^(?:break|case|chan|const|continue|default|defer|else|fallthrough|for|func|go|goto|if|import|interface|map|package|range|return|select|struct|switch|type|var)$/;
const GO_LIT = /^(?:true|false|nil|iota)$/;
const GO_TYPE = /^(?:any|bool|byte|comparable|complex64|complex128|error|float32|float64|int|int8|int16|int32|int64|rune|string|uint|uint8|uint16|uint32|uint64|uintptr)$/;
const GO_BUILTIN = /^(?:append|cap|clear|close|complex|copy|delete|imag|len|make|max|min|new|panic|print|println|real|recover)$/;
const GO_IDENT = /^[A-Za-z_][A-Za-z0-9_]*/;

// Scans one chunk of Go. `st.block` carries an open /* */ across lines;
// `st.raw` an open backquoted string.
function goChunk(text, st) {
  let out = '', i = 0;
  let prevWord = '';
  const n = text.length;
  while (i < n) {
    const rest = text.slice(i);
    let m;
    if (st.block) {
      const end = rest.indexOf('*/');
      if (end < 0) { out += span('t-com', rest); break; }
      out += span('t-com', rest.slice(0, end + 2));
      st.block = false; i += end + 2; continue;
    }
    if (st.raw) {
      const end = rest.indexOf('`');
      if (end < 0) { out += span('t-str', rest); break; }
      out += span('t-str', rest.slice(0, end + 1));
      st.raw = false; i += end + 1; continue;
    }
    if (rest.startsWith('/*')) { st.block = true; continue; }
    if (rest.startsWith('//')) { out += span('t-com', rest); break; }
    if (rest[0] === '`') { st.raw = true; out += span('t-str', '`'); i++; continue; }
    if ((m = /^"(?:\\.|[^"\\])*"?/.exec(rest))) { out += span('t-str', m[0]); i += m[0].length; prevWord = ''; continue; }
    if ((m = /^'(?:\\.|[^'\\])*'?/.exec(rest))) { out += span('t-str', m[0]); i += m[0].length; prevWord = ''; continue; }
    if ((m = /^(?:0[xXbBoO][0-9a-fA-F_]+|(?:\d[\d_]*)(?:\.[\d_]*)?(?:[eE][+-]?\d+)?i?)/.exec(rest))) {
      out += span('t-num', m[0]); i += m[0].length; prevWord = ''; continue;
    }
    if ((m = GO_IDENT.exec(rest))) {
      const w = m[0];
      const afterDot = text[i - 1] === '.';
      let cls;
      if (!afterDot && GO_KW.test(w)) cls = 't-kw';
      else if (!afterDot && GO_LIT.test(w)) cls = 't-num';
      else if (!afterDot && GO_TYPE.test(w)) cls = 't-cls';
      else if (rest[w.length] === '(')
        cls = !afterDot && GO_BUILTIN.test(w) ? 't-kw' : 't-fn';
      else if (prevWord === 'func' || prevWord === 'type') cls = 't-fn';
      else if (!afterDot && text[i + w.length] === '.') cls = 't-var'; // package qualifier
      else cls = 't-id';
      out += span(cls, w);
      prevWord = w; i += w.length; continue;
    }
    if ((m = /^\s+/.exec(rest))) { out += m[0]; i += m[0].length; continue; }
    out += span('t-op', rest[0]); i++; prevWord = '';
  }
  return out;
}

function go(src) {
  const st = { block: false, raw: false };
  return src.split('\n').map(line => goChunk(line, st)).join('\n');
}

// --- CSS value scanner (for <style> bodies; ported from go-styl) -------------

const CSS_KW = /^(?:and|or|not|only|from|to|important)$/;

function cssValue(text, st, swatch) {
  let out = '', i = 0;
  const n = text.length;
  while (i < n) {
    const rest = text.slice(i);
    let m;
    if (st.comment) {
      const end = rest.indexOf('*/');
      if (end < 0) { out += span('t-com', rest); break; }
      out += span('t-com', rest.slice(0, end + 2));
      st.comment = false; i += end + 2; continue;
    }
    if (rest.startsWith('/*')) { st.comment = true; continue; }
    if ((m = /^(['"])(?:\\.|(?!\1).)*\1?/.exec(rest))) { out += span('t-str', m[0]); i += m[0].length; continue; }
    if ((m = /^#[0-9a-fA-F]{3,8}\b/.exec(rest))) {
      out += swatch
        ? '<span class="t-col"><i style="background:' + m[0] + '"></i>' + m[0] + '</span>'
        : '<span class="t-col" style="text-decoration-color:' + m[0] + '">' + m[0] + '</span>';
      i += m[0].length; continue;
    }
    if ((m = /^(\d*\.)?\d+[a-zA-Z%]*/.exec(rest))) { out += span('t-num', m[0]); i += m[0].length; continue; }
    if ((m = /^![a-zA-Z]+/.exec(rest))) { out += span('t-kw', m[0]); i += m[0].length; continue; }
    if ((m = /^-{0,2}[A-Za-z_$][-\w$]*/.exec(rest))) {
      const w = m[0];
      out += span(rest[w.length] === '(' ? 't-fn' : CSS_KW.test(w) ? 't-kw' : 't-id', w);
      i += w.length; continue;
    }
    if ((m = /^\s+/.exec(rest))) { out += m[0]; i += m[0].length; continue; }
    out += span('t-op', rest[0]); i++;
  }
  return out;
}

function cssSelector(text, st) {
  let out = '', i = 0;
  const n = text.length;
  while (i < n) {
    const rest = text.slice(i);
    let m;
    if (st.comment || rest.startsWith('/*')) return out + cssValue(rest, st);
    if ((m = /^(['"])(?:\\.|(?!\1).)*\1?/.exec(rest))) { out += span('t-str', m[0]); i += m[0].length; continue; }
    if ((m = /^[.#][-\w]+/.exec(rest))) { out += span('t-cls', m[0]); i += m[0].length; continue; }
    if ((m = /^::?[-\w]+/.exec(rest))) { out += span('t-pse', m[0]); i += m[0].length; continue; }
    if ((m = /^(\d*\.)?\d+%?/.exec(rest))) { out += span('t-num', m[0]); i += m[0].length; continue; }
    if ((m = /^[-\w]+/.exec(rest))) { out += span('t-sel', m[0]); i += m[0].length; continue; }
    if ((m = /^\s+/.exec(rest))) { out += m[0]; i += m[0].length; continue; }
    out += span('t-op', rest[0]); i++;
  }
  return out;
}

// css tokenizes a stylesheet (a <style> body). Context stack: 'sel' expects
// selectors / at-rules; 'body' expects `prop: value;` declarations.
function css(text, swatch) {
  let out = '', i = 0;
  const stack = ['sel'];
  const st = { comment: false };
  while (i < text.length) {
    const rest = text.slice(i);
    let m;
    if ((m = /^\s+/.exec(rest))) { out += m[0]; i += m[0].length; continue; }
    if ((m = /^\/\*[^]*?(?:\*\/|$)/.exec(rest))) { out += span('t-com', m[0]); i += m[0].length; continue; }
    if (rest[0] === '}') {
      if (stack.length > 1) stack.pop();
      out += span('t-op', '}'); i++; continue;
    }
    if (stack[stack.length - 1] === 'sel') {
      if (rest[0] === '@') {
        m = /^@[-\w]+/.exec(rest);
        out += span('t-at', m[0]); i += m[0].length;
        const p = /^[^{;]*/.exec(text.slice(i))[0];
        out += cssValue(p, st, swatch); i += p.length;
        if (text[i] === '{') {
          stack.push(/^@(media|supports|keyframes|-[-\w]+-keyframes|document|layer|container)$/.test(m[0]) ? 'sel' : 'body');
          out += span('t-op', '{'); i++;
        } else if (text[i] === ';') { out += span('t-op', ';'); i++; }
        continue;
      }
      m = /^[^{}@]+/.exec(rest);
      if (m) {
        out += cssSelector(m[0], st); i += m[0].length;
        if (text[i] === '{') { stack.push('body'); out += span('t-op', '{'); i++; }
        continue;
      }
      if (rest[0] === '{') { stack.push('body'); out += span('t-op', '{'); i++; continue; }
      out += span('t-op', rest[0]); i++; continue;
    }
    if ((m = /^(--[\w-]+|\*?[-\w$]+)(\s*)(:)/.exec(rest))) {
      out += span('t-prop', m[1]) + m[2] + span('t-op', ':');
      i += m[0].length;
      const v = /^[^;}]*/.exec(text.slice(i))[0];
      out += cssValue(v, st, swatch); i += v.length;
      if (text[i] === ';') { out += span('t-op', ';'); i++; }
      continue;
    }
    out += span('t-op', rest[0]); i++;
  }
  return out;
}

// --- HTML ----------------------------------------------------------------------
// Tokenizes serialized HTML (the playground's output pane). <style> bodies are
// handed to the css tokenizer; <script> bodies stay plain.
function html(text, swatch) {
  let out = '', i = 0;
  const n = text.length;
  while (i < n) {
    const rest = text.slice(i);
    let m;
    if ((m = /^<!--[^]*?(?:-->|$)/.exec(rest))) { out += span('t-com', m[0]); i += m[0].length; continue; }
    if ((m = /^<!(?:DOCTYPE|doctype)[^>]*>?/.exec(rest))) { out += span('t-at', m[0]); i += m[0].length; continue; }
    if ((m = /^(<\/?)([A-Za-z][-\w]*)/.exec(rest))) {
      const closing = m[1] === '</';
      const tag = m[2].toLowerCase();
      out += span('t-op', m[1]) + span('t-sel', m[2]);
      i += m[0].length;
      // attributes until '>'
      while (i < n && text[i] !== '>') {
        const arest = text.slice(i);
        let am;
        if ((am = /^\s+/.exec(arest))) { out += am[0]; i += am[0].length; continue; }
        if (arest[0] === '/') { out += span('t-op', '/'); i++; continue; }
        if ((am = /^(['"])(?:(?!\1).)*\1?/.exec(arest))) { out += span('t-str', am[0]); i += am[0].length; continue; }
        if ((am = /^[^\s=>'"\/]+/.exec(arest))) { out += span('t-prop', am[0]); i += am[0].length; continue; }
        if (arest[0] === '=') { out += span('t-op', '='); i++; continue; }
        out += span('t-op', arest[0]); i++;
      }
      if (text[i] === '>') { out += span('t-op', '>'); i++; }
      if (!closing && (tag === 'style' || tag === 'script')) {
        const close = '</' + tag;
        let end = text.toLowerCase().indexOf(close, i);
        if (end < 0) end = n;
        const body = text.slice(i, end);
        out += tag === 'style' ? css(body, swatch) : span('t-id', body);
        i = end;
      }
      continue;
    }
    if ((m = /^&[#\w]+;/.exec(rest))) { out += span('t-num', m[0]); i += m[0].length; continue; }
    if ((m = /^[^<&]+/.exec(rest))) { out += esc(m[0]); i += m[0].length; continue; }
    out += span('t-op', rest[0]); i++;
  }
  return out;
}

// --- overlay editor wiring ---------------------------------------------------
// The <textarea> sits on top with transparent text (its wrapper's CSS handles
// that); the highlighted copy lives in `code` inside an overflow-hidden <pre>
// behind it. `enabled()` lets the host toggle highlighting off cheaply.
function editor(ta, code, enabled) {
  const pre = code.parentElement;
  const paint = () => {
    code.innerHTML = (enabled ? enabled() : true) ? go(ta.value) + '\n' : esc(ta.value) + '\n';
  };
  const sync = () => { pre.scrollTop = ta.scrollTop; pre.scrollLeft = ta.scrollLeft; };
  ta.addEventListener('input', paint);
  ta.addEventListener('scroll', sync);
  paint();
  return () => { paint(); sync(); };
}

const api = { go, html, css, escape: esc, editor };
if (typeof window !== 'undefined') window.eleHi = api;
if (typeof module !== 'undefined') module.exports = api;
})();
