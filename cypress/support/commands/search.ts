import { logCommand, updateLogMessage } from "./helpers";

const NO_LOG = { log: false };

export default function search(query: string): Cypress.Chainable<number> {
  const log = logCommand("search", { query }, query);

  // Give elastic some extra time to index changes
  cy.wait(1000, NO_LOG);

  cy.get('input[placeholder="Search..."]', NO_LOG)
    .clear(NO_LOG)
    .type(query, NO_LOG)
    .closest(".input-group", NO_LOG)
    .contains(".btn", "Search", NO_LOG)
    .click(NO_LOG);

  const REGEX = /^Showing (\d+-\d+ of )?(?<count>\d+) publications$/;

  return cy
    .contains(".card-header .c-body-small", REGEX, NO_LOG)
    .then(($el) => {
      const text = $el.text().trim();
      const count = parseInt(text.match(REGEX).groups.count);

      updateLogMessage(log, count);

      return count;
    })
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      search(query: string): Cypress.Chainable<number>;
    }
  }
}
