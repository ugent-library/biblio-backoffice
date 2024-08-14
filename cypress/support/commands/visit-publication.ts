import { logCommand, updateLogMessage, updateConsoleProps } from "./helpers";

export default function visitPublication(
  biblioIDAlias: Cypress.Alias = "@biblioId",
): void {
  const log = logCommand("visitPublication", {
    "Biblio ID alias": biblioIDAlias,
  });

  cy.get(biblioIDAlias, { log: false }).then((biblioId) => {
    updateLogMessage(log, biblioId);
    updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

    cy.intercept({ url: `/publication/${biblioId}/description*`, times: 1 }).as(
      "editPublication",
    );

    cy.visit(`/publication/${biblioId}`, { log: false });

    cy.wait("@editPublication", { log: false });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitPublication(biblioIDAlias?: Alias): Chainable<void>;
    }
  }
}
