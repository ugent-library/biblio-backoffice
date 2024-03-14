import { logCommand, updateLogMessage, updateConsoleProps } from "./helpers";

export default function visitPublication(): void {
  const log = logCommand("visitPublication");

  cy.get("@biblioId", { log: false }).then((biblioId) => {
    updateLogMessage(log, biblioId);
    updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

    cy.intercept(`/publication/${biblioId}/description*`).as("editPublication");

    cy.visit(`/publication/${biblioId}`, { log: false });

    cy.wait("@editPublication", { log: false });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitPublication(): Chainable<void>;
    }
  }
}
