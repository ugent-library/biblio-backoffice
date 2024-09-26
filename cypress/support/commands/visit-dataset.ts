import { logCommand, updateLogMessage, updateConsoleProps } from "./helpers";

export default function visitDataset(
  biblioIDAlias: Cypress.Alias = "@biblioId",
): void {
  const log = logCommand("visitDataset", { "Biblio ID alias": biblioIDAlias });

  cy.get(biblioIDAlias, { log: false }).then((biblioId) => {
    updateLogMessage(log, biblioId);
    updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

    cy.intercept({ url: `/dataset/${biblioId}/description*`, times: 1 }).as(
      "editDataset",
    );

    cy.visit(`/dataset/${biblioId}`, { log: false });

    cy.wait("@editDataset", { log: false });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitDataset(biblioIDAlias?: Alias): Chainable<void>;
    }
  }
}
