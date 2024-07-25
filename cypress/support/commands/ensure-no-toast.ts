import { logCommand } from "./helpers";

type EnsureNoToastOptions = {
  timeout?: number;
};

export default function ensureNoToast(
  options: EnsureNoToastOptions = { timeout: 6000 }, // Toast automatically disappear after 5 seconds
): Cypress.Chainable<JQuery<HTMLElement>> {
  logCommand("ensureNoToast", { options });

  const { timeout } = options;

  return cy.then(() => {
    // First check if there are any toasts and assert accordingly.
    // If not, either cy.get() or the chained assertion will fail.
    const $toasts = Cypress.$(".toast");
    if ($toasts.length > 0) {
      // expect(...) will not work here as that assertion would fail if toast is still hiding.
      // Cypress retry-ability fixes that.
      cy.get(".toast", { log: false, timeout }).should("not.be.visible");
    } else {
      cy.get(".toast", { log: false, timeout }).should("not.exist");
    }
  });
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
