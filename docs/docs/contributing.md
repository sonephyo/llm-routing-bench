---
title: Contributing
sidebar_position: 99
---

Contributions are welcome — whether that's reporting a bug, suggesting a new routing strategy, or submitting a pull request.

## Reporting Issues

If you run into a bug or have a feature request, open an issue on GitHub:

[github.com/sonephyo/llm-routing-bench/issues](https://github.com/sonephyo/llm-routing-bench/issues)

Please include:
- What you were trying to do
- What happened vs. what you expected
- Relevant logs or error output
- Your OS, Docker version, and GPU setup (if applicable)

## Submitting a Pull Request

1. Fork the repository on GitHub
2. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
3. Make your changes and commit them
4. Push your branch and open a pull request against `main`
5. Describe what your PR does and why in the PR description

## Ideas for Contributions

- New load balancing strategies (implement the `loadbalancer.Router` interface in Go)
- Additional benchmark scripts
- Documentation improvements
- Bug fixes
