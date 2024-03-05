import { logCommand } from "./helpers";

export default function switchMode(mode: "Researcher" | "Librarian"): void {
  const currentMode = mode === "Researcher" ? "Librarian" : "Researcher";

  let log: Cypress.Log;

  cy.visit("/", { log: false });

  cy.intercept({ method: "PUT", url: "/role/*" }, { log: false }).as(
    "switch-role",
  );

  cy.contains(`.c-sidebar > .dropdown > button`, currentMode, { log: false })
    .click({ log: false })
    .next(".dropdown-menu", { log: false })
    .contains(mode, { log: false })
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
    .click({ log: false });

  cy.wait("@switch-role", { log: false });

  cy.contains(`.c-sidebar > .dropdown > button`, mode, { log: false }).then(
    ($el) => {
      log.set("$el", $el).snapshot("after").end();
    },
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      switchMode(mode: "Researcher" | "Librarian"): Chainable<void>;
    }
  }
}
