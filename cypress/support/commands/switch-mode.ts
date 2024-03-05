import { logCommand } from "./helpers";

const NO_LOG = { log: false };

export default function switchMode(mode: "Researcher" | "Librarian"): void {
  const currentMode = mode === "Researcher" ? "Librarian" : "Researcher";

  let log: Cypress.Log;

  cy.visit("/", NO_LOG);

  cy.intercept({ method: "PUT", url: "/role/*" }, NO_LOG).as("switch-role");

  cy.contains(`.c-sidebar > .dropdown > button`, currentMode, NO_LOG)
    .click(NO_LOG)
    .next(".dropdown-menu", NO_LOG)
    .contains(".dropdown-item", mode, NO_LOG)
    .then(($el) => {
      log = logCommand(
        "switchMode",
        { "Current mode": currentMode, "New mode": mode },
        mode,
        $el,
      );
      log.set("type", "parent");
      log.snapshot("before");
    })
    .click(NO_LOG);

  cy.wait("@switch-role", NO_LOG);

  cy.contains(`.c-sidebar > .dropdown > button`, mode, NO_LOG).then(($el) => {
    log.set("$el", $el).snapshot("after").end();
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      switchMode(mode: "Researcher" | "Librarian"): Chainable<void>;
    }
  }
}
