import Controller from '@ember/controller';

export default Controller.extend({
  auth: Ember.inject.service(),

  actions: {
    logout() {
      this.get('auth').logout();
    }
  }
});
