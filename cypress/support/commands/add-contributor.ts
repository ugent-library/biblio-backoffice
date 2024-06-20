import { logCommand } from "./helpers";

const NO_LOG = { log: false };

const CONTRIBUTOR_MAP = {
  author: "Authors",
  editor: "Editors",
  supervisor: "Supervisors",
  creator: "Creators",
} as const;

type ContributorType = keyof typeof CONTRIBUTOR_MAP;

export function addAuthor(
  firstName: string,
  lastName: string,
  external = false,
  role?: string,
): void {
  addContributor("author", firstName, lastName, external, role);
}

export function addEditor(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("editor", firstName, lastName, external);
}

export function addSupervisor(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("supervisor", firstName, lastName, external);
}

export function addCreator(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("creator", firstName, lastName, external);
}

function addContributor(
  contributorType: ContributorType,
  firstName: string,
  lastName: string,
  external: boolean,
  role?: string,
): void {
  const consoleProps = {
    "Contributor type": contributorType,
    "First name": firstName,
    "Last name": lastName,
    "External contributor": external,
  };

  if (contributorType === "author") {
    consoleProps["Role"] = role;
  }

  const typeName = Cypress._.capitalize(contributorType);
  logCommand(`add${typeName}`, consoleProps, [
    `${firstName || "[missing]"} ${lastName || "[missing]"} ${external ? "(external)" : ""}`.trim(),
  ]);

  cy.updateFields(
    CONTRIBUTOR_MAP[contributorType],
    function () {
      const pathname = `/+(publication|dataset)/*/contributors/${contributorType === "creator" ? "author" : contributorType}/suggestions`;

      if (firstName) {
        cy.intercept(
          { pathname, query: { first_name: firstName, last_name: "" } },
          NO_LOG,
        ).as(`suggest${typeName}`);

        cy.setFieldByLabel("First name", firstName);
        cy.wait(`@suggest${typeName}`, NO_LOG);
      }

      if (lastName) {
        // Redefine intercept, now including last name
        cy.intercept(
          { pathname, query: { first_name: firstName, last_name: lastName } },
          NO_LOG,
        ).as(`suggest${typeName}`);

        cy.setFieldByLabel("Last name", lastName);
        cy.wait(`@suggest${typeName}`, NO_LOG);
      }

      cy.intercept(
        `/+(publication|dataset)/*/contributors/${contributorType === "creator" ? "author" : contributorType}/confirm-create*`,
        NO_LOG,
      ).as(`confirmCreate${typeName}`);
      cy.get(
        `.btn:contains("Add ${external ? "external " : ""}${contributorType}")`,
        NO_LOG,
      ).click(NO_LOG);

      cy.wait(`@confirmCreate${typeName}`, NO_LOG);
    },
    true,
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      addAuthor(
        firstName: string,
        lastName: string,
        external?: boolean,
        role?: string,
      ): Chainable<void>;

      addEditor(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;

      addSupervisor(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;

      addCreator(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;
    }
  }
}
