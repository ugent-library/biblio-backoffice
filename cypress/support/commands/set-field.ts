import { logCommand } from "./helpers";

type SetFieldOptions = {
  log?: boolean;
};

export default function setField(
  subject: JQuery<HTMLInputElement | HTMLSelectElement>,
  value: string,
  options: SetFieldOptions = { log: true },
): Cypress.Chainable<JQuery<HTMLInputElement | HTMLSelectElement>> {
  const log =
    options.log === true &&
    logCommand("setField", { subject, value }, value).snapshot("before");

  const field = cy.wrap(subject, { log: false });

  switch (subject.prop("tagName")) {
    case "INPUT":
      field.clear({ log: false }).type(value, { delay: 0, log: false });
      break;

    case "SELECT":
      field.select(value, { log: false });
      break;

    default:
      throw new Error(
        `Field of type '${subject.prop("tagName")}' is not supported.`,
      );
  }

  cy.then(() => log && log.snapshot("after"));

  return field.finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable<Subject> {
      setField(value: string, options?: SetFieldOptions): Chainable<Subject>;
    }
  }
}
