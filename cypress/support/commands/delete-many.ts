import { logCommand, updateConsoleProps } from "./helpers";

const NO_LOG = { log: false };

export function deletePublications(query: string) {
  const log = logCommand("deletePublications", { query }, query);

  deleteMany(log, "publication", query);
}

export function deleteDatasets(query: string) {
  const log = logCommand("deleteDatasets", { query }, query);

  deleteMany(log, "dataset", query);
}

function deleteMany(
  log: Cypress.Log,
  type: "publication" | "dataset",
  query: string,
) {
  if (query.trim().length === 0) {
    throw new Error("Query string is empty");
  }

  cy.getAllCookies(NO_LOG).then((originalCookies) => {
    cy.loginAsLibrarian("librarian1");

    cy.visit(`/${type}`, {
      qs: { q: query, "page-size": 1000 },
      ...NO_LOG,
    })
      .then(() => {
        const items = Cypress.$("button:contains('Biblio ID') + code").map(
          (_, i) => i.textContent.trim(),
        );

        updateConsoleProps(log, (cp) => (cp["Items to delete"] = items.get()));

        items.each((_, item) => {
          if (type === "publication") {
            cy.deletePublication(item);
          } else if (type === "dataset") {
            cy.deleteDataset(item);
          } else {
            throw new Error(`Unknown type: ${type}`);
          }
        });
      })
      .then(() => {
        // Restore original session
        cy.clearAllCookies(NO_LOG);

        for (const cookie of originalCookies) {
          cy.setCookie(cookie.name, cookie.value, NO_LOG);
        }
      });
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      deletePublications(query: string): Chainable<void>;

      deleteDatasets(query: string): Chainable<void>;
    }
  }
}
