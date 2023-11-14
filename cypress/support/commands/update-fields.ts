import { logCommand, updateConsoleProps } from './helpers'

const NO_LOG = { log: false }

const SECTIONS = [
  {
    tab: 'Description',
    sectionTitle: 'Publication details',
    buttonText: 'Edit',
    modalTitle: 'Edit publication details',
  },
  {
    tab: 'Description',
    sectionTitle: 'Dataset details',
    buttonText: 'Edit',
    modalTitle: 'Edit dataset details',
  },
  {
    tab: 'Description',
    sectionTitle: 'Conference details',
    buttonText: 'Edit',
    modalTitle: 'Edit conference details',
  },
  {
    tab: 'Description',
    sectionTitle: 'Abstract',
    buttonText: 'Add abstract',
    modalTitle: 'Add abstract',
  },
  {
    tab: 'Description',
    sectionTitle: 'Link',
    buttonText: 'Add link',
    modalTitle: 'Add link',
  },
  {
    tab: 'Description',
    sectionTitle: 'Lay summary',
    buttonText: 'Add lay summary',
    modalTitle: 'Add lay summary',
  },
  {
    tab: 'Description',
    sectionTitle: 'Additional information',
    buttonText: 'Edit',
    modalTitle: 'Edit additional information',
  },

  {
    tab: 'People & Affiliations',
    sectionTitle: 'Authors',
    buttonText: 'Add author',
    modalTitle: 'Add author',
  },
  {
    tab: 'People & Affiliations',
    sectionTitle: 'Editors',
    buttonText: 'Add editor',
    modalTitle: 'Add editor',
  },
  {
    tab: 'People & Affiliations',
    sectionTitle: 'Supervisors',
    buttonText: 'Add supervisor',
    modalTitle: 'Add supervisor',
  },
  {
    tab: 'People & Affiliations',
    sectionTitle: 'Creators',
    buttonText: 'Add creator',
    modalTitle: 'Add creator',
  },

  {
    tab: 'Biblio Messages',
    sectionTitle: 'Librarian tags',
    buttonText: 'Edit',
    modalTitle: 'Edit librarian tags',
  },
  {
    tab: 'Biblio Messages',
    sectionTitle: 'Librarian note',
    buttonText: 'Edit',
    modalTitle: 'Edit librarian note',
  },
  {
    tab: 'Biblio Messages',
    sectionTitle: 'Messages from and for Biblio team',
    buttonText: 'Edit',
    modalTitle: 'Edit messages from and for Biblio team',
  },
] as const

type FieldsSection = (typeof SECTIONS)[number]['sectionTitle']

export default function updateFields(section: FieldsSection, callback: () => void, persist = false): void {
  const log = logCommand('updateFields', { callback, persist }, section).snapshot('before')

  cy.location('pathname', NO_LOG).then(pathname => {
    if (!pathname.match(/^\/(publication|dataset)\/([A-Z0-9]+|add|add-single\/import\/confirm)$/)) {
      throw new Error('The updateFields command can only be called from the details page of a publication or dataset.')
    }

    const $lockLabel = Cypress.$('.bc-toolbar :contains("Locked")')
    if ($lockLabel.length > 0) {
      throw new Error('The updateFields command can only be called from an unlocked publication or dataset.')
    }
  })

  const SECTION = SECTIONS.find(s => s.sectionTitle === section)
  if (!section) {
    throw new Error(`Unknown fields section "${section}".`)
  }

  updateConsoleProps(log, cp => {
    cp.tab = SECTION.tab
    cp.section = SECTION.sectionTitle
    cp['Button text'] = SECTION.buttonText
    cp['Modal title'] = SECTION.modalTitle
  })

  cy.get('.nav.nav-tabs .nav-item .nav-link.active', NO_LOG).then($navLink => {
    if ($navLink.length > 1) {
      throw new Error('Multiple active ".nav-link" elements found.')
    }

    if ($navLink.text() !== SECTION.tab) {
      cy.contains('.nav-tabs .nav-item', SECTION.tab, NO_LOG).click(NO_LOG)
    }
  })

  cy.contains('.card-header', SECTION.sectionTitle, NO_LOG)
    .find('.btn', NO_LOG)
    .then($button => {
      if (!$button.text().includes(SECTION.buttonText)) {
        expect($button).to.contain.text(SECTION.buttonText)
      }

      return $button
    })
    .click(NO_LOG)

  cy.ensureModal(SECTION.modalTitle, NO_LOG)
    .within(NO_LOG, callback)
    .then(modal => {
      log.snapshot('during')

      return modal
    })
    .closeModal(persist, NO_LOG)

  cy.ensureNoModal(NO_LOG).then(() => log.snapshot('after').finish())
}

declare global {
  namespace Cypress {
    interface Chainable {
      updateFields(section: FieldsSection, callback: () => void, persist?: boolean): Chainable<void>
    }
  }
}
