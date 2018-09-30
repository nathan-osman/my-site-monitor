import Route from '@ember/routing/route';

export default Route.extend({
  model(params) {
    return this.store.query('site', {id: params.site_id});
  }
});
