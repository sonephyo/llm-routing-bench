import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Multiple Load Balancing Strategies',
    description: (
      <>
        Benchmark and compare Round Robin, Consistent Hashing, Least KV Cache and Least Queue
        routing strategies to find what works best for your LLM workload.
      </>
    ),
  },
  {
    title: 'Real-Time Metrics with Prometheus',
    description: (
      <>
        Built-in Prometheus integration tracks request counts and latency per
        backend, giving to evaluate routing performance.
      </>
    ),
  },
  {
    title: 'Pluggable & Extensible',
    description: (
      <>
        The router is built around <code>loadbalancer.Router</code> interface
        in Go, making it straightforward to add new routing strategies and run
        them against the same benchmark harness.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center padding-horiz--md" style={{paddingTop: '1.5rem'}}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
