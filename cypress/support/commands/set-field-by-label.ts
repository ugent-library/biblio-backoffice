import { logCommand, updateConsoleProps } from "./helpers";

const NO_LOG = { log: false };

export default function setFieldByLabel(
  labelCaption: string | RegExp,
  value: string,
): Cypress.Chainable<JQuery<HTMLElement>> {
  const log = logCommand(
    "setFieldByLabel",
    { "Label caption": labelCaption, value },
    `${labelCaption} = ${value}`,
  ).snapshot("before");

  cy.getLabel(labelCaption, NO_LOG)
    .then((label) => {
      updateConsoleProps(log, (cp) => (cp["Label element"] = label.get(0)));

      return label;
    })
    .click(NO_LOG);

  return cy
    .focused(NO_LOG)
    .then((field) => {
      updateConsoleProps(log, (cp) => {
        cp["Field element"] = field.get(0);
        cp["Old value"] = field.val();
      });

      return field;
    })
    .setField(value, NO_LOG)
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
