import { logCommand, updateLogMessage, updateConsoleProps } from "./helpers";

export default function visitDataset(
  bilbioIdAlias: Cypress.Alias = "@biblioId",
): void {
  const log = logCommand("visitDataset", { "Biblio ID alias": bilbioIdAlias });

  cy.get(bilbioIdAlias, { log: false }).then((biblioId) => {
    updateLogMessage(log, biblioId);
    updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

    cy.intercept(`/dataset/${biblioId}/description*`).as("editDataset");

    cy.visit(`/dataset/${biblioId}`, { log: false });

    cy.wait("@editDataset", { log: false });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitDataset(bilbioIdAlias?: Alias): Chainable<void>;
    }
  }
}
