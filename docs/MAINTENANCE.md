# ZINTENT Maintenance and Feedback Loop

This document outlines the commitment of the ZINTENT core team to the community and the internal maintenance rhythm.

## Issue Triage SLA (Service Level Agreement)

To ensure a responsive and healthy community, we aim for the following response times:

| Category | Response Time | Goal |
| :--- | :--- | :--- |
| **Critical Bugs** | < 24 Hours | Patch or workaround provided within 48 hours. |
| **Standard Bugs** | < 3 Days | Reproduced and labeled within 5 days. |
| **Feature Requests** | < 1 Week | Categorized (Backlog/Planned/Rejected) within 10 days. |
| **Discussions / Questions** | < 5 Days | Initial response or community redirection. |

## Weekly Maintenance Rhythm

Every **Friday**, the core team performs the following tasks:

1. **Issue Sweep**: Review all open issues, update labels, and ping stagnant discussions.
2. **PR Review**: Prioritize community PRs and provide constructive feedback.
3. **Snapshot Refresh**: Run `npm run test:snapshots:update` to ensure the registry and compiler output remain in sync with any intentional changes.
4. **Benchmark Audit**: Verify that the latest changes haven't introduced performance regressions.
5. **Roadmap Sync**: Update `ROADMAP.md` and `LAUNCH_PLAN_30D.md` to reflect the week's progress.

## Contribution Channels

- **GitHub Issues**: Primary for bugs and specific feature requests.
- **GitHub Discussions**: For broader design questions, RFCs, and community support.
- **X (Twitter)**: For announcements and quick updates.

## Release Cadence

- **Stable**: Monthly releases (e.g., v2.1.0, v2.2.0).
- **Canary/RC**: Weekly or as needed for testing major architectural changes.
