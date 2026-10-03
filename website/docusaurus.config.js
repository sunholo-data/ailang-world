// @ts-check
// AILANG World documentation site. The theme, fonts and navbar/footer shape
// follow the AILANG docs site (sunholo-data/ailang, docs/) so the two read as
// one family.

import {themes as prismThemes} from 'prism-react-renderer';

// Where the site is served. Both are configurable at build time:
//   SITE_URL=https://example.org BASE_URL=/ailang-world/ npm run build
const SITE_URL = process.env.SITE_URL || 'https://sunholo-data.github.io';
const BASE_URL_RAW = process.env.BASE_URL || '/';
const BASE_URL = BASE_URL_RAW.endsWith('/') ? BASE_URL_RAW : `${BASE_URL_RAW}/`;

const GITHUB_URL = 'https://github.com/sunholo-data/ailang-world';
const AILANG_URL = 'https://ailang.sunholo.com';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'AILANG World',
  tagline:
    'A semantic operating environment whose transactions are AILANG programs: agents propose, a verifier checks, and only verified, authorized changes commit.',
  favicon: 'https://ailang.sunholo.com/img/favicon.ico',

  url: SITE_URL,
  baseUrl: BASE_URL,

  // Used by `npm run deploy` (gh-pages). Publishing is the owner's call.
  organizationName: 'sunholo-data',
  projectName: 'ailang-world',
  deploymentBranch: 'gh-pages',
  trailingSlash: false,

  onBrokenLinks: 'throw',
  onBrokenAnchors: 'warn',

  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: 'throw',
    },
  },

  themes: ['@docusaurus/theme-mermaid'],

  // Google Fonts for Sunholo brand styling (same set as the AILANG site).
  headTags: [
    {
      tagName: 'link',
      attributes: {rel: 'preconnect', href: 'https://fonts.googleapis.com'},
    },
    {
      tagName: 'link',
      attributes: {
        rel: 'preconnect',
        href: 'https://fonts.gstatic.com',
        crossorigin: 'anonymous',
      },
    },
    {
      tagName: 'link',
      attributes: {
        rel: 'stylesheet',
        href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&family=Montserrat:wght@600;700;800&display=swap',
      },
    },
    {
      tagName: 'link',
      attributes: {
        rel: 'icon',
        type: 'image/svg+xml',
        href: `${BASE_URL}img/ailang-logo.svg`,
      },
    },
  ],

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: './sidebars.js',
          routeBasePath: '/docs',
          editUrl: `${GITHUB_URL}/tree/dev/website/`,
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      image: 'https://ailang.sunholo.com/img/ailang-social-card.jpg',
      navbar: {
        title: 'AILANG World',
        logo: {
          alt: 'AILANG logo',
          src: 'img/ailang-logo.svg',
        },
        items: [
          {to: '/docs/intro', label: 'What is World', position: 'left'},
          {
            type: 'docSidebar',
            sidebarId: 'docsSidebar',
            position: 'left',
            label: 'Documentation',
          },
          {to: '/docs/getting-started', label: 'Get started', position: 'left'},
          {to: '/docs/agents', label: 'For AI agents', position: 'left'},
          {to: '/docs/roadmap', label: 'Roadmap', position: 'left'},
          {href: AILANG_URL, label: 'AILANG', position: 'right'},
          {href: GITHUB_URL, label: 'GitHub', position: 'right'},
        ],
      },
      footer: {
        style: 'dark',
        logo: {
          alt: 'AILANG logo',
          src: 'img/ailang-logo.svg',
          href: '/',
          width: 48,
          height: 48,
        },
        links: [
          {
            title: 'Docs',
            items: [
              {label: 'What is AILANG World', to: '/docs/intro'},
              {label: 'Get started', to: '/docs/getting-started'},
              {label: 'Concepts', to: '/docs/concepts'},
              {label: 'Security', to: '/docs/security'},
              {label: 'Reference', to: '/docs/reference'},
            ],
          },
          {
            title: 'For agents',
            items: [
              {label: 'Agents guide', to: '/docs/agents'},
              {label: 'Guides', to: '/docs/guides'},
              {label: 'Roadmap to 1.0', to: '/docs/roadmap'},
            ],
          },
          {
            title: 'Project',
            items: [
              {label: 'GitHub', href: GITHUB_URL},
              {label: 'Issues', href: `${GITHUB_URL}/issues`},
              {
                label: 'Mission charter',
                href: `${GITHUB_URL}/blob/dev/design_docs/world-mission.md`,
              },
              {
                label: 'Design (DESIGN.md)',
                href: `${GITHUB_URL}/blob/dev/design_docs/DESIGN.md`,
              },
            ],
          },
          {
            title: 'AILANG family',
            items: [
              {label: 'AILANG language', href: AILANG_URL},
              {label: 'AILANG on GitHub', href: 'https://github.com/sunholo-data/ailang'},
              {label: 'llms.txt (AILANG)', href: `${AILANG_URL}/llms.txt`},
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Sunholo. Apache-2.0. Built with Docusaurus.`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
        additionalLanguages: ['bash', 'json', 'go', 'python'],
      },
      colorMode: {
        defaultMode: 'dark',
        disableSwitch: false,
        respectPrefersColorScheme: false,
      },
    }),
};

export default config;
