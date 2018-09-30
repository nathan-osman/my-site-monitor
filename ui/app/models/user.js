import DS from 'ember-data';

export default DS.Modle.extend({
  username: DS.attr('string'),
  password: DS.attr('string')
});
