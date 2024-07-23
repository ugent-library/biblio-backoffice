import { logCommand } from "./helpers";

export default function switchMode(mode: "Researcher" | "Librarian"): void {
  const currentMode = mode === "Researcher" ? "Librarian" : "Researcher";

  const log = logCommand(
    "switchMode",
    { "Current mode": currentMode, "New mode": mode },
    mode,
  );
  log.set("type", "parent");

  cy.session(
    mode,
    () => {
      cy.login("librarian1");

      cy.htmxRequest({
        method: "PUT",
        url: "/role/" + (mode === "Librarian" ? "curator" : "user"),
      });
    },
    { cacheAcrossSpecs: true },
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      switchMode(mode: "Researcher" | "Librarian"): Chainable<void>;
    }
  }
}
