import { Generator } from "@/components/Generator";
import { Install } from "@/components/Install";
import { repoUrl } from "@/lib/site";

const COMMANDS = [
  { command: "qrgen", description: "Start the interactive wizard" },
  { command: "qrgen history", description: "List your last 50 codes" },
  { command: "qrgen regen <id>", description: "Make a code from history again" },
  { command: "qrgen update", description: "Update to the latest release" },
  { command: "qrgen check-update", description: "See whether a newer release exists" },
  { command: "qrgen --version", description: "Print version information" },
  { command: "qrgen --help", description: "Show help" },
];

const KEYS = [
  { keys: ["Enter"], action: "Confirm" },
  { keys: ["↑", "↓"], action: "Move through lists (or j, k)" },
  { keys: ["←", "→"], action: "Switch options (or h, l)" },
  { keys: ["Tab"], action: "Next field, or open the file browser on the output step" },
  { keys: ["Space"], action: "Select or toggle" },
  { keys: ["c"], action: "Type a custom hex color" },
  { keys: ["Esc"], action: "Go back a step" },
  { keys: ["r"], action: "Make another code when done" },
  { keys: ["Ctrl", "C"], action: "Quit" },
];

export default function Home() {
  return (
    <>
      <header className="site-header">
        <a className="wordmark" href="#generate" aria-label="qrgen home">
          <FinderMark />
          qrgen
        </a>
        <nav aria-label="Main">
          <a href="#install">Install</a>
          <a href="#cli">Commands</a>
          <a href={repoUrl}>GitHub</a>
        </nav>
      </header>

      <main>
        <div className="intro">
          <h1>QR codes for links, WiFi, contacts, email and texts.</h1>
          <p>
            Fill in the form and download a PNG or SVG. Prefer the terminal? The same generator ships as{" "}
            <a href="#install">a CLI</a>.
          </p>
        </div>

        <Generator />
        <Install />

        <section className="cli" id="cli" aria-labelledby="cli-title">
          <div className="section-head">
            <h2 id="cli-title">Using the CLI</h2>
            <p>Run <code>qrgen</code> and follow the prompts. A few extra commands cover history and updates.</p>
          </div>
          <div className="cli-grid">
            <table>
              <caption>Commands</caption>
              <tbody>
                {COMMANDS.map((c) => (
                  <tr key={c.command}>
                    <th scope="row"><code>{c.command}</code></th>
                    <td>{c.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
            <table>
              <caption>Keys in the wizard</caption>
              <tbody>
                {KEYS.map((k) => (
                  <tr key={k.action}>
                    <th scope="row">
                      {k.keys.map((key) => <kbd key={key}>{key}</kbd>)}
                    </th>
                    <td>{k.action}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="meta">
            On Linux and macOS, if qrgen lives in <code>/usr/local/bin</code>, run <code>sudo qrgen update</code>.
            On Windows, update from an Administrator terminal.
          </p>
        </section>
      </main>

      <footer className="site-footer">
        <p>
          qrgen is open source under the MIT license. <a href={repoUrl}>Source, issues and releases on GitHub</a>.
        </p>
      </footer>
    </>
  );
}

function FinderMark() {
  return (
    <svg viewBox="0 0 7 7" width="22" height="22" aria-hidden="true">
      <path fillRule="evenodd" d="M0 0h7v7H0zM1 1v5h5V1z" />
      <rect x="2" y="2" width="3" height="3" />
    </svg>
  );
}
