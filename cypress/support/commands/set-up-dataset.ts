import { waitForIndex, extractSnapshotId } from "support/util";
import { logCommand, updateConsoleProps } from "./helpers";

type SetUpDatasetOptions = {
  prepareForPublishing?: boolean;
  title?: string;
  biblioIDAlias?: string;
  shouldWaitForIndex?: boolean;
};

export default function setUpDataset({
  prepareForPublishing = false,
  title = "The dataset title",
  biblioIDAlias = "biblioId",
  shouldWaitForIndex = false,
}: SetUpDatasetOptions = {}): void {
  const log = logCommand("setUpDataset", {
    "Prepare for publishing": prepareForPublishing,
    title,
    "Biblio ID alias": biblioIDAlias,
    "Should wait for index": shouldWaitForIndex,
  });

  cy.htmxRequest({
    method: "POST",
    url: "/add-dataset",
    form: true,
    body: { method: "manual" },
  })
    .then(extractBiblioId)
    .as(biblioIDAlias, { type: "static" })
    .then((biblioId) => {
      updateConsoleProps(log, (cp) => (cp["Biblio ID"] = biblioId));

      // Load the edit form to retrieve the snapshot ID
      cy.htmxRequest({
        url: `/dataset/${biblioId}/details/edit`,
      })
        .then(extractSnapshotId)
        .then((snapshotId) => {
          // Then update details
          const body = {
            title: `${title} [CYPRESSTEST]`,
            identifier_type: "DOI",
            identifier: "10.7202/1041023ar",
          };

          if (prepareForPublishing) {
            body["year"] = new Date().getFullYear().toString();
            body["publisher"] = "UGent";
            body["format"] = ["text/csv"];
            body["license"] = "CC0-1.0";
            body["access_level"] = "info:eu-repo/semantics/openAccess";
          }

          cy.htmxRequest({
            method: "PUT",
            url: `/dataset/${biblioId}/details`,
            headers: {
              "If-Match": snapshotId,
            },
            form: true,
            body,
          });
        });

      if (prepareForPublishing) {
        cy.addCreator("John", "Doe");
      }

      if (shouldWaitForIndex) {
        waitForIndex("dataset", biblioId);
      }
    });
}

function extractBiblioId(response: Cypress.Response<string>) {
  const { biblioId } = response.body.match(
    /href="\/dataset\/(?<biblioId>[A-Z0-9]+)\/add\/confirm"/,
  ).groups;

  return biblioId;
}

declare global {
  namespace Cypress {
    interface Chainable {
      setUpDataset(options?: SetUpDatasetOptions): Chainable<void>;
    }
  }
}
