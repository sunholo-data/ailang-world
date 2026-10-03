#!/usr/bin/env node
// Generates website/docs/agents/tools.md from packages/se-tools/transitions.json.
//
// Run from the repo root:   node website/scripts/gen-tool-reference.mjs
// Check without writing:    node website/scripts/gen-tool-reference.mjs --check
//
// Every tool id, title, description, argument name/type/requirement, effect,
// scope and cost on the page comes from transitions.json, so the reference
// cannot drift from the published manifest. The NOTES table below adds what
// the manifest does not carry: the plan-phase refusal texts and the finish
// behaviour, transcribed from packages/se-tools/se_tools/*.ail (each module's
// named tests pin those exact strings). Re-run after either source changes;
// --check exits 1 when the committed page is stale.

import { readFileSync, writeFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const repo = resolve(here, "..", "..");
const manifestPath = join(repo, "packages", "se-tools", "transitions.json");
const outPath = join(repo, "website", "docs", "agents", "tools.md");

const tools = JSON.parse(readFileSync(manifestPath, "utf8"));

// MDX-safe inline text: escape the characters MDX would treat as JSX or
// expressions, and the table delimiter.
const mdx = (s) =>
  String(s)
    .replace(/\\/g, "\\\\")
    .replace(/\|/g, "\\|")
    .replace(/\{/g, "\\{")
    .replace(/\}/g, "\\}")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

const PATH_RULE =
  'path "<p>" refused: it must be relative and non-empty, with no ".." segment and no leading "-"';

// Hand-maintained per-tool notes, sourced from the .ail modules (plan phase,
// finish phase) and the handler (host/broker/handlers_ailang.go).
const NOTES = {
  "ailang-read": {
    handlerOp: "policy-tool `read`",
    finish: false,
    example: { path: "hello.ail" },
    result: "`ok`, `content`, `tool` (the tool binary ref), `policy_digest`, `world`.",
    refusals: [PATH_RULE, 'missing required argument "path"', 'argument "path" must be a string'],
    policy: [
      'read /etc/hosts: path "/etc/hosts" escapes sandbox "<worktree>"',
      "read nope.ail: openat nope.ail: no such file or directory",
    ],
  },
  "ailang-write": {
    handlerOp: "policy-tool `write`",
    finish: false,
    example: { path: "notes/todo.txt", content: "first line\n" },
    result: "`ok`, `tool`, `policy_digest`, `world`.",
    refusals: [PATH_RULE, 'missing required argument "content"', 'unknown argument "old_text"; ailang-write admits only: path, content'],
    policy: [
      "write .git/config: .git/ is read-only to the lane's tools",
      'write .ailang/x.ail: matches fs_deny_write ".ailang/**" — read-only under this policy',
    ],
  },
  "ailang-edit": {
    handlerOp: "policy-tool `edit`",
    finish: false,
    example: {
      path: "hello.ail",
      old_text: 'println("hello from ailang-run")',
      new_text: 'println("hello from pi")',
    },
    result: "`ok`, `tool`, `policy_digest`, `world`.",
    refusals: [PATH_RULE, 'missing required argument "new_text"', 'unknown argument "old"; ailang-edit admits only: path, old_text, new_text'],
    policy: [
      "edit hello.ail: old_text occurs 8 times; include more context so it is unique",
      "edit hello.ail: old_text not found — the file may have changed; read it again",
    ],
  },
  "ailang-check": {
    handlerOp: "policy-tool `ai_check`",
    finish: true,
    finishText:
      "The finish phase parses `ai_check`'s JSON report into `ok: true`, `passed`, `error_count` and `errors` (each `code`, `message`, `file`; line and column stay inside `message`). `ok: true` means the check ran; `passed` is the verdict. When no report is produced the result is `ok: false` with `status` and `refused` (for example `the Ailang.Check effect ended denied`).",
    example: { path: "hello.ail" },
    result: "`ok`, `passed`, `error_count`, `errors`, `world` (the finish output; no `tool` or `policy_digest`).",
    refusals: [PATH_RULE, 'unknown argument "limit"; ailang-check admits only: path'],
    policy: [],
  },
  "ailang-run": {
    handlerOp: "`ailang run --policy <policy> [--args-json J] -- <path>` (cwd = the worktree root)",
    finish: false,
    example: { path: "hello.ail" },
    result:
      "`admitted`, `exit_code`, `decision`, `limit`, `stdout`, `stderr`, `world`. See [Results and errors](./results-and-errors.md#ailang-run-outcomes) for the four outcome shapes.",
    refusals: [
      PATH_RULE,
      'argument "args_json" must be a string holding JSON',
      'argument "args_json" must be a string',
      'unknown argument "caps"; ailang-run admits only: path, args_json',
    ],
    policy: [],
  },
  "builtins-search": {
    handlerOp: "policy-tool `builtins_list` (text form, no flags)",
    finish: true,
    finishText:
      "The finish phase parses the `Total: N builtins` header and one `name [effect] module` line per builtin, then keeps entries where `module` contains the `module` argument and `query` is a substring of name, module or effect (both case-insensitive). With a `query`, at most 10 matches are returned; `count` is the number returned, not the inventory total. A truncated or unparseable inventory is `ok: false` with `status` and `refused`, never an empty list. Matches carry no signature or description: use `ailang-cli` op `builtins_show` for those.",
    example: { query: "println" },
    result: "`count`, `matches` (each `name`, `module`, `effect`), `world`.",
    refusals: ['unknown argument "limit"; builtins-search admits only: query, module', 'argument "query" must be a string'],
    policy: [],
  },
  "examples-search": {
    handlerOp: "policy-tool `examples_search` over the corpus named by `serve --examples-dir`",
    finish: false,
    example: { query: "fold" },
    result: "`ok`, `argv`, `exit_code`, `stdout`, `stderr`, `tool`, `policy_digest`, `world`. Read a hit with `ailang-cli` op `examples_show`.",
    refusals: ['missing required argument "query"', 'unknown argument "limit"; examples-search admits only: query'],
    policy: [
      "no examples corpus configured: start ailang-worldd serve with --examples-dir DIR (an AILANG examples corpus: manifest.json plus runnable/; `ailang examples download` makes one at ~/.ailang/examples)",
    ],
  },
  "ailang-cli": {
    handlerOp: "policy-tool, `op` from the binary's `cli` list",
    finish: false,
    example: { op: "iface", module: "std/io", flags: { compact: "" } },
    result: "`ok`, `argv`, `exit_code`, `stdout`, `stderr`, `tool`, `policy_digest`, `world`.",
    refusals: [
      'op "read" is not an ailang-cli op; reads, writes and edits have their own tools, and run is ailang-run',
      'op "fmt" flag "write" writes the worktree; ailang-cli never writes (fmt without it returns the formatted text; write it with ailang-write)',
      'argument "flags" must be an object of string values',
      'missing required argument "op"',
      PATH_RULE,
    ],
    policy: ["op check does not admit flag --zz (admitted: --json --quiet --strict-syntax)"],
  },
};

// Flags each ailang-cli op admits on the v0.52.1 tool binary, from the V67/V78
// audit table in host/broker/handlers_ailang_cliwrite_test.go
// (measuredCLIFlags). policy-tool enforces it; fmt's --write is refused by World.
const CLI_FLAGS = {
  agent_prompt: "none", ai_check: "--timeout", axioms: "none", builtins_list: "--by-effect --by-module --json --module --query --verbose",
  builtins_show: "none", check: "--json --quiet --strict-syntax", devtools_prompt: "none",
  docs_search: "--json --limit", examples_list: "--status --tags", examples_search: "none",
  examples_show: "none", examples_tags: "none", fmt: "--check (--write is refused)", iface: "--compact",
  pkg_docs: "none", policy_check: "none", prompt: "none",
  test: "--allow-skips --json --no-color --package", tree: "none", version: "none",
};

const lines = [];
const out = (s = "") => lines.push(s);

out("---");
out("title: Tool reference");
out("sidebar_position: 3");
out("description: The eight AILANG World software-engineering tools, generated from packages/se-tools/transitions.json.");
out("---");
out();
out("# Tool reference");
out();
out(
  `AILANG World serves ${tools.length} tools. Each one is a pure AILANG transition that plans exactly one brokered effect, ` +
    "runs it inside your episode's worktree, and commits one World log entry. The MCP tool name and the A2A skill id are both the tool id below.",
);
out();
out("Rules that apply to every tool:");
out();
out("- **Arguments are strict.** A key outside the schema is refused by the plan, never ignored. The refusal commits a log entry with zero effects.");
out("- **Paths are relative to the worktree root**: non-empty, no `..` segment, no leading `-`. AILANG's policy layer refuses anything that resolves outside the worktree (absolute paths, symlinks leading out), keeps `.git/` read-only, and keeps the operator's deny list read-only.");
out("- **One call costs 1** from the session's budget for the tool's effect. Budgets count calls, not time or tokens.");
out("- **Every result carries `world`**: `{effects:[{id, status, record}], plan}`. See [Results and errors](./results-and-errors.md).");
out();
out("| Tool | Effect (grant needed) | Arguments | Finish phase |");
out("|---|---|---|---|");
for (const t of tools) {
  const props = Object.keys(t.inputSchema.properties || {});
  const req = new Set(t.inputSchema.required || []);
  const args = props.map((p) => (req.has(p) ? `\`${p}\`` : `\`${p}?\``)).join(", ") || "none";
  out(`| [\`${t.id}\`](#${t.id}) | \`${t.access.effect}\` (scope \`${t.access.scope}\`) | ${args} | ${NOTES[t.id]?.finish ? "yes" : "no"} |`);
}
out();
out(
  "Six grants cover the eight tools: `ailang-write` and `ailang-edit` share `Workspace.Write`, and the two searches share `Ailang.Discover`. " +
    "`tools/list` shows only the tools whose effect your session holds a grant for. A grant with budget 0 still lists the tool, but every call is `denied:budget`.",
);
out();

for (const t of tools) {
  const n = NOTES[t.id];
  if (!n) throw new Error(`no NOTES entry for ${t.id}; add one before regenerating`);
  out(`## ${t.id}`);
  out();
  out(`**${mdx(t.title)}.** Effect \`${t.declaredEffects.map((e) => e.effect).join(", ")}\`, scope \`${t.access.scope}\`, cost ${t.declaredEffects[0].cost} per call. Handler: ${n.handlerOp}. Module: \`${t.transitionFnFile}\`.`);
  out();
  out("Description served in `tools/list`:");
  out();
  out("```text");
  out(t.description);
  out("```");
  out();
  out("### Arguments");
  out();
  const props = t.inputSchema.properties || {};
  const req = new Set(t.inputSchema.required || []);
  if (Object.keys(props).length === 0) {
    out("None.");
  } else {
    out("| Name | Type | Required | Description |");
    out("|---|---|---|---|");
    for (const [name, s] of Object.entries(props)) {
      let type = s.type;
      if (s.type === "object" && s.additionalProperties && s.additionalProperties.type) {
        type = `object of ${s.additionalProperties.type}`;
      }
      out(`| \`${name}\` | ${type} | ${req.has(name) ? "yes" : "no"} | ${mdx(s.description || "")} |`);
    }
  }
  out();
  out(`\`additionalProperties\` is \`${t.inputSchema.additionalProperties}\`: any other key is refused.`);
  out();
  if (t.id === "ailang-cli") {
    out("Admitted ops and their flags (v0.52.1 tool binary). Pass flags as `{\"name\": \"value\"}`; boolean flags take `\"\"`.");
    out();
    out("| `op` | Fields | Admitted flags |");
    out("|---|---|---|");
    const fields = {
      check: "`path`", ai_check: "`path`", fmt: "`path`", tree: "`path`", policy_check: "`path`",
      iface: "`module`", pkg_docs: "`module`", docs_search: "`query`", examples_search: "`query`",
      examples_show: "`module` (a name)", builtins_show: "`module` (a name)", test: "`path`, `package` (both optional)",
    };
    for (const [op, flags] of Object.entries(CLI_FLAGS)) {
      out(`| \`${op}\` | ${fields[op] || "none"} | ${mdx(flags)} |`);
    }
    out();
  }
  out("### Result");
  out();
  out(n.result);
  out();
  if (n.finish) {
    out(n.finishText);
    out();
  }
  const outProps = Object.keys(t.outputSchema.properties || {});
  out(`Output schema properties: ${outProps.map((p) => `\`${p}\``).join(", ")} (required: ${(t.outputSchema.required || []).map((p) => `\`${p}\``).join(", ")}).`);
  out();
  out("### Refusals you will see");
  out();
  out("From the plan (zero effects, `world.effects` is `[]`, still committed):");
  out();
  out("```text");
  for (const r of n.refusals) out(r);
  out("```");
  out();
  if (n.policy.length) {
    out("From the handler or AILANG's policy layer (the effect ran with status `ok`; the output says `ok: false`):");
    out();
    out("```text");
    for (const r of n.policy) out(r);
    out("```");
    out();
  }
  out("### Example call");
  out();
  out("```json");
  out(
    JSON.stringify({
      jsonrpc: "2.0",
      id: 1,
      method: "tools/call",
      params: { name: t.id, arguments: n.example },
    }),
  );
  out("```");
  out();
}

out("---");
out();
out(
  "This page is generated by `website/scripts/gen-tool-reference.mjs` from `packages/se-tools/transitions.json` " +
    "(refusal texts transcribed from `packages/se-tools/se_tools/*.ail`). Do not edit it by hand: change the sources and run " +
    "`node website/scripts/gen-tool-reference.mjs` from the repo root.",
);

const page = lines.join("\n") + "\n";
if (process.argv.includes("--check")) {
  const current = readFileSync(outPath, "utf8");
  if (current !== page) {
    console.error(`${outPath} is stale; run: node website/scripts/gen-tool-reference.mjs`);
    process.exit(1);
  }
  console.log("tools.md is up to date");
} else {
  writeFileSync(outPath, page);
  console.log(`wrote ${outPath} (${tools.length} tools)`);
}
