import { logCommand } from "./helpers";

type PublicationType =
  | "Journal Article"
  | "Book Chapter"
  | "Book"
  | "Conference contribution"
  | "Dissertation"
  | "Miscellaneous"
  | "Issue";

const NO_LOG = { log: false };

export default function setUpPublication(
  publicationType: PublicationType,
  prepareForPublishing = false,
): void {
  logCommand(
    "setUpPublication",
    {
      "Publication type": publicationType,
      "Prepare for publishing": prepareForPublishing,
    },
    publicationType,
  );

  cy.visit("/publication/add", NO_LOG);

  cy.contains("Enter a publication manually", NO_LOG)
    .find(":radio", NO_LOG)
    .click(NO_LOG);
  cy.contains(".btn", "Add publication(s)", NO_LOG).click(NO_LOG);

  cy.intercept("/publication/*/description*").as("completeDescription");

  cy.contains(new RegExp(`^${publicationType}$`), NO_LOG).click(NO_LOG);
  cy.contains(".btn", "Add publication(s)", NO_LOG).click(NO_LOG);

  // Extract biblioId at this point
  cy.get("#show-content", NO_LOG)
    .attr("hx-get")
    .then((hxGet) => {
      const biblioId = hxGet.match(
        /\/publication\/(?<biblioId>.*)\/description/,
      )?.groups["biblioId"];

      if (!biblioId) {
        throw new Error("Could not extract biblioId.");
      }

      return biblioId;
    })
    .as("biblioId", { type: "static" });

  cy.wait("@completeDescription", NO_LOG);

  cy.updateFields(
    "Publication details",
    () => {
      cy.setFieldByLabel("Title", `The ${publicationType} title [CYPRESSTEST]`);

      if (prepareForPublishing) {
        cy.setFieldByLabel(
          "Publication year",
          new Date().getFullYear().toString(),
        );
      }
    },
    true,
  );

  if (prepareForPublishing) {
    cy.updateFields(
      "Authors",
      () => {
        cy.setFieldByLabel("First name", "Dries");
        cy.setFieldByLabel("Last name", "Moreels");

        cy.contains(".btn", "Add author", NO_LOG).click(NO_LOG);
      },
      /^Save$/,
    );
  }

  cy.contains(".btn", "Complete Description", NO_LOG).click(NO_LOG);
}

declare global {
  namespace Cypress {
    interface Chainable {
      setUpPublication(
        publicationType: PublicationType,
        prepareForPublishing?: boolean,
      ): Chainable<void>;
    }
  }
}
