// Code adapted from @4tw/cypress-drag-drop plugin

import { logCommand } from "./helpers";

const NO_LOG = { log: false };

const dataTransfer = new DataTransfer();

class DragSimulator {
  MAX_TRIES = 5;
  DELAY_INTERVAL_MS = 10;

  source: HTMLElement = null;
  initialSourcePosition: DOMRect = null;
  counter = 0;
  targetElement: HTMLElement = null;
  log: Cypress.Log = null;

  rectsEqual(r1: DOMRect, r2: DOMRect) {
    return (
      r1.top === r2.top &&
      r1.right === r2.right &&
      r1.bottom === r2.bottom &&
      r1.left === r2.left
    );
  }

  getTarget() {
    return cy.wrap(this.targetElement, NO_LOG);
  }

  isDropped() {
    return !this.rectsEqual(
      this.initialSourcePosition,
      this.source.getBoundingClientRect(),
    );
  }

  hasTriesLeft() {
    return this.counter < this.MAX_TRIES;
  }

  drag(subject: JQuery<HTMLElement>, targetSelector: string) {
    return this.init(subject, targetSelector)
      .then(() => this.dragstart())
      .then(() => this.dragover())
      .then((success) => {
        if (success) {
          return this.drop().then(() => true);
        } else {
          return cy.wrap(false, NO_LOG);
        }
      })
      .then((result) => {
        this.log.snapshot("after");

        return result;
      })
      .finishLog(this.log);
  }

  init($source: JQuery<HTMLElement>, target: string) {
    this.counter = 0;
    this.source = $source.get(0);
    this.initialSourcePosition = this.source.getBoundingClientRect();

    return cy.get(target, NO_LOG).then(($target) => {
      this.targetElement = $target.get(0);

      this.log = logCommand(
        "drag",
        {
          subject: this.source,
          "Target selector": target,
          target: this.targetElement,
        },
        this.source,
        $source,
      );

      this.log.snapshot("before");
    });
  }

  dragstart(clientPosition = {}) {
    return cy
      .wrap(this.source, NO_LOG)
      .trigger("pointerdown", {
        which: 1,
        button: 0,
        ...clientPosition,
        eventConstructor: "PointerEvent",
        ...NO_LOG,
      })
      .trigger("dragstart", {
        dataTransfer,
        eventConstructor: "DragEvent",
        scrollBehavior: "center",
        ...NO_LOG,
      });
  }

  dragover(): Cypress.Chainable<boolean> {
    if (!this.counter || (!this.isDropped() && this.hasTriesLeft())) {
      this.counter += 1;
      return this.getTarget()
        .trigger("dragover", {
          dataTransfer,
          eventConstructor: "DragEvent",
          ...NO_LOG,
        })
        .wait(this.DELAY_INTERVAL_MS, NO_LOG)
        .then(() => this.dragover());
    }

    if (!this.isDropped()) {
      console.error(`Exceeded maximum tries of: ${this.MAX_TRIES}, aborting`);
      return cy.wrap(false, NO_LOG);
    } else {
      return cy.wrap(true, NO_LOG);
    }
  }

  drop() {
    return this.getTarget().trigger("drop", {
      dataTransfer,
      eventConstructor: "DragEvent",
      ...NO_LOG,
    });
  }
}

export default function drag(
  subject: JQuery<HTMLElement>,
  targetSelector: string,
) {
  const dragSimulator = new DragSimulator();

  return dragSimulator
    .init(subject, targetSelector)
    .then(() => dragSimulator.dragstart())
    .then(() => dragSimulator.dragover())
    .then((success) => {
      if (success) {
        return dragSimulator.drop().then(() => true);
      } else {
        return cy.wrap(false, NO_LOG);
      }
    })
    .then((result) => {
      dragSimulator.log.snapshot("after");

      return result;
    })
    .finishLog(dragSimulator.log);
}

declare global {
  namespace Cypress {
    interface Chainable<Subject> {
      drag(targetSelector: string): Chainable<boolean>;
    }
  }
}
