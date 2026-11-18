// Takes the cards that X, Slack and other sites show for a shared link:
// the home page's hero, in every language of the site and in light, at
// 1200 by 630, into src/assets/share/<lang>.jpg. It needs Google Chrome
// and macOS's sips.
//
// The page is served by `astro dev` and driven through Chrome's DevTools
// protocol; a style sheet trims it to the mark, the title, the lead and the
// platforms' board, in a still frame.
import { execFileSync, spawn } from 'node:child_process';
import { mkdirSync, mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { langs } from '../src/site.ts';

const here = path.dirname(fileURLToPath(import.meta.url));
const website = path.resolve(here, '..');
const outDir = path.join(website, 'src', 'assets', 'share');
const chrome = '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';
const sitePort = 4329;
const debugPort = 9340;
const width = 1200;
const height = 630;
const locales = { en: 'en-US', zh: 'zh-CN' };
const missing = langs.filter((lang) => !locales[lang]);
if (missing.length) throw new Error(`no locale for ${missing.join(', ')}: add it to scripts/share-cards.mjs`);

const card = `
  .skip-link, .header__links, .header__actions, .pill, .hero__actions,
  .platform-lines, .stage, main > section:not(.hero), .cta-wrap, footer { display: none !important; }
  html, body { overflow: hidden; }
  .header { position: absolute; inset: 0 0 auto; border: 0; background: none; backdrop-filter: none; }
  .header__inner { padding-top: 24px; }
  .container { width: calc(100% - 112px); }
  .header::after { content: 'nsgo.dev'; position: absolute; top: 42px; right: 56px; color: var(--muted); font-size: 18px; font-weight: 500; }
  .hero { display: flex; align-items: center; height: ${height}px; padding: 40px 0 0; }
  .hero__inner { gap: 56px; }
  .hero h1 { margin-top: 0; font-size: 50px; }
  .hero__lead { margin-top: 24px; font-size: 20px; }
`;

async function capture(url, locale, out) {
  const profile = mkdtempSync(path.join(tmpdir(), 'nsgo-card-'));
  const browser = spawn(
    chrome,
    ['--headless=new', `--lang=${locale}`, `--accept-lang=${locale}`, `--remote-debugging-port=${debugPort}`, `--user-data-dir=${profile}`, 'about:blank'],
    { stdio: 'ignore' },
  );
  try {
    let page;
    for (let i = 0; i < 50 && !page; i++) {
      try {
        const targets = await (await fetch(`http://127.0.0.1:${debugPort}/json`)).json();
        page = targets.find((target) => target.type === 'page');
      } catch {
        await Bun.sleep(200);
      }
    }
    const socket = new WebSocket(page.webSocketDebuggerUrl);
    await new Promise((resolve) => (socket.onopen = resolve));
    let id = 0;
    const pending = new Map();
    socket.onmessage = (event) => {
      const message = JSON.parse(String(event.data));
      const resolve = pending.get(message.id);
      if (typeof resolve === 'function') {
        pending.delete(message.id);
        resolve(message.result);
      }
    };
    const send = (method, params = {}) =>
      new Promise((resolve) => {
        pending.set(++id, resolve);
        socket.send(JSON.stringify({ id, method, params }));
      });

    await send('Emulation.setLocaleOverride', { locale });
    await send('Emulation.setDeviceMetricsOverride', { width, height, deviceScaleFactor: 2, mobile: false });
    await send('Emulation.setEmulatedMedia', { features: [{ name: 'prefers-reduced-motion', value: 'reduce' }] });
    await send('Page.enable');
    await send('Page.navigate', { url });
    await Bun.sleep(4000);
    await send('Runtime.evaluate', {
      expression: `(() => {
        const style = document.createElement('style');
        style.textContent = ${JSON.stringify(card)};
        document.head.append(style);
        window.dispatchEvent(new Event('resize'));
      })()`,
    });
    await Bun.sleep(1500);
    const { data } = await send('Page.captureScreenshot', { format: 'png' });
    const png = `${out}.png`;
    await Bun.write(png, Buffer.from(data, 'base64'));
    execFileSync('sips', ['-s', 'format', 'jpeg', '-s', 'formatOptions', '90', '-Z', String(width), png, '--out', out], { stdio: 'ignore' });
    rmSync(png);
    socket.close();
  } finally {
    const exited = new Promise((resolve) => browser.once('exit', resolve));
    browser.kill();
    await exited;
    rmSync(profile, { recursive: true, force: true });
  }
}

const server = spawn('bunx', ['astro', 'dev', '--port', String(sitePort), '--ignore-lock'], { cwd: website, stdio: 'ignore' });
try {
  let ready = false;
  for (let i = 0; i < 150 && !ready; i++) {
    try {
      ready = (await fetch(`http://localhost:${sitePort}/`)).ok;
    } catch {
      await Bun.sleep(200);
    }
  }
  if (!ready) throw new Error('astro dev did not start');
  mkdirSync(outDir, { recursive: true });
  for (const lang of langs) {
    const out = path.join(outDir, `${lang}.jpg`);
    await capture(`http://localhost:${sitePort}/${lang === 'en' ? '' : `${lang}/`}`, locales[lang], out);
    console.log(path.relative(website, out));
  }
} finally {
  server.kill();
}
