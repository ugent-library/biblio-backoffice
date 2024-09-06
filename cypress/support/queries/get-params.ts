import { logCommand, updateConsoleProps } from "support/commands/helpers";

export default function (...names: string[]) {
  const log = logCommand("getParams", getInitialConsoleProps(names), names);

  const urlFn = cy.now("url", { log: false }) as () => string;

  return () => {
    const result = getParamsResult(urlFn(), names);

    if (cy.state("current") === this) {
      updateConsoleProps(log, (cp) => {
        cp.yielded = result;
      });
    }

    return result;
  };
}

function getInitialConsoleProps(names: string[]): Cypress.ObjectLike {
  switch (names.length) {
    case 0:
      return {};

    case 1: {
      return { name: names.at(0) };
    }

    default:
      return { names };
  }
}

function getParamsResult(url: string, names: string[]) {
  const params = getParamsObject(url);

  switch (names.length) {
    case 0:
      return params;

    case 1:
      return params[names.at(0)];

    default:
      return Cypress._.pick(params, ...names);
  }
}

function getParamsObject(url: string): Record<string, string | string[]> {
  const { searchParams } = new URL(url);

  const init: Record<string, string | string[]> = {};

  return Array.from(searchParams).reduce((previous, [key, value]) => {
    if (key in previous) {
      if (!Array.isArray(previous[key])) {
        previous[key] = [previous[key]];
      }

      previous[key].push(value);
    } else {
      previous[key] = value;
    }

    return previous;
  }, init);
}

declare global {
  namespace Cypress {
    interface Chainable {
      getParams(...names: string[]): Chainable<Record<string, string>>;
      getParams(name: string): Chainable<string>;
      getParams(): Chainable<Record<string, string>>;
    }
  }
}
