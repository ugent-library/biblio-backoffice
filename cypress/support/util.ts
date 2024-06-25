const NO_LOG = { log: false };

export function getRandomText() {
  return crypto.randomUUID().replace(/-/g, "").toUpperCase();
}

export function testFocusForLabel(
  labelText: string,
  fieldSelector: string,
  autoFocus = false,
) {
  cy.getLabel(labelText)
    .as("theLabel")
    .should("have.length", 1)
    .parent(NO_LOG)
    .find(fieldSelector)
    .should("have.length", 1)
    .first({ log: false })
    .as("theField")
    .should(autoFocus ? "be.focused" : "not.be.focused");

  if (autoFocus) {
    cy.focused().blur();

    cy.get("@theField").should("not.be.focused");
  }

  cy.get("@theLabel").click();

  cy.get("@theField").should("be.focused");
}

export function getCSRFToken() {
  return cy.state("ctx").CSRFToken;
}

export function extractHtmxJsonAttribute<T extends object>(
  response: Cypress.Response<string>,
  selector: string,
  hxAttributeName: `hx-${string}`,
): T {
  const $partial = Cypress.$(response.body);
  const $node = $partial.find(selector);

  const json = decodeEntities($node.attr(hxAttributeName));

  return JSON.parse(json) as T;
}

export function extractSnapshotId(response: Cypress.Response<string>): string {
  const hxHeaders = extractHtmxJsonAttribute<{ "If-Match": string }>(
    response,
    ".btn:Contains('Save'):not(:contains('Save and add next'))",
    "hx-headers",
  );

  return hxHeaders["If-Match"];
}

function decodeEntities(encodedString) {
  var textArea = document.createElement("textarea");

  textArea.innerHTML = encodedString;

  return textArea.value;
}

export function waitForIndex(
  scope: "publication" | "dataset",
  biblioId: string,
) {
  let count = 0;

  const doCheckIndex = () => {
    count++;

    const url = new URL(
      `/biblio_${scope}s/_search`,
      Cypress.env("ELASTICSEARCH_ORIGIN"),
    ).toString();

    cy.request({
      url,
      qs: {
        size: 0,
        q: `_id:${biblioId}`,
      },
      ...NO_LOG,
    }).then((r) => {
      if (r.body.hits.total < 1) {
        if (count < 20) {
          cy.wait(100, NO_LOG).then(doCheckIndex);
        } else {
          throw new Error("Timed out waiting for index to be ready");
        }
      }
    });
  };

  doCheckIndex();
}
