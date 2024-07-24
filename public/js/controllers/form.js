import { Controller } from 'https://unpkg.com/@hotwired/stimulus@3.2.2/dist/stimulus.js'

export default class extends Controller {
  static targets = [ 'required', 'submit' ];

  connect() {
    this.requiredTargets.forEach((element) => {
      element.addEventListener('input', this.validate.bind(this));
    });
    this.validate();
  }

  disconnect() {
    this.requiredTargets.forEach((element) => {
      element.removeEventListener('input', this.validate.bind(this));
    });
  }

  validate() {
    this.submitTarget.disabled = !this.#valid;
  }

  get #valid() {
    return this.requiredTargets.every((element) => element.value.trim() !== '');
  }
}
