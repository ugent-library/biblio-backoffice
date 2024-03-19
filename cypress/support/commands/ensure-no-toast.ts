import { logCommand } from "./helpers";

type EnsureNoToastOptions = {
  timeout?: number;
};

export default function ensureNoToast(
  options: EnsureNoToastOptions = { timeout: 6000 }, // Toast automatically disappear after 5 seconds
): Cypress.Chainable<JQuery<HTMLElement>> {
  logCommand("ensureNoToast", { options });

  const { timeout } = options;

  return cy.get(".toast", { log: false, timeout }).should("not.exist");
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureNoToast(
        options?: EnsureNoToastOptions,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
