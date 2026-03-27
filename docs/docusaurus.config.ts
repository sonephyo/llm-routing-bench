import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'LLM Routing Benchmark',
  tagline: 'Measuring how routing strategies affect tail latency in LLM inference serving',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: 'https://sonephyo.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/llm-routing-bench/',

  // GitHub pages deployment config.
  organizationName: 'sonephyo',
  projectName: 'llm-routing-bench',

  onBrokenLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  plugins: [
    [
      require.resolve('@easyops-cn/docusaurus-search-local'),
      {
        hashed: true,
        indexDocs: true,
        indexPages: true,
        docsRouteBasePath: '/docs',
      },
    ],
  ],

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/sonephyo/llm-routing-bench/tree/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'LLM Routing Benchmark',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docSideBar',
          position: 'left',
          label: 'Documentation',
        },
        {type: 'docSidebar', sidebarId: 'researchSideBar', position: 'left', label: 'Research'},
        {
          href: 'https://github.com/sonephyo/llm-routing-bench',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Documentation',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/intro',
            },
          ],
        },
        {
          title: 'Research',
          items: [
            {
              label: 'Motivation',
              to: '/docs/research/motivation',
            },
          ],
        },
        {
          title: 'Project',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/sonephyo/llm-routing-bench',
            },
            {
              label: 'Issues',
              href: 'https://github.com/sonephyo/llm-routing-bench/issues',
            },
          ],
        },
        {
          title: 'Author',
          items: [
            {
              label: 'GitHub @sonephyo',
              href: 'https://github.com/sonephyo',
            },
            {
              label: 'LinkedIn',
              href: 'https://www.linkedin.com/in/soney7/',
            },
            {
              label: 'Email',
              href: 'mailto:sonephyo7777777@gmail.com',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Sone Phyo. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
