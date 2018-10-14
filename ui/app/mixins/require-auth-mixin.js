import Mixin from '@ember/object/mixin';

export default Mixin.create({
  auth: Ember.inject.service(),

  beforeModel() {
    if (!this.get('auth').get('isAuthenticated')) {
      this.transitionTo('login');
    }
  }
});
