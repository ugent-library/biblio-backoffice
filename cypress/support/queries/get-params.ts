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

  return [...searchParams].reduce((previous, [name, value]) => {
    if (name in previous) {
      if (!Array.isArray(previous[name])) {
        previous[name] = [previous[name]];
      }

      previous[name].push(value);
    } else {
      previous[name] = value;
    }

    return previous;
  }, {});
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
