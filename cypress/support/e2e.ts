import "cypress-common";

import "./commands";
import "./queries";

Cypress.Keyboard.defaults({ keystrokeDelay: 0 });

beforeEach(() => {
  // Keep dashboard-icon polling from Cypress command log
  cy.intercept({ method: "GET", url: "/dashboard-icon" }, { log: false });
});
