import Route from '@ember/routing/route';
import moment from 'moment';

export default Route.extend({
  moment: Ember.inject.service(),

  beforeModel() {
    moment.relativeTimeThreshold('ss', 2);
  }
});
