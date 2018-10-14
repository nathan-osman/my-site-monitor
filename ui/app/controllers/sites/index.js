import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    delete(site) {
      // TODO: check for error
      site.destroyRecord();
    }
  }
});
