import Route from '@ember/routing/route';
import moment from 'moment';

export default Route.extend({
  auth: Ember.inject.service(),

  beforeModel() {
    moment.relativeTimeThreshold('ss', 2);
  },

  model() {
    return this.get('auth').restore();
  }
});
