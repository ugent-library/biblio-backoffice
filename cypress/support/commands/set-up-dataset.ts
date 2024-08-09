import { waitForIndex, extractSnapshotId } from "support/util";
import { logCommand, updateConsoleProps } from "./helpers";

type SetUpDatasetOptions = {
  title?: string;
  otherFields?: Record<string, unknown>;
  biblioIDAlias?: string;
  prepareForPublishing?: boolean;
  publish?: boolean;
  shouldWaitForIndex?: boolean;
};

export default function setUpDataset({
  title = "The dataset title",
  otherFields = {},
  biblioIDAlias = "biblioId",
  prepareForPublishing = false,
  publish = false,
  shouldWaitForIndex = false,
}: SetUpDatasetOptions = {}): void {
  const log = logCommand("setUpDataset", {
    title,
    "Other fields": otherFields,
    "Biblio ID alias": biblioIDAlias,
    "Prepare for publishing": prepareForPublishing,
    Publish: publish,
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
            title,
            identifier_type: "DOI",
            identifier: "10.7202/1041023ar",
          };

          if (prepareForPublishing || publish) {
            body["year"] = new Date().getFullYear().toString();
            body["publisher"] = "UGent";
            body["format"] = ["text/csv"];
            body["license"] = "CC0-1.0";
            body["access_level"] = "info:eu-repo/semantics/openAccess";
          }

          Object.assign(body, otherFields);

          if (Object.values(body).filter((v) => v !== null).length > 0) {
            cy.htmxRequest({
              method: "PUT",
              url: `/dataset/${biblioId}/details`,
              headers: {
                "If-Match": snapshotId,
              },
              form: true,
              body,
            });
          }
        });

      if (prepareForPublishing || publish) {
        cy.addCreator("John", "Doe", { biblioIDAlias: `@${biblioIDAlias}` });
      }

      if (publish) {
        publishDataset(biblioId);
      }

      if (shouldWaitForIndex) {
        waitForIndex("dataset", biblioId);
      }
    });
}

function publishDataset(biblioId: string) {
  cy.htmxRequest({
    url: `/dataset/${biblioId}/add/confirm`,
  })
    .then((r) =>
      extractSnapshotId(
        r,
        `button[hx-post="/dataset/${biblioId}/add/publish"]:contains("Publish to Biblio")`,
      ),
    )
    .then((snapshotId) => {
      cy.htmxRequest({
        method: "POST",
        url: `/dataset/${biblioId}/add/publish`,
        headers: {
          "If-Match": snapshotId,
        },
      });
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
