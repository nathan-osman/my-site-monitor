import DS from 'ember-data';
import { computed } from '@ember/object';

export default DS.Model.extend({
  url: DS.attr('string'),
  name: DS.attr('string'),
  pollInterval: DS.attr('number'),
  lastPoll: DS.attr('date'),
  nextPoll: DS.attr('date'),
  status: DS.attr('number'),
  statusTime: DS.attr('date'),
  user: DS.belongsTo('user'),

  isUnknown: computed('status', function() {
    return this.get('status') == 0;
  }),

  isOffline: computed('status', function() {
    return this.get('status') == 2;
  }),

  statusClass: computed('status', function() {
    return {
      1: 'online',
      2: 'offline'
    }[this.get('status')] || 'unknown';
  }),

  statusText: computed('status', function() {
    return {
      1: 'Online',
      2: 'Offline',
    }[this.get('status')] || 'Unknown';
  })
});
