import DS from 'ember-data';

export default DS.Model.extend({
  startTime: DS.attr('date'),
  endTime: DS.attr('date'),
  description: DS.attr('string'),
  site: DS.belongsTo('site')
});
