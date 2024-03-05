import { logCommand } from "./helpers";

export default function login(username, password): void {
  // WARNING: Whenever you change the code of the session setup, Cypress will throw an error:
  //   This session already exists. You may not create a new session with a previously used identifier.
  //   If you want to create a new session with a different setup function, please call cy.session() with
  //   a unique identifier other than...
  //
  // Temporarily uncomment the following line to clear the sessions if this happens
  // Cypress.session.clearAllSavedSessions()

  logCommand("login", { username }, username);

  cy.session(
    username,
    () => {
      cy.request("/login", { log: false })
        .then((response) => {
          const action = response.body.match(/action\=\"(.*)\" /)[1];

          return action.replace(/&amp;/g, "&");
        })
        .then((actionUrl) =>
          cy.request({
            method: "POST",
            url: actionUrl,
            form: true,

            body: {
              username,
              password,
            },

            // Make sure we redirect back and get the cookie we need
            followRedirect: true,

            // Make sure we don't leak passwords in the Cypress log
            log: false,
          }),
        );
    },
    {
      cacheAcrossSpecs: true,
    },
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      login(username: string, password: string): Chainable<void>;
    }
  }
}
