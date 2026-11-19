"use client";

import { useEffect, useState } from "react";
import { repo, rawUrl, releasesUrl } from "@/lib/site";
import { CopyCommand } from "./CopyCommand";
import { SegmentedControl } from "./SegmentedControl";

type OS = "macos" | "linux" | "windows";

const OS_OPTIONS = [
  { value: "macos", label: "macOS" },
  { value: "linux", label: "Linux" },
  { value: "windows", label: "Windows" },
] as const satisfies readonly { value: OS; label: string }[];

interface Method {
  title: string;
  command: string;
  prompt?: string;
  note?: string;
}

const goInstall: Method = {
  title: "With Go",
  command: "go install github.com/DalyChouikh/cmd/qrgen@latest",
  note: "Needs Go 1.25 or newer.",
};

const METHODS: Record<OS, Method[]> = {
  macos: [
    { title: "Install script", command: `curl -fsSL ${rawUrl}/install.sh | sh` },
    { title: "Homebrew", command: "brew install DalyChouikh/tap/qrgen" },
    goInstall,
  ],
  linux: [
    { title: "Install script", command: `curl -fsSL ${rawUrl}/install.sh | sh` },
    { title: "Homebrew", command: "brew install DalyChouikh/tap/qrgen" },
    goInstall,
  ],
  windows: [
    { title: "PowerShell", command: `irm ${rawUrl}/install.ps1 | iex`, prompt: ">", note: "Run in PowerShell." },
    goInstall,
  ],
};

// GoReleaser names archives qrgen_<version>_<os>_<arch>; the updater relies on the same scheme.
const GOOS: Record<OS, string> = { macos: "darwin", linux: "linux", windows: "windows" };
const ARCH_LABELS: Record<string, string> = { amd64: "x86-64", arm64: "ARM64" };

interface Asset {
  name: string;
  browser_download_url: string;
  size: number;
}

interface Release {
  tag_name: string;
  html_url: string;
  assets: Asset[];
}

function detectOS(): OS {
  const platform = navigator.userAgent.toLowerCase();
  if (platform.includes("windows")) return "windows";
  if (platform.includes("mac os") || platform.includes("macintosh")) return "macos";
  return "linux";
}

export function Install() {
  const [os, setOS] = useState<OS>("linux");
  const [release, setRelease] = useState<Release | null>(null);

  useEffect(() => {
    setOS(detectOS());
    const controller = new AbortController();
    fetch(`https://api.github.com/repos/${repo}/releases/latest`, { signal: controller.signal })
      .then((r) => (r.ok ? r.json() : null))
      .then(setRelease)
      .catch(() => setRelease(null));
    return () => controller.abort();
  }, []);

  const binaries =
    release?.assets.filter((a) => a.name.includes(`_${GOOS[os]}_`) && !a.name.endsWith(".txt")) ?? [];

  return (
    <section className="install" id="install" aria-labelledby="install-title">
      <div className="section-head">
        <h2 id="install-title">Get the terminal version</h2>
        <p>
          The CLI walks you through the same steps in an interactive wizard, previews the code right in
          your terminal and keeps your last 50 codes so you can regenerate them.
          {release && <> The latest release is <a href={release.html_url}>{release.tag_name}</a>.</>}
        </p>
      </div>

      <SegmentedControl
        label="Operating system"
        showLabel={false}
        variant="tabs"
        options={OS_OPTIONS}
        value={os}
        onChange={setOS}
      />

      <div className="install-panel">
        {METHODS[os].map((m) => (
          <div className="method" key={m.title}>
            <h3>{m.title}</h3>
            <CopyCommand command={m.command} prompt={m.prompt} />
            {m.note && <p className="meta">{m.note}</p>}
          </div>
        ))}

        <div className="method">
          <h3>Download a binary</h3>
          {binaries.length > 0 ? (
            <ul className="downloads">
              {binaries.map((a) => (
                <li key={a.name}>
                  <a href={a.browser_download_url}>{a.name}</a>
                  <span className="meta">
                    {ARCH_LABELS[archOf(a.name)] ?? archOf(a.name)}, {(a.size / 1024 / 1024).toFixed(1)} MB
                  </span>
                </li>
              ))}
            </ul>
          ) : (
            <p>
              Pick the archive for your system on the <a href={releasesUrl}>releases page</a>, extract it and
              move <code>qrgen</code> onto your PATH.
            </p>
          )}
        </div>
      </div>
    </section>
  );
}

function archOf(name: string) {
  return name.replace(/\.(tar\.gz|zip)$/, "").split("_").pop() ?? "";
}
