# TypeForm Slack Inviter

Automatically send Slack invitations for new responses in TypeForm.

## What It Does

TypeForm Slack Inviter checks TypeForm responses periodically and sends Slack
invites to emails in new responses.

Used in production at [RemoteMesh](https://www.remotemesh.com/), a chat community
for remote work.

## Installation

    go get github.com/sungwoncho/typeform-slack-inviter

## Usage

    typeform_slack_inviter --typeformAPIKey=### --formUID=### --interval=### slackAPIToken=###

* `typeformAPIKey`: API key for TypeForm
* `formUID`: a UID for your TypeForm
* `interval`: interval to check TypeForm responses, in minutes
* `SlackAPIToken`: API token for Slack

## License

MIT
