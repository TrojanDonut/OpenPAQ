# Contributing to OpenPAQ

Thank you for your interest in contributing to **OpenPAQ**!

We welcome all kinds of contributions — whether you're reporting bugs, suggesting features, or submitting pull requests.
However, please note that **GitHub Issues and Pull Requests are currently disabled** while we finalize an initial version of the project. 
Once we're ready, we'll announce the opening of community contributions on the project page.

---

## General Guidelines

Once contributions are open, please follow the guidelines below.

### 1. Creating Issues

When submitting an issue (bug report, feature request, or refactoring proposal), please make sure to:

- Provide a **clear rationale** — explain *why* this change or suggestion is meaningful.
- Include **examples or use cases** to help understand the context and benefit.
- Use one of the following tags in the issue title:
    - `[Bug]` — for bugs or regressions
    - `[Feature]` — for new capabilities
    - `[Refactor]` — for internal structural improvements

#### Example Issue

Title: [Feature] Better matches for Belgian addresses

Description: In order to better match some Belgian addresses a new normalizer for "be" is added and registered

"Abbreviations in streetnames"

Straat --> Str.

Laan --> Ln


### 2. Submitting Pull Requests (PRs)

Once pull requests are enabled:

- Make sure your PR has a **clear, focused scope**.
- Keep code **clean and readable**.
- Include **tests** and **documentation** if applicable.
- Ensure **existing tests pass**.
- Reference any related issue (e.g., `Fixes #42`).

