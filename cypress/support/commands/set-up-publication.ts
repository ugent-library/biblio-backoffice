import { extractSnapshotId, waitForIndex } from "support/util";
import { logCommand, updateConsoleProps } from "./helpers";

const publicationTypes = {
  "Journal Article": "journal_article",
  "Book Chapter": "book_chapter",
  Book: "book",
  "Conference contribution": "conference",
  Dissertation: "dissertation",
  Miscellaneous: "miscellaneous",
  "Book (editor)": "book_editor",
  "Issue (editor)": "issue_editor",
} as const;

type PublicationTypes = typeof publicationTypes;
export type PublicationType = keyof PublicationTypes;

type SetUpPublicationOptions = {
  prepareForPublishing?: boolean;
  title?: string;
  biblioIDAlias?: string;
  shouldWaitForIndex?: boolean;
};

export default function setUpPublication(
  publicationType: PublicationType = "Miscellaneous",
  options: SetUpPublicationOptions = {},
): void {
  const {
    prepareForPublishing = false,
    title = `The ${publicationType} title`,
    biblioIDAlias = "biblioId",
    shouldWaitForIndex = false,
  } = options;

  const log = logCommand(
    "setUpPublication",
    {
      "Publication type": publicationType,
      "Prepare for publishing": prepareForPublishing,
      title,
      "Biblio ID alias": biblioIDAlias,
      "Should wait for index": shouldWaitForIndex,
    },
    publicationType,
  );

  cy.htmxRequest({
    method: "POST",
    url: "/add-publication/import/single/confirm",
    form: true,
    body: { publication_type: publicationTypes[publicationType] },
  })
    .then(extractBiblioId)
    .as(biblioIDAlias, { type: "static" })
    .then((biblioId) => {
      updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

      // Load the edit form to retrieve the snapshot ID
      cy.htmxRequest({
        url: `/publication/${biblioId}/details/edit`,
      })
        .then(extractSnapshotId)
        .then((snapshotId) => {
          // Then update details

          const body = {
            title,
            classification: "U",
          };

          if (prepareForPublishing) {
            body["year"] = new Date().getFullYear().toString();
          }

          cy.htmxRequest({
            method: "PUT",
            url: `/publication/${biblioId}/details`,
            headers: {
              "If-Match": snapshotId,
            },
            form: true,
            body,
          });
        });

      if (prepareForPublishing) {
        cy.addAuthor("John", "Doe");
      }

      if (shouldWaitForIndex) {
        waitForIndex("publication", biblioId);
      }
    });
}

function extractBiblioId(response: Cypress.Response<string>) {
  const { biblioId } = response.body.match(
    /href="\/publication\/(?<biblioId>[A-Z0-9]+)\/add\/confirm"/,
  ).groups;

  return biblioId;
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
