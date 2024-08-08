import { logCommand } from "./helpers";

export function loginAsResearcher(username: string): void {
  login("loginAsResearcher", username, null);
}

export function loginAsLibrarian(
  username: string,
  reviewerMode = "Librarian",
): void {
  login("loginAsLibrarian", username, reviewerMode);
}

function login(
  commandName: string,
  username: string,
  reviewerMode?: string,
): void {
  // WARNING: Whenever you change the code of the session setup, Cypress will throw an error:
  //   This session already exists. You may not create a new session with a previously used identifier.
  //   If you want to create a new session with a different setup function, please call cy.session() with
  //   a unique identifier other than...
  //
  // Temporarily uncomment the following line to clear the sessions if this happens
  // Cypress.session.clearAllSavedSessions()
  const sessionName = username + (reviewerMode ? ` (${reviewerMode})` : "");

  const consoleProps = { username };
  if (reviewerMode) {
    consoleProps["Reviewer mode"] = reviewerMode;
  }

  logCommand(commandName, consoleProps, sessionName);

  // First clear the current CSRF token
  cy.state("ctx").CSRFToken = "";

  cy.session(
    sessionName,
    () => {
      cy.request({ url: "/login", log: false }).then((response) => {
        const form = new DOMParser()
          .parseFromString(response.body, "text/html")
          .querySelector("form");

        const body = Object.fromEntries(new FormData(form));
        if ("username" in body) {
          body.username = username;
        }
        if ("password" in body) {
          throw new Error(
            "Current Cypres setup doesn't support password login.",
          );
        }

        cy.request({
          method: "POST",
          url: form.action,
          form: true,
          body,

          // Make sure we redirect back and get the cookie we need
          followRedirect: true,

          // Make sure we don't leak passwords in the Cypress log
          log: false,
        });

        if (reviewerMode) {
          cy.htmxRequest({
            method: "PUT",
            url: "/role/" + (reviewerMode === "Librarian" ? "curator" : "user"),
          });
        }
      });
    },
    {
      cacheAcrossSpecs: true,
    },
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsResearcher(
        username: "researcher1" | "researcher2",
      ): Chainable<void>;

      loginAsLibrarian(
        username: "librarian1" | "librarian2",
        reviewerMode?: "Librarian" | "Researcher",
      ): Chainable<void>;
    }
  }
}
