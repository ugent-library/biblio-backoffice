import { logCommand, updateConsoleProps } from "./helpers";

export default function setFieldByLabel(
  labelCaption: string | RegExp,
  value: string,
): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand(
    "setFieldByLabel",
    { "Label caption": labelCaption, value },
    `${labelCaption} = ${value}`,
  ).snapshot("before");

  cy.getLabel(labelCaption, { log: false })
    .then((label) => {
      updateConsoleProps(log, (cp) => (cp["Label element"] = label.get(0)));

      return label;
    })
    .click({ log: false });

  return cy
    .focused({ log: false })
    .then((field) => {
      updateConsoleProps(log, (cp) => {
        cp["Field element"] = field.get(0);
        cp["Old value"] = field.val();
      });

      return field;
    })
    .setField(value, { log: false })
    .then((field) => {
      log.snapshot("after");

      return field;
    });
}

declare global {
  namespace Cypress {
    interface Chainable {
      setFieldByLabel(
        fieldLabel: string | RegExp,
        value: string,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
