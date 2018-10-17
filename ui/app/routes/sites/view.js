import Route from '@ember/routing/route';

export default Route.extend({
  model(params) {
    return Ember.RSVP.hash({
      site: this.store.findRecord('site', params.site_id),
      outages: this.store.query('outage', {site_id: params.site_id})
    });
  }
});
