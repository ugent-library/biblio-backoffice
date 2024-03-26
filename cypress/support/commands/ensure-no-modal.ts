import { logCommand } from "./helpers";

type EnsureNoModalOptions = {
  log?: boolean;
};

export default function ensureNoModal(
  options: EnsureNoModalOptions = { log: true },
): void {
  if (options.log === true) {
    logCommand("ensureNoModal");
  }

  cy.get("#modals > *, #modal, #modal-backdrop", { log: false }).should(
    "not.exist",
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureNoModal(options?: EnsureNoModalOptions): Chainable<void>;
    }
  }
}
