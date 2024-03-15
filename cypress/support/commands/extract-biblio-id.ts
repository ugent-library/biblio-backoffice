import { logCommand, updateConsoleProps } from "./helpers";

export default function extractBiblioId(
  subject: undefined | JQuery<HTMLElement>,
  alias: string = "biblioId",
) {
  const log = logCommand("extractBiblioId", { alias });

  let cySubject: Cypress.Chainable;
  if (subject) {
    cySubject = cy.wrap(subject, { log: false });
  } else {
    cySubject = cy.get(".card > .card-body > .list-group > .list-group-item", {
      log: false,
    });
  }

  cySubject.then((el) => {
    updateConsoleProps(log, (cp) => (cp.subject = el));
  });

  cySubject.then((s) => {
    if (s.length !== 1) {
      expect(s).to.have.length(1);
    }
  });

  cySubject
    .contains("Biblio ID", { log: false })
    .next(".c-code", { log: false })
    .invoke({ log: false }, "text")
    .as(alias, { type: "static" })
    .finishLog(log, true);
}

declare global {
  namespace Cypress {
    interface Chainable {
      extractBiblioId(alias?: string): Chainable<string> | never;
    }
  }
}
