# Contributing

Thank you for for taking the time to contribute to the project!

All contributions are welcome - minor fixes, new features, any improvement.

If you'd like to contribute to Autonity, please follow our process for reporting issues and submitting code changes.

## Log an issue

Before logging a new issue check existing [issues](https://github.com/autonity/autonity/issues) to make sure it hasn't already been reported and fixed.

If creating an issue, please:

- Provide a descriptive title
- Describe the problem and if possible how to replicate it
- Include any error messages or codes that were returned if relevant
- State the version of Autonity software you were using


## Submit a code change

Changes are proposed and submitted using a Git pull request workflow for submitting work to the project's `master` branch. 

For a minor change, simply submit a pull request.

If you want to propose a major change, then please raise a Git issue and discuss with the core devs to get some early feedback on the idea. For example, a new feature request could already be met by some other means.

### Contributor workflow

The basic developer workflow is:

- Fork the Autonity repo and create a topic branch off [`master`](https://github.com/autonity/autonity/tree/master)

- If your change adds code that should be tested, then add tests and make sure that your tests pass before committing

- Follow commit best practice:
    - Use incremental and [atomic commits](https://en.wikipedia.org/wiki/Atomic_commit) that affect a specific area of functionality. A clean git history helps review and debugging, making it easier for the community and code maintainers to review a change now and trace change history in the future.
    - Make sure tests pass before you commit
    - Reference any issues the commit resolves by ID
    - Use simple descriptive commit messages
    - If a change is breaking, make this clear in the commit message

- Create a [pull request from the fork](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork). Please give the PR a descriptive title and summarise the changes being proposed. Add references to any relevant Git Issues the PR resolves.

- Review. PR's will be reviewed by code maintainers before merge into the codebase. Discuss and action any review feedback by comment and discussion in the PR. Accepted changes are merged in by project code maintainers after PR approval.

- Thank you! Please accept our thanks for your contribution to the project!

## Setup and contributing guidelines

### Environment setup

For configuring your environment and managing project dependencies to build and run Autonity Go Client please see the [README](https://github.com/autonity/autonity#readme).

### Coding guidelines

Please ensure your contributions follow our coding guidelines:

- Go code must follow official Golang coding guidelines for:
    - [formatting](https://golang.org/doc/effective_go.html#formatting) and use [gofmt](https://golang.org/cmd/gofmt/))
    - documentation and follow the [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines
- Solidity code must follow the official Solidity coding [Style Guide](https://docs.soliditylang.org/en/latest/style-guide.html)
- JavaScript must adhere to Mozilla Developer Network [JavaScript guidelines](https://developer.mozilla.org/en-US/docs/mdn/guidelines/code_guidelines/javascript)
