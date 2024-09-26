import { logCommand } from "./helpers";

export function deletePublication(biblioId: string) {
  logCommand(`deletePublication`, { biblioId }, biblioId);

  deleteOne("publication", biblioId);
}

export function deleteDataset(biblioId: string) {
  logCommand(`deleteDataset`, { biblioId }, biblioId);

  deleteOne("dataset", biblioId);
}

function deleteOne(scope: Biblio.Scope, biblioId: string) {
  cy.htmxRequest({
    url: `/${scope}/${biblioId}/confirm-delete`,
  }).then((response) => {
    const dangerButton = Cypress.$(response.body).find(".btn-danger");
    cy.htmxRequest({
      method: "DELETE",
      url: dangerButton.attr("hx-delete"),
      headers: JSON.parse(dangerButton.attr("hx-headers")),
    });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      deletePublication(biblioId: string): Chainable<void>;

      deleteDataset(biblioId: string): Chainable<void>;
    }
  }
}
