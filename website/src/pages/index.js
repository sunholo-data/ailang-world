import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import {
  ArrowRight,
  Bot,
  BookOpen,
  CheckCircle2,
  FileCode2,
  GitBranch,
  History,
  KeyRound,
  Network,
  Play,
  ShieldCheck,
  Waypoints,
} from 'lucide-react';
import styles from './index.module.css';

const GITHUB_URL = 'https://github.com/sunholo-data/ailang-world';

function GitHubMark({size = 18}) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" aria-hidden="true" fill="currentColor">
      <path d="M12 .5C5.65.5.5 5.65.5 12a11.5 11.5 0 0 0 7.86 10.92c.58.1.79-.25.79-.56v-2c-3.2.7-3.87-1.37-3.87-1.37-.53-1.33-1.28-1.69-1.28-1.69-1.05-.72.08-.7.08-.7 1.16.08 1.77 1.19 1.77 1.19 1.03 1.77 2.71 1.26 3.37.96.1-.75.4-1.26.73-1.55-2.56-.29-5.25-1.28-5.25-5.7 0-1.26.45-2.29 1.19-3.1-.12-.29-.52-1.47.11-3.06 0 0 .97-.31 3.17 1.18a11 11 0 0 1 5.77 0c2.2-1.49 3.17-1.18 3.17-1.18.63 1.59.23 2.77.11 3.06.74.81 1.19 1.84 1.19 3.1 0 4.43-2.7 5.4-5.26 5.69.41.36.78 1.06.78 2.14v3.17c0 .31.21.67.8.56A11.5 11.5 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5Z" />
    </svg>
  );
}

function Hero() {
  const {siteConfig} = useDocusaurusContext();
  const logo = useBaseUrl('/img/ailang-logo.svg');
  return (
    <header className={styles.hero}>
      <div className={styles.heroBg} aria-hidden="true">
        <div className={clsx(styles.orb, styles.orb1)} />
        <div className={clsx(styles.orb, styles.orb2)} />
        <div className={styles.grid} />
      </div>
      <div className={styles.heroContent}>
        <div className={styles.lockup}>
          <img src={logo} alt="AILANG logo" width="96" height="96" className={styles.heroLogo} />
        </div>
        <h1 className={styles.heroTitle}>
          AILANG <span className={styles.gradient}>World</span>
        </h1>
        <p className={styles.heroLead}>
          Unix made everything a file. World makes everything a{' '}
          <strong>typed state transition</strong>.
        </p>
        <p className={styles.heroSubtitle}>{siteConfig.tagline}</p>
        <div className={styles.heroActions}>
          <Link to="/docs/getting-started" className={clsx(styles.btn, styles.btnPrimary)}>
            <Play size={18} /> Get started
          </Link>
          <Link to="/docs/agents" className={clsx(styles.btn, styles.btnAccent)}>
            <Bot size={18} /> For AI agents
          </Link>
          <a href={GITHUB_URL} className={clsx(styles.btn, styles.btnSecondary)}>
            <GitHubMark /> GitHub
          </a>
        </div>
      </div>
    </header>
  );
}

function StatusBanner() {
  return (
    <section className={styles.statusWrap} aria-label="Project status">
      <Link to="/docs/roadmap" className={styles.status}>
        <span className={styles.statusPill}>Pre-1.0</span>
        <span className={styles.statusText}>
          <strong>5 of 7 release clauses met.</strong> The two left measure World against the
          status quo: a non-inferiority floor for resident agents, and provenance answers timed
          against log archaeology.
        </span>
        <span className={styles.statusLink}>
          See the roadmap <ArrowRight size={16} />
        </span>
      </Link>
    </section>
  );
}

const STEPS = [
  {
    icon: FileCode2,
    title: 'Propose',
    text: 'An agent proposes a change as an AILANG program. It never writes to the world directly.',
  },
  {
    icon: ShieldCheck,
    title: 'Verify',
    text: 'Types and Z3 contracts are checked before anything runs. Unverified proposals stop here.',
  },
  {
    icon: CheckCircle2,
    title: 'Commit',
    text: 'Verified, authorized, budgeted proposals commit to an append-only, content-addressed log.',
  },
  {
    icon: History,
    title: 'Replay',
    text: 'A pure transition plus its recorded effect results rebuilds any past state bit-for-bit.',
  },
];

function HowItWorks() {
  return (
    <section className={styles.section}>
      <div className={styles.container}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>How it works</h2>
          <p className={styles.sectionSubtitle}>
            The kernel is an immutable world graph. State changes only through transitions:
            pure AILANG functions, checked before they run. Every effect passes through a
            capability-checked broker, and MCP and A2A are the native boundary.
          </p>
        </div>
        <ol className={styles.steps}>
          {STEPS.map(({icon: Icon, title, text}, i) => (
            <li key={title} className={styles.step}>
              <div className={styles.stepHead}>
                <span className={styles.stepNum}>{i + 1}</span>
                <Icon size={22} className={styles.stepIcon} aria-hidden="true" />
              </div>
              <h3 className={styles.stepTitle}>{title}</h3>
              <p className={styles.stepText}>{text}</p>
            </li>
          ))}
        </ol>
      </div>
    </section>
  );
}

