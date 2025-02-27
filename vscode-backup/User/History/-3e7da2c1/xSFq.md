# How we work

## Introduction

The Authorization team builds and provides foundations, frameworks, services and APIs that provide secure authorization for external customers and internal teams and services. This document describes how we work. It is closely related to the GitHub Engineering [how we work](https://thehub.github.com/engineering/how-we-work/) doc and inspired from the Authentication's team's [how-we-work](https://github.com/github/authentication/blob/main/docs/how_we_work.md) document.

## Guiding Principles

Everyone in the team trusts themselves and each other to do the right thing. Everyone on the team holds themselves and others accountable to our values and our mission. Mechanics and processes are simply tools we use to help us do our best work efficiently. If a process or practice isn't doing this, we can change it and iterate.

See our [Design Principles](/docs/design-principles.md) for the principles that govern the Authorization platform.

## Key Points

- We are running a version of the [6-week cycle](https://basecamp.com/shapeup). In essence, we run 6-week iterations where we focus on shipping planned work followed by 2-weeks of cooldown where we work on un-planned work — including tech debt. [This diagram](https://miro.com/app/board/uXjVPWC0O8Q=/) helps illustrate the process.
- Our primary customer base is internal teams. We recognize that empowering internal customers, e.g. feature teams, pays off more than doing one-off feature work ourselves.
- We deliver the highest priority customer-facing improvement in collaboration with the product team. We also take the time to invest in our foundations so we can easily extend and build on top of those foundations to meet customer needs.
- We focus on quality while maintaining a healthy tempo. If we are at risk of slipping a deadline, we will first attempt to cut scope. Extending deadlines is not our first option, but something we evaluate with our Product partners and leadership chain when we need to.
- We limit the number of work streams in progress to a manageable size based on team members' availability. For example, accounting for the FR rotation.
- We set quarterly OKRs in collaboration with the Identity, and Platform organization goals. These goals are tied to our [top-level objectives](/planning/objectives/README.md) to help give purpose
- We use a combination of Initiatives, Epics, Batches, Tasks, Bugs to manage work.
- We mitigate an incident first, and then we investigate and create a post-mortem. Returning a service to a fully working state and restoring SLOs is the paramount goal during incidents.
- When considering what to prioritize, we follow the larger Engineering's prioritization of Security > Reliability > Performance > Features. As our focus is Authorization, most if not all of our potential projects will have a Security impact.
- We run [6-week cycles](https://miro.com/app/board/uXjVPWC0O8Q=/) followed by 2-week cooldowns

## Process flow

Initiative ⇒ Epic ⇒ Batch ⇒ Task

## Initiatives (May span multiple quarters)

An initiative is a goal set in collaboration with the product team that provides the most important value to our customers.
The intent is to clarify the "Why?" behind the work we do and its ties to the desired customer-facing outcome.
An initiative can have one or more Epics.

## Epic (Roughly 4-6 weeks — or 1 iteration — of work)

An epic is part of a business initiative that, potentially combined with other epics, provides an impactful customer solution. An epic could directly or indirectly affect the customer. Examples of Epics are

- https://github.com/github/authorization/issues/2224
- https://github.com/github/authorization/issues/2221

Other Epic requirements:
- Weekly status update (green/yellow/red) provided by Epic DRI
- Epic Target Date calculated based on Batch Estimates and any deadline

## Batch (Roughly 1-2 weeks of work)

A batch is a shippable software component that, combined with other batches, delivers a solution defined in a parent Epic. Batches roll up to epics. Examples of batches are:

- https://github.com/github/authorization/issues/2229
- https://github.com/github/authorization/issues/2230

Other Batch requirements:
- Batch Issue should be assigned to Engineers working on it
- Estimate should be provided by the Engineers working in the Batch based on the Tasks outlined

## Task (<1 week of work)

A task is the smallest unit of work that is implemented with one or, preferably, a pair of developers. We use one or more PRs to execute tasks.

## Key Players

- **Epic DRI**: Responsible for Epic Kick-off, Weekly Status Updates on the Epic as well as any coordination needed across Batches within the Epic, and leading weekly sync (if one is needed)
- **Batch Lead**: DRI of the batch; Kicks off the batch with a meeting to go over the plan and gets feedback from the team. Advises engineers in picking the next highest priority task. Coordinates releases, creates (or delegates) demo presentations, and writes feature-complete posts.
- **Engineer**: Takes one task at a time, preferably with a pair, and completes the task set in the Acceptance Criteria. Responsible for the end to end development of the task up to release, monitoring, dashboards, and post-deployment diagnosis.
- **Product Manager**: Plays an essential role as the interface between customers and the team. Analyzes customer requirements and articulate them to the team. Leads the effort in prioritizing the Product backlog.
- **Engineering Manager**: Facilitates the interaction between PM and Engineers. Responsible for assigning the right engineer to do Spikes, lead batches, and manage team work in progress (WIP).

## Process Guidelines

- Initiatives are long-running, but all Epics within an Initiative aren't defined upfront. Ideally we plan 2-3 Epics ahead.
- No more than 2 Epics in parallel
- Engineers working on the Epic define its Batches
- Engineers working on a Batch define its Tasks
- Batch Engineers have autonomy and ownership over the feature implementation
- We pair as much as possible
- We work on the highest priority item in-order
- If it is unclear what the next work item is, consult with Batch Lead, then Epic Lead, then EM.
- Any major technical/architectural decisions with consequences are documented in an ADR
