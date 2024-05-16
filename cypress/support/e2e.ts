import "cypress-common";

import "./commands";
import "./queries";

Cypress.Keyboard.defaults({ keystrokeDelay: 0 });

beforeEach(() => {
  // Store latest CSRFToken whenever we find one
  cy.on("window:load", function (win) {
    const metaTag = Cypress.$("meta[name='csrf-token']");
    if (metaTag.length > 0) {
      // Cannot easily store aliased value from an event handler so we keep it on the test context
      const ctx = cy.state("ctx");
      ctx["CSRFToken"] = metaTag.prop("content");
    }
  });

  // Keep dashboard-icon polling from Cypress command log
  cy.intercept({ method: "GET", url: "/dashboard-icon" }, { log: false });
  cy.intercept(
    { method: "GET", url: "/candidate-records-icon" },
    { log: false },
  );
});
