import { logCommand } from "./helpers";

export default function closeToast(
  subject: JQuery<HTMLElement>,
): Cypress.Chainable<JQuery<HTMLElement>> {
  logCommand("closeToast", { subject });

  if (!subject.is(".toast")) {
    throw new Error("Command subject is not a toast.");
  }

  return cy.wrap(subject, { log: false }).within({ log: false }, () => {
    cy.get(".btn-close", { log: false }).click({ log: false });
  });
}

declare global {
  namespace Cypress {
    interface Chainable<Subject> {
      closeToast(): Chainable<Subject>;
    }
  }
}
