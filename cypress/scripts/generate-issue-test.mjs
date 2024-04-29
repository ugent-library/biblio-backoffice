import fs from "fs";
import path from "path";
import { execSync } from "child_process";

run();

async function run() {
  const issueId = process.argv[2];

  if (!issueId) {
    error("Please provide a valid issue ID.");
  }

  const specPath = path.resolve(`./cypress/e2e/issues/issue-${issueId}.cy.ts`);
  if (fs.existsSync(specPath)) {
    error(`Issue spec ${specPath} already exists.`);
  }

  const response = await fetch(
    `https://api.github.com/repos/ugent-library/biblio-backoffice/issues/${issueId}`,
  );
  if (response.ok) {
    const { title, pull_request } = await response.json();

    if (pull_request) {
      error(`#${issueId} is a pull request, not an issue.`);
    }

    fs.writeFileSync(
      specPath,
      `// https://github.com/ugent-library/biblio-backoffice/issues/${issueId}

describe('Issue #${issueId}: ${title}', () => {
  it('should ...', () => {

  })
})`,
    );

    execSync(`open ${specPath}`);
  } else {
    const { message } = await response.json();
    error(message);
  }
}

function error(message) {
  console.error(`ERROR: ${message}`);
  process.exit(1);
}
