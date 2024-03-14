import { logCommand } from "./helpers";

export default function clocked(callback: () => void) {
  const log = logCommand("clocked");
  log.set("type", "parent");

  cy.clock({ log: false });

  callback();

  cy.then(() => {
    this.clock.restore();
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      clocked(callback: () => void): Chainable<void>;
    }
  }
}
