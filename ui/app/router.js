import EmberRouter from '@ember/routing/router';
import config from './config/environment';

const Router = EmberRouter.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('login');
  this.route('sites', function() {
    this.route('new');
    this.route('view', {path: ':site_id'});
  });
});

export default Router;
