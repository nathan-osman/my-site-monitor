import DS from 'ember-data';

export default DS.Modle.extend({
  startTime: DS.attr('date'),
  endTime: DS.attr('date'),
  description: DS.attr('string'),
  site: DS.belongsTo('site')
});
