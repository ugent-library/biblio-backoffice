import { logCommand } from "./helpers";

export type PublicationType =
  | "Journal Article"
  | "Book Chapter"
  | "Book"
  | "Conference contribution"
  | "Dissertation"
  | "Miscellaneous"
  | "Issue";

const NO_LOG = { log: false };

type SetUpPublicationOptions = {
  prepareForPublishing?: boolean;
  title?: string;
  biblioIDAlias?: string;
};

export default function setUpPublication(
  publicationType: PublicationType = "Miscellaneous",
  options: SetUpPublicationOptions = {},
): void {
  const {
    prepareForPublishing = false,
    title = `The ${publicationType} title`,
    biblioIDAlias = "biblioId",
  } = options;

  logCommand(
    "setUpPublication",
    {
      "Publication type": publicationType,
      "Prepare for publishing": prepareForPublishing,
      title,
      "Biblio ID alias": biblioIDAlias,
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
    .as(biblioIDAlias, { type: "static" });

  cy.wait("@completeDescription", NO_LOG);

  cy.updateFields(
    "Publication details",
    () => {
      cy.setFieldByLabel("Title", `${title} [CYPRESSTEST]`);

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
      true,
    );
  }

  cy.contains(".btn", "Complete Description", NO_LOG).click(NO_LOG);
}

declare global {
  namespace Cypress {
    interface Chainable {
      setUpPublication(
        publicationType?: PublicationType,
        options?: SetUpPublicationOptions,
      ): Chainable<void>;
    }
  }
}
