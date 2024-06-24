import { extractSnapshotId } from "support/util";
import { logCommand } from "./helpers";

type SetUpDatasetOptions = {
  prepareForPublishing?: boolean;
  title?: string;
  biblioIDAlias?: string;
};

export default function setUpDataset({
  prepareForPublishing = false,
  title = "The dataset title",
  biblioIDAlias = "biblioId",
}: SetUpDatasetOptions = {}): void {
  logCommand("setUpDataset", {
    "Prepare for publishing": prepareForPublishing,
    title,
    "Biblio ID alias": biblioIDAlias,
  });

  cy.htmxRequest<string>({
    method: "POST",
    url: "/add-dataset",
    form: true,
    body: { method: "manual" },
  })
    .then(extractBiblioId)
    .as(biblioIDAlias, { type: "static" })
    .then((biblioId) => {
      // Load the edit form to retrieve the snapshot ID
      cy.htmxRequest<string>({
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
    });

  if (prepareForPublishing) {
    cy.addCreator("John", "Doe");
  }
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
