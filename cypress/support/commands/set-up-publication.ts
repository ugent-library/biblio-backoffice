import { extractSnapshotId, waitForIndex } from "support/util";
import { logCommand, updateConsoleProps } from "./helpers";

const publicationTypes = {
  "Journal Article": "journal_article",
  "Book Chapter": "book_chapter",
  Book: "book",
  "Conference contribution": "conference",
  Dissertation: "dissertation",
  Miscellaneous: "miscellaneous",
  "Book editor": "book_editor",
  "Issue editor": "issue_editor",
} as const;

type PublicationTypes = typeof publicationTypes;
export type PublicationType = keyof PublicationTypes;

type SetUpPublicationOptions = {
  title?: string;
  otherFields?: Record<string, unknown>;
  biblioIDAlias?: string;
  prepareForPublishing?: boolean;
  publish?: boolean;
  shouldWaitForIndex?: boolean;
};

export default function setUpPublication(
  publicationType: PublicationType = "Miscellaneous",
  options: SetUpPublicationOptions = {},
): void {
  const {
    title = `The ${publicationType} title`,
    otherFields = {},
    biblioIDAlias = "biblioId",
    prepareForPublishing = false,
    publish = false,
    shouldWaitForIndex = false,
  } = options;

  const log = logCommand(
    "setUpPublication",
    {
      title,
      "Other fields": otherFields,
      "Biblio ID alias": biblioIDAlias,
      "Prepare for publishing": prepareForPublishing,
      Publish: publish,
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

          if (prepareForPublishing || publish) {
            body["year"] = new Date().getFullYear().toString();
          }

          Object.assign(body, otherFields);

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

      if (prepareForPublishing || publish) {
        cy.addAuthor("John", "Doe");
      }

      if (publish) {
        publishPublication(biblioId);
      }

      if (shouldWaitForIndex) {
        waitForIndex("publication", biblioId);
      }
    });
}

function publishPublication(biblioId: string) {
  cy.htmxRequest({
    url: `/publication/${biblioId}/add/confirm`,
  })
    .then((r) =>
      extractSnapshotId(
        r,
        `button[hx-post="/publication/${biblioId}/add/publish"]:contains("Publish to Biblio")`,
      ),
    )
    .then((snapshotId) => {
      cy.htmxRequest({
        method: "POST",
        url: `/publication/${biblioId}/add/publish`,
        headers: {
          "If-Match": snapshotId,
        },
      });
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
