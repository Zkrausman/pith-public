package install

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetupPiHook installs Pith's global Pi extension. Project policy stays in the
// trusted project's .pi/pith.json; no consumer repository policy is embedded.
func SetupPiHook() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".pi", "agent", "extensions", "pith")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, "index.ts")
	if err := os.WriteFile(path, []byte(piExtension), 0644); err != nil {
		return err
	}
	fmt.Printf("Installed Pi hook at %s\n", path)
	return nil
}

const piExtension = `import { spawn } from "node:child_process";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

type Config = { enabled?: boolean; thresholdBytes?: number; telemetryEnabled?: boolean; rawBypass?: boolean };
async function config(cwd: string, trusted: boolean): Promise<Config> {
  if (!trusted) return {};
  try { return JSON.parse(await readFile(join(cwd, ".pi", "pith.json"), "utf8")); } catch { return {}; }
}
function transform(binary: string, request: unknown, signal?: AbortSignal): Promise<any> {
  return new Promise((resolve, reject) => {
    const child = spawn(binary, ["pi", "transform"], { stdio: ["pipe", "pipe", "ignore"] });
    let out = "";
    child.stdout.on("data", (chunk) => { out += String(chunk); });
    child.once("error", reject);
    child.once("close", (code) => code === 0 ? resolve(JSON.parse(out)) : reject(new Error("pith exited " + code)));
    signal?.addEventListener("abort", () => child.kill(), { once: true });
    child.stdin.end(JSON.stringify(request));
  });
}
export default function (pi: ExtensionAPI) {
  pi.on("tool_result", async (event, ctx) => {
    const command = (event.input as { command?: unknown }).command;
    if (event.toolName !== "bash" || typeof command !== "string" || event.isError) return;
    const blocks = event.content as Array<{ type: string; text?: string }>;
    if (!Array.isArray(blocks) || blocks.some((b) => b.type !== "text" || typeof b.text !== "string")) return;
    const output = blocks.map((b) => b.text!).join("\n");
    const cfg = await config(ctx.cwd, ctx.isProjectTrusted());
    if (cfg.enabled === false || cfg.rawBypass || output.length < (cfg.thresholdBytes ?? 8000)) return;
    try {
      const response = await transform(process.env.PITH_BIN || "pith", { command, output, exitCode: Number((event.details as any)?.exitCode ?? 0), thresholdBytes: cfg.thresholdBytes, telemetryEnabled: cfg.telemetryEnabled === true }, ctx.signal);
      if (typeof response?.output !== "string") return;
      return { content: [{ type: "text", text: response.output }], details: { ...event.details, pith: { parser: response.parser, passthrough: response.passthrough } } };
    } catch { return; } // Pith failure always preserves Pi's original result.
  });
}
`