const FEATURES = [
  {
    icon: Waypoints,
    title: 'Provenance you can walk',
    text: 'Every commit leaves evidence and a trace. Ask the world why something happened and walk the log to the answer, instead of grepping for it.',
    to: '/docs/concepts',
  },
  {
    icon: KeyRound,
    title: 'Explicit authority',
    text: 'No ambient authority. Every effect goes through the broker with a capability and a budget check, and every allowed, denied or failed effect is recorded.',
    to: '/docs/security',
  },
  {
    icon: GitBranch,
    title: 'Deterministic replay',
    text: 'Determinism comes from the language, not from reconstructed logs. Recorded effect results are replay input, and replay never dispatches a live handler.',
    to: '/docs/concepts',
  },
  {
    icon: Bot,
    title: 'Agents as residents',
    text: 'Eight AILANG coding tools (read, write, edit, check, run, two searches and the CLI) are served over MCP and A2A. Each call runs one brokered effect and commits one log entry.',
    to: '/docs/agents',
  },
];

function Features() {
  return (
    <section className={clsx(styles.section, styles.sectionAlt)}>
      <div className={styles.container}>
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>What World gives you</h2>
          <p className={styles.sectionSubtitle}>
            World manages goals, typed state, capabilities, effects, evidence, budgets,
            contracts, provenance and AI proposals. It runs on one machine, over SQLite,
            with no cloud in the core.
          </p>
        </div>
        <div className={styles.features}>
          {FEATURES.map(({icon: Icon, title, text, to}) => (
            <Link key={title} to={to} className={styles.feature}>
              <div className={styles.featureIcon}>
                <Icon size={24} aria-hidden="true" />
              </div>
              <h3 className={styles.featureTitle}>{title}</h3>
              <p className={styles.featureText}>{text}</p>
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}

// Copied from docs/QUICKSTART.md §9 (software-engineering tools). Do not
// invent flags here: that section's flags are bound to the CLI by tests.
const SAMPLE = `# List the tools this session may call (MCP over the local daemon)
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \\
  -H 'Content-Type: application/json' \\
  -H 'Accept: application/json, text/event-stream' \\
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' \\
  http://127.0.0.1:7644/mcp/

# Call one: a brokered effect, one committed log entry
curl -s -H "Authorization: Bearer $(cat /tmp/se-session)" \\
  -H 'Content-Type: application/json' \\
  -H 'Accept: application/json, text/event-stream' \\
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/call",
       "params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}' \\
  http://127.0.0.1:7644/mcp/`;

function Sample() {
  return (
    <section className={styles.section}>
      <div className={clsx(styles.container, styles.sampleGrid)}>
        <div className={styles.sampleCopy}>
          <h2 className={styles.sectionTitle}>MCP native, session scoped</h2>
          <p>
            An agent connects to the daemon's <code>/mcp/</code> endpoint with a minted
            session. It sees only the tools whose effects that session holds a grant for.
          </p>
          <p>
            The call result carries the tool's output plus the <code>world</code> block:
            the effect records it produced and the plan that ran. Resolve any record by its
            content address.
          </p>
          <Link to="/docs/agents" className={styles.textLink}>
            Connect an agent <ArrowRight size={16} />
          </Link>
        </div>
        <div className={styles.terminal}>
          <div className={styles.terminalHeader}>
            <span className={clsx(styles.dot, styles.dotRed)} />
            <span className={clsx(styles.dot, styles.dotYellow)} />
            <span className={clsx(styles.dot, styles.dotGreen)} />
            <span className={styles.terminalName}>mcp: tools/list</span>
          </div>
          <pre className={styles.terminalBody}>
            <code>{SAMPLE}</code>
          </pre>
        </div>
      </div>
    </section>
  );
}

function Closing() {
  return (
    <section className={clsx(styles.section, styles.closing)}>
      <div className={styles.container}>
        <h2 className={styles.closingTitle}>
          The underlying OS executes bytes. <span className={styles.gradient}>World executes intent.</span>
        </h2>
        <div className={styles.closingLinks}>
          <Link to="/docs/intro" className={clsx(styles.btn, styles.btnPrimary)}>
            <BookOpen size={18} /> What is AILANG World
          </Link>
          <a href="https://ailang.sunholo.com" className={clsx(styles.btn, styles.btnGhost)}>
            <Network size={18} /> The AILANG language
          </a>
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <Hero />
      <main>
        <StatusBanner />
        <HowItWorks />
        <Features />
        <Sample />
        <Closing />
      </main>
    </Layout>
  );
}
