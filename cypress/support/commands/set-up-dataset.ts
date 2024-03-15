import { logCommand } from "./helpers";

const NO_LOG = { log: false };

export default function setUpDataset(
  prepareForPublishing = false,
  title?: string,
): void {
  title ??= "The dataset title";

  logCommand("setUpDataset", {
    "Prepare for publishing": prepareForPublishing,
    title,
  });

  cy.visit("/dataset/add", NO_LOG);

  cy.intercept("/dataset/*/description*").as("completeDescription");

  cy.contains("Register a dataset manually", NO_LOG)
    .find(":radio", NO_LOG)
    .click(NO_LOG);
  cy.contains(".btn", "Add dataset", NO_LOG).click(NO_LOG);

  // Extract biblioId at this point
  cy.get("#show-content", NO_LOG)
    .attr("hx-get")
    .then((hxGet) => {
      const biblioId = hxGet.match(/\/dataset\/(?<biblioId>.*)\/description/)
        ?.groups["biblioId"];

      if (!biblioId) {
        throw new Error("Could not extract biblioId.");
      }

      return biblioId;
    })
    .as("biblioId", { type: "static" });

  cy.wait("@completeDescription", NO_LOG);

  cy.updateFields(
    "Dataset details",
    () => {
      cy.setFieldByLabel("Title", `${title} [CYPRESSTEST]`);

      cy.setFieldByLabel("Persistent identifier type", "DOI");
      cy.setFieldByLabel("Identifier", "10.5072/test/t");

      if (prepareForPublishing) {
        cy.setFieldByLabel("Access level", "Open access");
        cy.setFieldByLabel("Data format", "text/csv")
          .next(".autocomplete-hits", NO_LOG)
          .contains(".badge", "text/csv", NO_LOG)
          .click(NO_LOG);
        cy.setFieldByLabel("Publisher", "UGent");
        cy.setFieldByLabel(
          "Publication year",
          new Date().getFullYear().toString(),
        );
        cy.setFieldByLabel("License", "CC0 (1.0)");
      }
    },
    true,
  );

  if (prepareForPublishing) {
    cy.updateFields(
      "Creators",
      () => {
        cy.setFieldByLabel("First name", "Dries");
        cy.setFieldByLabel("Last name", "Moreels");

        cy.contains(".btn", "Add creator", NO_LOG).click(NO_LOG);
      },
      /^Save$/,
    );
  }

  cy.contains(".btn", "Complete Description", NO_LOG).click(NO_LOG);
}

declare global {
  namespace Cypress {
    interface Chainable {
      setUpDataset(
        prepareForPublishing?: boolean,
        title?: string,
      ): Chainable<void>;
    }
  }
}
