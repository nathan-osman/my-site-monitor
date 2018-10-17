import DS from 'ember-data';
import { computed } from '@ember/object';
import moment from 'moment';

export default DS.Model.extend({
  startTime: DS.attr('date'),
  endTime: DS.attr('date'),
  description: DS.attr('string'),
  site: DS.belongsTo('site'),

  duration: computed('startTime', 'endTime', function() {
    let startTime = moment(this.get('startTime'));
    return startTime.diff(this.get('endTime') || moment());
  })
});
