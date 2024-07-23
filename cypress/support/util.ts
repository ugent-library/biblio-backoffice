const NO_LOG = { log: false };
const ACTUAL_FIELDS_SELECTOR =
  "input:not([type=hidden], [type=submit], [type=reset], [type=button], [type=image]), textarea, select";

export function getRandomText() {
  return crypto.randomUUID().replace(/-/g, "").toUpperCase();
}

export function testFormAccessibility(
  form: Record<string, string | RegExp>,
  autoFocusLabel?: string,
  ignoredFields: string[] = [],
) {
  testFocusForLabels(ignoredFields, form, autoFocusLabel);

  testAriaDescriptionAttributes();
}

function testFocusForLabels(
  ignoredFields: string[],
  form: Record<string, string | RegExp>,
  autoFocusLabel: string,
) {
  cy.get<HTMLFormElement>(ACTUAL_FIELDS_SELECTOR, NO_LOG).then(
    (allFormFields) => {
      ignoredFields.forEach((f) => (allFormFields = allFormFields.not(f)));

      for (const selector of Object.keys(form)) {
        testFocusForLabel(
          form[selector],
          selector,
          autoFocusLabel === form[selector],
        );

        allFormFields = allFormFields.not(selector);
      }

      expect(allFormFields.length).to.eq(
        0,
        "Not all fields were checked: " +
          allFormFields
            .get()
            .map((f) => `${f.tagName.toLowerCase()}[name=${f.name}]`)
            .join(", "),
      );
    },
  );
}

function testAriaDescriptionAttributes() {
  cy.root(NO_LOG).then(($context) => {
    // Add 2 extra fields for each repeated form value
    $context.find("button.form-value-add").trigger("click");
    $context.find("button.form-value-add").trigger("click");

    $context.find(".form-text").each((_, formText) => {
      const $fields = $context
        .find(formText)
        .parent()
        .find(ACTUAL_FIELDS_SELECTOR);
      expect($fields).to.have.length.above(0);

      $fields.each((_, field) => {
        expect(field).to.satisfy(
          (f: HTMLElement) =>
            // aria-details references id attribute of .form-text element
            f.getAttribute("aria-details") === formText.id ||
            // aria-describedby references id attribute of .form-text element
            f.getAttribute("aria-describedby") === formText.id ||
            // aria-description contains the same text content as .form-text element
            f.getAttribute("aria-description") === formText.innerText,
          "Field is missing accessibility information referencing its .form-text node.",
        );
      });
    });
  });
}

function testFocusForLabel(
  labelText: string | RegExp,
  fieldSelector: string,
  autoFocus = false,
) {
  cy.getLabel(labelText)
    .as("theLabel")
    .then(($label) => {
      if ($label.length !== 1) {
        expect($label).to.have.lengthOf(1);
      }
    })
    .parent(NO_LOG)
    .find(fieldSelector, NO_LOG)
    .then(($field) => {
      if ($field.length !== 1) {
        expect($field).to.have.lengthOf(1);
      }
    })
    .first(NO_LOG)
    .as("theField")
    .should(autoFocus ? "be.focused" : "not.be.focused");

  if (autoFocus) {
    cy.focused(NO_LOG).blur();
    cy.get("@theField", NO_LOG).should("not.be.focused");
  }

  cy.get("@theLabel", NO_LOG).click();
  cy.get("@theField", NO_LOG).should("be.focused");
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

  if ($node.length === 0) {
    return null;
  }

  const hxAttribute = $node.attr(hxAttributeName);
  return JSON.parse(decodeEntities(hxAttribute)) as T;
}

export function extractSnapshotId(
  response: Cypress.Response<string>,
  selector = ".btn:Contains('Save'):not(:contains('Save and add next'))",
): string {
  const hxHeaders = extractHtmxJsonAttribute<{ "If-Match": string }>(
    response,
    selector,
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
