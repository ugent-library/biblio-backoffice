import { logCommand } from "./helpers";

const NO_LOG = { log: false };

type EnsureModalOptions = {
  log?: boolean;
};

export default function ensureModal(
  expectedTitle: string | RegExp | null,
  options: EnsureModalOptions = { log: true },
): Cypress.Chainable<JQuery<HTMLElement>> {
  let log: Cypress.Log | null = null;
  if (options.log === true) {
    log = logCommand(
      "ensureModal",
      { "Expected title": expectedTitle },
      expectedTitle,
    );
  }

  // Assertion "be.visible" doesn't work here because it is behind the dialog
  cy.get("#modal-backdrop", NO_LOG).then((modalBackdrop) => {
    if (!modalBackdrop.get(0).classList.contains("show")) {
      cy.wrap(modalBackdrop, NO_LOG).should("have.class", "show");
    }
  });

  return cy
    .get("#modal", NO_LOG)
    .should("be.visible")
    .within(NO_LOG, () => {
      if (expectedTitle === null) {
        cy.get(".modal-title", NO_LOG).should("not.exist");
      } else if (expectedTitle instanceof RegExp) {
        cy.get(".modal-title", NO_LOG)
          .invoke(NO_LOG, "text")
          .should("match", expectedTitle);
      } else {
        cy.get(".modal-title", NO_LOG).should("have.text", expectedTitle);
      }
    })
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureModal(
        expectedTitle: string | RegExp | null,
        options?: EnsureModalOptions,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
