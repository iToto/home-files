# Weekly meeting

## Desired outcomes

* A prioritized list of issues, in the current project or milestone, so that the team has shared understanding of what to work on next.
* An understanding of what we're each working on next, so that we can proceed efficiently - whether in one or parallel paths.
* A list of any blockers, and an agreement on who can assist in unblocking and when.
* A list of any newly-identified action items, captured via [meeting notes][meeting-notes].
* An understanding of our progress toward the upcoming milestone (if applicable).
* An awareness of the project's *next* milestone (if applicable).
* A list of interactions with other teams (e.g., Identity, Security, etc.) that we’ll face in the next two weeks.

## Stakeholders

* Engineering manager: @iToto
* Product manager: @hpsin

## Meeting roles

* Facilitator: first responder from the previous week
* Notetaker: rotated

## Decision making

We prefer to make decisions through consensus. If the team cannot reach consensus, the EM will identify the appropriate person to make the decision at hand, and we will move forward with that person's decision.

## Preparation

* Merge the weekly deploy branch
* Go through 'In progress' column
  * Mark any issues as 'Done' with a cross reference to the PR that closed it
  * Ping responsible team member to see if they have an update or need anything to move forward

## Process

0. Identify notetaker for our [meeting notes][meeting-notes].
0. Are there any announcements before we go through the board?
0. Review our SLOs: [authzd]( https://catalog.githubapp.com/services/authzd), [Roles and Permissions](https://catalog.githubapp.com/services/github/roles_and_permissions)
0. Check [action items from previous week][weekly-action-items], and identify how to resolve any open action items.
0. Share your screen with the [Auhtorization board][triage] open
0. Triage column
    * Give a brief 1-2 sentence summary of the issue
    * Explain the recommended priority. Team members then need to :+1: / :-1: to agree / disagree with the assessment. Move card to appropriate column.
    * If an issue requires additional input, leave it in the 'Triage' column for follow up during the week. Assign a person to do the follow up
0. Review in-progress and blocked items in the [WIP Tab on the Auhorization board][in-progress]
    * Reminder: Ready For Work is where we're pulling work for the current iteration
    * Is this still in progress? What do you need here?
    * Is there any work that needs to be transferred to this week's FR?
    * Is there any work that needs to be included to in-progress
    * Is there anything that needs to be brought in/out of the current iteration?

#### Within Iteration?

0. Open the Project Prioritization in the [Iteration Planning Authorization board][project-prio]
    * Review projects in current iteration and give a quick status update (on schedule/at risk)
    * Briefly go over any changes to iteration priority and discuss if needed. Raise any concerns or questions.
0. Check if there's anything coming up this/next week that will affect our time.

#### Within Cooldown?

1. Open the Cooldown tab in the [Iteration Planning Authorization board][cooldown-tab]
    * Review projects in current iteration and give a quick status update (on schedule/at risk)
    * Briefly go over any changes to iteration priority and discuss if needed. Raise any concerns or questions.
2. Check if there's anything coming up this/next week that will affect our time.



0. Each person answers:
    * "Is there anything you're nervous about, or think isn't going well?"
    * "Is there anything you're excited about, or think is going well?"
0. Publish [notes][meeting-notes].
0. Notetaker rotates the Topic in `#authorization`for the new FR.

[meeting-notes]: https://github.com/github/authorization/discussions
[weekly-action-items]: https://github.com/github/authorization/labels/weekly-ai
[triage]: https://github.com/orgs/github/projects/7025/views/6
[in-progress]: https://github.com/orgs/github/projects/7025/views/4
[project-prio]: https://github.com/orgs/github/projects/7025/views/13
[identity-board]: https://github.com/orgs/github/projects/3765/views/13
[cooldown-tab]:https://github.com/orgs/github/projects/7025/views/15
