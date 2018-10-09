import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    create() {
      this.store.createRecord('site', {
        url: this.get('url'),
        name: this.get('name'),
        pollInterval: this.get('pollInterval')
      }).save().then(
        () => { this.transitionToRoute('index'); }
      );
    }
  }
});
