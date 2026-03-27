import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docSideBar: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: [
        'getting-started/prerequisites',
        'getting-started/dev-mode',
        'getting-started/production-mode',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'reference/routing-strategies',
        'reference/benchmarking',
        'reference/metrics',
        'reference/makefile-reference',
      ],
    },
    'contributing',
    'about',
  ],

  researchSideBar: [
    {
      type: 'category',
      label: 'Research',
      collapsible: false,
      items: [
        'research/motivation',
        'research/routing-strategies',
        'research/experiment-setup',
        'research/results',
        'research/related-work',
      ],
    },
  ],
};

export default sidebars;
