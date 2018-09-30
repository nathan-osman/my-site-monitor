import DS from 'ember-data';

export default DS.Modle.extend({
  url: DS.attr('string'),
  name: DS.attr('string'),
  pollInterval: DS.attr('number'),
  pollTime: DS.attr('date'),
  status: DS.attr('status'),
  statusTime: DS.attr('date'),
  user: DS.belongsTo('user')
});
